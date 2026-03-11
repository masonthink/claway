# OpenClaw 插件能力详解

## 插件本质

插件是 **in-process 运行的代码模块**，与 Gateway 同一个 Node.js 进程，拥有完整的运行时能力。

> "A plugin is a small code module that extends OpenClaw with extra features (commands, tools, and Gateway RPC). Plugins run in-process with the Gateway."

## 插件结构

### Manifest (openclaw.plugin.json)

```json
{
  "id": "claway-collab",
  "configSchema": {
    "type": "object",
    "properties": {
      "platformUrl": { "type": "string" },
      "apiKey": { "type": "string" }
    }
  },
  "uiHints": {
    "apiKey": { "label": "API Key", "sensitive": true }
  }
}
```

### 入口文件

```typescript
// 函数格式
export default function(api) { /* ... */ }

// 对象格式
export default {
  id: "claway-collab",
  name: "Claway Collab",
  configSchema: { /* ... */ },
  register(api) { /* ... */ }
}
```

## 核心 API 能力

### 1. registerTool — 注册 Agent 工具

```typescript
api.registerTool({
  name: "collab_submit",
  description: "提交协作产出到平台",
  parameters: {
    type: "object",
    properties: {
      taskId: { type: "string" },
      content: { type: "string" },
      fileChanges: { type: "array", items: { type: "object" } }
    },
    required: ["taskId", "content"]
  },
  execute: async (_execId, params) => {
    // 可以发网络请求！
    const res = await fetch("https://platform.example.com/api/submit", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(params)
    });
    const data = await res.json();
    return { content: [{ type: "text", text: JSON.stringify(data) }] };
  }
});
```

**关键**: execute() 内可做任何事，包括网络请求。工具分两类：
- **Required tools**: 自动对所有 Agent 可用
- **Optional tools**: 需要 `tools.allow` 白名单（推荐用于有副作用的工具）

### 2. registerHttpRoute — 注册 HTTP 端点

```typescript
api.registerHttpRoute({
  path: "/collab/task",
  auth: "plugin",      // "plugin" 用插件自己的认证，"gateway" 用 Gateway token
  match: "exact",      // "exact" 或 "prefix"
  handler: async (req, res) => {
    // 标准 Node.js HTTP handler
    let body = "";
    req.on("data", chunk => body += chunk);
    req.on("end", () => {
      const data = JSON.parse(body);
      // 处理逻辑...
      res.statusCode = 200;
      res.end(JSON.stringify({ ok: true }));
    });
    return true;
  }
});
```

### 3. registerService — 后台服务

```typescript
api.registerService({
  id: "collab-bridge",
  start: () => {
    // 启动后台任务，如 WebSocket 长连接
    api.logger.info("Collab bridge started");
  },
  stop: () => {
    // 清理资源
    api.logger.info("Collab bridge stopped");
  }
});
```

### 4. on() — 生命周期钩子

```typescript
// 注入系统上下文
api.on("before_prompt_build", (event, ctx) => {
  return {
    prependSystemContext: "协作任务上下文...",
    appendSystemContext: "其他 Agent 的最新产出..."
  };
}, { priority: 10 });

// 模型选择覆盖
api.on("before_model_resolve", (event, ctx) => {
  return { modelOverride: "anthropic/claude-opus-4-6" };
});
```

可用钩子：
| 钩子 | 时机 | 能力 |
|------|------|------|
| `before_model_resolve` | session 加载前 | 覆盖模型/provider |
| `before_prompt_build` | session 加载后 | 注入 system context |
| `before_tool_call` | 工具调用前 | 拦截/修改参数 |
| `after_tool_call` | 工具调用后 | 修改结果 |
| `tool_result_persist` | 结果持久化时 | 修改存储内容 |

### 5. registerGatewayMethod — 自定义 RPC

```typescript
api.registerGatewayMethod("collab.status", ({ respond }) => {
  respond(true, { connected: true, activeTasks: 3 });
});
```

### 6. registerCommand — 自动回复命令

```typescript
api.registerCommand({
  name: "collabstatus",
  description: "显示协作状态",
  acceptsArgs: false,
  requireAuth: true,
  handler: (ctx) => ({
    text: `当前协作任务: 3, 连接状态: 正常`
  })
});
```

### 7. registerCli — CLI 命令

```typescript
api.registerCli(({ program }) => {
  program.command("collab")
    .description("Manage collaboration")
    .action(() => { console.log("Collab status..."); });
}, { commands: ["collab"] });
```

### 8. registerContextEngine — 上下文引擎（高级）

```typescript
api.registerContextEngine("collab-context", () => ({
  info: { id: "collab-context", name: "Collab Context", ownsCompaction: false },
  async ingest() { return { ingested: true }; },
  async assemble({ messages }) {
    // 在消息中插入协作上下文
    return { messages: [...collabContext, ...messages], estimatedTokens: 0 };
  },
  async compact() { return { ok: true, compacted: false }; }
}));
```

## 插件安装与分发

### 安装方式
```bash
# 从 npm
openclaw plugins install @claway/collab

# 从本地路径（开发）
openclaw plugins install -l ./extensions/collab

# 固定版本
openclaw plugins install @claway/collab --pin
```

### 配置启用
```json5
{
  plugins: {
    enabled: true,
    allow: ["claway-collab"],
    entries: {
      "claway-collab": {
        enabled: true,
        config: {
          platformUrl: "https://platform.claway.com",
          apiKey: "user-api-key"
        }
      }
    }
  }
}
```

### 加载顺序
1. `plugins.load.paths` 配置路径
2. `<workspace>/.openclaw/extensions/`
3. `~/.openclaw/extensions/`
4. 内置 extensions

### npm 发布要求
```json
{
  "name": "@claway/collab",
  "openclaw": {
    "extensions": ["./dist/index.js"]
  }
}
```

## SDK 导入

```typescript
import { SomeType } from "openclaw/plugin-sdk/core";
```

## 安全注意事项

- 插件 in-process 运行，等同于受信代码
- 不能路径穿越（entry 必须在插件目录内）
- 非 bundled 插件检查文件所有权
- 可通过 `plugins.allow/deny` 控制加载
- `hooks.allowPromptInjection: false` 可禁止插件注入 prompt
