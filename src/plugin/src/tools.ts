// Tool registrations for the Claway OpenClaw plugin.

import { ClawayClient } from "./client";
import { runAuthFlow, loadToken } from "./auth";

// Helper to create a standard tool response
function textResult(text: string) {
  return { content: [{ type: "text", text }] };
}

// Helper to safely execute a tool action with error handling
async function safeExecute(fn: () => Promise<string>): Promise<any> {
  try {
    return textResult(await fn());
  } catch (err: any) {
    return textResult(`操作失败: ${err.message || String(err)}`);
  }
}

export function registerTools(api: any, client: ClawayClient) {
  // ========== Auth Tools ==========

  api.registerTool({
    name: "claway_auth",
    description:
      "Authenticate with Claway using your X (Twitter) account. Opens a browser window for OAuth login. Run this first before using any other Claway tools.",
    parameters: {
      type: "object",
      properties: {
        action: {
          type: "string",
          enum: ["login", "status", "logout"],
          description: "Action: 'login' to authenticate, 'status' to check current auth, 'logout' to clear saved token",
        },
      },
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const action = params.action || "login";
        const platformUrl = (client as any).baseUrl || "";

        if (action === "status") {
          const token = loadToken();
          if (token) {
            try {
              const me = await client.getMe();
              return [
                `已登录 Claway`,
                `  用户名: ${me.username}`,
                `  积分余额: ${me.credits_balance}`,
              ].join("\n");
            } catch {
              return "Token 已保存但可能已过期，请重新运行 claway_auth login";
            }
          }
          return "未登录。请运行 claway_auth login 进行认证。";
        }

        if (action === "logout") {
          const fs = await import("fs");
          const path = await import("path");
          const os = await import("os");
          const authFile = path.join(os.homedir(), ".config", "claway", "auth.json");
          try {
            fs.unlinkSync(authFile);
          } catch {}
          return "已退出登录，Token 已清除。";
        }

        // Login flow
        const authPromise = runAuthFlow(platformUrl);
        const pending = (runAuthFlow as any)._pending;

        if (!pending) {
          return "无法启动认证服务器，请重试。";
        }

        const lines = [
          `请在浏览器中打开以下链接完成 X 账号授权:`,
          ``,
          `  ${pending.authUrl}`,
          ``,
          `等待授权完成... (2 分钟超时)`,
        ];

        // Wait for the auth to complete
        try {
          const result = await authPromise;
          lines.push(``);
          lines.push(`认证成功! Token 已保存到 ~/.config/claway/auth.json`);
          lines.push(`现在可以使用所有 Claway 工具了。`);
        } catch (err: any) {
          lines.push(``);
          lines.push(`认证失败: ${err.message}`);
        }

        return lines.join("\n");
      }),
  });

  // ========== Initiator Tools ==========

  api.registerTool({
    name: "claway_create_idea",
    description:
      "Create a new idea on the Claway platform. An idea is a product concept that gets broken down into research tasks for contributors to work on.",
    parameters: {
      type: "object",
      properties: {
        title: { type: "string", description: "Idea title" },
        description: { type: "string", description: "Detailed description of the idea" },
        target_user_hint: { type: "string", description: "Hint about the target user group" },
        problem_definition: { type: "string", description: "The problem this idea solves" },
        initiator_cut_percent: {
          type: "number",
          description: "Initiator's revenue share percentage (10-30)",
        },
        package_type: {
          type: "string",
          enum: ["light", "standard"],
          description: "Package type: 'light' (5 tasks) or 'standard' (9 tasks)",
        },
      },
      required: ["title", "description", "initiator_cut_percent", "package_type"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const idea = await client.createIdea(params);
        return [
          `创建成功!`,
          `  Idea ID: ${idea.id}`,
          `  标题: ${idea.title}`,
          `  套餐: ${idea.package_type}`,
          `  状态: ${idea.status}`,
          ``,
          `系统已自动为该 Idea 生成任务，使用 claway_view_idea 查看任务列表。`,
        ].join("\n");
      }),
  });

  api.registerTool({
    name: "claway_my_ideas",
    description:
      "List ideas that I initiated on the Claway platform. Returns ideas with their current status.",
    parameters: {
      type: "object",
      properties: {
        status: {
          type: "string",
          description: "Filter by status: active, completed, etc. Leave empty for all.",
        },
      },
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const resp = await client.listIdeas(params.status);
        const ideas = resp.ideas || [];
        if (ideas.length === 0) {
          return "暂无 Idea。";
        }
        const lines = [`共 ${resp.total} 个 Idea:\n`];
        for (const idea of ideas) {
          lines.push(`  [${idea.id}] ${idea.title}`);
          lines.push(`      状态: ${idea.status} | 套餐: ${idea.package_type}`);
          lines.push("");
        }
        return lines.join("\n");
      }),
  });

  api.registerTool({
    name: "claway_review_task",
    description:
      "Review a submitted task (approve or reject). Only the idea initiator can review tasks. Approving awards credits to the contributor.",
    parameters: {
      type: "object",
      properties: {
        task_id: { type: "number", description: "Task ID to review" },
        action: {
          type: "string",
          enum: ["approve", "reject"],
          description: "Review action",
        },
        quality_score: {
          type: "number",
          enum: [1.0, 1.2, 1.5],
          description: "Quality multiplier (required for approve): 1.0=standard, 1.2=good, 1.5=excellent",
        },
        reject_reason: {
          type: "string",
          description: "Reason for rejection (required for reject)",
        },
      },
      required: ["task_id", "action"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        await client.reviewTask(
          params.task_id,
          params.action,
          params.quality_score,
          params.reject_reason
        );
        if (params.action === "approve") {
          return `任务 #${params.task_id} 已批准 (质量系数: ${params.quality_score})，贡献者已获得积分奖励。`;
        }
        return `任务 #${params.task_id} 已驳回。原因: ${params.reject_reason}`;
      }),
  });

  api.registerTool({
    name: "claway_publish_prd",
    description:
      "Merge all approved task outputs into a final PRD document and publish it. Only the idea initiator can publish.",
    parameters: {
      type: "object",
      properties: {
        idea_id: { type: "number", description: "Idea ID to publish PRD for" },
      },
      required: ["idea_id"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const prd = await client.publishPRD(params.idea_id);
        return [
          `PRD 发布成功!`,
          `  PRD ID: ${prd.id}`,
          `  Idea: ${prd.idea_id}`,
          `  版本: ${prd.version}`,
        ].join("\n");
      }),
  });

  // ========== Contributor Tools ==========

  api.registerTool({
    name: "claway_browse_ideas",
    description:
      "Browse available ideas on the Claway platform that you can contribute to. Shows ideas with open tasks.",
    parameters: {
      type: "object",
      properties: {
        status: {
          type: "string",
          description: "Filter by status (default: active)",
        },
        limit: { type: "number", description: "Max results (default: 20)" },
        offset: { type: "number", description: "Pagination offset" },
      },
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const resp = await client.listIdeas(
          params.status || "active",
          params.limit,
          params.offset
        );
        const ideas = resp.ideas || [];
        if (ideas.length === 0) {
          return "当前没有可参与的 Idea。";
        }
        const lines = [`共 ${resp.total} 个可参与的 Idea:\n`];
        for (const idea of ideas) {
          lines.push(`  [${idea.id}] ${idea.title}`);
          lines.push(`      ${idea.description.slice(0, 100)}${idea.description.length > 100 ? "..." : ""}`);
          lines.push(`      套餐: ${idea.package_type} | 发起人分成: ${idea.initiator_cut_percent}%`);
          lines.push("");
        }
        return lines.join("\n");
      }),
  });

  api.registerTool({
    name: "claway_view_idea",
    description:
      "View detailed information about an idea, including its task list with status and descriptions.",
    parameters: {
      type: "object",
      properties: {
        idea_id: { type: "number", description: "Idea ID to view" },
      },
      required: ["idea_id"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const [idea, tasksResp] = await Promise.all([
          client.getIdea(params.idea_id),
          client.getIdeaTasks(params.idea_id),
        ]);
        const tasks = tasksResp.tasks || [];

        const lines = [
          `Idea #${idea.id}: ${idea.title}`,
          `  状态: ${idea.status}`,
          `  套餐: ${idea.package_type}`,
          `  描述: ${idea.description}`,
          idea.target_user_hint ? `  目标用户: ${idea.target_user_hint}` : "",
          idea.problem_definition ? `  问题定义: ${idea.problem_definition}` : "",
          `  发起人分成: ${idea.initiator_cut_percent}%`,
          ``,
          `任务列表 (${tasks.length} 个):`,
        ];

        for (const task of tasks) {
          const claimed = task.claimed_by ? ` | 认领者: ${task.claimed_by}` : "";
          lines.push(`  [${task.id}] ${task.type} - ${task.title}`);
          lines.push(`      状态: ${task.status}${claimed}`);
          lines.push(`      依赖: ${task.dependencies || "无"}`);
          lines.push(`      验收标准: ${task.acceptance_criteria}`);
          lines.push("");
        }

        return lines.filter(Boolean).join("\n");
      }),
  });

  api.registerTool({
    name: "claway_claim_task",
    description:
      "Claim an open task so you can work on it. The task's dependency tasks must all be approved before claiming.",
    parameters: {
      type: "object",
      properties: {
        task_id: { type: "number", description: "Task ID to claim" },
      },
      required: ["task_id"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        await client.claimTask(params.task_id);
        const task = await client.getTask(params.task_id);
        return [
          `任务已认领!`,
          `  任务 #${task.id}: ${task.title}`,
          `  类型: ${task.type}`,
          `  描述: ${task.description}`,
          `  验收标准: ${task.acceptance_criteria}`,
          ``,
          `使用 claway_get_task_context 获取依赖产出，然后开始工作。`,
          `完成后使用 claway_submit_task 提交成果。`,
        ].join("\n");
      }),
  });

  api.registerTool({
    name: "claway_unclaim_task",
    description:
      "Release a claimed task back to open status. Only the current claimer can unclaim.",
    parameters: {
      type: "object",
      properties: {
        task_id: { type: "number", description: "Task ID to unclaim" },
      },
      required: ["task_id"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        await client.unclaimTask(params.task_id);
        return `任务 #${params.task_id} 已释放，其他人可以认领。`;
      }),
  });

  api.registerTool({
    name: "claway_get_task_context",
    description:
      "Get full context for a task, including task details, acceptance criteria, and all approved dependency outputs. Use this before starting work on a claimed task.",
    parameters: {
      type: "object",
      properties: {
        task_id: { type: "number", description: "Task ID to get context for" },
      },
      required: ["task_id"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const task = await client.getTask(params.task_id);
        const context = await client.getIdeaContext(task.idea_id);

        const lines = [
          `===== 任务上下文 =====`,
          ``,
          `任务 #${task.id}: ${task.title} (${task.type})`,
          `  描述: ${task.description}`,
          `  验收标准: ${task.acceptance_criteria}`,
          `  Token 限制提示: ${task.token_limit_hint || "无"}`,
          `  依赖: ${task.dependencies || "无"}`,
          ``,
        ];

        // Include dependency outputs
        const entries = context.entries || [];
        const depEntries = entries.filter(
          (e: any) => e.status === "approved" && e.content
        );

        if (depEntries.length > 0) {
          lines.push(`===== 已完成的依赖产出 =====`);
          lines.push(``);
          for (const entry of depEntries) {
            lines.push(`--- ${entry.task_type}: ${entry.title} ---`);
            lines.push(entry.content);
            lines.push(``);
          }
        } else {
          lines.push(`(暂无已完成的依赖产出)`);
        }

        return lines.join("\n");
      }),
  });

  api.registerTool({
    name: "claway_submit_task",
    description:
      "Submit your completed work for a claimed task. The content should be the full task output. The initiator will review your submission.",
    parameters: {
      type: "object",
      properties: {
        task_id: { type: "number", description: "Task ID to submit" },
        content: { type: "string", description: "Full task output content" },
        note: {
          type: "string",
          description: "Short note about this submission (max 200 chars)",
        },
      },
      required: ["task_id", "content"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        await client.submitTask(params.task_id, params.content, params.note || "");
        return [
          `任务 #${params.task_id} 已提交!`,
          `等待发起人审核...`,
          ``,
          `提交备注: ${params.note || "(无)"}`,
        ].join("\n");
      }),
  });

  api.registerTool({
    name: "claway_update_document",
    description:
      "Update the document content for a task you're working on. Use this for saving in-progress work before final submission.",
    parameters: {
      type: "object",
      properties: {
        task_id: { type: "number", description: "Task ID" },
        content: { type: "string", description: "Updated document content" },
      },
      required: ["task_id", "content"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        await client.updateDocument(params.task_id, params.content);
        return `任务 #${params.task_id} 的文档已更新。`;
      }),
  });

  // ========== Shared Tools ==========

  api.registerTool({
    name: "claway_my_compute",
    description:
      "View your compute (LLM token) usage on the Claway platform, including total cost and breakdown by idea.",
    parameters: { type: "object", properties: {} },
    execute: async (_execId: string, _params: any) =>
      safeExecute(async () => {
        const compute = await client.getMyCompute();
        return [
          `我的算力使用:`,
          `  总消耗: $${compute.total_cost_usd?.toFixed(4) || "0.0000"}`,
          `  总 Token 数: ${compute.total_tokens || 0}`,
          `  请求次数: ${compute.total_requests || 0}`,
        ].join("\n");
      }),
  });

  api.registerTool({
    name: "claway_my_credits",
    description:
      "View your credits balance and transaction history on the Claway platform.",
    parameters: { type: "object", properties: {} },
    execute: async (_execId: string, _params: any) =>
      safeExecute(async () => {
        const resp = await client.getMyCredits();
        const lines = [`我的积分:`];
        lines.push(`  余额: ${resp.balance || 0} 积分`);

        const txs = resp.transactions || [];
        if (txs.length > 0) {
          lines.push(``);
          lines.push(`最近交易:`);
          for (const tx of txs.slice(0, 10)) {
            const sign = tx.amount >= 0 ? "+" : "";
            lines.push(`  ${sign}${tx.amount} | ${tx.type} | ${tx.description}`);
          }
        }

        return lines.join("\n");
      }),
  });

  api.registerTool({
    name: "claway_my_contributions",
    description:
      "View your contribution history on the Claway platform, showing tasks you completed and credits earned.",
    parameters: { type: "object", properties: {} },
    execute: async (_execId: string, _params: any) =>
      safeExecute(async () => {
        const resp = await client.getMyContributions();
        const contribs = resp.contributions || [];
        if (contribs.length === 0) {
          return "暂无贡献记录。";
        }
        const lines = [`我的贡献 (${contribs.length} 条):\n`];
        for (const c of contribs) {
          lines.push(`  Idea #${c.idea_id} / Task #${c.task_id}`);
          lines.push(`    消耗: $${c.cost_usd?.toFixed(4) || "0"} | 质量: ${c.quality_score} | 加权: ${c.weighted_score?.toFixed(4) || "0"}`);
          lines.push("");
        }
        return lines.join("\n");
      }),
  });
}
