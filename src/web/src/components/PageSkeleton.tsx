export default function PageSkeleton({ variant = "detail" }: { variant?: "home" | "detail" }) {
  if (variant === "home") {
    return (
      <div className="animate-pulse" role="status" aria-label="Loading">
        <div className="px-7 pb-16 pt-20">
          <div className="mx-auto max-w-[720px] space-y-4 text-center">
            <div className="mx-auto h-4 w-40 rounded" style={{ background: "var(--surface-muted)" }} />
            <div className="mx-auto h-12 w-96 rounded" style={{ background: "var(--surface-muted)" }} />
            <div className="mx-auto h-6 w-80 rounded" style={{ background: "var(--surface-muted)" }} />
          </div>
        </div>
        <div className="px-7">
          <div className="mx-auto grid max-w-[1200px] gap-5 sm:grid-cols-2 lg:grid-cols-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-40 rounded-[16px]" style={{ background: "var(--surface)" }} />
            ))}
          </div>
        </div>
        <span className="sr-only">Loading</span>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-[860px] animate-pulse px-7 py-8" role="status" aria-label="Loading">
      <div className="mb-6 h-4 w-20 rounded" style={{ background: "var(--surface-muted)" }} />
      <div className="mb-8 h-48 rounded-[20px]" style={{ background: "var(--surface)" }} />
      <div className="space-y-4">
        <div className="h-32 rounded-[16px]" style={{ background: "var(--surface)" }} />
        <div className="h-32 rounded-[16px]" style={{ background: "var(--surface)" }} />
      </div>
      <span className="sr-only">Loading</span>
    </div>
  );
}
