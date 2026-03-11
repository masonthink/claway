"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Lock, ShoppingCart } from "lucide-react";
import MarkdownRenderer from "@/components/MarkdownRenderer";
import { getPRD, purchasePRD, type PRD } from "@/lib/api";
import { isLoggedIn } from "@/lib/auth";

export default function PRDViewPage() {
  const { id } = useParams<{ id: string }>();
  const [prd, setPrd] = useState<PRD | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [purchasing, setPurchasing] = useState(false);

  useEffect(() => {
    if (!id) return;
    getPRD(id)
      .then(setPrd)
      .catch((err) => setError(err.message));
  }, [id]);

  const handlePurchase = async () => {
    if (!id || !isLoggedIn()) {
      alert("请先登录");
      return;
    }
    setPurchasing(true);
    try {
      await purchasePRD(id);
      const updated = await getPRD(id);
      setPrd(updated);
    } catch (err) {
      setError(err instanceof Error ? err.message : "购买失败");
    } finally {
      setPurchasing(false);
    }
  };

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
              <p className="mb-1 text-sm text-ink-soft">查看完整 PRD 需要</p>
              <p className="mb-5 font-display text-2xl font-bold text-accent-deep">
                {prd.price} Credits
              </p>
              <button
                onClick={handlePurchase}
                disabled={purchasing}
                className="inline-flex items-center gap-2 rounded-[10px] px-6 py-2.5 text-sm font-semibold text-white hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-60 disabled:hover:translate-y-0"
                style={{ background: "linear-gradient(135deg, var(--accent), var(--accent-deep))" }}
              >
                <ShoppingCart className="h-4 w-4" />
                {purchasing ? "处理中..." : "购买"}
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
