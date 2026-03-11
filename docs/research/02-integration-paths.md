# 间接实现路径分析

OpenClaw 不原生支持跨 Gateway Agent 协作。以下是 5 条间接实现路径，从简单到复杂排列。

---

## 路径一：Webhook 双向桥接（最快可验证）

### 原理
每个用户的 Gateway 有 `/hooks/agent` 端点，平台作为中间人向各 Gateway 发 webhook。

```
用户A Gateway                协调平台                  用户B Gateway
/hooks/agent  ◄── HTTP POST ── 协调服务 ── HTTP POST ──► /hooks/agent
```

### 关键 API

```bash
# 向 Agent 发送任务
curl -X POST http://gateway:18789/hooks/agent \
  -H 'Authorization: Bearer SECRET' \
  -H 'Content-Type: application/json' \
  -d '{
    "message": "实现用户登录 API",
    "agentId": "backend",
    "sessionKey": "hook:collab:task-001",
    "deliver": false,
    "timeoutSeconds": 120
  }'
```

### 参数说明
- `message`: 发给 Agent 的指令（必填）
- `agentId`: 指定 Agent
- `sessionKey`: 隔离不同任务上下文
- `deliver`: false 表示不发到聊天频道
- `wakeMode`: "now" 立即执行
- `model`: 可覆盖模型
- `timeoutSeconds`: 超时

### 问题
- 回复是异步的，没有直接 HTTP 回调
- 需要 Gateway 可达（公网或 Tailscale）
- 需要配合其他方式读取 Agent 回复

### 用户配置
```json5
{
  hooks: {
    enabled: true,
    token: "shared-secret",
    allowedAgentIds: ["main", "backend", "frontend"],
    defaultSessionKey: "hook:ingress",
    allowRequestSessionKey: true,
    allowedSessionKeyPrefixes: ["hook:collab:"]
  }
}
```

---

## 路径二：Plugin + 出站回调（核心方案）⭐

### 原理
OpenClaw 插件 **in-process 运行**，可以：
1. `api.registerTool()` — 注册 Agent 可调用的工具（execute 内可发网络请求）
2. `api.registerHttpRoute()` — 注册自定义 HTTP 端点
3. `api.registerService()` — 运行后台服务（如 WebSocket 长连接）
4. `api.on("before_prompt_build")` — 注入系统上下文

### 插件核心设计

```typescript
export default function(api) {
  // 1. 后台服务：维持到平台的 WebSocket 长连接
  api.registerService({
    id: "collab-bridge",
    start: () => {
      // 建立 WebSocket 到 platform.example.com
      // Gateway 不需要暴露端口！
    },
    stop: () => { /* 断开 */ }
  });

  // 2. Agent 工具：提交协作产出
  api.registerTool({
    name: "collab_submit",
    description: "提交协作任务的产出到平台",
    parameters: { /* TypeBox schema */ },
    execute: async (_id, params) => {
      // 插件内可以 fetch() 外部 URL
      await fetch("https://platform/api/v1/tasks/submit", {
        method: "POST",
        body: JSON.stringify(params)
      });
      return { content: [{ type: "text", text: "已提交" }] };
    }
  });

  // 3. Agent 工具：获取协作上下文
  api.registerTool({
    name: "collab_context",
    description: "获取当前协作任务的最新上下文",
    parameters: {},
    execute: async () => {
      const ctx = await fetch("https://platform/api/v1/tasks/current/context");
      return { content: [{ type: "text", text: await ctx.text() }] };
    }
  });

  // 4. 生命周期钩子：自动注入协作上下文
  api.on("before_prompt_build", (event, ctx) => {
    return {
      prependSystemContext: "你正在参与协作任务 task-001，你的角色是后端开发...",
      appendSystemContext: "前端 Agent 已完成登录页面，API 接口定义见..."
    };
  }, { priority: 10 });

  // 5. HTTP 端点：接收平台下发的任务（兜底）
  api.registerHttpRoute({
    path: "/collab/task",
    auth: "plugin",
    handler: async (req, res) => {
      res.statusCode = 200;
      res.end(JSON.stringify({ ok: true }));
      return true;
    }
  });
}
```

### 插件 Manifest

```json
{
  "id": "claway-collab",
  "configSchema": {
    "type": "object",
    "properties": {
      "platformUrl": { "type": "string" },
      "apiKey": { "type": "string" }
    },
    "required": ["platformUrl", "apiKey"]
  },
  "uiHints": {
    "platformUrl": { "label": "Platform URL" },
    "apiKey": { "label": "API Key", "sensitive": true }
  }
}
```

### 优势
- Gateway 不需要暴露端口（插件主动出站连接）
- Agent 通过工具自然地与平台交互
- 上下文注入让 Agent 理解协作背景
- 插件安装简单：`openclaw plugins install @claway/collab`

---

## 路径三：OpenAI 兼容 API 作为统一接口

### 原理
Gateway 的 `POST /v1/chat/completions` 是标准 OpenAI 格式。

```python
from openai import OpenAI

client = OpenAI(
    base_url="https://user-a-gateway.tailnet/v1",
    api_key="user-a-gateway-token"
)

response = client.chat.completions.create(
    model="openclaw:frontend-agent",  # 选择具体 agent
    messages=[
        {"role": "system", "content": "你是协作任务的前端开发者..."},
        {"role": "user", "content": "实现登录页面组件"}
    ],
    stream=True,
    user="collab:task-001"  # 持久 session key
)
```

### 关键细节
- `model` 字段用 `openclaw:<agentId>` 选择 Agent
- `user` 字段创建持久 session
- 支持 SSE 流式
- 需要显式启用，默认禁用
- Token 是 operator 级权限（安全风险）

### 适用场景
- 快速原型验证
- 平台需要同步获取 Agent 回复时

---

## 路径四：ACP 协议反向利用

### 原理
把平台伪装成 ACP 后端，让 Gateway 主动连接。

```
用户 Gateway ──sessions_spawn(runtime:"acp")──► 平台（伪装 ACP 后端）
                                                    │
                                                    ├── 接收任务
                                                    ├── 路由到其他用户 Agent
                                                    └── 回传结果 (streamTo:"parent")
```

### 关键细节
- `persistent` 模式可跨多轮保持状态
- `streamTo: "parent"` 实时回传进度
- 可绑定到 Discord/Telegram 线程
- 运行在宿主系统，不受沙盒限制

### 风险
- acpx 协议规范文档不完整，需逆向研究
- 实现复杂度高

---

## 路径五：共享聊天频道作为消息总线（最简但最 hacky）

### 原理
多个用户 Gateway 接入同一个 Telegram Group 或 Discord Channel。

```
User A Gateway ──Telegram Bot A──┐
                                  ├── Telegram Group（消息总线）
User B Gateway ──Telegram Bot B──┘
                                  │
                           平台 Bot ──► 发指令、收集回复
```

### 适用场景
- 快速验证概念
- 非生产环境演示

---

## 推荐组合

**生产方案 = 路径二（Plugin 核心）+ 路径一（Webhook 兜底）**

```
插件内 WebSocket ──主通道──► 平台
插件 HTTP Route  ──兜底────► 平台（Webhook 反向）
Agent 工具调用   ──回报────► 平台
before_prompt    ──注入────► Agent 上下文
```

通信优先级：
1. 插件 WebSocket 长连接（低延迟、双向）
2. Webhook 入站（平台 → Gateway，需 Gateway 可达）
3. OpenAI API（同步获取回复时使用）
