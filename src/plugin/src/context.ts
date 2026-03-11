// Context injection for active tasks.
// When the user has claimed tasks, injects task context into the agent's prompt.

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

        // Fetch the user's active ideas to find claimed tasks
        const ideasResp = await client.listIdeas("active", 100);
        const ideas = ideasResp.ideas || [];

        if (ideas.length === 0) {
          return {};
        }

        const claimedTasks: any[] = [];
        const depOutputs: Map<number, any[]> = new Map();

        // Check each idea for tasks claimed by this user
        for (const idea of ideas) {
          const tasksResp = await client.getIdeaTasks(idea.id);
          const tasks = tasksResp.tasks || [];

          for (const task of tasks) {
            if (task.claimed_by === me.id && (task.status === "claimed" || task.status === "rejected")) {
              claimedTasks.push({ ...task, idea_title: idea.title });

              // Fetch dependency outputs for this task's idea
              if (!depOutputs.has(idea.id)) {
                try {
                  const ctx = await client.getIdeaContext(idea.id);
                  const approved = (ctx.entries || []).filter(
                    (e: any) => e.status === "approved" && e.content
                  );
                  depOutputs.set(idea.id, approved);
                } catch {
                  depOutputs.set(idea.id, []);
                }
              }
            }
          }
        }

        if (claimedTasks.length === 0) {
          return {};
        }

        // Build context string
        const lines: string[] = [
          "## Claway - 你当前认领的任务",
          "",
        ];

        for (const task of claimedTasks) {
          lines.push(`### 任务 #${task.id}: ${task.title} (${task.type})`);
          lines.push(`- Idea: ${task.idea_title}`);
          lines.push(`- 状态: ${task.status}`);
          lines.push(`- 描述: ${task.description}`);
          lines.push(`- 验收标准: ${task.acceptance_criteria}`);

          if (task.status === "rejected" && task.reject_reason) {
            lines.push(`- **驳回原因: ${task.reject_reason}**`);
          }

          // Show available dependency outputs
          const deps = depOutputs.get(task.idea_id) || [];
          if (deps.length > 0) {
            lines.push(`- 可用依赖产出: ${deps.map((d: any) => d.task_type).join(", ")}`);
          }

          lines.push("");
        }

        lines.push(
          "提示: 使用 claway_get_task_context 获取完整的任务上下文和依赖产出内容。"
        );
        lines.push(
          "完成后使用 claway_submit_task 提交成果。"
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
