"use client";

import { useState } from "react";
import { X, Eye, Edit3, Send } from "lucide-react";
import { submitTask } from "@/lib/api";
import { useToast } from "./Toast";

interface TaskSubmitModalProps {
  taskId: string;
  taskName: string;
  onClose: () => void;
  onSuccess: () => void;
}

export default function TaskSubmitModal({
  taskId,
  taskName,
  onClose,
  onSuccess,
}: TaskSubmitModalProps) {
  const [content, setContent] = useState("");
  const [note, setNote] = useState("");
  const [preview, setPreview] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const { toast } = useToast();

  const handleSubmit = async () => {
    if (!content.trim()) {
      toast("error", "请输入产出内容");
      return;
    }
    setSubmitting(true);
    try {
      await submitTask(taskId, {
        content,
        note,
      });
      toast("success", "提交成功");
      onSuccess();
    } catch (err) {
      toast("error", err instanceof Error ? err.message : "提交失败");
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
        className="w-full max-w-[640px] rounded-[20px] p-6"
        style={{
          background: "var(--surface)",
          border: "1px solid var(--line)",
          boxShadow: "var(--shadow-lg)",
        }}
      >
        {/* Header */}
        <div className="mb-5 flex items-center justify-between">
          <h3 className="font-display text-lg tracking-[-0.02em]">
            提交产出 — {taskName}
          </h3>
          <button
            onClick={onClose}
            className="rounded-lg p-1.5 text-ink-soft hover:text-ink"
            style={{ border: "1px solid var(--line)" }}
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        {/* Content area */}
        <div className="mb-4">
          <div className="mb-2 flex items-center gap-2">
            <label className="text-sm font-medium">产出内容</label>
            <button
              onClick={() => setPreview(!preview)}
              className="ml-auto flex items-center gap-1 rounded-lg px-2 py-1 text-xs text-ink-soft hover:text-ink"
              style={{ border: "1px solid var(--line)" }}
            >
              {preview ? (
                <>
                  <Edit3 className="h-3 w-3" />
                  编辑
                </>
              ) : (
                <>
                  <Eye className="h-3 w-3" />
                  预览
                </>
              )}
            </button>
          </div>
          {preview ? (
            <pre
              className="min-h-[200px] whitespace-pre-wrap rounded-[12px] p-4 font-body text-sm leading-relaxed"
              style={{
                background: "var(--surface-muted)",
                border: "1px solid var(--line)",
              }}
            >
              {content || "（无内容）"}
            </pre>
          ) : (
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="输入 Markdown 格式的产出内容..."
              className="w-full resize-y rounded-[12px] p-4 text-sm leading-relaxed outline-none"
              style={{
                background: "var(--surface-muted)",
                border: "1px solid var(--line)",
                minHeight: "200px",
                color: "var(--ink)",
              }}
            />
          )}
        </div>

        {/* Note */}
        <div className="mb-6">
          <label className="mb-2 block text-sm font-medium">
            备注 <span className="text-ink-soft font-normal">（{note.length}/200）</span>
          </label>
          <input
            type="text"
            value={note}
            onChange={(e) => setNote(e.target.value.slice(0, 200))}
            placeholder="简要说明产出内容..."
            className="w-full rounded-[12px] px-4 py-2.5 text-sm outline-none"
            style={{
              background: "var(--surface-muted)",
              border: "1px solid var(--line)",
              color: "var(--ink)",
            }}
          />
        </div>

        {/* Actions */}
        <div className="flex items-center justify-end gap-3">
          <button
            onClick={onClose}
            className="rounded-[10px] px-5 py-2.5 text-sm font-medium text-ink-soft hover:text-ink"
            style={{ border: "1px solid var(--line)" }}
          >
            取消
          </button>
          <button
            onClick={handleSubmit}
            disabled={submitting || !content.trim()}
            className="inline-flex items-center gap-2 rounded-[10px] px-5 py-2.5 text-sm font-semibold text-white hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-60 disabled:hover:translate-y-0"
            style={{
              background: "linear-gradient(135deg, var(--accent), var(--accent-deep))",
            }}
          >
            <Send className="h-3.5 w-3.5" />
            {submitting ? "提交中..." : "提交"}
          </button>
        </div>
      </div>
    </div>
  );
}
