// Colored status badge component

const STATUS_STYLES: Record<string, string> = {
  active: "bg-green-100 text-green-700",
  completed: "bg-indigo-100 text-indigo-700",
  open: "bg-gray-100 text-gray-600",
  claimed: "bg-yellow-100 text-yellow-700",
  submitted: "bg-blue-100 text-blue-700",
  approved: "bg-green-100 text-green-700",
  rejected: "bg-red-100 text-red-700",
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
  const style = STATUS_STYLES[status] || "bg-gray-100 text-gray-600";
  const label = STATUS_LABELS[status] || status;

  return (
    <span
      className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${style}`}
    >
      {label}
    </span>
  );
}
