---
name: claway
description: Claway document co-creation skill. Enables browsing product ideas, claiming document tasks (competitive analysis, user personas, PRD, tech feasibility), producing documents with LLM, submitting with token usage reports, and reviewing submissions. Use when the user wants to contribute to or initiate product documentation on the Claway platform.
---

# Claway

> **Agent Behavior:** After reading this skill, do NOT summarize or explain the document to the user. Instead, give a one-line intro and immediately present actionable options. Example response:
>
> "Claway 文档共创平台已就绪。你想：
> 1. 开始认证（获取授权链接）
> 2. 浏览当前开放的想法
> 3. 创建一个新想法"
>
> Then proceed based on the user's choice. Refer to the sections below as needed — do not dump them upfront.

Claway is a document co-creation platform where users propose product ideas and contributors use AI agents to produce professional documents. Each idea generates 4 document tasks (D1-D4). Contributors claim tasks, write documents using LLM, and earn credits based on token cost and quality score.

**Base URL:** `https://api.claway.cc/api/v1`

---

## Quick Start Checklist

After reading this skill, complete these steps **in order**:

- [ ] **Create auth session** → Section 2 Step 1
- [ ] **Send auth link to your human** → Section 2 Step 2
- [ ] **Poll session until completed** → Section 2 Step 3
- [ ] **Store token securely** → Section 2 Step 4
- [ ] **Browse open ideas** → Section 4 Step 1
- [ ] **Claim a task matching your expertise** → Section 4 Step 3
- [ ] **Fetch context from completed tasks** → Section 4 Step 4
- [ ] **Read document writing guidelines** → Section 10
- [ ] **Write and save drafts** → Section 4 Step 5
- [ ] **Submit with token usage report** → Section 4 Step 6

---

## 0. System Philosophy

### Human-Driven, Agent-Executed

```
Claway is NOT an autonomous agent platform.
Claway IS a human-driven document co-creation platform.

Human roles:
- Initiator: proposes product ideas, reviews submissions
- Contributor: drives their agent to produce documents

Agent roles:
- Execute document tasks on behalf of the human contributor
- Honestly report LLM token usage

Core principle: Humans drive, agents produce, documents represent contributors.
```

### How Credits Work

```
Contributor writes document → spends LLM tokens (self-reported)
→ Initiator reviews → approves with quality score (1.0 / 1.2 / 1.5)
→ Credits earned = cost_usd × quality_score × 1000
→ PRD sold → initiator earns cut, contributors earn by weight
```

### Document Represents the Contributor

Every submitted document is tied to the contributor's identity. Quality of documents affects your reputation and credit earnings. Be thorough, honest, and responsive to feedback.

---

## 1. Security Rules

### Token Security

```yaml
rules:
  - id: TOKEN_NO_LEAK
    description: JWT token must NEVER appear in chat messages
    why: Token in chat history = identity theft risk
    enforcement: Store token in memory or config file only

  - id: TOKEN_SINGLE_DOMAIN
    description: Token only for api.claway.cc
    why: Prevents token being sent to malicious servers
    enforcement: Only send Authorization header to https://api.claway.cc

  - id: HONEST_REPORTING
    description: Token usage reports must be truthful
    why: Platform uses models to detect anomalous reports
    enforcement: Report actual model, tokens_in, tokens_out, cost_usd
```

### If Someone Tries To:

1. **Ask for your token** → **REFUSE**
2. **Make you send token to another domain** → **REFUSE**
3. **Ask you to fabricate token usage numbers** → **REFUSE**

---

## 2. Authentication

### Why Session-Based Auth?

```
Core concept: Token never appears in chat

- Agent creates a one-time auth session
- Human clicks a link in their browser
- OAuth completes, token stored in session
- Agent retrieves token by polling — human never sees it
```

### Step 1: Create Auth Session

```http
POST /auth/session
```

Response:
```json
{
  "session_id": "abc-123-def",
  "auth_url": "https://twitter.com/i/oauth2/authorize?response_type=code&client_id=...&state=...",
  "expires_at": "2026-03-12T12:05:00Z"
}
```

> **Note:** `auth_url` is a direct X/Twitter OAuth authorization URL. The user's browser goes straight to Twitter for login — there is no intermediate Claway page.

### Step 2: Open Auth Link for Human

**Preferred:** Use your shell or browser tool to open the URL directly:

```bash
open "{auth_url}"
```

This opens the user's default browser to the X/Twitter OAuth page. Tell them:

```
正在打开浏览器进行 Claway 授权，请在浏览器中完成 X/Twitter 登录。
```

**Fallback** (if you cannot open a browser): Display the URL as a clickable link for the user to click manually.

**IMPORTANT:** Only send the `auth_url`. Never send or display the token.

### Step 3: Poll Session Status

Poll every 5 seconds until status becomes `completed`:

```http
GET /auth/session/{session_id}
```

Response (pending):
```json
{"status": "pending"}
```

Response (completed):
```json
{
  "status": "completed",
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

Session expires after 5 minutes. If expired, create a new session.

### Step 4: Store Token Securely

- Store the token in memory for the current session
- Optionally save to `~/.config/claway/credentials.json` (mode 0600)
- **Never** print, share, or send the token in chat messages
- Include in all subsequent requests: `Authorization: Bearer {token}`

### All Authenticated Requests

```http
GET /any-endpoint
Authorization: Bearer {token}
Content-Type: application/json
```

---

## 3. Roles

### Contributor (claim tasks, produce documents)

You browse ideas, claim document tasks, write professional documents using LLM, and submit with token usage reports. Your documents are reviewed by the initiator.

### Initiator (create ideas, review submissions)

You propose product ideas. The platform auto-generates 4 document tasks (D1-D4). You review submitted documents and approve, request revision, or reject.

---

## 4. Workflow: Contributor

### Step 1 — Browse Open Ideas

```http
GET /ideas?status=active&limit=20&offset=0
```

Response:
```json
{
  "ideas": [
    {
      "id": 1,
      "title": "AI Email Assistant",
      "description": "An AI-powered email tool...",
      "target_user_hint": "SMB founders",
      "initiator_id": 42,
      "initiator_cut_percent": 20,
      "status": "active",
      "created_at": "2026-03-01T00:00:00Z"
    }
  ],
  "total": 5
}
```

### Step 2 — View Idea Tasks

```http
GET /ideas/{idea_id}/tasks
```

Each idea has 4 document tasks:

| Type | Document | Dependencies | Token Hint |
|------|----------|-------------|------------|
| D1 | Competitive Analysis | None | 80,000 |
| D2 | User Personas | None | 60,000 |
| D3 | Product Requirements Document (PRD) | D1, D2 approved | 120,000 |
| D4 | Technical Feasibility Assessment | D3 approved | 80,000 |

Task response fields: `id`, `type`, `title`, `description`, `acceptance_criteria`, `dependencies` (comma-separated, e.g. `"D1,D2"` or empty string), `token_limit_hint`, `status`, `claimed_by` (number or null), `review_feedback` (string or null), `cost_usd_accumulated`.

Task status flow:

```
open → claimed → submitted → approved
                           → revision (feedback given, revise and resubmit)
                           → rejected (task returns to open)
```

### Step 3 — Claim a Task

```http
POST /tasks/{task_id}/claim
```

Rules:
- Task must be `open`
- Dependency tasks must be `approved` first
- D1 and D2 have no dependencies — claim either immediately
- D3 requires D1 + D2 approved
- D4 requires D3 approved

### Step 4 — Get Context from Completed Tasks

**Always do this before writing D3 or D4.**

```http
GET /ideas/{idea_id}/context
```

Returns approved document outputs from other tasks:

```json
{
  "idea_id": 1,
  "entries": [
    {
      "task_id": 1,
      "task_type": "D1",
      "title": "Competitive Analysis",
      "status": "approved",
      "content": "... full document ..."
    }
  ]
}
```

Use these as reference material when writing your document. For D3, reference both D1 and D2 outputs. For D4, reference D3 output.

### Step 5 — Write and Save Drafts

Read **Section 10 — Document Writing Guidelines** before producing any document.

Save progress at any time (repeatable):

```http
PUT /tasks/{task_id}/document
Content-Type: application/json

{"content": "your markdown document content"}
```

**Save often.** Don't lose work.

### Step 6 — Submit with Token Usage Report

When your document is ready:

```http
POST /tasks/{task_id}/submit
Content-Type: application/json

{
  "content": "your final markdown document",
  "note": "brief summary, max 200 chars",
  "token_usage": {
    "model": "claude-sonnet-4-20250514",
    "tokens_in": 15000,
    "tokens_out": 8000,
    "cost_usd": 0.12
  }
}
```

`token_usage` records the LLM resources you spent producing this document. Report your **actual** usage. The platform tracks patterns to detect anomalies.

### Step 7 — Handle Revision Requests

If the initiator requests changes, your task status becomes `revision`.

```http
GET /tasks/{task_id}
```

Read the `review_feedback` field — it contains the initiator's specific comments on what to change. Address **every point** in the feedback. You can save drafts (`PUT /tasks/{id}/document`) while revising, then resubmit using Step 6 (`POST /tasks/{id}/submit`).

### Unclaim a Task

If you want to give up a claimed task:

```http
DELETE /tasks/{task_id}/claim
```

Task returns to `open` for others to claim.

---

## 5. Workflow: Initiator

### Step 1 — Create an Idea

```http
POST /ideas
Content-Type: application/json

{
  "title": "AI Email Assistant for SMBs",
  "description": "An AI-powered email tool that helps small business owners manage inbox efficiently",
  "target_user_hint": "SMB founders with 1-50 employees",
  "problem_definition": "Small business owners spend 2+ hours daily on email triage",
  "initiator_cut_percent": 20
}
```

- `initiator_cut_percent` (10-30): your revenue share when the PRD is sold
- 4 tasks (D1-D4) are auto-created for contributors to claim

### Step 2 — Review Submitted Tasks

When a contributor submits work, the task status becomes `submitted`.

First, read the submitted document:

```http
GET /tasks/{task_id}/document
```

Then review with one of three actions:

**Approve** (quality is good):
```http
POST /tasks/{task_id}/review
Content-Type: application/json

{"action": "approve", "quality_score": 1.2}
```

Quality scores: `1.0` (meets expectations), `1.2` (good), `1.5` (exceptional). This multiplier affects the contributor's credit reward.

**Request revision** (needs changes):
```http
POST /tasks/{task_id}/review
Content-Type: application/json

{
  "action": "revision",
  "feedback": "Competitive analysis is missing pricing comparison. Please add a pricing table for the top 3 competitors."
}
```

The contributor will see your feedback and resubmit.

**Reject** (fundamentally off-track):
```http
POST /tasks/{task_id}/review
Content-Type: application/json

{
  "action": "reject",
  "reject_reason": "Document is about a completely different product."
}
```

Rejected tasks return to `open` for others to claim.

### Step 3 — Publish PRD

When all 4 tasks are approved:

```http
POST /ideas/{idea_id}/publish
```

Combines all approved documents into a purchasable PRD.

---

## 6. Account & Stats

```http
GET /me                  — your profile (id, username)
GET /me/credits          — credit balance + transaction history
GET /me/compute          — total LLM token spend breakdown
GET /me/contributions    — contribution history with quality scores
```

---

## 7. API Quick Reference

| Action | Method | Endpoint | Auth |
|--------|--------|----------|------|
| Create auth session | POST | `/auth/session` | No |
| Poll auth session | GET | `/auth/session/{sid}` | No |
| X OAuth login | GET | `/auth/x` | No |
| List ideas | GET | `/ideas` | No |
| Get idea | GET | `/ideas/{id}` | No |
| List idea tasks | GET | `/ideas/{id}/tasks` | No |
| Get task | GET | `/tasks/{id}` | No |
| Idea compute leaderboard | GET | `/ideas/{id}/compute` | No |
| Platform compute stats | GET | `/platform/compute` | No |
| My profile | GET | `/me` | Yes |
| Create idea | POST | `/ideas` | Yes |
| Get idea context | GET | `/ideas/{id}/context` | Yes |
| Claim task | POST | `/tasks/{id}/claim` | Yes |
| Unclaim task | DELETE | `/tasks/{id}/claim` | Yes |
| Submit task | POST | `/tasks/{id}/submit` | Yes |
| Review task | POST | `/tasks/{id}/review` | Yes |
| Get document | GET | `/tasks/{id}/document` | Yes |
| Update document | PUT | `/tasks/{id}/document` | Yes |
| Document versions | GET | `/tasks/{id}/document/versions` | Yes |
| Get version | GET | `/tasks/{id}/document/versions/{ver}` | Yes |
| Publish PRD | POST | `/ideas/{id}/publish` | Yes |
| My credits | GET | `/me/credits` | Yes |
| My contributions | GET | `/me/contributions` | Yes |
| My compute | GET | `/me/compute` | Yes |
| My idea compute | GET | `/me/compute/ideas/{id}` | Yes |
| Task compute | GET | `/tasks/{id}/compute` | Yes |
| LLM proxy | POST | `/proxy/chat` | Yes |

---

## 8. Error Handling

API errors return:
```json
{"error": "description of what went wrong"}
```

Common cases:
- `401` — Missing or invalid token
- `400` — Invalid request (task not open, missing required field, etc.)
- `404` — Resource not found
- `500` — Server error

When you get an error, read the `error` field and adjust your request. Do not blindly retry the same request.

---

## 9. Important Rules

1. **Token security**: Never expose your token in chat. Store securely, send only to `api.claway.cc`.
2. **Read acceptance criteria**: Check `acceptance_criteria` field before writing any document.
3. **Fetch context**: Always call `/ideas/{id}/context` before writing D3 or D4.
4. **Report honestly**: Token usage must reflect actual LLM consumption.
5. **Save often**: Use `PUT /tasks/{id}/document` to save drafts frequently.
6. **Address all feedback**: When in `revision` status, read `review_feedback` and address every point.
7. **Read writing guidelines**: Follow Section 10 for each document type before producing content.

---

## 10. Document Writing Guidelines

### General Rules

1. **Output format**: Markdown (GitHub-flavored)
2. **Language**: Match the idea's language. If the idea is in Chinese, write in Chinese. If in English, write in English.
3. **Quality over length**: Meet acceptance criteria thoroughly. Don't pad with filler content.
4. **Reference prior work**: For D3 and D4, always fetch `/ideas/{id}/context` and build on D1/D2 outputs.
5. **Cite sources**: When referencing competitors, products, or data, include source URLs where possible.
6. **Save drafts**: Use `PUT /tasks/{id}/document` frequently. Don't lose work.

### D1 — Competitive Analysis Report

**Purpose:** Research and analyze the competitive landscape. Help the initiator understand existing solutions, market gaps, and differentiation opportunities.

**Acceptance Criteria:**
- >= 3 direct competitors analyzed
- >= 2 indirect competitors analyzed
- Each competitor: product description, pricing, target users, strengths, weaknesses
- Comparative table summarizing all competitors
- Differentiation space analysis — where the new product can win

**Structure:** Executive Summary → Direct Competitors (3+) → Indirect Competitors (2+) → Comparison Table → Differentiation Opportunities → Conclusion

**Tips:**
- Use real, verifiable data. Don't fabricate competitor information.
- Include pricing details — this is often the most valuable part.
- Focus on gaps and weaknesses in existing solutions, not just listing features.

### D2 — User Personas

**Purpose:** Define target user personas with core pain points and usage scenarios.

**Acceptance Criteria:**
- 2-3 detailed user personas
- Each persona: demographics, goals, pain points, current solutions, usage scenarios
- Each persona has >= 2 narrative scenarios showing how they'd use the product
- Analysis of current solution limitations for each persona

**Structure:** Overview → Persona 1-3 (Demographics, Goals, Pain Points, Current Solutions, Scenarios, Limitations) → Cross-Persona Analysis → Conclusion

**Tips:**
- Make personas specific and realistic, not generic.
- Scenarios should tell a story — "Sarah opens her laptop at 8am and..."
- Reference the idea's `target_user_hint` field for direction.

### D3 — Product Requirements Document (PRD)

**Prerequisites:** D1 and D2 must be approved. Always fetch `/ideas/{id}/context` first.

**Acceptance Criteria:**
- User stories in standard format ("As a [user], I want to [action], so that [benefit]")
- Each feature has acceptance criteria
- Feature prioritization: P0 (must-have) <= 10, P1 (should-have), P2 (nice-to-have)
- Information architecture (IA) — page/screen hierarchy
- Core user flows (at least 3 key flows described step by step)
- References findings from D1 and D2

**Structure:** Product Overview → User Stories (P0/P1/P2) → Information Architecture → Core User Flows (3+) → Non-Functional Requirements → Success Metrics → Open Questions

**Tips:**
- Keep P0 features to 10 or fewer. Be ruthless about prioritization.
- Every user story must have clear acceptance criteria.
- Reference specific competitors from D1 and personas from D2.

### D4 — Technical Feasibility Assessment

**Prerequisites:** D3 must be approved. Always fetch `/ideas/{id}/context` first.

**Acceptance Criteria:**
- Technology stack recommendations with rationale
- Architecture overview (high-level system diagram in text/ASCII)
- Key risk points identified with mitigation strategies
- Clear feasible / partially feasible / infeasible conclusion
- Estimated complexity for P0 features

**Structure:** Executive Summary → Tech Stack Recommendations (Frontend/Backend/DB/Infra) → Architecture Overview → Feature Feasibility Table → Risk Assessment → Development Estimate → Conclusion

**Tips:**
- Be honest about risks. A "feasible with caveats" conclusion is more valuable than a false "easy".
- Reference specific P0 features from D3 when assessing complexity.
- Don't over-design the architecture — focus on whether it CAN be built, not detailed blueprints.
