import Link from "next/link";
import { Cpu } from "lucide-react";
import StatusBadge from "./StatusBadge";
import type { Idea } from "@/lib/api";

export default function IdeaCard({ idea }: { idea: Idea }) {
  return (
    <Link
      href={`/ideas/${idea.id}`}
      className="block rounded-xl border border-gray-200 bg-white p-5 shadow-sm transition-shadow hover:shadow-md"
    >
      <div className="mb-3 flex items-start justify-between">
        <h3 className="text-lg font-semibold text-gray-900">{idea.title}</h3>
        <StatusBadge status={idea.status} />
      </div>

      <p className="mb-4 line-clamp-2 text-sm text-gray-500">
        {idea.description}
      </p>

      <div className="flex items-center justify-between text-xs text-gray-400">
        <div className="flex items-center gap-3">
          <span className="rounded bg-indigo-50 px-2 py-0.5 font-medium text-indigo-600">
            {idea.package_type}
          </span>
          <span>
            {idea.tasks_completed}/{idea.tasks_total} 已完成
          </span>
        </div>
        <div className="flex items-center gap-1">
          <Cpu className="h-3.5 w-3.5" />
          <span>{idea.total_compute_cost.toFixed(2)} tokens</span>
        </div>
      </div>
    </Link>
  );
}
