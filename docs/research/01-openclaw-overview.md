# OpenClaw 平台概述

## 平台定位

OpenClaw 是一个**开源、自托管的个人 AI 助手网关平台**，充当中心化网关，连接 20+ 聊天渠道（WhatsApp、Telegram、Discord、Slack、Signal、iMessage 等）到 AI 编码代理（Claude Code、Codex、Gemini CLI 等）。

- GitHub: `github.com/openclaw/openclaw`
- 文档: `docs.openclaw.ai`
- 许可证: MIT

## 核心架构特征

### Gateway 是单用户信任边界
> "not a hostile multi-tenant security boundary for multiple adversarial users sharing one agent/gateway"

每个 Gateway 是一个用户的私有领地，一个受信操作者。

### Agent 间通信默认关闭
`tools.agentToAgent` 需显式启用 + 白名单，且**仅限同一 Gateway 内的 Agent**。

### 无跨 Gateway 通信机制
文档未提及任何 Gateway 间原生通信能力。**这是我们需要解决的核心问题。**

## 核心协议

| 协议 | 端点 | 用途 |
|------|------|------|
| WebSocket + JSON | `ws://127.0.0.1:18789` | Gateway 核心 RPC 通信 |
| HTTP Webhook | `POST /hooks/wake`, `/hooks/agent` | 外部系统触发 Agent |
| OpenAI 兼容 API | `POST /v1/chat/completions` | 标准对话接口（默认禁用） |
| Tools Invoke API | `POST /tools/invoke` | HTTP 调用工具 |
| ACP | `sessions_spawn(runtime:"acp")` | 集成外部编码工具 |
| SSE | `text/event-stream` | 流式响应 |

## 技术栈

| 组件 | 技术 |
|------|------|
| 核心语言 | TypeScript |
| 运行时 | Node.js >= 22 |
| 包管理 | pnpm (monorepo) |
| 测试 | Vitest |
| 类型定义 | `@sinclair/typebox` |
| 默认端口 | 18789 |
| 配置文件 | `~/.openclaw/openclaw.json` (JSON5) |

## 三种扩展机制

### 1. Skills（SKILL.md）
- Markdown + YAML frontmatter
- 无需编写代码
- 通过 ClawHub 分发

### 2. Plugins（TypeScript）
- **in-process 运行**，与 Gateway 同一进程
- 可注册：工具、HTTP 路由、后台服务、CLI 命令、生命周期钩子
- 可发起网络请求（不受 Lobster 沙盒限制）
- 分发：npm 包或本地路径

### 3. 外部 HTTP API
- OpenAI 兼容 API
- Webhook
- Tools Invoke API

## 认证方式

| 方式 | 说明 | 适合场景 |
|------|------|---------|
| Token | Bearer Token，operator 级权限 | 本地/可信环境 |
| Password | 密码认证 | 同上 |
| Trusted Proxy | 反向代理身份传递 | Tailscale 等 |

**注意**: 所有认证方式都是"全有或全无"，没有细粒度 scope。

## 多 Agent 架构

- 每个 Agent 独立的 workspace、state directory、session store
- Bindings 配置做消息路由（channel/accountId/peer 匹配）
- `sessions_send` 工具做 Agent 间消息传递（同 Gateway）
- Broadcast Groups 让多个 Agent 同时处理同一消息（实验性，仅 WhatsApp）

## 关键限制

1. 不支持跨 Gateway Agent 通信
2. 认证无细粒度权限控制
3. Gateway 默认绑定 loopback
4. Broadcast Groups 仅 WhatsApp
5. ACP 后端规范文档不完整
