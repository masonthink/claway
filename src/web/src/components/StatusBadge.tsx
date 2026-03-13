const STATUS_STYLES: Record<string, { bg: string; color: string }> = {
  open: { bg: "rgba(43,198,164,0.14)", color: "rgb(26,107,91)" },
  closed: { bg: "rgba(255,107,74,0.14)", color: "var(--accent-deep)" },
  cancelled: { bg: "rgba(42,31,25,0.06)", color: "var(--ink-soft)" },
  draft: { bg: "rgba(231,187,103,0.18)", color: "#92700a" },
  submitted: { bg: "rgba(59,130,246,0.1)", color: "#2563eb" },
};

const STATUS_LABELS: Record<string, string> = {
  open: "进行中",
  closed: "已揭榜",
  cancelled: "已取消",
  draft: "草稿",
  submitted: "已提交",
};

export default function StatusBadge({ status }: { status: string }) {
  const style = STATUS_STYLES[status] || STATUS_STYLES.open;
  const label = STATUS_LABELS[status] || status;

  return (
    <span
      className="inline-flex shrink-0 items-center rounded-[8px] px-2.5 py-0.5 text-[0.75rem] font-semibold leading-tight"
      style={{ background: style.bg, color: style.color }}
    >
      {label}
    </span>
  );
}
