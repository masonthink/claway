"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Lock } from "lucide-react";
import MarkdownRenderer from "@/components/MarkdownRenderer";
import { getPRD, type PRD } from "@/lib/api";

export default function PRDViewPage() {
  const { id } = useParams<{ id: string }>();
  const [prd, setPrd] = useState<PRD | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;
    getPRD(id)
      .then(setPrd)
      .catch((err) => setError(err.message));
  }, [id]);

  if (error) {
    return (
      <div className="mx-auto max-w-[860px] px-7 py-12">
        <div className="rounded-[12px] p-4 text-sm" style={{ background: "rgba(239,68,68,0.08)", color: "#dc2626", border: "1px solid rgba(239,68,68,0.15)" }}>
          {error}
        </div>
      </div>
    );
  }

  if (!prd) {
    return (
      <div className="mx-auto max-w-[860px] px-7 py-12 text-center text-ink-soft">
        Loading...
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-[860px] px-7 py-8">
      <Link href={`/ideas/${prd.idea_id}`} className="mb-6 inline-flex items-center gap-1.5 text-sm text-ink-soft hover:text-ink">
        <ArrowLeft className="h-4 w-4" />
        返回 Idea
      </Link>

      <div
        className="rounded-[20px] p-6"
        style={{ background: "var(--surface)", border: "1px solid var(--line)", boxShadow: "var(--shadow-sm)" }}
      >
        <h1 className="mb-6 font-display text-2xl tracking-[-0.02em]">{prd.title}</h1>

        {prd.purchased ? (
          <MarkdownRenderer content={prd.content} />
        ) : (
          <div>
            <div className="relative">
              <MarkdownRenderer content={prd.preview} />
              <div className="absolute inset-x-0 bottom-0 h-32" style={{ background: "linear-gradient(transparent, var(--surface))" }} />
            </div>

            <div
              className="mt-6 rounded-[16px] p-8 text-center"
              style={{ background: "var(--surface-muted)", border: "1px solid var(--line)" }}
            >
              <Lock className="mx-auto mb-3 h-7 w-7 text-ink-soft opacity-50" />
              <p className="mb-1 text-sm font-medium" style={{ color: "var(--ink)" }}>
                完整内容即将开放
              </p>
              <p className="text-sm text-ink-soft">
                PRD 购买功能正在开发中，敬请期待
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
