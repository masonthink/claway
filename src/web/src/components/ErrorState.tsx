"use client";

export default function ErrorState({ message, onRetry }: { message?: string; onRetry?: () => void }) {
  return (
    <div
      className="rounded-[12px] p-6 text-center"
      role="alert"
      style={{ background: "rgba(239,68,68,0.08)", border: "1px solid rgba(239,68,68,0.15)" }}
    >
      <p className="mb-1 text-sm font-medium" style={{ color: "#dc2626" }}>
        加载失败
      </p>
      {message && (
        <p className="mb-3 text-xs text-ink-soft">{message}</p>
      )}
      {onRetry && (
        <button
          onClick={onRetry}
          className="rounded-[8px] px-4 py-2 text-sm font-medium hover:bg-surface-muted"
          style={{ border: "1px solid var(--line)" }}
        >
          重试
        </button>
      )}
    </div>
  );
}
