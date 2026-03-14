// Tool registrations for the Claway OpenClaw plugin.
// v3: contribution-based bidding model with blind voting.

import { ClawayClient } from "./client";
import { runAuthFlow, loadToken } from "./auth";

// Helper to create a standard tool response
function textResult(text: string) {
  return { content: [{ type: "text", text }] };
}

// Web URL helpers — always generate links so users can view content in browser
const WEB_BASE = "https://claway.cc";
function ideaUrl(id: string | number) { return `${WEB_BASE}/idea/${id}`; }
function draftUrl(id: string | number) { return `${WEB_BASE}/draft/${id}`; }
function resultUrl(id: string | number) { return `${WEB_BASE}/idea/${id}/result`; }
function userUrl(username: string) { return `${WEB_BASE}/user/${username}`; }

// Helper to safely execute a tool action with error handling
async function safeExecute(fn: () => Promise<string>): Promise<any> {
  try {
    return textResult(await fn());
  } catch (err: any) {
    const msg = err.message || String(err);
    let hint = "";
    if (msg.includes("401") || msg.includes("Unauthorized")) {
      hint = "\n提示: 认证已过期或未登录，请先运行 claway_auth login。";
    } else if (msg.includes("403")) {
      hint = "\n提示: 权限不足，可能在操作他人的资源。";
    } else if (msg.includes("404")) {
      hint = "\n提示: 资源不存在，请确认 ID 是否正确。";
    } else if (msg.includes("409") || msg.includes("conflict") || msg.includes("already")) {
      hint = "\n提示: 操作冲突，可能是重复投票或重复提交。";
    } else if (msg.includes("429")) {
      hint = "\n提示: 操作频率过高，请稍后再试。";
    }
    return textResult(`操作失败: ${msg}${hint}`);
  }
}

// Truncate text to a max length with ellipsis
function truncate(text: string, max: number): string {
  if (!text || text.length <= max) return text || "";
  return text.slice(0, max) + "...";
}

export function registerTools(api: any, client: ClawayClient) {
  // ========== Auth ==========

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
          description:
            "Action: 'login' to authenticate, 'status' to check current auth, 'logout' to clear saved token",
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
                `  显示名: ${me.display_name || me.username}`,
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
          const authFile = path.join(
            os.homedir(),
            ".config",
            "claway",
            "auth.json"
          );
          try {
            fs.unlinkSync(authFile);
          } catch {}
          return "已退出登录。如需重新登录，请运行 claway_auth login。";
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

        try {
          await authPromise;
          lines.push(``);
          lines.push(`认证成功! 你已登录 Claway，可以开始浏览和参与想法了。`);
        } catch (err: any) {
          lines.push(``);
          lines.push(`认证失败: ${err.message}`);
        }

        return lines.join("\n");
      }),
  });

  // ========== Idea ==========

  api.registerTool({
    name: "claway_create_idea",
    description:
      "Create a new idea on the Claway platform. An idea is a product concept that contributors compete to write the best proposal for. After creation, a 7-day bidding period begins.",
    parameters: {
      type: "object",
      properties: {
        title: {
          type: "string",
          description: "Idea title (max 50 characters)",
        },
        description: {
          type: "string",
          description: "Detailed description of the product idea",
        },
        target_user: {
          type: "string",
          description: "Target user group (one sentence)",
        },
        core_problem: {
          type: "string",
          description: "The core problem this idea solves (one sentence)",
        },
        out_of_scope: {
          type: "string",
          description: "What is explicitly out of scope (optional)",
        },
      },
      required: ["title", "description", "target_user", "core_problem"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const idea = await client.createIdea(params);
        return [
          `Idea 创建成功!`,
          `  ID: ${idea.id}`,
          `  标题: ${idea.title || params.title}`,
          `  截止时间: ${idea.deadline}（7 天后揭榜）`,
          `  查看页面: ${ideaUrl(idea.id)}`,
          ``,
          `你可以把链接分享给其他人，邀请他们提交方案和投票。`,
        ].join("\n");
      }),
  });

  api.registerTool({
    name: "claway_list_ideas",
    description:
      "Browse ideas on the Claway platform. Shows ideas with contribution count and voting stats. Filter by status: 'open' (bidding in progress), 'closed' (revealed), or leave empty for all.",
    parameters: {
      type: "object",
      properties: {
        status: {
          type: "string",
          enum: ["open", "closed"],
          description: "Filter by status: 'open' (bidding), 'closed' (revealed)",
        },
        page: { type: "number", description: "Page number (default: 1)" },
        limit: { type: "number", description: "Results per page (default: 20)" },
      },
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const resp = await client.listIdeas(
          params.status,
          params.page,
          params.limit
        );
        const ideas = resp.ideas || resp.data || [];
        if (ideas.length === 0) {
          if (params.status) {
            return `没有${params.status === "open" ? "进行中" : "已揭榜"}的 Idea。`;
          }
          return `平台暂无 Idea。使用 claway_create_idea 发起第一个!`;
        }
        const total = resp.total || ideas.length;
        const lines = [`共 ${total} 个 Idea:\n`];
        for (const idea of ideas) {
          lines.push(`  [${idea.id}] ${idea.title}`);
          lines.push(
            `      ${truncate(idea.description, 100)}`
          );
          lines.push(
            `      状态: ${idea.status === "open" ? "进行中" : "已揭榜"} | 贡献数: ${idea.contribution_count ?? "?"} | 截止: ${idea.deadline || "?"}`
          );
          lines.push(`      链接: ${ideaUrl(idea.id)}`);
          lines.push("");
        }
        return lines.join("\n");
      }),
  });

  api.registerTool({
    name: "claway_get_idea",
    description:
      "View detailed information about an idea, including description, target user, core problem, and contribution/voting stats.",
    parameters: {
      type: "object",
      properties: {
        idea_id: { type: "string", description: "Idea ID" },
      },
      required: ["idea_id"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const idea = await client.getIdea(params.idea_id);

        const lines = [
          `Idea: ${idea.title}`,
          `  ID: ${idea.id}`,
          `  状态: ${idea.status === "open" ? "进行中" : "已揭榜"}`,
          `  描述: ${idea.description}`,
          `  目标用户: ${idea.target_user}`,
          `  核心问题: ${idea.core_problem}`,
          idea.out_of_scope ? `  不做: ${idea.out_of_scope}` : "",
          `  发起人: ${idea.initiator?.username || idea.initiator_id || "?"}`,
          `  贡献数: ${idea.contribution_count ?? "?"}`,
          `  投票人数: ${idea.voter_count ?? "?"}`,
          `  截止时间: ${idea.deadline || "?"}`,
          `  创建时间: ${idea.created_at}`,
        ];

        lines.push("");
        lines.push(`  查看页面: ${ideaUrl(idea.id)}`);
        if (idea.status === "open") {
          lines.push(`\n下一步: 提交方案 (claway_create_contribution) 或查看已有方案 (claway_list_contributions)`);
        } else if (idea.status === "closed") {
          lines.push(`  揭榜结果: ${resultUrl(idea.id)}`);
          lines.push(`\n使用 claway_get_result 查看排名详情。`);
        }

        return lines.filter(Boolean).join("\n");
      }),
  });

  // ========== Contribution ==========

  api.registerTool({
    name: "claway_create_contribution",
    description:
      "Create a draft contribution for an idea. The content is a full Markdown proposal document. The decision_log records key choices made during the agent-guided process. The draft can be edited before final submission.",
    parameters: {
      type: "object",
      properties: {
        idea_id: { type: "string", description: "Idea ID to contribute to" },
        content: {
          type: "string",
          description: "Full Markdown document content",
        },
        decision_log: {
          type: "array",
          description:
            "Key decisions made during the process, e.g. [{question: '...', choice: '...'}]",
          items: { type: "object" },
        },
      },
      required: ["idea_id", "content", "decision_log"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const result = await client.createContribution(
          params.idea_id,
          params.content,
          params.decision_log || []
        );
        const cid = result.contribution_id || result.id;
        return [
          `草稿已创建!`,
          `  贡献 ID: ${cid}`,
          `  状态: draft (草稿，未提交)`,
          `  网页预览: ${draftUrl(cid)} (仅你可见)`,
          ``,
          `在浏览器中查看渲染效果，回来告诉我要改什么。`,
          `确认无误后使用 claway_submit_contribution 提交（提交后不可修改）。`,
        ].join("\n");
      }),
  });

  api.registerTool({
    name: "claway_update_contribution",
    description:
      "Update a draft contribution's content. Only works while the contribution is still in draft status (not yet submitted).",
    parameters: {
      type: "object",
      properties: {
        contribution_id: {
          type: "string",
          description: "Contribution ID",
        },
        content: {
          type: "string",
          description: "Updated full Markdown document content",
        },
        decision_log: {
          type: "array",
          description: "Updated decision log (optional)",
          items: { type: "object" },
        },
      },
      required: ["contribution_id", "content"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const result = await client.updateContribution(
          params.contribution_id,
          params.content,
          params.decision_log
        );
        return [
          `草稿已更新!`,
          `  贡献 ID: ${result.contribution_id || params.contribution_id}`,
          `  更新时间: ${result.updated_at || "now"}`,
          `  网页预览: ${draftUrl(result.contribution_id || params.contribution_id)}`,
          ``,
          `刷新浏览器查看最新内容。还要改别的吗？`,
          `确认无误后使用 claway_submit_contribution 提交。`,
        ].join("\n");
      }),
  });

  api.registerTool({
    name: "claway_submit_contribution",
    description:
      "Submit a draft contribution, locking it permanently. IMPORTANT: This action is irreversible — confirm with the user before calling this tool.",
    parameters: {
      type: "object",
      properties: {
        contribution_id: {
          type: "string",
          description: "Contribution ID to submit",
        },
      },
      required: ["contribution_id"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const result = await client.submitContribution(params.contribution_id);
        return [
          `方案已提交!`,
          `  贡献 ID: ${result.contribution_id || params.contribution_id}`,
          `  状态: submitted (已锁定，不可修改)`,
          ``,
          `方案将匿名参与评选，揭榜前不公开署名。祝好运!`,
        ].join("\n");
      }),
  });

  api.registerTool({
    name: "claway_get_contribution",
    description:
      "View the full content of a contribution, including the Markdown document and metadata.",
    parameters: {
      type: "object",
      properties: {
        contribution_id: {
          type: "string",
          description: "Contribution ID",
        },
      },
      required: ["contribution_id"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const c = await client.getContribution(params.contribution_id);
        const lines = [
          `贡献 #${c.id}`,
          `  Idea ID: ${c.idea_id}`,
          `  状态: ${c.status === "draft" ? "草稿" : "已提交"}`,
          `  作者: ${c.author?.username || c.author_id || "(匿名)"}`,
          c.status === "draft"
            ? `  网页预览: ${draftUrl(c.id)}`
            : `  Idea 页面: ${ideaUrl(c.idea_id)}`,
          `  提交时间: ${c.submitted_at || "(草稿未提交)"}`,
          `  创建时间: ${c.created_at}`,
          ``,
          `===== 文档内容 =====`,
          ``,
          c.content || "(空)",
        ];
        if (c.status === "draft") {
          lines.push("");
          lines.push(`可使用 claway_update_contribution 修改，或 claway_submit_contribution 提交。`);
        }
        return lines.join("\n");
      }),
  });

  api.registerTool({
    name: "claway_list_contributions",
    description:
      "List contributions for a specific idea. Before reveal, contributions are anonymous and randomly ordered. Only submitted contributions are shown publicly.",
    parameters: {
      type: "object",
      properties: {
        idea_id: { type: "string", description: "Idea ID" },
      },
      required: ["idea_id"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const resp = await client.listContributions(params.idea_id);
        const contributions = resp.contributions || resp.data || [];
        if (contributions.length === 0) {
          return `该 Idea 暂无已提交的方案。使用 claway_create_contribution 提交你的方案!`;
        }
        const lines = [`共 ${contributions.length} 份方案:\n`];
        for (let i = 0; i < contributions.length; i++) {
          const c = contributions[i];
          lines.push(`  方案 #${i + 1} [${c.id}]`);
          lines.push(`    提交时间: ${c.submitted_at || c.created_at}`);
          if (c.preview) {
            lines.push(`    摘要: ${c.preview}`);
          } else if (c.content) {
            lines.push(`    摘要: ${truncate(c.content, 200)}`);
          }
          if (c.author) {
            lines.push(`    作者: ${c.author.username || c.author}`);
          }
          lines.push(`    查看: claway_get_contribution ${c.id}`);
          lines.push("");
        }
        return lines.join("\n");
      }),
  });

  api.registerTool({
    name: "claway_my_contributions",
    description:
      "View your own contributions across all ideas, including drafts and submitted proposals.",
    parameters: { type: "object", properties: {} },
    execute: async (_execId: string, _params: any) =>
      safeExecute(async () => {
        const resp = await client.getMyContributions();
        const contributions = resp.contributions || resp.data || [];
        if (contributions.length === 0) {
          return `暂无贡献记录。使用 claway_list_ideas 浏览想法，找到感兴趣的参与!`;
        }
        const lines = [`我的贡献 (${contributions.length} 份):\n`];
        for (const c of contributions) {
          lines.push(`  [${c.id}] Idea: ${c.idea_title || c.idea_id}`);
          lines.push(
            `    状态: ${c.status === "draft" ? `草稿 | 预览: ${draftUrl(c.id)}` : `已提交于 ${c.submitted_at}`}`
          );
          lines.push("");
        }
        return lines.join("\n");
      }),
  });

  // ========== Vote ==========

  api.registerTool({
    name: "claway_vote",
    description:
      "Cast a vote for a contribution. IMPORTANT: Votes are irreversible — confirm the user's choice before calling this tool.",
    parameters: {
      type: "object",
      properties: {
        idea_id: { type: "string", description: "Idea ID" },
        contribution_id: {
          type: "string",
          description: "Contribution ID to vote for",
        },
      },
      required: ["idea_id", "contribution_id"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const result = await client.vote(
          params.idea_id,
          params.contribution_id
        );
        return [
          `投票成功!`,
          `  Idea: ${params.idea_id}`,
          `  投给方案: ${params.contribution_id}`,
          `  投票时间: ${result.voted_at || "now"}`,
          ``,
          `投票不可撤回。结果将在截止后揭晓。`,
          `  Idea 页面: ${ideaUrl(params.idea_id)}`,
        ].join("\n");
      }),
  });

  api.registerTool({
    name: "claway_my_votes",
    description:
      "View your voting history across all ideas.",
    parameters: { type: "object", properties: {} },
    execute: async (_execId: string, _params: any) =>
      safeExecute(async () => {
        const resp = await client.getMyVotes();
        const votes = resp.votes || resp.data || [];
        if (votes.length === 0) {
          return `暂无投票记录。使用 claway_list_ideas 查看可投票的 Idea。`;
        }
        const lines = [`我的投票 (${votes.length} 条):\n`];
        for (const v of votes) {
          lines.push(`  Idea: ${v.idea_title || v.idea_id}`);
          lines.push(
            `    投给: ${v.contribution_id} | 时间: ${v.voted_at}`
          );
          lines.push(`    链接: ${ideaUrl(v.idea_id)}`);
          lines.push("");
        }
        return lines.join("\n");
      }),
  });

  // ========== Result ==========

  api.registerTool({
    name: "claway_get_result",
    description:
      "View the reveal results for an idea. Only available after the idea's bidding period has ended (status: closed). Shows ranked contributions with vote counts and featured status.",
    parameters: {
      type: "object",
      properties: {
        idea_id: { type: "string", description: "Idea ID" },
      },
      required: ["idea_id"],
    },
    execute: async (_execId: string, params: any) =>
      safeExecute(async () => {
        const resp = await client.getResult(params.idea_id);
        const results = resp.ranked_results || resp.results || resp.data || [];
        const totalVotes = resp.total_votes ?? "?";

        if (results.length === 0) {
          return "暂无揭榜结果（可能 Idea 尚未截止或无贡献）。";
        }

        const lines = [
          `揭榜结果 (总票数: ${totalVotes})`,
          `  揭榜时间: ${resp.revealed_at || "?"}`,
          `  完整结果: ${resultUrl(params.idea_id)}`,
          ``,
        ];

        for (const r of results) {
          const featured = r.is_featured ? " [精选]" : "";
          const author = r.author?.username || r.author || "?";
          lines.push(
            `  #${r.rank} ${r.vote_count} 票${featured}  作者: ${author}`
          );
          lines.push(`    贡献 ID: ${r.contribution_id}`);
          lines.push("");
        }

        return lines.join("\n");
      }),
  });

  // ========== Personal ==========

  api.registerTool({
    name: "claway_my_ideas",
    description:
      "List ideas that you initiated on the Claway platform.",
    parameters: { type: "object", properties: {} },
    execute: async (_execId: string, _params: any) =>
      safeExecute(async () => {
        const resp = await client.getMyIdeas();
        const ideas = resp.ideas || resp.data || [];
        if (ideas.length === 0) {
          return `你还没有发起过 Idea。使用 claway_create_idea 创建你的第一个想法!`;
        }
        const lines = [`我发起的 Idea (${ideas.length} 个):\n`];
        for (const idea of ideas) {
          lines.push(`  [${idea.id}] ${idea.title}`);
          lines.push(
            `      状态: ${idea.status === "open" ? "进行中" : "已揭榜"} | 贡献数: ${idea.contribution_count ?? "?"} | 截止: ${idea.deadline || "?"}`
          );
          lines.push(`      链接: ${ideaUrl(idea.id)}`);
          lines.push("");
        }
        return lines.join("\n");
      }),
  });

  api.registerTool({
    name: "claway_whoami",
    description:
      "Show current authenticated user information on the Claway platform.",
    parameters: { type: "object", properties: {} },
    execute: async (_execId: string, _params: any) =>
      safeExecute(async () => {
        const me = await client.getMe();
        return [
          `当前用户:`,
          `  用户名: ${me.username}`,
          `  显示名: ${me.display_name || me.username}`,
          me.avatar_url ? `  头像: ${me.avatar_url}` : "",
          `  主页: ${userUrl(me.username)}`,
          `  注册时间: ${me.created_at || "?"}`,
        ]
          .filter(Boolean)
          .join("\n");
      }),
  });
}
