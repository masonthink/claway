"use client";

import { ChevronLeft, ChevronRight } from "lucide-react";

interface PaginationProps {
  total: number;
  limit: number;
  offset: number;
  onChange: (offset: number) => void;
}

export default function Pagination({
  total,
  limit,
  offset,
  onChange,
}: PaginationProps) {
  const totalPages = Math.ceil(total / limit);
  const currentPage = Math.floor(offset / limit) + 1;

  if (totalPages <= 1) return null;

  const pages: (number | "...")[] = [];
  for (let i = 1; i <= totalPages; i++) {
    if (i === 1 || i === totalPages || Math.abs(i - currentPage) <= 1) {
      pages.push(i);
    } else if (pages[pages.length - 1] !== "...") {
      pages.push("...");
    }
  }

  return (
    <div className="flex items-center justify-center gap-1.5 pt-8">
      <button
        disabled={currentPage === 1}
        onClick={() => onChange((currentPage - 2) * limit)}
        aria-label="上一页"
        className="flex h-9 w-9 items-center justify-center rounded-[8px] disabled:opacity-30"
        style={{ border: "1px solid var(--line)" }}
      >
        <ChevronLeft className="h-4 w-4" aria-hidden="true" />
      </button>

      {pages.map((p, i) =>
        p === "..." ? (
          <span key={`e${i}`} className="px-1 text-ink-soft">
            ...
          </span>
        ) : (
          <button
            key={p}
            onClick={() => onChange((p - 1) * limit)}
            className="flex h-9 w-9 items-center justify-center rounded-[8px] text-sm font-medium"
            style={{
              background:
                p === currentPage ? "var(--accent)" : "transparent",
              color: p === currentPage ? "#fff" : "var(--ink)",
              border:
                p === currentPage ? "none" : "1px solid var(--line)",
            }}
          >
            {p}
          </button>
        )
      )}

      <button
        disabled={currentPage === totalPages}
        onClick={() => onChange(currentPage * limit)}
        aria-label="下一页"
        className="flex h-9 w-9 items-center justify-center rounded-[8px] disabled:opacity-30"
        style={{ border: "1px solid var(--line)" }}
      >
        <ChevronRight className="h-4 w-4" aria-hidden="true" />
      </button>
    </div>
  );
}
