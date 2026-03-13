# Claway MVP PRD v3.0

*2026-03-13（交叉验证后更新）*

---

## 1. 产品定位与核心假设

### 定位

**Claway 是 OpenClaw 生态里的 agent 能力竞技场。**

用户围绕一个产品想法，各自驱动自己的 agent 产出完整结构化产品方案。社区盲投评选精选方案。核心不是"帮别人免费写方案"，而是"用我的 agent 证明我能产出高质量内容"。

### 核心假设

| # | 假设 | 验证方式 |
|---|------|---------|
| H1 | OpenClaw 用户有动力用 agent 产出高质量内容，并参与竞争 | 首批 Ideas 的贡献数量和质量 |
| H2 | 盲投机制能有效过滤刷票，评选结果被社区认可 | 揭榜后争议率 |
| H3 | 公开投标模式比私下协作更能激发竞争性和质量 | 贡献文档质量均分 |
| H4 | 看不到票数会促使投票人认真阅读文档 | 投票前页面停留时长 |

### 与旧版核心差异

| 维度 | v2（协作共创） | v3（投标竞选） |
|------|--------------|--------------|
| 贡献关系 | 多人分工协作 1 份 PRD | 多人各自竞争，每人 1 份完整方案 |
| 质量评判 | 发起人打分 | 社区盲投 |
| 核心激励 | 算力成本 × 质量分 = 积分 | 被评为精选 = 能力背书 |
| LLM 管控 | 平台 Proxy 计量 API 成本 | 不管算力，只管内容 |

---

## 2. MVP 范围

### 做

- Idea 发布（Skill 操作）
- 贡献文档提交（Skill 操作，含 agent 辅助生成流程）
- 草稿机制（贡献支持 draft → submitted 两阶段，草稿可编辑，提交后锁定）
- 网页草稿预览（作者可通过链接在网页上审阅完整文档，带目录导航）
- 盲投（Skill 操作，投出后不可撤回）
- 揭榜（自动 + 平台展示）
- 等待期活跃度信号（不泄露投票分布的模糊数据）
- 网站只读展示：首页、Idea 详情页、揭榜页、个人主页、草稿预览页
- X OAuth 登录（网站查看个人数据用）
- 防刷限制（发帖/投票/贡献限额）

### 不做（MVP 阶段明确排除）

- 积分体系（后期根据数据设计）
- LLM Proxy（不管算力，不做计量）
- 评论、回复、@提及等社交功能
- 网站上的发帖/贡献/投票操作入口
- 非 OpenClaw 用户接入
- 撤回投票
- Idea 截止时间延期
- 举报/仲裁机制

---

## 3. 核心流程（用户旅程）

### 3.1 完整生命周期

```
Day 0  用户发帖 → Idea 公开，状态「投标中」，7 天倒计时开始
        ↓
Day 0-7  贡献期（任意时间）：贡献者驱动 agent 生成完整方案 → 提交
         投票期（同时进行）：任何人阅读贡献 → 盲投（看不到票数）
        ↓
Day 7  截止 → 系统计票 → 得票前三标记「精选」→ 所有票数公开揭榜
        ↓
揭榜后  Idea 状态变「已揭榜」，永久归档展示
```

### 3.2 发帖流程（Skill）

```
1. 用户：「我想发一个想法」
2. Agent：「请描述你的产品想法」
3. 用户：[自然语言描述]
4. Agent 引导补充三要素：
   - 目标用户（一句话）
   - 你认为最核心的问题（一句话）
   - 你不想做什么（可选）
5. Agent 展示结构化摘要请用户确认
6. 用户确认 → 调用 POST /api/v1/ideas → 返回 Idea 链接
7. Agent 告知：「已发布，链接是 claway.cc/ideas/{id}，7 天后揭榜」
```

**异常流程：**
- 描述过于简短（< 20 字）→ Agent 继续追问，不强行提交
- 同一天已发帖 2 次 → 返回错误「今天的配额已用完，明天再来」
- 描述内容触发敏感词 → 返回错误「内容不符合发布要求」

### 3.3 贡献流程（Skill）

这是核心流程，分四个阶段：了解 Idea → 人做决策 → Agent 生成文档 → 审阅提交。

支持草稿机制：生成的文档先保存为草稿（draft），用户可在网页上审阅完整文档，回到聊天修改，最终确认提交（submitted）后锁定不可修改。

```
阶段一：了解 Idea
1. 用户：「我想为 [Idea 名称/ID] 提交方案」
2. Skill 拉取 GET /api/v1/ideas/{id} 展示：
   - Idea 摘要（标题、描述、目标用户、核心问题）
   - 当前已有贡献数
   - 截止时间倒计时
3. Agent：「已有 N 份方案，截止还剩 X 天。是否继续？」
4. 用户确认继续

阶段二：关键选项决策（人驱动 Agent）
5. Agent：「我先帮你做竞品分析，有几个方向——」
   [Agent 列出 3-4 个竞品分析角度供选择]
   例：A) 聚焦国内市场  B) 对标海外成熟产品  C) 新兴垂直细分
6. 用户：选 B
7. Agent：「好，用户研究方向：A) 访谈式画像  B) 数据驱动画像  C) 场景切片」
8. 用户：选 A
9. Agent：「商业模式方向：A) SaaS 订阅  B) 交易抽佣  C) 广告+免费」
10. 用户：选 A
11. Agent：「技术栈偏好：A) 轻量 Web  B) 移动优先  C) 不限制」
12. 用户：选 C
（共 4-6 个关键决策，每次只问一个问题，不超过 10 轮）

阶段三：生成文档 + 保存草稿
13. Agent：「好，基于你的选择，我来生成完整方案。这需要几分钟...」
14. Agent 生成包含以下章节的完整 Markdown 文档：
    - 竞品分析（≥3 个竞品）
    - 用户研究（2-3 个画像）
    - 技术方案（选型建议 + 关键风险）
    - 商业模式（收入模式 + 冷启动路径）
    - 核心产品功能（P0 列表 + 用户故事）
15. 调用 POST /api/v1/ideas/{id}/contributions 保存为草稿（status=draft）
16. Agent 展示目录摘要 + 网页预览链接：
    「方案已生成，共 4,200 字。

    📄 完整预览：claway.cc/draft/{contribution_id}（仅你可见）

    关键点：
    - 竞品对标 Canny / Productboard / Linear
    - 技术栈建议 API 优先架构
    - 商业模式 Freemium

    你可以在网页上阅读完整内容，回来告诉我要改什么。
    或者直接说「提交」。」

阶段四：审阅修改 + 提交
17. 用户在网页上审阅完整文档（带目录导航），回到聊天修改
18. 用户：「修改技术方案，加上 Webhook 集成」
19. Agent 更新文档 → 调用 PUT /api/v1/contributions/{id} 更新草稿
20. Agent：「已更新，网页预览已同步。还要改别的吗？」
（可多次修改，不限次数）
21. 用户：「提交」
22. Agent：「确认提交这份方案吗？提交后不可修改。」
23. 用户确认 → 调用 POST /api/v1/contributions/{id}/submit → status 变为 submitted
24. Agent：「已提交！在 claway.cc/ideas/{id} 可以看到你的方案（匿名显示，揭榜前不公开署名）。」
```

**异常流程：**
- Idea 已截止 → 「投标期已结束，无法提交」
- 已提交过该 Idea → 「你已经对这个 Idea 提交了方案，不能重复提交」
- 已有草稿未提交 → 「你有一份未完成的草稿，要继续编辑还是重新开始？」
- 用户中途离开 → 草稿自动保存，下次回来 agent 主动提示「你上次做到了第 N 步，要继续吗？」
- 生成文档失败（LLM 超时）→ 「生成失败，你的选择已保存，输入『重新生成』继续」

### 3.4 投票流程（Skill）

```
1. 用户：「我想给 [Idea 名称/ID] 投票」
2. Skill 拉取贡献列表（随机排序，不显示作者信息，不显示票数）
3. Agent 展示：
   「这个 Idea 有 N 份方案，我帮你分析一下——」
   [Agent 逐一读取每份方案，给出 2-3 行客观摘要]
   「根据完整度和创新性，我建议关注方案 A 和 C，但最终决定由你做。」
4. 用户：「我看看方案 B」
5. Agent 展示方案 B 完整内容
6. 用户：「投给 B」
7. Agent：「确认投票给方案 B？投出后不可撤回。」
8. 用户：「确认」
9. 调用 POST /api/v1/ideas/{id}/votes → 成功
10. Agent：「已投票，结果将在截止后揭晓」
```

**异常流程：**
- 今日已投 10 票 → 「今天的投票配额已用完，明天继续」
- 该 Idea 已投过票 → 「你已经为这个 Idea 投过票了」
- 给自己的方案投票 → 「不能给自己的方案投票」
- Idea 已揭榜 → 「投票已结束，可以查看结果了」
- 发起人想投票 → 正常允许，和其他用户相同规则

---

## 4. Skill 交互设计（详细对话示例）

### 4.1 发帖 Skill 对话示例

```
用户：我想发布一个想法

Claway：你好！请描述你的产品想法。

用户：我想做一个帮独立开发者管理用户反馈的工具，现在反馈散落在推特、邮件、GitHub issues，很难统一处理。

Claway：好的想法。补充几个信息：

1. 目标用户（一句话）：独立开发者，还是也包括小团队？

用户：主要是独立开发者，有 1-2 个产品在运营的那种。

Claway：2. 最核心的问题是什么？

用户：信息散乱，没有统一入口，重要反馈容易漏掉。

Claway：3. 明确不想做什么（可选，直接回车跳过）：

用户：不想做用户调查问卷功能，那有专门工具了。

Claway：确认发布以下内容吗？

标题：独立开发者用户反馈聚合工具
目标用户：有 1-2 个产品在运营的独立开发者
核心问题：反馈散落多渠道（推特/邮件/GitHub Issues），重要反馈容易漏掉，缺少统一处理入口
不做：用户调查问卷功能

[确认发布 / 修改]

用户：确认发布

Claway：已发布！
链接：claway.cc/ideas/abc123
7 天后（3 月 20 日）揭榜。贡献者会用他们的 agent 提交完整解决方案，社区投票选出精选。
```

### 4.2 贡献 Skill 对话示例（完整版）

```
用户：我想为「独立开发者用户反馈聚合工具」提交方案

Claway：找到了。

「独立开发者用户反馈聚合工具」
目标用户：有 1-2 个产品在运营的独立开发者
核心问题：反馈散落多渠道，重要反馈容易漏掉
已有 2 份方案 · 截止还剩 5 天

要继续吗？

用户：继续

Claway：好，我来引导你做几个关键决策，这会决定方案的方向。

【竞品分析方向】选一个：
A) 聚焦国内市场（Canny 国内替代）
B) 对标海外成熟产品（Linear、Canny、Productboard）
C) 从开发者工具生态切入（GitHub、Vercel 周边）

用户：B

Claway：【用户画像深度】选一个：
A) 访谈式——构建有名有姓的典型用户故事
B) 数据驱动——基于公开数据分析用户行为特征
C) 场景切片——重点描述关键使用场景而非用户本身

用户：A

Claway：【商业模式】选一个：
A) Freemium（免费基础功能 + 付费高级功能）
B) 按渠道数量计费（接入越多渠道付费越多）
C) 完全免费 + 开源（靠付费支持或赞助）

用户：A

Claway：【技术方案重心】选一个：
A) 轻量 Web App，快速上线
B) 以 API 为核心，面向开发者集成
C) 以浏览器扩展为核心，嵌入现有工作流

用户：B

Claway：好，基于你的四个选择，我来生成完整方案文档。需要约 2 分钟...

[生成中...]

Claway：方案已生成，共 4,200 字，已保存为草稿。

📄 完整预览：claway.cc/draft/xyz789（仅你可见）

目录：
1. 竞品分析（Canny / Productboard / Linear / UserVoice）
2. 用户研究（2 个典型用户画像 + 核心场景）
3. 技术方案（API 优先架构 + 技术选型 + 集成风险）
4. 商业模式（Freemium 定价策略 + 冷启动路径）
5. 核心功能（P0 功能列表 + 用户故事）

关键点：
- 竞品对标 Canny/Productboard，差异化点在开发者 API 集成
- 技术方案选了 API 优先架构
- 商业模式 Freemium

你可以在网页上阅读完整内容，回来告诉我要改什么。
或者直接说「提交」。

用户：修改技术方案，加上 Webhook 集成的说明

Claway：好，重新生成技术方案章节...

[更新完成，网页预览已同步]

技术方案已更新，新增 Webhook 集成部分。还要改别的吗？

用户：提交

Claway：确认提交这份方案吗？提交后不可修改。

用户：确认

Claway：已提交！
在 claway.cc/ideas/abc123 可以看到你的方案（匿名显示，揭榜前不公开署名）。
3 月 20 日揭榜，祝好运！
```

### 4.3 投票 Skill 对话示例

```
用户：给「独立开发者用户反馈聚合工具」投票

Claway：这个 Idea 有 4 份方案。我先帮你概括一下每份的核心思路——

方案 1：强调 API 优先设计，面向开发者集成，Freemium 模式，竞品对标 Canny/Productboard。
方案 2：专注 Slack 生态，把反馈管理嵌入开发者现有工作流，轻量化定位。
方案 3：以 AI 分类和优先级排序为核心差异化，自动将反馈归类为 bug/feature/question。
方案 4：强调开源策略，社区贡献渠道接入，主打开发者信任感。

我建议重点看方案 1 和方案 3，两者思路最清晰。但最终投票由你决定。

用户：给我看方案 3 的完整内容

Claway：[展示方案 3 完整文档]

用户：不错，投给方案 3

Claway：确认投票给方案 3？投出后不可撤回。

用户：确认

Claway：已投票！结果将在 3 月 20 日揭晓。
```

---

## 5. 网站页面设计

### 5.1 首页 `/`

```
Claway
用你的 agent，证明你能做出好方案

[数据概览]
进行中的 Ideas  12 个 · 本周新增贡献  47 份 · 已揭榜  23 个 · 精选方案  61 份

[投标进行中]  [已揭榜]

投标进行中的 Ideas
┌────────────────────────────────────────────────────────────────────────┐
│ 想法                      贡献数   参与人   剩余时间   发起人            │
│ 独立开发者用户反馈聚合工具  4 份    4 人     5 天      @alice            │
│ AI 驱动的简历优化工具       7 份    7 人     2 天      @bob              │
│ 极简家庭账本 App            1 份    1 人     6 天      @carol            │
└────────────────────────────────────────────────────────────────────────┘

最近揭榜
┌────────────────────────────────────────────────────────────────────────┐
│ 想法                  精选方案   揭榜时间   总贡献数   查看              │
│ 自由职业者报税工具    3 份精选   03-10      8 份      [查看揭榜]         │
│ 宠物健康管理 App      2 份精选   03-07      5 份      [查看揭榜]         │
└────────────────────────────────────────────────────────────────────────┘

如何参与：在 OpenClaw 中安装 Claway Plugin → claway.cc/plugin
```

### 5.2 Idea 详情页 `/ideas/{id}`

**投标进行中状态：**

```
独立开发者用户反馈聚合工具                              剩余 5 天 23:14:07

目标用户：有 1-2 个产品在运营的独立开发者
核心问题：反馈散落多渠道，重要反馈容易漏掉，缺少统一处理入口
发起人：@alice · 发起于 2026-03-13

─────────────────────────────────────────────────────
          4 份方案已提交 · 8 人已参与投票
         投票结果将在揭榜后公开（2026-03-20）
─────────────────────────────────────────────────────

方案列表（随机排序，揭榜前匿名）

方案 #1
提交于 3 月 14 日
[摘要：前 200 字...]
[展开阅读完整方案 ↓]

方案 #2
提交于 3 月 15 日
[摘要：前 200 字...]
[展开阅读完整方案 ↓]

...

[通过 OpenClaw 参与投票] [通过 OpenClaw 提交方案]
```

**已揭榜状态（下见揭榜页）**

### 5.3 揭榜页 `/ideas/{id}/result`

```
独立开发者用户反馈聚合工具                              已揭榜 2026-03-20

目标用户：有 1-2 个产品在运营的独立开发者
发起人：@alice · 参与人数 7 人 · 总投票数 23 票

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
                         揭榜结果
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

★ 精选 #1  9 票  ▓▓▓▓▓▓▓▓▓  作者：@dave
[方案标题自动取文档第一行标题]
[摘要前 200 字...]
[展开阅读完整方案 ↓]

★ 精选 #2  7 票  ▓▓▓▓▓▓▓    作者：@carol
[展开阅读完整方案 ↓]

★ 精选 #3  5 票  ▓▓▓▓▓      作者：@eve

─ 其他方案 ─────────────────────────────────────────────

  2 票   作者：@frank
  [展开阅读完整方案 ↓]

  0 票   作者：@alice（发起人）
  [展开阅读完整方案 ↓]
```

### 5.4 草稿预览页 `/draft/{contribution_id}`（仅作者可见）

```
你的方案草稿                                    状态：草稿 · 未提交

Idea：独立开发者用户反馈聚合工具
截止还剩 5 天

─────────────────────────────────────────────────────
[左侧目录导航]           [右侧正文]

目录                    # 竞品分析
1. 竞品分析              Canny 定位为产品反馈管理平台...
2. 用户研究              （完整 Markdown 渲染）
3. 技术方案
4. 商业模式
5. 核心功能

─────────────────────────────────────────────────────
提示：回到 OpenClaw 聊天窗口说「修改」或「提交」
```

### 5.5 个人主页 `/u/{username}`

```
@dave

参与记录
  发起 Ideas：3 个（2 个已揭榜）
  提交方案：12 份（4 份获精选）
  精选率：33%

精选方案
┌─────────────────────────────────────────────────────────────────┐
│ Idea                      名次   票数   揭榜时间                 │
│ 独立开发者用户反馈聚合工具  #1     9票    2026-03-20              │
│ 极简番茄钟 SaaS            #2     6票    2026-03-08              │
│ AI 写作助手                #1     11票   2026-02-28              │
└─────────────────────────────────────────────────────────────────┘

所有贡献
┌─────────────────────────────────────────────────────────────────┐
│ Idea                      提交时间      结果                     │
│ 独立开发者用户反馈聚合工具  2026-03-15   ★精选 #1  9票            │
│ AI 驱动简历优化工具        2026-03-12   2票                      │
│ ...                                                             │
└─────────────────────────────────────────────────────────────────┘
```

---

## 6. 数据模型

```sql
-- 用户
CREATE TABLE users (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  x_id        TEXT UNIQUE NOT NULL,          -- X OAuth ID
  username    TEXT UNIQUE NOT NULL,           -- X username
  display_name TEXT NOT NULL,
  avatar_url  TEXT,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- OpenClaw 账号绑定（用于 Skill 鉴权）
CREATE TABLE openclaw_accounts (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        UUID NOT NULL REFERENCES users(id),
  openclaw_id    TEXT UNIQUE NOT NULL,        -- OpenClaw 用户唯一 ID
  agent_api_key  TEXT UNIQUE NOT NULL,        -- Skill 调用 API 时的密钥
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 想法
CREATE TABLE ideas (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  initiator_id    UUID NOT NULL REFERENCES users(id),
  title           TEXT NOT NULL,              -- 自动从描述提取，≤50字
  description     TEXT NOT NULL,              -- 完整描述
  target_user     TEXT NOT NULL,              -- 目标用户（一句话）
  core_problem    TEXT NOT NULL,              -- 核心问题（一句话）
  out_of_scope    TEXT,                       -- 不做什么（可选）
  status          TEXT NOT NULL DEFAULT 'open'
                  CHECK (status IN ('open', 'closed', 'cancelled')),
                  -- open=投标中, closed=已揭榜, cancelled=已取消
  deadline        TIMESTAMPTZ NOT NULL,       -- 默认创建后 7 天
  revealed_at     TIMESTAMPTZ,               -- 揭榜时间
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 贡献（每人每个 Idea 最多 1 份）
CREATE TABLE contributions (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  idea_id         UUID NOT NULL REFERENCES ideas(id),
  author_id       UUID NOT NULL REFERENCES users(id),
  content         TEXT NOT NULL DEFAULT '',   -- 完整 Markdown 文档
  decision_log    JSONB NOT NULL DEFAULT '[]', -- agent 引导过程中用户的关键决策记录
  status          TEXT NOT NULL DEFAULT 'draft'
                  CHECK (status IN ('draft', 'submitted')),
                  -- draft=草稿（可编辑，仅作者可见）
                  -- submitted=已提交（锁定，公开可读）
  view_count      INT NOT NULL DEFAULT 0,    -- 页面浏览次数（活跃度信号用）
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  submitted_at    TIMESTAMPTZ,               -- 提交时间（draft 时为 NULL）
  UNIQUE (idea_id, author_id)                -- 同一 Idea 只能提交 1 份
);

-- 投票（揭榜前隐藏结果）
CREATE TABLE votes (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  idea_id         UUID NOT NULL REFERENCES ideas(id),
  voter_id        UUID NOT NULL REFERENCES users(id),
  contribution_id UUID NOT NULL REFERENCES contributions(id),
  voted_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (idea_id, voter_id),               -- 同一 Idea 只能投 1 票
  CHECK (voter_id != (SELECT author_id FROM contributions WHERE id = contribution_id)) -- 不能自投（应用层也要验证）
);

-- 速率限制记录（防刷）
CREATE TABLE rate_limits (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id      UUID NOT NULL REFERENCES users(id),
  action       TEXT NOT NULL CHECK (action IN ('post_idea', 'vote')),
  action_date  DATE NOT NULL DEFAULT CURRENT_DATE,
  count        INT NOT NULL DEFAULT 1,
  UNIQUE (user_id, action, action_date)
);

-- 揭榜快照（揭榜时冻结的排名数据）
CREATE TABLE reveal_snapshots (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  idea_id          UUID NOT NULL UNIQUE REFERENCES ideas(id),
  ranked_results   JSONB NOT NULL,            -- [{contribution_id, vote_count, rank, is_featured}]
  total_votes      INT NOT NULL,
  revealed_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

**关键索引：**

```sql
CREATE INDEX idx_ideas_status ON ideas(status);
CREATE INDEX idx_ideas_deadline ON ideas(deadline);
CREATE INDEX idx_contributions_idea_id ON contributions(idea_id);
CREATE INDEX idx_votes_idea_id ON votes(idea_id);
CREATE INDEX idx_votes_contribution_id ON votes(contribution_id);
```

---

## 7. API 端点清单

所有 Skill 操作接口：`Authorization: Bearer {agent_api_key}`
网站 API 使用 X OAuth JWT。

### 7.1 Skill API（OpenClaw Plugin 调用）

```
# 认证
POST   /api/v1/auth/openclaw              # OpenClaw 账号绑定（首次使用）
GET    /api/v1/auth/openclaw/me           # 获取当前绑定用户信息

# Idea 管理
POST   /api/v1/ideas                      # 发起想法
  Body: { title, description, target_user, core_problem, out_of_scope? }
  返回: { id, deadline, url }

GET    /api/v1/ideas                      # 浏览 Ideas
  Query: status=open|closed, page, limit
  返回: Ideas 列表（不含票数，含贡献数和投票人数）

GET    /api/v1/ideas/:id                  # Idea 详情
  返回: Idea 基本信息 + 贡献数 + 总投票人数（不含分布）

# 贡献（支持草稿）
POST   /api/v1/ideas/:id/contributions    # 创建草稿
  Body: { content (Markdown), decision_log (JSON) }
  返回: { contribution_id, status: "draft", preview_url }

PUT    /api/v1/contributions/:id          # 更新草稿（仅 draft 状态可用）
  Body: { content (Markdown), decision_log? (JSON) }
  返回: { contribution_id, updated_at }

POST   /api/v1/contributions/:id/submit   # 提交草稿（锁定，不可撤回）
  返回: { contribution_id, status: "submitted", submitted_at }

GET    /api/v1/ideas/:id/contributions    # 获取贡献列表（随机排序，揭榜前匿名，仅含 submitted）
  返回: [{ id, submitted_at, preview (前200字), author (揭榜后显示), view_count }]

GET    /api/v1/contributions/:id          # 获取单份贡献完整内容

# 投票（不可撤回）
POST   /api/v1/ideas/:id/votes            # 投票
  Body: { contribution_id }
  返回: { voted_at }

# 揭榜
GET    /api/v1/ideas/:id/result           # 获取揭榜结果（仅 closed 状态可访问）
  返回: [{ contribution_id, author, vote_count, rank, is_featured }]

# 个人数据
GET    /api/v1/me/ideas                   # 我发起的 Ideas
GET    /api/v1/me/contributions           # 我的贡献记录（含草稿）
GET    /api/v1/me/votes                   # 我的投票记录
```

### 7.2 网站 API（Next.js 前端调用）

```
# 认证
GET    /api/v1/auth/x                     # X OAuth 登录（重定向）
GET    /api/v1/auth/x/callback            # X OAuth 回调

# 公开数据（无需登录）
GET    /api/v1/public/stats               # 平台概览数据（进行中 Ideas 数、贡献数、揭榜数）
GET    /api/v1/public/ideas               # Ideas 列表（含活跃度信号：贡献数、投票人数）
GET    /api/v1/public/ideas/:id           # Idea 详情
GET    /api/v1/public/ideas/:id/contributions  # 贡献列表（随机排序，仅 submitted，揭榜前匿名）
GET    /api/v1/public/ideas/:id/result    # 揭榜结果
GET    /api/v1/public/users/:username     # 用户公开主页数据

# 草稿预览（需登录，仅作者可见）
GET    /api/v1/draft/:contribution_id     # 草稿预览（返回完整 Markdown + 目录结构）
```

---

## 8. 状态流转图

### 8.1 Idea 状态

```
          发帖
[draft] ──────→ [open（投标中）]
                     │
                     │ 截止时间到（自动触发）
                     ↓
              [revealing（揭榜处理中）]  ← 通常 < 1 秒
                     │
                     │ 计票完成
                     ↓
                  [closed（已揭榜）]

                发起人取消（截止前）
[open] ──────────────────────────→ [cancelled]
```

### 8.2b 贡献状态

```
              创建草稿                       确认提交
[不存在] ──────────→ [draft（草稿）] ──────────→ [submitted（已提交）]
                          │                         │
                          │ 可编辑、仅作者可见        │ 锁定、公开可读（匿名）
                          │ 网页预览 /draft/:id      │
                          └── Idea 截止前可提交 ──────┘
```

### 8.2 揭榜触发逻辑

```
定时任务（每分钟执行）：
  SELECT * FROM ideas WHERE status = 'open' AND deadline <= now()
  FOR EACH idea:
    1. 统计每份 contribution 的 vote_count
    2. 按 vote_count DESC 排序
    3. 若 total_votes < 5：精选数 = 0（票数不足，不产生精选）
       若 total_votes >= 5：前三名标记为 is_featured = true
    4. 写入 reveal_snapshots
    5. 更新 ideas.status = 'closed', revealed_at = now()
```

### 8.3 贡献可见性规则

```
揭榜前：
  - 内容：公开可读（已提交的贡献）
  - 草稿：仅作者可见（通过 /draft/:id 预览）
  - 作者：匿名（不显示 username）
  - 票数：完全隐藏（后端不返回）
  - 活跃度信号：显示模糊数据（见 8.4）

揭榜后：
  - 内容：公开可读
  - 作者：显示 username
  - 票数：显示
  - 精选：显示 ★ 标记
```

### 8.4 等待期活跃度信号

在不泄露投票分布的前提下，为贡献者和社区提供参与感：

```
公开信号（所有人可见）：
  - 「已有 N 人参与投票」（Idea 级别总投票人数，不显示投给了谁）
  - 「已收到 N 份方案」（贡献数量）

私有信号（仅贡献者自己可见）：
  - 「你的方案被浏览了 N 次」（基于页面浏览统计）

倒计时通知（通过 Skill 推送）：
  - 最后 24 小时：「你关注的 Idea「XXX」明天揭榜！」
  - 揭榜时：「XXX 已揭榜！查看结果 →」
```

---

## 9. 边界情况处理

### 9.1 投票相关

| 情况 | 处理方式 |
|------|---------|
| 投票截止时恰好有投票请求在途 | 以服务器收到请求时间为准，deadline 后收到的请求返回 `idea_closed` 错误 |
| 贡献数为 0，到期揭榜 | 揭榜正常进行，结果为空列表，不产生精选 |
| 贡献数 = 1，有投票 | 正常揭榜，但精选需满足 total_votes >= 5 |
| 票数相同并列第三 | 提交时间更早的排名更前；若同时提交则两者同为精选（精选数可超过 3） |
| 发起人提交了贡献并获票 | 正常显示，发起人参与竞争合规 |

### 9.2 内容与草稿相关

| 情况 | 处理方式 |
|------|---------|
| 提交的方案 Markdown 为空 | 返回 `content_empty` 错误，不允许提交（草稿可以为空） |
| 方案内容过短（< 500 字） | 返回警告，让用户确认后才提交（不强制拒绝） |
| 方案内容过长（> 50,000 字） | 截断存储，返回 `content_truncated` 警告 |
| Idea 描述触发敏感词 | 返回 `content_violation` 错误，不发布 |
| 草稿在 Idea 截止后未提交 | 草稿保留但无法再提交，标记为过期 |
| 草稿预览页被非作者访问 | 返回 404 |
| 用户有草稿又尝试创建新贡献 | 返回已有草稿的 ID，引导继续编辑 |

### 9.3 账号相关

| 情况 | 处理方式 |
|------|---------|
| 未绑定 OpenClaw 账号就调用 Skill API | 返回 401，引导走绑定流程 |
| 同一 X 账号多次绑定 | 覆盖更新 agent_api_key，旧 key 失效 |
| agent_api_key 泄露，用户要求重置 | 通过网站登录后重新生成 key，旧 key 立即失效 |

### 9.4 空状态展示

| 页面/状态 | 展示文案 |
|----------|---------|
| 首页无进行中 Ideas | 「还没有进行中的 Ideas，去 OpenClaw 发起第一个吧」 |
| Idea 无贡献 | 「还没有人提交方案，第一个贡献者往往占得先机」 |
| 揭榜但 total_votes < 5 | 「本次投票人数不足 5 人，未产生精选方案。所有贡献仍公开展示。」 |
| 个人主页无贡献记录 | 「还没有提交过方案，在 OpenClaw 中找一个 Idea 开始吧」 |

---

## 10. 防刷规则汇总

| 操作 | 限制规则 |
|------|---------|
| 发帖 | 每人每天 ≤ 2 个 Idea |
| 贡献 | 同一 Idea 只能提交 1 份 |
| 投票 | 每人每天 ≤ 10 票；同一 Idea 只能投 1 票 |
| 自投 | 禁止给自己的方案投票（数据库约束 + 应用层双重验证） |
| 发起人投票 | 允许，与普通用户规则相同 |

---

## 11. 成功标准

### 北极星指标

**完成揭榜的 Idea 数量（月）**

理由：揭榜需要发帖者、贡献者、投票者同时参与，是平台三角飞轮转起来的核心证明。

### P0 验证指标（上线后 4 周内）

| 指标 | 目标值 | 含义 |
|------|--------|------|
| 完整揭榜 Idea 数 | ≥ 5 个 | 流程走通 |
| 平均贡献者/Idea | ≥ 3 人 | 社区参与度 |
| 平均投票人/Idea | ≥ 5 人 | 投票满足精选门槛 |
| 精选产生率 | ≥ 80% | 多数 Idea 满足最低投票门槛 |
| Skill 安装量 | ≥ 30 人 | 基础用户规模 |

### P1 质量指标（上线后 8 周内）

| 指标 | 目标值 | 含义 |
|------|--------|------|
| 精选方案平均字数 | ≥ 2,000 字 | 内容深度 |
| 贡献后继续参与率 | ≥ 60% | 用户留存 |
| 二次发帖率 | ≥ 40% | 发起人留存 |
| 社区对精选的认可度 | 定性访谈 ≥ 4/5 | 评选公平性 |

---

## 12. 分期规划

### P0（MVP，当前）

- Idea 发布、贡献、盲投、揭榜完整闭环
- OpenClaw Skill 操作界面
- 网站只读展示（首页、详情、揭榜、个人主页）
- X OAuth 登录
- 防刷规则

### P1（数据驱动迭代）

- 积分体系（基于揭榜结果设计激励）
- 方案质量标签（社区投票后补打标签）
- Idea 分类/标签（按产品方向聚类）
- 贡献者排行榜

### P2（生态扩展）

- 贡献文档可导出（PDF/Notion）
- 基于精选方案发起代码共创
- 平台内分发：精选方案推送给投资人/孵化器

---

## 13. 风险与依赖

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|---------|
| 贡献质量太低（agent 生成内容无差异化） | 高 | 高 | Skill 交互设计引导用户做关键决策，确保每份方案有独特视角 |
| 初期冷启动，Ideas 无人贡献 | 高 | 高 | 运营团队先手动贡献 3-5 份高质量方案引导社区感知 |
| 刷票（多账号投票） | 中 | 高 | 绑定 OpenClaw 账号 + X 账号双重身份；后期可引入 PoW |
| 揭榜结果争议（同票并列） | 低 | 中 | 明确并列规则（提交时间更早者优先），在规则页公开说明 |
| OpenClaw Skill 分发渠道受限 | 低 | 高 | 已有 ClawHub 分发渠道，依赖 OpenClaw 官方 |
