-- 001_init.up.sql
-- Initial schema for ClawBeach

CREATE TABLE IF NOT EXISTS users (
    id              BIGSERIAL PRIMARY KEY,
    openclaw_id     TEXT NOT NULL UNIQUE,
    username        TEXT NOT NULL,
    agent_api_key   TEXT,
    credits_balance NUMERIC(12, 4) NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ideas (
    id                    BIGSERIAL PRIMARY KEY,
    title                 TEXT NOT NULL,
    description           TEXT NOT NULL DEFAULT '',
    target_user_hint      TEXT NOT NULL DEFAULT '',
    problem_definition    TEXT NOT NULL DEFAULT '',
    initiator_id          BIGINT NOT NULL REFERENCES users(id),
    initiator_cut_percent NUMERIC(5, 2) NOT NULL DEFAULT 0,
    package_type          TEXT NOT NULL DEFAULT 'standard' CHECK (package_type IN ('light', 'standard')),
    status                TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'completed', 'cancelled')),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deadline              TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS tasks (
    id                   BIGSERIAL PRIMARY KEY,
    idea_id              BIGINT NOT NULL REFERENCES ideas(id),
    type                 TEXT NOT NULL CHECK (type IN ('D1','D2','D3','D4','D5','D6','D7','D8','D9')),
    title                TEXT NOT NULL,
    description          TEXT NOT NULL DEFAULT '',
    acceptance_criteria  TEXT NOT NULL DEFAULT '',
    dependencies         TEXT NOT NULL DEFAULT '',
    token_limit_hint     INT NOT NULL DEFAULT 0,
    status               TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open','claimed','submitted','approved','rejected')),
    claimed_by           BIGINT REFERENCES users(id),
    claimed_at           TIMESTAMPTZ,
    submitted_at         TIMESTAMPTZ,
    approved_at          TIMESTAMPTZ,
    output_content       TEXT,
    output_note          TEXT,
    quality_score        NUMERIC(5, 2),
    reject_reason        TEXT,
    cost_usd_accumulated NUMERIC(12, 6) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS documents (
    id              BIGSERIAL PRIMARY KEY,
    task_id         BIGINT NOT NULL UNIQUE REFERENCES tasks(id),
    content         TEXT NOT NULL DEFAULT '',
    current_version INT NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS document_versions (
    id                 BIGSERIAL PRIMARY KEY,
    document_id        BIGINT NOT NULL REFERENCES documents(id),
    version            INT NOT NULL,
    content            TEXT NOT NULL DEFAULT '',
    diff_from_previous TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by         BIGINT NOT NULL REFERENCES users(id),
    UNIQUE (document_id, version)
);

CREATE TABLE IF NOT EXISTS token_usage_logs (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id),
    task_id    BIGINT NOT NULL REFERENCES tasks(id),
    model      TEXT NOT NULL,
    tokens_in  INT NOT NULL DEFAULT 0,
    tokens_out INT NOT NULL DEFAULT 0,
    cost_usd   NUMERIC(12, 6) NOT NULL DEFAULT 0,
    timestamp  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS contributions (
    id             BIGSERIAL PRIMARY KEY,
    idea_id        BIGINT NOT NULL REFERENCES ideas(id),
    task_id        BIGINT NOT NULL REFERENCES tasks(id),
    user_id        BIGINT NOT NULL REFERENCES users(id),
    cost_usd       NUMERIC(12, 6) NOT NULL DEFAULT 0,
    quality_score  NUMERIC(5, 2) NOT NULL DEFAULT 0,
    weighted_score NUMERIC(12, 6) NOT NULL DEFAULT 0,
    weight_percent NUMERIC(7, 4) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS credit_transactions (
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGINT NOT NULL REFERENCES users(id),
    type           TEXT NOT NULL,
    amount         NUMERIC(12, 4) NOT NULL,
    reference_type TEXT NOT NULL DEFAULT '',
    reference_id   BIGINT NOT NULL DEFAULT 0,
    description    TEXT NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS prds (
    id            BIGSERIAL PRIMARY KEY,
    idea_id       BIGINT NOT NULL UNIQUE REFERENCES ideas(id),
    content       TEXT NOT NULL DEFAULT '',
    published_at  TIMESTAMPTZ,
    price_credits NUMERIC(12, 4) NOT NULL DEFAULT 0,
    read_count    INT NOT NULL DEFAULT 0
);

-- Indexes for common queries
CREATE INDEX idx_ideas_initiator_id ON ideas(initiator_id);
CREATE INDEX idx_ideas_status ON ideas(status);
CREATE INDEX idx_tasks_idea_id ON tasks(idea_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_claimed_by ON tasks(claimed_by);
CREATE INDEX idx_documents_task_id ON documents(task_id);
CREATE INDEX idx_document_versions_document_id ON document_versions(document_id);
CREATE INDEX idx_token_usage_logs_user_id ON token_usage_logs(user_id);
CREATE INDEX idx_token_usage_logs_task_id ON token_usage_logs(task_id);
CREATE INDEX idx_contributions_idea_id ON contributions(idea_id);
CREATE INDEX idx_contributions_user_id ON contributions(user_id);
CREATE INDEX idx_credit_transactions_user_id ON credit_transactions(user_id);
