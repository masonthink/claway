const STATUS_STYLES: Record<string, { bg: string; color: string }> = {
  active: { bg: "rgba(43,198,164,0.14)", color: "rgb(26,107,91)" },
  completed: { bg: "rgba(255,107,74,0.14)", color: "var(--accent-deep)" },
  open: { bg: "rgba(42,31,25,0.06)", color: "var(--ink-soft)" },
  claimed: { bg: "rgba(231,187,103,0.18)", color: "#92700a" },
  submitted: { bg: "rgba(59,130,246,0.1)", color: "#2563eb" },
  approved: { bg: "rgba(34,197,94,0.1)", color: "rgb(22,163,74)" },
  rejected: { bg: "rgba(239,68,68,0.1)", color: "rgb(220,38,38)" },
};

const STATUS_LABELS: Record<string, string> = {
  active: "进行中",
  completed: "已完成",
  open: "待认领",
  claimed: "已认领",
  submitted: "已提交",
  approved: "已通过",
  rejected: "已拒绝",
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
