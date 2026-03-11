import Link from "next/link";
import StatusBadge from "./StatusBadge";
import type { Idea } from "@/lib/api";

export default function IdeaCard({ idea }: { idea: Idea }) {
  const completed = idea.tasks_completed ?? 0;
  const total = idea.tasks_total ?? 1;
  const progress = total > 0 ? (completed / total) * 100 : 0;

  return (
    <Link
      href={`/ideas/${idea.id}`}
      className="card-hover flex flex-col rounded-[16px] p-5"
      style={{
        background: "var(--surface)",
        border: "1px solid var(--line)",
        boxShadow: "var(--shadow-sm)",
      }}
    >
      <div className="mb-2 flex items-start justify-between gap-3">
        <h3 className="font-display text-[1.05rem] leading-snug tracking-[-0.01em]">
          {idea.title}
        </h3>
        <StatusBadge status={idea.status} />
      </div>

      <p className="mb-4 line-clamp-2 flex-1 text-[0.88rem] leading-[1.5] text-ink-soft">
        {idea.description}
      </p>

      {/* Progress bar */}
      <div className="mb-3">
        <div className="progress-bar">
          <div className="progress-bar-fill" style={{ width: `${progress}%` }} />
        </div>
        <p className="mt-1.5 text-xs text-ink-soft">
          {completed}/{total} 任务完成
        </p>
      </div>

      <div className="flex items-center justify-between text-xs text-ink-soft">
        <span>by @{idea.initiator}</span>
        <span
          className="rounded-[8px] px-2 py-0.5 text-[0.75rem] font-medium"
          style={{ background: "rgba(43,198,164,0.12)", color: "rgb(26,107,91)" }}
        >
          {idea.package_type}
        </span>
      </div>
    </Link>
  );
}
