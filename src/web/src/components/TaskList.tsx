import { Cpu, User } from "lucide-react";
import StatusBadge from "./StatusBadge";
import type { Task } from "@/lib/api";

export default function TaskList({ tasks }: { tasks: Task[] }) {
  return (
    <div
      className="overflow-hidden rounded-[16px]"
      style={{ border: "1px solid var(--line)", background: "var(--surface)" }}
    >
      {tasks.map((task, i) => (
        <div
          key={task.id}
          className="flex items-center justify-between gap-3 px-5 py-3.5"
          style={{
            borderBottom: i === tasks.length - 1 ? "none" : "1px solid var(--line)",
          }}
        >
          <div className="flex items-center gap-3 min-w-0">
            <span
              className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg text-[0.7rem] font-bold"
              style={{ background: "rgba(255,107,74,0.1)", color: "var(--accent-deep)" }}
            >
              {task.code}
            </span>
            <div className="min-w-0">
              <p className="truncate text-sm font-medium">{task.name}</p>
              {task.claimed_by && (
                <p className="flex items-center gap-1 text-xs text-ink-soft">
                  <User className="h-3 w-3" />
                  @{task.claimed_by}
                </p>
              )}
            </div>
          </div>

          <div className="flex shrink-0 items-center gap-2.5">
            {task.token_cost > 0 && (
              <span className="flex items-center gap-1 text-xs text-ink-soft">
                <Cpu className="h-3 w-3" />
                {task.token_cost.toFixed(2)}
              </span>
            )}
            <StatusBadge status={task.status} />
          </div>
        </div>
      ))}
    </div>
  );
}
