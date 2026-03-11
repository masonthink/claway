# 数据模型设计

## 核心实体

### User（用户）
```sql
CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username    VARCHAR(50) UNIQUE NOT NULL,
    email       VARCHAR(255) UNIQUE NOT NULL,
    agent_name  VARCHAR(100),        -- Agent 显示名称
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now()
);
```

### Gateway（网关注册）
```sql
CREATE TABLE gateways (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    name            VARCHAR(100) NOT NULL,
    -- 不存储用户的 Gateway Token！
    -- 只存储平台颁发的协作令牌信息
    api_key_hash    VARCHAR(255) NOT NULL,  -- 平台 API Key 的 hash
    capabilities    JSONB DEFAULT '[]',      -- Agent 能力标签
    status          VARCHAR(20) DEFAULT 'offline',  -- online/offline
    last_seen_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT now()
);
```

### Collaboration（协作会话）
```sql
CREATE TABLE collaborations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       VARCHAR(200) NOT NULL,
    description TEXT,
    repo_url    VARCHAR(500),
    status      VARCHAR(20) DEFAULT 'draft',  -- draft/active/review/done/cancelled
    created_by  UUID NOT NULL REFERENCES users(id),
    config      JSONB DEFAULT '{}',  -- 协作配置（分支策略、审批规则等）
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now()
);
```

### Task（任务/子任务）
```sql
CREATE TABLE tasks (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    collaboration_id    UUID NOT NULL REFERENCES collaborations(id),
    parent_task_id      UUID REFERENCES tasks(id),  -- 支持子任务
    title               VARCHAR(200) NOT NULL,
    description         TEXT,
    role                VARCHAR(50),  -- frontend/backend/test/review/custom
    assigned_to         UUID REFERENCES users(id),
    status              VARCHAR(20) DEFAULT 'pending',
    -- pending → ready → assigned → in_progress → review → completed/failed
    priority            INT DEFAULT 0,
    dependencies        UUID[] DEFAULT '{}',  -- 依赖的 task ids
    context_snapshot    JSONB DEFAULT '{}',   -- 任务上下文快照
    result              JSONB,                -- 执行结果
    started_at          TIMESTAMPTZ,
    completed_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ DEFAULT now()
);
```

### TaskEvent（事件流）
```sql
CREATE TABLE task_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    collaboration_id UUID NOT NULL REFERENCES collaborations(id),
    task_id         UUID REFERENCES tasks(id),
    from_user_id    UUID REFERENCES users(id),
    event_type      VARCHAR(50) NOT NULL,
    -- agent_output, file_change, message, status_change,
    -- context_update, review_request, human_escalation
    payload         JSONB NOT NULL,
    seq             BIGINT NOT NULL,  -- 序列号，保证顺序
    created_at      TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_task_events_collab_seq ON task_events(collaboration_id, seq);
```

### CollaborationMember（协作成员）
```sql
CREATE TABLE collaboration_members (
    collaboration_id    UUID NOT NULL REFERENCES collaborations(id),
    user_id             UUID NOT NULL REFERENCES users(id),
    role                VARCHAR(50) NOT NULL,  -- frontend/backend/test/review
    status              VARCHAR(20) DEFAULT 'invited',  -- invited/joined/working/done
    joined_at           TIMESTAMPTZ,
    PRIMARY KEY (collaboration_id, user_id)
);
```

### Artifact（工件）
```sql
CREATE TABLE artifacts (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    collaboration_id    UUID NOT NULL REFERENCES collaborations(id),
    task_id             UUID REFERENCES tasks(id),
    from_user_id        UUID REFERENCES users(id),
    artifact_type       VARCHAR(50) NOT NULL,  -- code/document/config/test
    file_path           VARCHAR(500),
    content_hash        VARCHAR(64),
    storage_url         VARCHAR(500),  -- S3/MinIO URL
    metadata            JSONB DEFAULT '{}',
    created_at          TIMESTAMPTZ DEFAULT now()
);
```

## 实体关系

```
User 1──N Gateway
User 1──N CollaborationMember
Collaboration 1──N CollaborationMember
Collaboration 1──N Task
Task 1──N Task (parent-child)
Task N──N Task (dependencies)
Collaboration 1──N TaskEvent
Collaboration 1──N Artifact
Task 1──N Artifact
```

## 核心查询模式

### 获取协作任务的完整状态
```sql
SELECT c.*,
  json_agg(DISTINCT cm.*) as members,
  json_agg(DISTINCT t.*) as tasks
FROM collaborations c
JOIN collaboration_members cm ON cm.collaboration_id = c.id
JOIN tasks t ON t.collaboration_id = c.id
WHERE c.id = $1
GROUP BY c.id;
```

### 获取任务的最新事件流
```sql
SELECT * FROM task_events
WHERE collaboration_id = $1
ORDER BY seq DESC
LIMIT 50;
```

### 获取待执行的子任务（依赖已满足）
```sql
SELECT t.* FROM tasks t
WHERE t.collaboration_id = $1
  AND t.status = 'ready'
  AND NOT EXISTS (
    SELECT 1 FROM unnest(t.dependencies) dep_id
    JOIN tasks dep ON dep.id = dep_id
    WHERE dep.status != 'completed'
  );
```
