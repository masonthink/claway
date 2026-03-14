"use client";

import { useEffect } from "react";

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error("Unhandled error:", error);
  }, [error]);

  return (
    <div className="flex min-h-[50vh] flex-col items-center justify-center px-7 py-20">
      <div
        className="w-full max-w-md rounded-[16px] p-6 text-center"
        style={{
          background: "var(--surface)",
          border: "1px solid var(--line)",
          boxShadow: "var(--shadow-sm)",
        }}
      >
        <h2 className="mb-2 font-display text-lg tracking-[-0.02em]">
          Something went wrong
        </h2>
        <p className="mb-4 text-sm text-ink-soft">
          An unexpected error occurred. Please try again.
        </p>
        <button
          onClick={reset}
          className="inline-flex items-center gap-2 rounded-[10px] px-4 py-2 text-sm font-medium text-white"
          style={{
            background: "linear-gradient(135deg, var(--accent), var(--accent-deep))",
          }}
        >
          Try again
        </button>
      </div>
    </div>
  );
}
