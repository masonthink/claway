"use client";

import { useEffect, useState, useMemo } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, FileText, List } from "lucide-react";
import MarkdownRenderer from "@/components/MarkdownRenderer";
import StatusBadge from "@/components/StatusBadge";
import { getDraftPreview, type Contribution } from "@/lib/api";
import ErrorState from "@/components/ErrorState";

interface TocItem {
  level: number;
  text: string;
  id: string;
}

function extractToc(content: string): TocItem[] {
  const headingRegex = /^(#{1,4})\s+(.+)$/gm;
  const items: TocItem[] = [];
  let match;
  while ((match = headingRegex.exec(content)) !== null) {
    const text = match[2].trim();
    const id = text
      .toLowerCase()
      .replace(/[^\w]+/g, "-")
      .replace(/^-|-$/g, "");
    items.push({
      level: match[1].length,
      text,
      id,
    });
  }
  return items;
}

export default function DraftPreviewPage() {
  const { id } = useParams<{ id: string }>();
  const [contrib, setContrib] = useState<Contribution | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [showToc, setShowToc] = useState(true);

  useEffect(() => {
    if (!id) return;
    getDraftPreview(id)
      .then(setContrib)
      .catch((err) => setError(err.message));
  }, [id]);

  const toc = useMemo(() => {
    if (!contrib?.content) return [];
    return extractToc(contrib.content);
  }, [contrib?.content]);

  const reload = () => {
    if (!id) return;
    setError(null);
    getDraftPreview(id).then(setContrib).catch((err) => setError(err.message));
  };

  if (error) {
    return (
      <div className="mx-auto max-w-[860px] px-7 py-12">
        <ErrorState message={error} onRetry={reload} />
      </div>
    );
  }

  if (!contrib) {
    return null; // loading.tsx handles this
  }

  return (
    <div className="mx-auto max-w-[1100px] px-7 py-8">
      <Link href="/" className="mb-6 inline-flex items-center gap-1.5 text-sm text-ink-soft hover:text-ink">
        <ArrowLeft className="h-4 w-4" />
        Back
      </Link>

      {/* Header */}
      <div
        className="mb-6 flex items-center justify-between rounded-[16px] p-5"
        style={{ background: "var(--surface)", border: "1px solid var(--line)" }}
      >
        <div className="flex items-center gap-3">
          <FileText className="h-5 w-5 text-accent" />
          <div>
            <h1 className="font-display text-lg tracking-[-0.02em]">Draft Preview</h1>
            <p className="text-xs text-ink-soft">
              Contribution #{contrib.id} - Last updated {new Date(contrib.updated_at).toLocaleDateString("en-US")}
            </p>
          </div>
        </div>
        <div className="flex items-center gap-3">
          <StatusBadge status={contrib.status} />
          {toc.length > 0 && (
            <button
              onClick={() => setShowToc(!showToc)}
              className="flex items-center gap-1.5 rounded-[8px] px-3 py-1.5 text-xs font-medium text-ink-soft hover:text-ink"
              style={{ border: "1px solid var(--line)" }}
            >
              <List className="h-3.5 w-3.5" />
              Contents
            </button>
          )}
        </div>
      </div>

      <div className="flex gap-6">
        {/* TOC sidebar */}
        {showToc && toc.length > 0 && (
          <aside className="hidden w-56 shrink-0 lg:block">
            <nav
              className="sticky top-20 rounded-[14px] p-4"
              style={{ background: "var(--surface)", border: "1px solid var(--line)" }}
            >
              <p className="mb-3 text-xs font-semibold uppercase tracking-wider text-ink-soft">
                Contents
              </p>
              <ul className="space-y-1">
                {toc.map((item) => (
                  <li
                    key={item.id}
                    style={{ paddingLeft: `${(item.level - 1) * 12}px` }}
                  >
                    <a
                      href={`#${item.id}`}
                      className="block truncate rounded-[6px] px-2 py-1 text-xs text-ink-soft hover:bg-surface-muted hover:text-ink"
                    >
                      {item.text}
                    </a>
                  </li>
                ))}
              </ul>
            </nav>
          </aside>
        )}

        {/* Content */}
        <div className="min-w-0 flex-1">
          <MarkdownRenderer content={contrib.content} />
        </div>
      </div>
    </div>
  );
}
