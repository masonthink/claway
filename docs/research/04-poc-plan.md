# POC 验证计划

## 概述

在正式开发前，需要验证 5 个关键技术点。总计约 3-4 天。

---

## POC-1：插件工具内网络请求（半天）

### 验证什么
`api.registerTool()` 的 `execute()` 函数内能否 `fetch()` 外部 URL。

### 验证步骤
1. 创建最小插件，注册一个 `test_fetch` 工具
2. 工具 execute 内 POST 到 httpbin.org/post
3. 安装插件到 Gateway
4. 通过聊天触发 Agent 调用该工具
5. 验证请求是否成功发出和接收

### 最小代码
```typescript
export default function(api) {
  api.registerTool({
    name: "test_fetch",
    description: "Test external HTTP request",
    parameters: { type: "object", properties: {} },
    execute: async () => {
      const res = await fetch("https://httpbin.org/post", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ test: true, timestamp: Date.now() })
      });
      const data = await res.json();
      return { content: [{ type: "text", text: JSON.stringify(data) }] };
    }
  });
}
```

### 风险等级：高
如果失败，需退化为 sidecar 模式。

### 通过标准
- fetch 请求成功发出
- 响应数据正确返回到 Agent

---

## POC-2：插件后台服务 WebSocket 长连接（半天）

### 验证什么
`api.registerService()` 内能否建立并维持 WebSocket 长连接。

### 验证步骤
1. 搭建简单 WebSocket echo 服务器
2. 创建插件，在 registerService.start() 中建立 WebSocket 连接
3. 验证连接建立、消息收发、断线重连

### 最小代码
```typescript
import WebSocket from "ws";

export default function(api) {
  let ws: WebSocket | null = null;

  api.registerService({
    id: "ws-test",
    start: () => {
      ws = new WebSocket("wss://echo.websocket.org");
      ws.on("open", () => api.logger.info("WS connected"));
      ws.on("message", (data) => api.logger.info(`WS received: ${data}`));
      ws.on("close", () => {
        api.logger.info("WS closed, reconnecting...");
        setTimeout(() => { /* reconnect logic */ }, 3000);
      });
    },
    stop: () => { ws?.close(); }
  });
}
```

### 风险等级：中
WebSocket 库在 Gateway 进程内是否可用需确认。

### 通过标准
- WebSocket 连接建立成功
- 消息双向传输正常
- 断线后能自动重连

---

## POC-3：before_prompt_build 上下文注入（半天）

### 验证什么
通过 `api.on("before_prompt_build")` 注入的内容是否对 Agent 可见并影响其行为。

### 验证步骤
1. 创建插件，注册 before_prompt_build 钩子
2. 注入特定指令（如"你的名字是 CollabAgent，回答任何问题时先说你的名字"）
3. 通过聊天向 Agent 提问，观察是否遵循注入的指令
4. 测试动态更新注入内容

### 最小代码
```typescript
export default function(api) {
  let collabContext = "你正在参与协作任务 task-001。你的角色是后端开发。";

  api.on("before_prompt_build", () => ({
    prependSystemContext: collabContext,
    appendSystemContext: "当前任务进度：前端已完成登录页面。"
  }), { priority: 10 });

  // 通过 HTTP 端点动态更新上下文
  api.registerHttpRoute({
    path: "/collab/update-context",
    auth: "plugin",
    handler: async (req, res) => {
      let body = "";
      req.on("data", chunk => body += chunk);
      req.on("end", () => {
        collabContext = JSON.parse(body).context;
        res.end(JSON.stringify({ ok: true }));
      });
      return true;
    }
  });
}
```

### 风险等级：低
文档明确支持此功能。

### 通过标准
- Agent 回复中体现了注入的上下文
- 动态更新上下文后，Agent 行为相应改变
- 不影响 Agent 的其他能力

---

## POC-4：Webhook + Session 历史读取（1天）

### 验证什么
通过 `/hooks/agent` 发任务后，能否通过 API 读取 Agent 的回复。

### 验证步骤
1. 配置 Gateway hooks
2. 发送 `/hooks/agent` 请求，设置 `deliver: false`, 指定 `sessionKey`
3. 等待 Agent 执行完毕
4. 通过 `/tools/invoke` 调用 `sessions_list` 查看 session
5. 尝试读取 session 内的消息历史

### 验证命令
```bash
# Step 1: 发送任务
curl -X POST http://127.0.0.1:18789/hooks/agent \
  -H 'Authorization: Bearer SECRET' \
  -H 'Content-Type: application/json' \
  -d '{
    "message": "列出当前目录的文件",
    "sessionKey": "hook:collab:test-001",
    "deliver": false,
    "timeoutSeconds": 30
  }'

# Step 2: 读取 session 列表
curl -X POST http://127.0.0.1:18789/tools/invoke \
  -H 'Authorization: Bearer SECRET' \
  -H 'Content-Type: application/json' \
  -d '{"tool": "sessions_list", "action": "json"}'

# Step 3: 尝试通过 OpenAI API 读取（用相同 session）
curl http://127.0.0.1:18789/v1/chat/completions \
  -H 'Authorization: Bearer SECRET' \
  -H 'Content-Type: application/json' \
  -d '{
    "model": "openclaw:main",
    "messages": [{"role": "user", "content": "总结你刚才做了什么"}],
    "user": "collab:test-001"
  }'
```

### 风险等级：中
回复捕获的具体机制文档不够明确。

### 通过标准
- 能从 session 历史中获取 Agent 的完整回复
- 回复内容可被程序化解析

---

## POC-5：端到端双 Gateway 消息中继（1天）

### 验证什么
两个 Gateway 通过外部服务中继消息的端到端可行性和延迟。

### 验证步骤
1. 搭建最小 WebSocket 中继服务器
2. 启动两个 Gateway 实例（不同端口），各装 collab 插件
3. Gateway A 的 Agent 通过 `collab_submit` 工具发消息
4. 中继服务器路由到 Gateway B
5. Gateway B 的 Agent 收到消息并处理
6. 测量端到端延迟

### 架构
```
Gateway A (port 18789)          中继服务器           Gateway B (port 18790)
  Agent A ──collab_submit──►   WS Server    ──webhook──► Agent B
                               (localhost:3000)
```

### 风险等级：低
各单项技术都是成熟方案，主要验证集成。

### 通过标准
- 消息端到端延迟 < 500ms
- 消息内容完整无损
- 断线重连后消息不丢失

---

## Go/No-Go 判定

| 条件 | 结果 | 决策 |
|------|------|------|
| POC-1 通过 | 插件可发网络请求 | ✅ 走插件路线 |
| POC-1 失败 | 插件网络受限 | ⚠️ 退化为 sidecar 模式 |
| POC-2 通过 | 可维持长连接 | ✅ 用 WebSocket 做主通道 |
| POC-2 失败 | 长连接不稳定 | ⚠️ 退化为轮询 + Webhook |
| POC-3 通过 | 上下文注入有效 | ✅ Agent 能理解协作背景 |
| POC-4 通过 | 能读取 Agent 回复 | ✅ Webhook 可作为兜底通道 |
| POC-5 通过 | 端到端集成可行 | ✅ 开始正式开发 |

**关键路径**: POC-1 → POC-2 → POC-5（必须全部通过）
