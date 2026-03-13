# Claway - 产品方案投标平台

## 产品定位
OpenClaw 用户的产品方案投标平台。用户发起产品想法（Idea），社区其他用户驱动 Agent 协作完成完整的产品方案文档（Contribution），通过盲投（blind voting）评选前三名精选方案。

## 核心机制 (v3)
- **投标模式**：用户各自独立提交完整产品方案，盲投期间不可见投票分布
- **人驱动 Agent (A mode)**：Agent 生成方案草稿，人做关键决策确认
- **草稿机制**：贡献经历 draft → submitted，提交后不可修改
- **揭榜规则**：截止后自动揭榜，需 ≥5 票才有精选，前 3 名按票数排名
- **随机排序**：盲投期间贡献随机排列，消除先发优势
- **每个想法每人只能投一票，禁止自投**

## 技术栈
| 层 | 技术 |
|---|---|
| Backend | Go 1.23 + Echo v4 + pgx/v5 |
| Database | PostgreSQL 16 |
| Frontend | Next.js 15 + TailwindCSS 4 |
| Plugin | TypeScript (OpenClaw Skill, 17 tools) |

## 项目结构
```
claway/
├── CLAUDE.md
├── docs/
│   ├── prd-mvp-v3.md                      # 产品需求文档 (v3)
│   ├── research/                           # 技术调研
│   └── architecture/                       # 架构设计
├── src/
│   ├── backend/                            # Go 后端 API
│   │   ├── cmd/server/main.go              # 入口 + 揭榜定时任务
│   │   ├── internal/
│   │   │   ├── config/                     # 环境变量配置
│   │   │   ├── model/                      # 数据模型 (User, Idea, Contribution, Vote, RateLimit, RevealSnapshot)
│   │   │   ├── store/                      # 数据库访问层 (idea, contribution, vote, rate_limit, reveal, user)
│   │   │   ├── service/                    # 业务逻辑层 (idea, contribution, vote, reveal, auth)
│   │   │   ├── handler/                    # HTTP 处理器 (idea, contribution, vote, user, stats, auth)
│   │   │   ├── middleware/                 # JWT 认证中间件
│   │   │   └── testutil/                   # 测试工具
│   │   ├── migrations/                     # SQL 迁移文件 (001-005)
│   │   └── scripts/                        # 辅助脚本
│   ├── web/                                # Next.js 前端
│   │   └── src/
│   │       ├── app/                        # 页面路由 (/, /idea/[id], /idea/[id]/result, /draft/[id], /user/[username])
│   │       ├── components/                 # 共享组件
│   │       └── lib/                        # API 客户端
│   └── plugin/                             # OpenClaw 插件
│       └── src/                            # 17 个 agent 工具
└── scripts/
```

## 后端 API
```
# 认证
GET    /api/v1/auth/x                       # X OAuth 登录
GET    /api/v1/auth/x/callback              # X OAuth 回调
POST   /api/v1/auth/openclaw/callback       # OpenClaw OAuth 回调

# 公开 API
GET    /api/v1/public/stats                 # 平台统计
GET    /api/v1/public/ideas                 # 想法列表
GET    /api/v1/public/ideas/:id             # 想法详情
GET    /api/v1/public/ideas/:id/contributions  # 贡献列表（盲投匿名）
GET    /api/v1/public/ideas/:id/result      # 揭榜结果
GET    /api/v1/public/users/:username       # 用户资料

# 认证 API
POST   /api/v1/ideas                        # 创建想法
GET    /api/v1/me/ideas                     # 我的想法
POST   /api/v1/ideas/:id/contributions      # 创建草稿
PUT    /api/v1/contributions/:id            # 更新草稿
POST   /api/v1/contributions/:id/submit     # 提交锁定
GET    /api/v1/contributions/:id            # 获取贡献
GET    /api/v1/me/contributions             # 我的贡献
POST   /api/v1/ideas/:id/vote              # 投票
GET    /api/v1/me/votes                     # 我的投票
GET    /api/v1/draft/:contribution_id       # 草稿预览（作者）
```

## 数据模型
- **User**: id, openclaw_id, username, display_name, avatar_url
- **Idea**: id, initiator_id, title, description, target_user, core_problem, out_of_scope, status (open/closed/cancelled), deadline, revealed_at
- **Contribution**: id, idea_id, author_id, content (markdown), decision_log (jsonb), status (draft/submitted), view_count, submitted_at
- **Vote**: id, idea_id, voter_id, contribution_id (UNIQUE idea_id+voter_id)
- **RateLimit**: user_id, action (post_idea/vote), action_date, count
- **RevealSnapshot**: idea_id, ranked_results (jsonb), total_votes, revealed_at

## 本地开发
```bash
# 后端
cd src/backend
DATABASE_URL="postgres://mason@localhost:5432/claway?sslmode=disable" \
JWT_SECRET="dev-secret" \
go run ./cmd/server/

# 前端
cd src/web
NEXT_PUBLIC_API_URL="http://localhost:8081/api/v1" \
npm run dev

# 测试
cd src/backend
go test -v ./internal/store/... ./internal/service/...
```

## 部署架构
- **后端**: GitHub Actions → Docker → GHCR → SSH 部署 VPS
- **前端**: Vercel 自动部署，自定义域名 claway.cc
- **VPS**: 45.32.57.146 (Vultr Tokyo)
- **API 域名**: api.claway.cc
- **前端域名**: claway.cc (Vercel)
- **反向代理**: Cloudflare → Caddy → localhost:8081
- **数据库**: PostgreSQL 16 (Docker 容器 claway-postgres)
- **部署路径**: /opt/claway/
- **SSH 密钥**: ~/.ssh/dtc_deploy_vps
- **Cloudflare SSL**: Flexible 模式
- **认证**: X (Twitter) OAuth 2.0 + PKCE + OpenClaw OAuth

## 开发规范
- 中文编写文档
- 代码注释用英文
- Git commit message 用英文
