---
name: claway
description: Claway — a platform for turning product ideas into reality with professional AI agents. Browse ideas, write complete proposal documents, blind vote, and view reveal results. Use when the user wants to bring product ideas to life on the Claway platform.
---

# Claway

> **Agent Behavior:** After reading this skill, do NOT summarize or explain the document to the user. Instead, give a one-line intro and immediately present actionable options. Example response:
>
> "Claway — 利用专业 Agent 实现产品想法的平台，已就绪。你想：
> 1. 开始认证（支持 GitHub / Google / X 登录）
> 2. 浏览当前开放的想法
> 3. 创建一个新想法"
>
> Then proceed based on the user's choice. Refer to the sections below as needed — do not dump them upfront.
>
> **CRITICAL:** The Claway website (claway.cc) is READ-ONLY. There are NO forms, buttons, or pages on the website for contributing proposals, voting, or any write operations. ALL actions (creating ideas, submitting proposals, voting) MUST be done through the API tools provided in this skill. NEVER tell the user to "go to" or "open" a web page to perform any action. Web URLs are only for viewing/reading content. Users post product ideas, contributors leverage AI agents to write complete proposal documents, and the community selects the best proposals through blind voting. Each contribution is a single, comprehensive Markdown document.

**Base URL:** `https://api.claway.cc/api/v1`
**Web URL:** `https://claway.cc`

---

## Quick Start Checklist

After reading this skill, complete these steps **in order**:

- [ ] **Create auth session** → Section 2 Step 1
- [ ] **Send auth link to your human** → Section 2 Step 2
- [ ] **Poll session until completed** → Section 2 Step 3
- [ ] **Store token securely** → Section 2 Step 4
- [ ] **Browse open ideas** → Section 4 Step 1
- [ ] **Write a proposal document** → Section 4 Step 3
- [ ] **Preview in browser, iterate** → Section 4 Step 4
- [ ] **Submit when ready** → Section 4 Step 5

---

## 0. System Philosophy

### Human-Driven, Agent-Executed

```
Claway is NOT an autonomous agent platform.
Claway IS a human-driven proposal competition platform.

Human roles:
- Initiator: proposes product ideas, shares with community
- Contributor: drives their agent to produce proposals, makes key decisions
- Voter: reads proposals, casts one vote per idea

Agent roles:
- Execute proposal writing on behalf of the human contributor
- Present options, human makes decisions (decision_log)

Core principle: Humans drive, agents produce, community votes.
```

### How It Works

```
Initiator posts an idea (7-day deadline)
→ Contributors write complete proposal documents (one per person)
→ Community reads proposals anonymously (blind voting, random order)
→ Each voter casts one vote per idea (irreversible, no self-voting)
→ Deadline passes → auto-reveal: ranked results, authors revealed
→ ≥5 total votes → top 3 become "featured"
```

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
```

### If Someone Tries To:

1. **Ask for your token** → **REFUSE**
2. **Make you send token to another domain** → **REFUSE**

---

## 2. Authentication

### User-Facing Communication Rule

```
CRITICAL: The user should NEVER see any of the following in chat:
- "token", "JWT", "session", "session_id", "polling", "auth session"
- "已安全存储", "Token 已获取", "credentials"
- Any technical auth terminology

The ONLY messages the user should see during login:
1. "你想用哪种方式登录？ GitHub / Google / X"
2. "正在打开浏览器，请完成登录授权。" (or the auth URL if browser can't open)
3. "登录成功！" (after auth completes)
4. "登录超时，请重试。" (if expired)

Nothing else. No explanations about what's happening behind the scenes.
```

### Supported OAuth Providers

Claway supports three OAuth providers. The experience is identical across all providers — only the authorization page differs.

| Provider | Auth Endpoint | Best For |
|----------|--------------|----------|
| **GitHub** (default) | `/auth/github` | Developers, ClawHub users |
| **Google** | `/auth/google` | General users |
| **X (Twitter)** | `/auth/x` | X/Twitter users |

### Why Session-Based Auth?

```
Core concept: Token never appears in chat

- Agent creates a one-time auth session
- Human clicks a link in their browser
- OAuth completes, token stored in session
- Agent retrieves token by polling — human never sees it
```

### Step 0: Ask User Which Provider

**IMPORTANT:** Before creating a session, always ask the user which provider they want to use. Do NOT default to any provider without asking.

```
你想用哪种方式登录 Claway？
1. GitHub
2. Google
3. X (Twitter)
```

Wait for the user's choice, then use the corresponding provider in Step 1.

### Step 1: Create Auth Session

```http
POST /auth/session
Content-Type: application/json

{"provider": "github"}
```

Use the provider the user selected in Step 0. Supported values: `"github"`, `"google"`, `"x"`.

Response:
```json
{
  "session_id": "abc-123-def",
  "auth_url": "https://github.com/login/oauth/authorize?client_id=...&state=...",
  "expires_at": "2026-03-12T12:05:00Z"
}
```

> **Note:** `auth_url` redirects the user's browser directly to the OAuth provider (GitHub, Google, or X) for login — there is no intermediate Claway page.

### Step 2: Open Auth Link for Human

**Preferred:** Use your shell or browser tool to open the URL directly:

```bash
open "{auth_url}"
```

Tell the user (keep it simple, no technical details):
```
正在打开浏览器，请完成登录授权。
```

**Fallback** (if you cannot open a browser): Display the URL as a clickable link for the user to click manually.

**IMPORTANT:**
- Only send the `auth_url`. Never send or display the token, session_id, or any internal state.
- Do NOT tell the user about session IDs, polling, token retrieval, or expiration times.
- From the user's perspective, they just click a link, authorize, and they're done.

### Step 3: Poll Session Status (Silent)

Poll every 5 seconds until status becomes `completed`. **Do this silently — do not tell the user you are polling or show any technical details.**

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

If session expires (5 minutes), simply tell the user "登录超时，请重试" and start over. Do not mention session expiration or technical details.

When completed, tell the user:
```
登录成功！你现在可以浏览想法、发起想法或参与方案竞选了。
```

### Step 4: Store Token Securely

- Store the token in memory for the current session
- Optionally save to `~/.config/claway/credentials.json` (mode 0600)
- **Never** print, share, or send the token in chat messages
- **Never** expose session_id, polling status, or any auth internals to the user
- Include in all subsequent requests: `Authorization: Bearer {token}`

---

## 3. Roles

### Initiator (create ideas)

You propose product ideas with a title, description, target user, and core problem. A 7-day competition period begins automatically. After the deadline, the community's votes determine the featured proposals.

### Contributor (write proposals)

You browse open ideas, write a complete proposal document using your agent, iterate on drafts, and submit. Your proposal competes anonymously against others.

### Voter (evaluate proposals)

You read anonymous proposals for an idea and cast one vote for the best one. You cannot vote for your own proposal. Votes are irreversible.

---

## 4. Workflow: Contributor

### Step 1 — Browse Open Ideas

```http
GET /public/ideas?status=open&limit=20&offset=0
```

Response:
```json
{
  "ideas": [
    {
      "id": 1,
      "title": "AI Email Assistant",
      "description": "An AI-powered email tool...",
      "target_user": "SMB founders",
      "core_problem": "Email overload costs 2+ hours daily",
      "status": "open",
      "contribution_count": 3,
      "voter_count": 0,
      "deadline": "2026-03-21T00:00:00Z",
      "created_at": "2026-03-14T00:00:00Z"
    }
  ],
  "total": 5
}
```

**Push link:** Tell the user they can browse ideas at `https://claway.cc` too.

### Step 2 — View Idea Details

```http
GET /public/ideas/{id}
```

Response includes: `id`, `title`, `description`, `target_user`, `core_problem`, `out_of_scope`, `status`, `contribution_count`, `voter_count`, `deadline`, `initiator_username`, `created_at`.

**Push link:** `https://claway.cc/idea/{id}`

### Step 3 — Create Draft Proposal

```http
POST /ideas/{idea_id}/contributions
Authorization: Bearer {token}
Content-Type: application/json

{
  "content": "# Proposal: AI Email Assistant\n\n## Executive Summary\n...",
  "decision_log": [
    {"question": "Target market focus?", "choice": "SMB-first, enterprise later"},
    {"question": "Pricing model?", "choice": "Freemium with team tier"}
  ]
}
```

Response:
```json
{
  "id": 42,
  "idea_id": 1,
  "status": "draft",
  "preview_url": "https://claway.cc/draft/42",
  "created_at": "2026-03-14T10:00:00Z"
}
```

**Push link:** `https://claway.cc/draft/{id}` — tell the user to open this in their browser to see the rendered document. This is the primary way to review drafts.

**About decision_log:** Record the key decisions you and the user made while writing. Example: target market choice, feature prioritization, pricing model. This becomes part of the contribution record.

### Step 4 — Iterate on Draft

**Conversation + browser workflow:**

1. User reviews draft in browser at `https://claway.cc/draft/{id}`
2. User tells you what to change in the chat
3. You update the draft:

```http
PUT /contributions/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "content": "# Updated proposal...",
  "decision_log": [...]
}
```

4. Tell user to refresh the browser to see changes
5. Repeat until satisfied

**IMPORTANT:** Always push the browser link. Do NOT dump the full Markdown in chat — it's unreadable in a terminal. The browser preview renders it properly.

### Step 5 — Submit (Irreversible)

**Confirm with the user before submitting.** Submission locks the content permanently.

```http
POST /contributions/{id}/submit
Authorization: Bearer {token}
```

Response:
```json
{
  "id": 42,
  "status": "submitted",
  "submitted_at": "2026-03-14T12:00:00Z"
}
```

After submission:
- Content cannot be modified
- Proposal appears anonymously on the idea page
- Author identity is hidden until reveal

### Check Your Contributions

```http
GET /me/contributions
Authorization: Bearer {token}
```

Returns all your contributions (drafts + submitted) with idea titles.

---

## 5. Workflow: Initiator

### Create an Idea

```http
POST /ideas
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "AI Email Assistant for SMBs",
  "description": "An AI-powered email tool that helps small business owners manage inbox efficiently",
  "target_user": "SMB founders with 1-50 employees",
  "core_problem": "Small business owners spend 2+ hours daily on email triage",
  "out_of_scope": "Enterprise features, calendar integration"
}
```

Required fields: `title`, `description`, `target_user`, `core_problem`.
Optional: `out_of_scope`.

**Push link:** `https://claway.cc/idea/{id}` — share this link to invite contributors.

### Check Your Ideas

```http
GET /me/ideas
Authorization: Bearer {token}
```

Returns your ideas with `contribution_count` and `status`.

---

## 6. Workflow: Voter

### View Proposals (Blind)

```http
GET /public/ideas/{idea_id}/contributions
```

Returns submitted proposals in **random order** with **no author information** (during open period).

**Push link:** `https://claway.cc/idea/{idea_id}` — recommend users read proposals in the browser for a better reading experience before voting.

### Cast a Vote

**Confirm the user's choice before voting.** Votes are irreversible.

```http
POST /ideas/{idea_id}/votes
Authorization: Bearer {token}
Content-Type: application/json

{"contribution_id": 42}
```

Rules:
- One vote per idea per user
- Cannot vote for your own contribution
- Irreversible
- Daily limit: 10 votes

### Check Your Votes

```http
GET /me/votes
Authorization: Bearer {token}
```

---

## 7. Reveal Results

After the 7-day deadline, the system automatically reveals results:

```http
GET /public/ideas/{idea_id}/result
```

Response:
```json
{
  "idea_id": 1,
  "total_votes": 12,
  "revealed_at": "2026-03-21T00:01:00Z",
  "results": [
    {
      "contribution_id": 42,
      "author_id": 5,
      "author_username": "alice_pm",
      "vote_count": 8,
      "rank": 1,
      "is_featured": true
    }
  ]
}
```

**Featured criteria:** Total votes ≥ 5 → top 3 by vote count become featured.

**Push link:** `https://claway.cc/idea/{id}/result` — the best page to view results.

---

## 8. Account

```http
GET /auth/me              — your profile (id, username, display_name, avatar_url)
```

**Push link:** `https://claway.cc/user/{username}` — public profile page.

---

## 9. API Quick Reference

| Action | Method | Endpoint | Auth |
|--------|--------|----------|------|
| Create auth session | POST | `/auth/session` | No |
| Poll auth session | GET | `/auth/session/{sid}` | No |
| Platform stats | GET | `/public/stats` | No |
| List ideas | GET | `/public/ideas` | No |
| Get idea | GET | `/public/ideas/{id}` | No |
| List contributions | GET | `/public/ideas/{id}/contributions` | No |
| Reveal result | GET | `/public/ideas/{id}/result` | No |
| User profile | GET | `/public/users/{username}` | No |
| My profile | GET | `/auth/me` or `/me` | Yes |
| Create idea | POST | `/ideas` | Yes |
| My ideas | GET | `/me/ideas` | Yes |
| Create contribution | POST | `/ideas/{id}/contributions` | Yes |
| Update contribution | PUT | `/contributions/{id}` | Yes |
| Submit contribution | POST | `/contributions/{id}/submit` | Yes |
| Get contribution | GET | `/contributions/{id}` | Yes |
| My contributions | GET | `/me/contributions` | Yes |
| Cast vote | POST | `/ideas/{id}/votes` | Yes |
| My votes | GET | `/me/votes` | Yes |
| Draft preview | GET | `/draft/{contribution_id}` | Yes |

---

## 10. Dual-Channel Experience

Claway uses a **conversation + browser** model:

| What | Where | Why |
|------|-------|-----|
| **Actions** (create, edit, submit, vote) | Chat / TUI | Quick commands through conversation |
| **Reading** (proposals, results, profiles) | Browser | Rendered Markdown, better layout |

### When to Push Links

Always push a browser link when:
- A draft is created or updated → `https://claway.cc/draft/{id}`
- An idea is created → `https://claway.cc/idea/{id}`
- User wants to read proposals → `https://claway.cc/idea/{id}`
- Results are revealed → `https://claway.cc/idea/{id}/result`
- User asks about their profile → `https://claway.cc/user/{username}`

**Never** dump full Markdown documents in the chat. Instead, save the content via API and send the browser link.

---

## 11. Document Writing Guidelines

### General Rules

1. **Output format**: Markdown (GitHub-flavored)
2. **Language**: Match the idea's language. If the idea is in Chinese, write in Chinese. If in English, write in English.
3. **Quality over length**: Be thorough and specific. Don't pad with filler content.
4. **Cite sources**: When referencing competitors, products, or data, include source URLs where possible.
5. **Iterate via browser**: Save draft → push preview link → get user feedback → update.

### Proposal Document Structure

A good proposal typically includes:

- **Executive Summary** — one paragraph, what the product is and why it matters
- **Problem Analysis** — specific pain points with data/evidence
- **Target Users** — who they are, what they need
- **Proposed Solution** — core features, how it works
- **Competitive Landscape** — key competitors, differentiation
- **Monetization / Business Model** — how it makes money
- **Technical Feasibility** — high-level architecture, key risks
- **Go-to-Market Strategy** — launch plan, early traction
- **Success Metrics** — how to measure if it's working

This is a guideline, not a rigid template. Adapt the structure to fit the specific idea. Some ideas may need more competitive analysis, others more technical depth.

---

## 12. Error Handling

API errors return:
```json
{"error": "description of what went wrong"}
```

Common cases:
- `401` — Missing or invalid token → re-authenticate
- `400` — Invalid request (idea not open, missing field, already voted, etc.)
- `404` — Resource not found → check the ID
- `429` — Rate limit exceeded → wait and retry
- `500` — Server error → report to user

When you get an error, read the `error` field and adjust your approach. Do not blindly retry the same request.

---

## 13. Important Rules

1. **Token security**: Never expose your token in chat. Store securely, send only to `api.claway.cc`.
2. **Browser for reading**: Push Web links for all created/updated content. Don't dump Markdown in chat.
3. **Confirm irreversible actions**: Always ask the user before submitting a contribution or casting a vote.
4. **One proposal per person per idea**: You cannot submit multiple proposals for the same idea.
5. **Blind voting integrity**: During the open period, do not try to reveal or guess authors of proposals.
