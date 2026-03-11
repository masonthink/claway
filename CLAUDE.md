# Claway - 文档共创平台

## 产品定位
为 OpenClaw 用户定制的文档协作平台。用户发起产品想法，社区其他用户驱动 agent 协作完成开发前文档（竞品分析、用户画像、PRD 等）。文档以积分形式销售，贡献者按算力投入获得分成。

## 核心机制
- **人驱动 Agent**：专业用户驱动 agent 完成任务，平台不做执行引擎
- **LLM Proxy 计量**：agent 通过平台代理调用 LLM，自动记录 token 消耗和美元成本
- **积分体系**：API 成本 × 质量系数 × 1000 = 积分；查看 PRD 消耗积分；积分按贡献权重分配
- **9 文档模板**：D1 竞品分析 → D2 用户画像 → D3 用户旅程 → D4 功能需求 → D5 信息架构 → D6 页面流程 → D7 交互设计 → D8 视觉设计 → D9 技术可行性

## 技术栈
| 层 | 技术 |
|---|---|
| Backend | Go 1.23 + Echo v4 + pgx/v5 |
| Database | PostgreSQL 16 |
| Frontend | Next.js 15 + Tailwind CSS |
| Plugin | TypeScript (OpenClaw Plugin) |

## 项目结构
```
claway/
├── CLAUDE.md
├── docs/
│   ├── prd-mvp-v2.md                      # 产品需求文档
│   ├── research/                           # 技术调研
│   └── architecture/                       # 架构设计
├── src/
│   ├── backend/                            # Go 后端 API
│   │   ├── cmd/server/main.go              # 入口
│   │   ├── internal/
│   │   │   ├── config/                     # 环境变量配置
│   │   │   ├── model/                      # 数据模型
│   │   │   ├── store/                      # 数据库访问层
│   │   │   ├── service/                    # 业务逻辑层
│   │   │   ├── handler/                    # HTTP 处理器
│   │   │   ├── middleware/                 # JWT 认证中间件
│   │   │   └── testutil/                   # 测试工具
│   │   ├── migrations/                     # SQL 迁移文件
│   │   └── scripts/                        # 辅助脚本
│   ├── web/                                # Next.js 前端
│   │   └── src/
│   │       ├── app/                        # 页面路由
│   │       ├── components/                 # 共享组件
│   │       └── lib/                        # API 客户端
│   └── plugin/                             # OpenClaw 插件
│       └── src/                            # 14 个 agent 工具
└── scripts/
```

## 后端 API（27 端点）
```
POST   /api/v1/ideas                  # 创建想法（自动生成子任务）
GET    /api/v1/ideas                  # 想法列表
GET    /api/v1/ideas/:id              # 想法详情
GET    /api/v1/ideas/:id/context      # 聚合上下文（给 agent 用）
GET    /api/v1/ideas/:id/tasks        # 子任务列表
GET    /api/v1/tasks/:id              # 任务详情
POST   /api/v1/tasks/:id/claim        # 认领任务
DELETE /api/v1/tasks/:id/claim        # 放弃任务
POST   /api/v1/tasks/:id/submit       # 提交产出
POST   /api/v1/tasks/:id/review       # 验收（发起人）
GET    /api/v1/tasks/:id/document     # 获取文档
PUT    /api/v1/tasks/:id/document     # 更新文档
POST   /api/v1/ideas/:id/publish      # 发布 PRD
POST   /api/v1/proxy/chat             # LLM 代理
GET    /api/v1/me/compute             # 我的算力
GET    /api/v1/me/credits             # 我的积分
POST   /api/v1/prd/:id/purchase       # 购买 PRD
```

## 本地开发
```bash
# 后端
cd src/backend
DATABASE_URL="postgres://mason@localhost:5432/claway?sslmode=disable" \
JWT_SECRET="dev-secret" \
go run ./cmd/server/

# 前端
cd src/web
NEXT_PUBLIC_API_URL="http://localhost:8080/api/v1" \
npm run dev

# 测试
cd src/backend
go test -v ./internal/
```

## 开发规范
- 中文编写文档
- 代码注释用英文
- Git commit message 用英文
