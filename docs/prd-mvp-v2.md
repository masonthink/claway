# AgentForge MVP PRD 最终版 v2.0

*2026-03-11*

---

## 1. 产品定位

**为 OpenClaw 用户定制的文档协作平台。用户发起产品想法，社区其他用户各自驱动 agent 从不同专业角度协作，完成开发前的所有文档设计工作——从竞品分析到 API 设计，产出一套可直接交付开发团队的完整文档集。**

---

## 2. 核心痛点

有想法的人很多，但独自把想法收敛成结构化、可执行的产品方案极难。原因是：

- 一个人的视角有限（缺竞品研究、缺用户洞察、缺技术视角）
- 从想法到 PRD 需要大量专业知识和时间
- AI 可以加速这个过程，但需要有专业背景的人驱动 agent 才有质量

AgentForge 把这个过程变成多人协作：每个参与者贡献自己最擅长的专业视角，用 agent 放大产出效率。

---

## 3. MVP 范围

### 做

- 想法发布与浏览（OpenClaw Skill + 只读网站）
- PRD 子任务认领与执行（Skill 操作）
- LLM Proxy 自动计量 API 成本
- 产出提交与发起人验收
- 贡献权重计算（基于 API 成本美元）
- 中心服务：文档存储、版本管理、Agent 读取与同步
- 最终 PRD 归档展示（网站只读）
- 积分体系（算力投入获积分、消耗积分查看 PRD、积分实时分配）

### 不做

- 网站上的任何操作入口
- 代码共创（第二期）
- 手动填报 token
- 社交功能（评论、点赞、关注）
- 非 OpenClaw 用户接入
- 实时结算（只做月结）
- 复杂仲裁机制

---

## 4. 核心流程

```
【发起】
  发起人在 OpenClaw 中描述产品想法（自然语言）
  Skill 引导补充：目标用户、大致方向、发起人抽成比例
  平台创建 Idea，生成标准子任务列表
  Idea 在网站公开，状态：招募中

【认领】
  贡献者在 OpenClaw 中浏览开放 Idea
  选择认领某个子任务（如"竞品分析"）
  Skill 返回：任务详情 + 验收标准 + LLM Proxy 接入参数

【执行】
  贡献者在本地驱动 agent，通过平台 LLM Proxy 调用 LLM
  平台自动记录：模型、tokens_in、tokens_out、换算美元成本
  执行过程中平台只感知 API 调用，不感知对话内容

【提交】
  贡献者在 OpenClaw 中提交产出
  内容：Markdown 文档 + 关键决策说明（200字以内）
  发起人收到通知

【验收】
  发起人在 OpenClaw 中查看提交内容
  打分：达标(1.0) / 良好(1.2) / 优秀(1.5)
  打回需填写原因，贡献者可修改重新提交（最多3次）

【合并归档】
  所有子任务通过验收
  平台将各子任务产出合并为完整 PRD（Markdown）
  PRD 在网站上公开展示，永久归档
  贡献权重快照冻结

【积分流通】
  PRD 在网站公开展示（免费预览摘要）
  其他用户花积分解锁完整 PRD
  积分实时按权重分配给贡献者
  贡献者的积分可用于阅读其他 PRD
```

---

## 5. 文档体系与任务模板

### 定位

系统覆盖**开发工作前的所有文档设计工作**。产出不是一份 PRD，而是一套完整的、可直接交付开发团队的文档集。

### 发起流程

发起人在 OpenClaw 中说出模糊想法（如"我想做一个线下预定系统"），Skill + agent 引导他完成问题定义：
- 核心问题陈述（≤50字）
- 目标用户方向
- 不做什么

问题定义完成后，作为所有后续文档的上下文输入，D1、D2 立即开放认领。

### 三个阶段、9 种文档

#### 阶段 1：探索期（可并行）

| # | 文档 | 目标 | 所需专业背景 | token上限 |
|---|------|------|-------------|-----------|
| D1 | 竞品分析 | 市场现状、竞品优劣势、差异化机会 | 产品/市场 | 80k |
| D2 | 用户研究 | 用户画像、核心痛点、使用场景 | 用户研究/产品 | 60k |

#### 阶段 2：定义期（D1+D2 完成后开始，D3-D5 可并行）

| # | 文档 | 目标 | 所需专业背景 | token上限 |
|---|------|------|-------------|-----------|
| D3 | 产品需求文档(PRD) | 用户故事、验收标准、P0/P1优先级 | 产品 | 100k |
| D4 | 商业模式设计 | 收入模式、定价策略、冷启动路径 | 产品/商业 | 60k |
| D5 | 成功指标体系 | 北极星指标、过程指标、度量方案 | 产品/数据 | 40k |

#### 阶段 3：设计期（D3 完成后开始，D6-D9 可并行）

| # | 文档 | 目标 | 所需专业背景 | token上限 |
|---|------|------|-------------|-----------|
| D6 | 信息架构 | 页面结构、导航体系、权限矩阵 | 产品/交互设计 | 60k |
| D7 | 核心用户流程 | 主流程+异常流程，覆盖所有P0功能 | 交互设计 | 80k |
| D8 | 设计规范 | 组件库选型、色板字体、定制组件清单 | UI设计 | 60k |
| D9 | 技术可行性评估 | 技术选型建议、关键风险点、可行性结论 | 技术/架构 | 60k |

### 文档依赖 DAG

```
发起想法 → 问题定义（发起流程自动完成）
                │
                ├→ D1 竞品分析 ─┐
                └→ D2 用户研究 ─┤
                                │
                                ├→ D3 PRD ──┬→ D6 信息架构
                                │           ├→ D7 用户流程
                                ├→ D4 商业模式  ├→ D8 设计规范
                                └→ D5 成功指标  └→ D9 技术可行性
                                                    │
                                                    ▼
                                            可以开始写代码
```

关键路径：想法 → D2 → D3 → D6 → **可以开始写代码**

### 验收标准

| 文档 | 核心验收要求 |
|------|-------------|
| D1 竞品分析 | ≥3直接竞品+≥2间接竞品，含差异化空间分析 |
| D2 用户研究 | 2-3个画像，每个≥2个叙述性场景，当前方案局限性 |
| D3 PRD | 用户故事格式，每个功能有验收标准，P0≤10个 |
| D4 商业模式 | 收入模式+定价依据+冷启动路径 |
| D5 成功指标 | 北极星指标1个+过程指标3-5个，每个有目标值和度量方案 |
| D6 信息架构 | 完整页面列表（路由+内容）+导航结构+权限矩阵 |
| D7 用户流程 | 覆盖所有P0功能，正常路径+≥2异常路径，含空/错/加载状态 |
| D8 设计规范 | 组件库选型+色板字体间距+定制组件清单 |
| D9 技术可行性 | 技术选型建议(附理由)+关键风险点+明确的可行/不可行结论 |

### 两种预设套餐

**轻量套餐**（5个文档）— 简单工具、验证想法

```
D1, D2 → D3 → D7, D9
```

**标准套餐**（9个文档）— 完整产品方案

全部 9 个文档。

### 收敛机制

- **必须选择套餐**，不允许完全自定义文档组合
- **平台强制依赖顺序**，前置文档未完成时后续文档不开放认领
- **验收时逐条确认**，Skill 提示发起人对照验收标准，不能只打"通过"

---

## 6. 中心服务：文档存储、版本管理与 Agent 同步

中心服务是平台的核心基础设施，负责所有产出文档的统一管理，让多个 agent 在协作过程中能读取彼此的产出。

### 核心职责

```
文档保存    每个子任务的产出（Markdown）统一存储在中心服务
版本管理    每次提交/修改生成新版本，可追溯历史，支持 diff 查看
Agent 读取  贡献者的 agent 通过 API 拉取已完成任务的产出作为上下文
状态同步    多个 agent 同时工作时，能读到最新的项目状态和其他人的产出
```

### 实际场景

```
贡献者 B 认领「核心功能设计」任务
  → Skill 自动拉取已完成的「竞品分析」和「用户场景」产出
  → 作为上下文注入 agent，agent 基于这些信息开展工作
  → B 的 agent 产出过程中可随时调用 API 查看最新状态
  → B 提交后，后续认领「信息架构」的贡献者 C 也能读到 B 的产出
```

### 文档 API

```
# 文档读取（agent 调用）
GET    /api/v1/tasks/{id}/document                # 获取任务最新产出文档
GET    /api/v1/tasks/{id}/document/versions        # 获取版本列表
GET    /api/v1/tasks/{id}/document/versions/{ver}   # 获取指定版本

# 项目上下文（agent 调用，获取所有已完成任务的产出摘要）
GET    /api/v1/ideas/{id}/context                  # 返回该 Idea 所有已完成任务的产出

# 文档写入（提交时由平台自动处理）
POST   /api/v1/tasks/{id}/submit                   # 提交产出时自动创建新版本
PUT    /api/v1/tasks/{id}/document                  # 更新文档（被打回后修改重提交）
```

### 数据模型补充

```
Document
  id, task_id             FK: Task
  content                 Markdown 全文
  version                 自增版本号
  created_at
  created_by              FK: User

DocumentVersion           版本历史
  id, document_id         FK: Document
  version                 版本号
  content                 该版本的完整内容
  diff_from_previous      与上一版本的 diff
  created_at
```

### 设计原则

- **只读优先**：agent 对其他任务的产出只有读权限，写权限仅限自己认领的任务
- **自动注入上下文**：认领任务时 Skill 自动拉取相关已完成产出，贡献者无需手动操作
- **版本不可删除**：所有版本永久保留，确保可追溯
- **最终一致性**：agent 读取到的可能不是最新提交（短暂延迟可接受），但保证最终一致

---

## 7. 交互界面

### 双界面模型

- **网站 = 展示 + 付费**（项目列表、进度、算力数据、PRD 展示、付费购买），无共创操作入口
- **OpenClaw 聊天窗口 = 共创操作界面**，所有共创操作通过 OpenClaw Skill 完成

### 网站登录

通过 **OpenClaw OAuth** 登录，用户用 OpenClaw 账号授权。

```
公开页面（无需登录）         登录后可见
─────────────────         ──────────────
项目列表                   我的算力投入与积分
项目详情（进度、算力数据）   花积分解锁完整 PRD
PRD 免费预览               积分流水明细
个人公开主页
平台算力总览
```

### 人机协作模型

- 人是主体，agent 是工具
- 有专业能力的人驱动 agent 调整方向，更好地完成工作
- 平台不感知人和 agent 之间的对话过程
- 人全程参与执行，持续纠偏和决策

### 网站页面

#### 首页

```
AgentForge
用 AI 算力协作，把想法变成产品方案

[平台算力总览]
  累计参与用户 42 人 · 累计算力投入 $128.7 · 完成项目 8 个 · 累计 LLM 调用 12,340 次

[招募中]                          [已完成]

招募中的 Idea
┌───────────────────────────────────────────────────────────────┐
│ 想法简述              开放任务   参与人数   已投入算力  发起人  │
│ 帮自由职业者报税的工具  4/7 个    2 人      $1.23     @alice   │
│ 独立开发者用户反馈系统  7/7 个    0 人      $0.00     @bob     │
└───────────────────────────────────────────────────────────────┘

已完成的 PRD
┌───────────────────────────────────────────────────────────────┐
│ 产品名称           完成时间   贡献者   总算力投入   查看        │
│ 极简番茄钟 SaaS     2026-03   5 人    $4.23       [阅读 PRD]  │
└───────────────────────────────────────────────────────────────┘

如何参与：在 OpenClaw 中安装 AgentForge Skill
```

#### Idea 详情页 `/ideas/{id}`

```
[想法标题]
发起人 @alice · 发起于 2026-03-08 · 截止 2026-03-22

[想法描述]
我想做一个帮自由职业者管理收入和自动计算税款的工具...

算力投入总览
  总投入 $1.57 · LLM 调用 286 次 · 模型分布：Sonnet 78% / Opus 22%

子任务进度
┌─────────────────────────────────────────────────────────────────────┐
│ 任务           状态      认领者   算力投入   调用次数  模型     质量分 │
│ 竞品分析       已完成    @carol   $0.82     142次    Sonnet   优秀1.5│
│ 用户场景       进行中    @dave    $0.34     89次     Sonnet   —      │
│ 核心功能设计   待认领    —        —          —        —        —      │
│ 信息架构       待认领    —        —          —        —        —      │
│ 数据模型       已完成    @alice   $0.41     55次     Opus     良好1.2│
│ 商业模式       待认领    —        —          —        —        —      │
│ 技术方案       待认领    —        —          —        —        —      │
└─────────────────────────────────────────────────────────────────────┘

贡献者算力分布（实时）
┌─────────────────────────────────────────────────────────────┐
│ 贡献者    算力投入   占比    加权权重（含质量分）   预估分成比 │
│ @carol    $0.82     52%     45%                   45%       │
│ @dave     $0.34     22%     20%（进行中）          —         │
│ @alice    $0.41     26%     35%                   35%       │
└─────────────────────────────────────────────────────────────┘
```

#### PRD 展示页 `/prd/{id}`

```
[产品名称] PRD                              [解锁完整版 580 积分]

贡献者：@alice @carol @dave @eve @frank
完成于 2026-03-18 · 总算力投入 $4.23 · 已被 12 人解锁

目录
  1. 竞品分析        @carol 完成  ★优秀
  2. 用户场景        @dave 完成   ★良好
  ...

[免费预览：竞品分析摘要]
...前300字内容展示...

[花 580 积分解锁完整 PRD]    积分不够？投入算力或购买积分
```

#### 个人主页 `/users/{username}`

```
@username

算力投入总览
  总算力投入 $12.34 · 累计 LLM 调用 2,180 次 · 参与项目 5 个 · 平均质量分 1.3

模型使用分布
  ████████████░░░  Claude Sonnet  68%  $8.39
  ████░░░░░░░░░░░  Claude Opus    25%  $3.09
  ██░░░░░░░░░░░░░  GPT-4o          7%  $0.86

算力投入趋势（最近 30 天）
  03/01 ▎$0.42
  03/02 ▎▎$0.81
  03/03 ▎▎▎$1.23
  ...

参与记录
┌─────────────────────────────────────────────────────────────┐
│ Idea 名称     角色     任务       算力投入  调用次数  质量分  │
│ 自由职业报税  贡献者   竞品分析   $0.82    142次    1.5     │
│ 用户反馈系统  发起人   数据模型   $0.41    55次     1.2     │
└─────────────────────────────────────────────────────────────┘
```

---

## 7. OpenClaw Skill 命令设计

### 发起人操作

```
「发起想法」→ 引导输入想法描述、目标用户、抽成比例 → 自动生成6个标准子任务 → 发布
「查看我发起的想法」→ 列出所有 Idea 及各任务状态
「验收提交」→ 展示待验收内容 → 打分（达标/良好/优秀）或打回
「合并发布 PRD」→ 所有任务通过后，合并产出发布
```

### 贡献者操作

```
「看看有什么想法」→ 列出招募中的 Idea
「查看 [Idea名称] 的任务」→ 展示子任务详情和验收标准
「认领 [任务名称]」→ 返回任务详情 + LLM Proxy 接入参数 + 已完成任务的产出摘要
「提交任务」→ 提交 Markdown 产出 + 关键决策说明
「我的算力」→ 总投入、模型分布、各项目投入明细、实时权重
「我的积分」→ 当前余额、获得/消耗明细、各项目分成记录
「我的贡献」→ 各 Idea 的 API 成本、权重、积分收入
```

### 首次使用

```
「配置 AgentForge」→ 绑定 OpenClaw 账号 → 生成 Agent API Key → 输出 LLM Proxy 配置
```

---

## 9. 平台 API 端点清单

所有请求：`Authorization: Bearer {agent_api_key}`

```
# Idea 管理
POST   /api/v1/ideas                           # 发起想法，自动生成子任务
GET    /api/v1/ideas                           # 列出 Idea（?status=recruiting）
GET    /api/v1/ideas/{id}                      # Idea 详情含任务列表

# 任务操作
GET    /api/v1/ideas/{id}/tasks                # 列出子任务
GET    /api/v1/tasks/{id}                      # 任务详情+验收标准+相关已完成产出
POST   /api/v1/tasks/{id}/claim                # 认领任务
DELETE /api/v1/tasks/{id}/claim                # 放弃任务
POST   /api/v1/tasks/{id}/submit               # 提交产出
POST   /api/v1/tasks/{id}/review               # 验收（仅发起人）

# PRD 发布
POST   /api/v1/ideas/{id}/publish              # 合并发布完整 PRD（仅发起人）

# 文档与上下文（agent 读取）
GET    /api/v1/tasks/{id}/document             # 获取任务最新产出文档
GET    /api/v1/tasks/{id}/document/versions    # 获取版本列表
GET    /api/v1/tasks/{id}/document/versions/{ver}  # 获取指定版本
GET    /api/v1/ideas/{id}/context              # 获取该 Idea 所有已完成任务的产出摘要
PUT    /api/v1/tasks/{id}/document             # 更新文档（打回后修改重提交）

# LLM Proxy
POST   /api/v1/proxy/chat                      # 代理 LLM 调用（Header: X-Task-ID）

# 算力投入查询
GET    /api/v1/me/compute                      # 我的算力总览（总投入、模型分布、趋势）
GET    /api/v1/me/compute/ideas/{id}           # 我在某个 Idea 的算力明细
GET    /api/v1/ideas/{id}/compute              # 某个 Idea 的算力投入总览（各贡献者）
GET    /api/v1/tasks/{id}/compute              # 某个任务的算力消耗明细
GET    /api/v1/platform/compute                # 平台算力总览（公开数据）

# 积分与贡献
GET    /api/v1/me/credits                      # 积分余额和流水明细
GET    /api/v1/me/contributions                # 我的贡献列表
POST   /api/v1/prd/{id}/purchase               # 花积分购买阅读 PRD
```

---

## 10. LLM Proxy 设计

```
调用链路：
  贡献者 agent
    → POST /api/v1/proxy/chat
      Header: Authorization: Bearer {agent_api_key}
      Header: X-Task-ID: {task_id}
      Body: 标准 OpenAI 格式
    → 平台 Proxy 层
      1. 验证 api_key，确认 task_id 归属当前用户
      2. 确认任务状态为「进行中」
      3. 转发到真实 LLM API
      4. 记录：user_id, task_id, model, tokens_in, tokens_out
         cost_usd = tokens_in × 输入单价 + tokens_out × 输出单价
      5. 累加到 task.cost_usd_accumulated
      6. 原始响应透传给 agent

模型定价表（内置，定期维护）：
  claude-opus-4       输入 $15/M token    输出 $75/M token
  claude-sonnet-4-5   输入 $3/M token     输出 $15/M token
  gpt-4o              输入 $2.5/M token   输出 $10/M token
  gpt-4o-mini         输入 $0.15/M token  输出 $0.6/M token

费用由贡献者自己承担，平台只做计量记录，不代付。
```

---

## 11. 数据模型

```
User
  id, openclaw_id, username
  agent_api_key
  credits_balance          当前积分余额

Idea
  id, title, description, target_user_hint
  initiator_id            FK: User
  initiator_cut_percent   10-30
  status                  recruiting / in_progress / completed / cancelled
  created_at, deadline

Task
  id, idea_id             FK: Idea
  type                    D1_competitive_analysis / D2_user_research /
                          D3_prd / D4_business_model / D5_success_metrics /
                          D6_information_architecture / D7_user_flows /
                          D8_design_spec / D9_tech_feasibility
  dependencies            依赖的文档 type 列表（平台自动填充）
  title, description, acceptance_criteria
  token_limit_hint        参考上限（不强制截断）
  status                  open / claimed / submitted / approved / rejected
  claimed_by              FK: User
  claimed_at, submitted_at, approved_at
  output_content          Markdown 全文
  output_note             贡献者说明（≤200字）
  quality_score           1.0 / 1.2 / 1.5
  reject_reason
  cost_usd_accumulated    实时累计 API 成本（美元）

TokenUsageLog
  id, user_id, task_id
  model
  tokens_in, tokens_out
  cost_usd
  timestamp

Document
  id, task_id             FK: Task
  content                 Markdown 全文（最新版）
  current_version         当前版本号
  created_at, updated_at

DocumentVersion           版本历史
  id, document_id         FK: Document
  version                 版本号（自增）
  content                 该版本完整内容
  diff_from_previous      与上一版本的 diff
  created_at
  created_by              FK: User

Contribution             任务通过验收后生成
  id, idea_id, task_id, user_id
  cost_usd
  quality_score
  weighted_score         cost_usd × quality_score
  weight_percent         归一化后的权重（月结时冻结）

PRD
  id, idea_id
  content                合并后的完整 Markdown
  published_at
  price_credits          阅读定价（积分）
  read_count             已被解锁次数

CreditTransaction        积分流水
  id
  user_id                FK: User
  type                   earn_contribute / earn_read_share / spend_read / earn_topup
  amount                 变动积分数（正=增加，负=消耗）
  reference_type         task / prd / topup
  reference_id           关联的 task_id 或 prd_id
  description            可读描述
  created_at
```

---

## 12. 积分体系

MVP 阶段不引入支付系统，用**积分**作为平台内的价值流通单位。

### 核心循环

```
投入算力 → 获得积分 → 消耗积分查看他人 PRD → 驱动更多人投入算力
```

### 积分获取

```
贡献算力获得积分：
  积分 = API 成本（美元）× 质量系数(1.0/1.2/1.5) × 1000

  示例：
    消耗 $0.82 API 成本，质量评分「优秀」1.5
    获得积分 = 0.82 × 1.5 × 1000 = 1,230 积分

花钱购买积分（后续开放）：
  $1 = 1,000 积分（与算力获取等价，不打折，保护贡献者价值）
```

### 积分消耗

```
查看完整 PRD 文档：消耗 N 积分（由发起人设定，平台给建议值）
  建议定价 = 项目总积分产出 × 10%（即贡献者集体投入的 10%）

  示例：
    一个 PRD 总共消耗 $4.23 API 成本，质量加权后产出 5,800 积分
    建议阅读定价 = 580 积分
```

### 积分分配（PRD 被购买阅读时）

```
每次有人花积分查看 PRD：
  平台抽成    = 10%
  发起人抽成  = 发起时设定（10%-30%）
  贡献者池    = 剩余部分，按贡献权重分配

贡献权重 = 该用户加权分数 / 项目总加权分数
加权分数 = API 成本（美元）× 质量系数

示例：
  某人花 580 积分查看 PRD
  平台 58 积分，发起人抽成 20% = 116 积分，贡献者池 406 积分

  贡献者 A（权重 45%）：获得 183 积分
  贡献者 B（权重 20%）：获得 81 积分
  贡献者 C（权重 35%）：获得 142 积分
```

### 数据模型

```
User 新增字段
  credits_balance         当前积分余额

CreditTransaction         积分流水（每笔变动一条记录）
  id
  user_id                 FK: User
  type                    earn_contribute / earn_purchase / spend_read / earn_share
  amount                  变动积分数（正数增加，负数消耗）
  reference_type          task / prd / topup
  reference_id            关联的 task_id 或 prd_id
  description             可读描述
  created_at

PRD 新增字段
  price_credits           阅读定价（积分）
  read_count              已被阅读次数
```

### 设计原则

- **算力投入 = 积分来源**，贡献越多积分越多，不花钱也能获取
- **花钱买积分是补充手段**，不是主要路径，保护贡献者的劳动价值
- **积分不可提现**（MVP 阶段），避免金融合规问题
- **所有积分流水可追溯**，用户随时查看明细

---

## 13. 分期规划

### 第一期（本期 MVP）：PRD 共创

- 交付物：结构化产品方案文档
- 验证：用户愿意参与、分成机制公平、流程走通
- 成功标准：≥3 个 Idea 完整闭环，≥3 人/Idea，公平性满意度 ≥ 4/5

### 第二期：代码共创

- 交付物：可运行的 Web 产品（前端 + 轻量后端 + 数据库）
- 新增：产品模板体系、技术栈标准化（Next.js + PostgreSQL）、平台自动构建部署、托管销售
- 复用：LLM Proxy、分成机制、OpenClaw Skill 框架
