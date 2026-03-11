"use client";

import { useState } from "react";
import { X, CheckCircle, XCircle } from "lucide-react";
import { reviewTask } from "@/lib/api";
import { useToast } from "./Toast";

interface TaskReviewModalProps {
  taskId: string;
  taskName: string;
  outputContent: string;
  onClose: () => void;
  onSuccess: () => void;
}

const QUALITY_OPTIONS = [
  { label: "达标", score: 1.0, desc: "符合要求" },
  { label: "良好", score: 1.2, desc: "超出预期" },
  { label: "优秀", score: 1.5, desc: "显著优秀" },
];

export default function TaskReviewModal({
  taskId,
  taskName,
  outputContent,
  onClose,
  onSuccess,
}: TaskReviewModalProps) {
  const [qualityScore, setQualityScore] = useState(1.0);
  const [rejectReason, setRejectReason] = useState("");
  const [mode, setMode] = useState<"approve" | "reject" | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const { toast } = useToast();

  const handleApprove = async () => {
    setSubmitting(true);
    try {
      await reviewTask(taskId, { quality_score: qualityScore });
      toast("success", "审核通过");
      onSuccess();
    } catch (err) {
      toast("error", err instanceof Error ? err.message : "审核失败");
    } finally {
      setSubmitting(false);
    }
  };

  const handleReject = async () => {
    if (!rejectReason.trim()) {
      toast("error", "请输入拒绝原因");
      return;
    }
    setSubmitting(true);
    try {
      await reviewTask(taskId, {
        quality_score: 0,
        reject_reason: rejectReason,
      });
      toast("success", "已拒绝该提交");
      onSuccess();
    } catch (err) {
      toast("error", err instanceof Error ? err.message : "审核失败");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div
      className="modal-overlay fixed inset-0 z-40 flex items-center justify-center p-4"
      style={{ background: "rgba(0,0,0,0.45)", backdropFilter: "blur(4px)" }}
      onClick={(e) => e.target === e.currentTarget && onClose()}
    >
      <div
        className="w-full max-w-[640px] max-h-[85vh] overflow-y-auto rounded-[20px] p-6"
        style={{
          background: "var(--surface)",
          border: "1px solid var(--line)",
          boxShadow: "var(--shadow-lg)",
        }}
      >
        {/* Header */}
        <div className="mb-5 flex items-center justify-between">
          <h3 className="font-display text-lg tracking-[-0.02em]">
            审核 — {taskName}
          </h3>
          <button
            onClick={onClose}
            className="rounded-lg p-1.5 text-ink-soft hover:text-ink"
            style={{ border: "1px solid var(--line)" }}
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        {/* Submitted content */}
        <div className="mb-5">
          <label className="mb-2 block text-sm font-medium">提交内容</label>
          <pre
            className="max-h-[240px] overflow-y-auto whitespace-pre-wrap rounded-[12px] p-4 font-body text-sm leading-relaxed"
            style={{
              background: "var(--surface-muted)",
              border: "1px solid var(--line)",
            }}
          >
            {outputContent || "（无内容）"}
          </pre>
        </div>

        {/* Quality score */}
        {mode !== "reject" && (
          <div className="mb-5">
            <label className="mb-2 block text-sm font-medium">质量评分</label>
            <div className="flex gap-2">
              {QUALITY_OPTIONS.map((opt) => (
                <button
                  key={opt.score}
                  onClick={() => {
                    setQualityScore(opt.score);
                    setMode("approve");
                  }}
                  className="flex-1 rounded-[12px] px-3 py-3 text-center text-sm font-medium"
                  style={{
                    background:
                      qualityScore === opt.score && mode === "approve"
                        ? "rgba(34,197,94,0.12)"
                        : "var(--surface-muted)",
                    border: `1px solid ${
                      qualityScore === opt.score && mode === "approve"
                        ? "rgba(34,197,94,0.3)"
                        : "var(--line)"
                    }`,
                    color:
                      qualityScore === opt.score && mode === "approve"
                        ? "rgb(22,163,74)"
                        : "var(--ink)",
                  }}
                >
                  <div className="font-semibold">{opt.label}</div>
                  <div className="mt-0.5 text-xs opacity-60">x{opt.score}</div>
                </button>
              ))}
            </div>
          </div>
        )}

        {/* Reject area */}
        {mode === "reject" && (
          <div className="mb-5">
            <label className="mb-2 block text-sm font-medium">拒绝原因</label>
            <textarea
              value={rejectReason}
              onChange={(e) => setRejectReason(e.target.value)}
              placeholder="请说明拒绝的具体原因..."
              className="w-full resize-y rounded-[12px] p-4 text-sm leading-relaxed outline-none"
              style={{
                background: "var(--surface-muted)",
                border: "1px solid var(--line)",
                minHeight: "100px",
                color: "var(--ink)",
              }}
            />
          </div>
        )}

        {/* Actions */}
        <div className="flex items-center justify-end gap-3">
          <button
            onClick={onClose}
            className="rounded-[10px] px-5 py-2.5 text-sm font-medium text-ink-soft hover:text-ink"
            style={{ border: "1px solid var(--line)" }}
          >
            取消
          </button>

          {mode !== "approve" && (
            <button
              onClick={() => {
                if (mode === "reject") {
                  handleReject();
                } else {
                  setMode("reject");
                }
              }}
              disabled={submitting}
              className="inline-flex items-center gap-2 rounded-[10px] px-5 py-2.5 text-sm font-semibold hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-60 disabled:hover:translate-y-0"
              style={{
                background: "rgba(239,68,68,0.1)",
                color: "rgb(220,38,38)",
                border: "1px solid rgba(239,68,68,0.2)",
              }}
            >
              <XCircle className="h-3.5 w-3.5" />
              {mode === "reject"
                ? submitting
                  ? "处理中..."
                  : "确认拒绝"
                : "拒绝"}
            </button>
          )}

          {mode !== "reject" && (
            <button
              onClick={() => {
                if (mode === "approve") {
                  handleApprove();
                } else {
                  setMode("approve");
                  setQualityScore(1.0);
                }
              }}
              disabled={submitting}
              className="inline-flex items-center gap-2 rounded-[10px] px-5 py-2.5 text-sm font-semibold text-white hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-60 disabled:hover:translate-y-0"
              style={{
                background: "linear-gradient(135deg, var(--seafoam), #1a8a7f)",
              }}
            >
              <CheckCircle className="h-3.5 w-3.5" />
              {mode === "approve"
                ? submitting
                  ? "处理中..."
                  : "确认通过"
                : "通过"}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
