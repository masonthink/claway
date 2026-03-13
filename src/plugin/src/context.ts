// Context injection for active contributions.
// When the user has draft contributions, injects context into the agent's prompt.

import { ClawayClient } from "./client";

export function registerContextInjection(api: any, client: ClawayClient) {
  api.on(
    "before_prompt_build",
    async () => {
      try {
        const me = await client.getMe();
        if (!me || !me.id) {
          return {};
        }

        // Fetch the user's contributions to find active drafts
        const contribResp = await client.getMyContributions();
        const contributions = contribResp.contributions || contribResp.data || [];

        const drafts = contributions.filter(
          (c: any) => c.status === "draft"
        );

        if (drafts.length === 0) {
          return {};
        }

        // Build context string
        const lines: string[] = [
          "## Claway - 你的未完成草稿",
          "",
        ];

        for (const draft of drafts) {
          lines.push(`### 草稿: ${draft.idea_title || draft.idea_id}`);
          lines.push(`- 贡献 ID: ${draft.id}`);
          lines.push(`- Idea ID: ${draft.idea_id}`);
          lines.push(`- 状态: 草稿 (未提交)`);
          if (draft.updated_at) {
            lines.push(`- 最后更新: ${draft.updated_at}`);
          }
          lines.push("");
        }

        lines.push(
          "提示: 使用 claway_update_contribution 修改草稿内容。"
        );
        lines.push(
          "确认无误后使用 claway_submit_contribution 提交（提交后不可修改）。"
        );

        return {
          prependSystemContext: lines.join("\n"),
        };
      } catch {
        // Silently fail - context injection is best-effort
        return {};
      }
    },
    { priority: 10 }
  );
}
