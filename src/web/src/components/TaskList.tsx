"use client";

import { useState } from "react";
import { Cpu, User, Hand, XCircle, Send, ClipboardCheck } from "lucide-react";
import StatusBadge from "./StatusBadge";
import TaskSubmitModal from "./TaskSubmitModal";
import TaskReviewModal from "./TaskReviewModal";
import { claimTask, unclaimTask, getTaskDocument, type Task } from "@/lib/api";
import { isLoggedIn } from "@/lib/auth";
import { useToast } from "./Toast";

interface TaskListProps {
  tasks: Task[];
  currentUsername?: string;
  isInitiator?: boolean;
  onRefresh?: () => void;
}

export default function TaskList({
  tasks,
  currentUsername,
  isInitiator,
  onRefresh,
}: TaskListProps) {
  const [loadingId, setLoadingId] = useState<string | null>(null);
  const [submitModal, setSubmitModal] = useState<Task | null>(null);
  const [reviewModal, setReviewModal] = useState<{
    task: Task;
    content: string;
  } | null>(null);
  const { toast } = useToast();

  const loggedIn = isLoggedIn();

  const handleClaim = async (taskId: string) => {
    setLoadingId(taskId);
    try {
      await claimTask(taskId);
      toast("success", "认领成功");
      onRefresh?.();
    } catch (err) {
      toast("error", err instanceof Error ? err.message : "认领失败");
    } finally {
      setLoadingId(null);
    }
  };

  const handleUnclaim = async (taskId: string) => {
    setLoadingId(taskId);
    try {
      await unclaimTask(taskId);
      toast("success", "已放弃任务");
      onRefresh?.();
    } catch (err) {
      toast("error", err instanceof Error ? err.message : "放弃失败");
    } finally {
      setLoadingId(null);
    }
  };

  return (
    <>
      <div
        className="overflow-hidden rounded-[16px]"
        style={{ border: "1px solid var(--line)", background: "var(--surface)" }}
      >
        {tasks.map((task, i) => {
          const isClaimedByMe =
            currentUsername && task.claimed_by === currentUsername;

          return (
            <div
              key={task.id}
              className="flex items-center justify-between gap-3 px-5 py-3.5"
              style={{
                borderBottom:
                  i === tasks.length - 1 ? "none" : "1px solid var(--line)",
              }}
            >
              <div className="flex items-center gap-3 min-w-0">
                <span
                  className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg text-[0.7rem] font-bold"
                  style={{
                    background: "rgba(255,107,74,0.1)",
                    color: "var(--accent-deep)",
                  }}
                >
                  {task.type}
                </span>
                <div className="min-w-0">
                  <p className="truncate text-sm font-medium">{task.title}</p>
                  {task.claimed_by && (
                    <p className="flex items-center gap-1 text-xs text-ink-soft">
                      <User className="h-3 w-3" />
                      @{task.claimed_by}
                    </p>
                  )}
                </div>
              </div>

              <div className="flex shrink-0 items-center gap-2.5">
                {task.cost_usd_accumulated > 0 && (
                  <span className="flex items-center gap-1 text-xs text-ink-soft">
                    <Cpu className="h-3 w-3" />
                    {task.cost_usd_accumulated.toFixed(2)}
                  </span>
                )}
                <StatusBadge status={task.status} />

                {/* Action buttons */}
                {loggedIn && task.status === "open" && (
                  <button
                    onClick={() => handleClaim(task.id)}
                    disabled={loadingId === task.id}
                    className="inline-flex items-center gap-1 rounded-[8px] px-2.5 py-1.5 text-xs font-semibold hover:-translate-y-0.5 disabled:opacity-60 disabled:hover:translate-y-0"
                    style={{
                      background: "rgba(43,198,164,0.12)",
                      color: "rgb(26,107,91)",
                      border: "1px solid rgba(43,198,164,0.2)",
                    }}
                  >
                    <Hand className="h-3 w-3" />
                    {loadingId === task.id ? "..." : "认领"}
                  </button>
                )}

                {loggedIn && task.status === "claimed" && isClaimedByMe && (
                  <>
                    <button
                      onClick={() => setSubmitModal(task)}
                      className="inline-flex items-center gap-1 rounded-[8px] px-2.5 py-1.5 text-xs font-semibold hover:-translate-y-0.5"
                      style={{
                        background: "rgba(59,130,246,0.1)",
                        color: "#2563eb",
                        border: "1px solid rgba(59,130,246,0.2)",
                      }}
                    >
                      <Send className="h-3 w-3" />
                      提交
                    </button>
                    <button
                      onClick={() => handleUnclaim(task.id)}
                      disabled={loadingId === task.id}
                      className="inline-flex items-center gap-1 rounded-[8px] px-2.5 py-1.5 text-xs font-semibold hover:-translate-y-0.5 disabled:opacity-60 disabled:hover:translate-y-0"
                      style={{
                        background: "rgba(239,68,68,0.08)",
                        color: "rgb(220,38,38)",
                        border: "1px solid rgba(239,68,68,0.15)",
                      }}
                    >
                      <XCircle className="h-3 w-3" />
                      {loadingId === task.id ? "..." : "放弃"}
                    </button>
                  </>
                )}

                {loggedIn &&
                  task.status === "submitted" &&
                  isInitiator && (
                    <button
                      onClick={async () => {
                        try {
                          const doc = await getTaskDocument(task.id);
                          setReviewModal({
                            task,
                            content: doc.content,
                          });
                        } catch {
                          setReviewModal({
                            task,
                            content: "（无法加载文档内容）",
                          });
                        }
                      }}
                      className="inline-flex items-center gap-1 rounded-[8px] px-2.5 py-1.5 text-xs font-semibold hover:-translate-y-0.5"
                      style={{
                        background: "rgba(231,187,103,0.15)",
                        color: "#92700a",
                        border: "1px solid rgba(231,187,103,0.25)",
                      }}
                    >
                      <ClipboardCheck className="h-3 w-3" />
                      审核
                    </button>
                  )}
              </div>
            </div>
          );
        })}
      </div>

      {/* Submit modal */}
      {submitModal && (
        <TaskSubmitModal
          taskId={submitModal.id}
          taskName={submitModal.title}
          onClose={() => setSubmitModal(null)}
          onSuccess={() => {
            setSubmitModal(null);
            onRefresh?.();
          }}
        />
      )}

      {/* Review modal */}
      {reviewModal && (
        <TaskReviewModal
          taskId={reviewModal.task.id}
          taskName={reviewModal.task.title}
          outputContent={reviewModal.content}
          onClose={() => setReviewModal(null)}
          onSuccess={() => {
            setReviewModal(null);
            onRefresh?.();
          }}
        />
      )}
    </>
  );
}
