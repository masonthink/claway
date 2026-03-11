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
  getMe,
  publishIdea,
  type Idea,
  type Task,
  type ComputeLeaderboard,
  type User,
} from "@/lib/api";
import { isLoggedIn } from "@/lib/auth";
import { useToast } from "@/components/Toast";

export default function IdeaDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [idea, setIdea] = useState<Idea | null>(null);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [leaderboard, setLeaderboard] = useState<ComputeLeaderboard | null>(null);
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [publishing, setPublishing] = useState(false);
  const { toast } = useToast();

  const loadData = () => {
    if (!id) return;
    Promise.all([
      getIdea(id).then(setIdea),
      getIdeaTasks(id).then((d) => setTasks(d.tasks)),
      getIdeaCompute(id).then(setLeaderboard),
    ]).catch((err) => setError(err.message));
  };

  useEffect(() => {
    loadData();
    if (isLoggedIn()) {
      getMe().then(setCurrentUser).catch(() => {});
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  const handlePublish = async () => {
    if (!id) return;
    setPublishing(true);
    try {
      await publishIdea(id);
      toast("success", "PRD 发布成功");
      loadData();
    } catch (err) {
      toast("error", err instanceof Error ? err.message : "发布失败");
    } finally {
      setPublishing(false);
    }
  };

  if (error) {
    return (
      <div className="mx-auto max-w-[860px] px-7 py-12">
        <div className="rounded-[12px] p-4 text-sm" style={{ background: "rgba(239,68,68,0.08)", color: "#dc2626", border: "1px solid rgba(239,68,68,0.15)" }}>
          Failed to load idea: {error}
        </div>
      </div>
    );
  }

  if (!idea) {
    return (
      <div className="mx-auto max-w-[860px] px-7 py-12 text-center text-ink-soft">
        Loading...
      </div>
    );
  }

  const completed = tasks.filter((t) => t.status === "approved").length;
  const total = tasks.length;
  const progress = total > 0 ? (completed / total) * 100 : 0;

  return (
    <div className="mx-auto max-w-[860px] px-7 py-8">
      <Link href="/" className="mb-6 inline-flex items-center gap-1.5 text-sm text-ink-soft hover:text-ink">
        <ArrowLeft className="h-4 w-4" />
        返回列表
      </Link>

      {/* Header card */}
      <div
        className="mb-8 rounded-[20px] p-6"
        style={{ background: "var(--surface)", border: "1px solid var(--line)", boxShadow: "var(--shadow-sm)" }}
      >
        <div className="mb-3 flex items-start justify-between gap-3">
          <div>
            <h1 className="font-display text-2xl tracking-[-0.02em]">{idea.title}</h1>
            <p className="mt-1.5 flex flex-wrap items-center gap-2 text-sm text-ink-soft">
              <span>by @{idea.initiator}</span>
              <span
                className="rounded-[8px] px-2 py-0.5 text-[0.75rem] font-medium"
                style={{ background: "rgba(43,198,164,0.12)", color: "rgb(26,107,91)" }}
              >
                {idea.package_type}
              </span>
            </p>
          </div>
          <StatusBadge status={idea.status} />
        </div>

        <p className="mb-5 text-[0.95rem] leading-relaxed text-ink-soft">{idea.description}</p>

        {/* Progress */}
        <div className="mb-4">
          <div className="progress-bar">
            <div className="progress-bar-fill" style={{ width: `${progress}%` }} />
          </div>
          <p className="mt-1.5 text-xs text-ink-soft">{completed}/{total} 任务完成</p>
        </div>

        {idea.status === "completed" && (
          <Link
            href={`/prd/${idea.id}`}
            className="inline-flex items-center gap-2 rounded-[10px] px-5 py-2.5 text-sm font-semibold text-white hover:-translate-y-0.5"
            style={{ background: "linear-gradient(135deg, var(--accent), var(--accent-deep))" }}
          >
            <FileText className="h-4 w-4" />
            查看 PRD
          </Link>
        )}
      </div>

      {/* Publish PRD button */}
      {currentUser &&
        idea.initiator === currentUser.username &&
        idea.status === "active" &&
        tasks.length > 0 &&
        tasks.every((t) => t.status === "approved") && (
          <div className="mb-8">
            <button
              onClick={handlePublish}
              disabled={publishing}
              className="inline-flex items-center gap-2 rounded-[10px] px-5 py-2.5 text-sm font-semibold text-white hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-60 disabled:hover:translate-y-0"
              style={{
                background: "linear-gradient(135deg, var(--accent), var(--accent-deep))",
              }}
            >
              <FileText className="h-4 w-4" />
              {publishing ? "发布中..." : "发布 PRD"}
            </button>
          </div>
        )}

      {/* Tasks */}
      <div className="mb-8">
        <h2 className="mb-4 font-display text-lg tracking-[-0.02em]">Tasks</h2>
        <TaskList
          tasks={tasks}
          currentUsername={currentUser?.username}
          isInitiator={currentUser?.username === idea.initiator}
          onRefresh={loadData}
        />
      </div>

      {/* Leaderboard */}
      {leaderboard && leaderboard.entries.length > 0 && (
        <div>
          <h2 className="mb-4 flex items-center gap-2 font-display text-lg tracking-[-0.02em]">
            <Trophy className="h-5 w-5 text-gold" />
            Compute Leaderboard
          </h2>
          <div className="overflow-hidden rounded-[16px]" style={{ border: "1px solid var(--line)", background: "var(--surface)" }}>
            {leaderboard.entries.map((entry, i) => (
              <div
                key={entry.username}
                className="flex items-center justify-between px-5 py-3"
                style={{ borderBottom: i === leaderboard.entries.length - 1 ? "none" : "1px solid var(--line)" }}
              >
                <div className="flex items-center gap-3">
                  <span
                    className="flex h-6 w-6 items-center justify-center rounded-lg text-[0.7rem] font-bold"
                    style={{
                      background: i === 0 ? "rgba(231,187,103,0.2)" : "rgba(42,31,25,0.06)",
                      color: i === 0 ? "#92700a" : "var(--ink-soft)",
                    }}
                  >
                    {i + 1}
                  </span>
                  <span className="font-mono text-sm">@{entry.username}</span>
                </div>
                <span className="text-sm text-ink-soft">
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
