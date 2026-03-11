import { Cpu, User } from "lucide-react";
import StatusBadge from "./StatusBadge";
import type { Task } from "@/lib/api";

export default function TaskList({ tasks }: { tasks: Task[] }) {
  return (
    <div className="space-y-3">
      {tasks.map((task) => (
        <div
          key={task.id}
          className="flex items-center justify-between rounded-lg border border-gray-200 bg-white p-4"
        >
          <div className="flex items-center gap-3">
            <span className="flex h-8 w-8 items-center justify-center rounded-full bg-indigo-100 text-xs font-bold text-indigo-700">
              {task.code}
            </span>
            <div>
              <p className="text-sm font-medium text-gray-900">{task.name}</p>
              {task.claimed_by && (
                <p className="flex items-center gap-1 text-xs text-gray-400">
                  <User className="h-3 w-3" />
                  {task.claimed_by}
                </p>
              )}
            </div>
          </div>

          <div className="flex items-center gap-3">
            {task.token_cost > 0 && (
              <span className="flex items-center gap-1 text-xs text-gray-400">
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
