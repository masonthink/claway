"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Trophy, FileText } from "lucide-react";
import StatusBadge from "@/components/StatusBadge";
import TaskList from "@/components/TaskList";
import {
  getIdea,
  getIdeaTasks,
  getIdeaCompute,
  type Idea,
  type Task,
  type ComputeLeaderboard,
} from "@/lib/api";

export default function IdeaDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [idea, setIdea] = useState<Idea | null>(null);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [leaderboard, setLeaderboard] = useState<ComputeLeaderboard | null>(
    null
  );
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;

    Promise.all([
      getIdea(id).then(setIdea),
      getIdeaTasks(id).then((d) => setTasks(d.tasks)),
      getIdeaCompute(id).then(setLeaderboard),
    ]).catch((err) => setError(err.message));
  }, [id]);

  if (error) {
    return (
      <div className="mx-auto max-w-4xl px-4 py-12">
        <div className="rounded-lg bg-red-50 p-4 text-sm text-red-600">
          Failed to load idea: {error}
        </div>
      </div>
    );
  }

  if (!idea) {
    return (
      <div className="mx-auto max-w-4xl px-4 py-12 text-center text-gray-400">
        Loading...
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-4xl px-4 py-8">
      {/* Back link */}
      <Link
        href="/"
        className="mb-6 inline-flex items-center gap-1 text-sm text-gray-500 hover:text-indigo-600 transition-colors"
      >
        <ArrowLeft className="h-4 w-4" />
        返回列表
      </Link>

      {/* Idea header */}
      <div className="mb-8 rounded-xl border border-gray-200 bg-white p-6">
        <div className="mb-4 flex items-start justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">{idea.title}</h1>
            <p className="mt-1 text-sm text-gray-400">
              by {idea.initiator} &middot;{" "}
              <span className="rounded bg-indigo-50 px-2 py-0.5 text-xs font-medium text-indigo-600">
                {idea.package_type}
              </span>
            </p>
          </div>
          <StatusBadge status={idea.status} />
        </div>
        <p className="text-gray-600">{idea.description}</p>

        {/* View PRD button */}
        {idea.status === "completed" && (
          <Link
            href={`/prd/${idea.id}`}
            className="mt-4 inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 transition-colors"
          >
            <FileText className="h-4 w-4" />
            查看 PRD
          </Link>
        )}
      </div>

      {/* Tasks */}
      <div className="mb-8">
        <h2 className="mb-4 text-lg font-semibold text-gray-900">
          Tasks ({tasks.filter((t) => t.status === "approved").length}/
          {tasks.length} 已完成)
        </h2>
        <TaskList tasks={tasks} />
      </div>

      {/* Compute leaderboard */}
      {leaderboard && leaderboard.entries.length > 0 && (
        <div>
          <h2 className="mb-4 flex items-center gap-2 text-lg font-semibold text-gray-900">
            <Trophy className="h-5 w-5 text-yellow-500" />
            Compute Leaderboard
          </h2>
          <div className="rounded-xl border border-gray-200 bg-white">
            {leaderboard.entries.map((entry, i) => (
              <div
                key={entry.username}
                className={`flex items-center justify-between px-5 py-3 ${
                  i > 0 ? "border-t border-gray-100" : ""
                }`}
              >
                <div className="flex items-center gap-3">
                  <span
                    className={`flex h-7 w-7 items-center justify-center rounded-full text-xs font-bold ${
                      i === 0
                        ? "bg-yellow-100 text-yellow-700"
                        : i === 1
                          ? "bg-gray-200 text-gray-600"
                          : i === 2
                            ? "bg-orange-100 text-orange-700"
                            : "bg-gray-100 text-gray-500"
                    }`}
                  >
                    {i + 1}
                  </span>
                  <span className="text-sm font-medium text-gray-800">
                    {entry.username}
                  </span>
                </div>
                <span className="text-sm text-gray-500">
                  {entry.total_cost.toFixed(2)} tokens
                </span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
