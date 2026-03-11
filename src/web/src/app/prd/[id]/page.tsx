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
      // Reload PRD data after purchase
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
      <div className="mx-auto max-w-4xl px-4 py-12">
        <div className="rounded-lg bg-red-50 p-4 text-sm text-red-600">
          {error}
        </div>
      </div>
    );
  }

  if (!prd) {
    return (
      <div className="mx-auto max-w-4xl px-4 py-12 text-center text-gray-400">
        Loading...
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-4xl px-4 py-8">
      <Link
        href={`/ideas/${prd.idea_id}`}
        className="mb-6 inline-flex items-center gap-1 text-sm text-gray-500 hover:text-indigo-600 transition-colors"
      >
        <ArrowLeft className="h-4 w-4" />
        返回 Idea
      </Link>

      <div className="rounded-xl border border-gray-200 bg-white p-6">
        <h1 className="mb-6 text-2xl font-bold text-gray-900">{prd.title}</h1>

        {prd.purchased ? (
          <MarkdownRenderer content={prd.content} />
        ) : (
          <div>
            {/* Preview */}
            <div className="relative">
              <MarkdownRenderer content={prd.preview} />
              <div className="absolute inset-x-0 bottom-0 h-32 bg-gradient-to-t from-white to-transparent" />
            </div>

            {/* Purchase prompt */}
            <div className="mt-4 rounded-lg border border-indigo-200 bg-indigo-50 p-6 text-center">
              <Lock className="mx-auto mb-3 h-8 w-8 text-indigo-400" />
              <p className="mb-1 text-sm text-gray-600">
                查看完整 PRD 需要
              </p>
              <p className="mb-4 text-2xl font-bold text-indigo-700">
                {prd.price} Credits
              </p>
              <button
                onClick={handlePurchase}
                disabled={purchasing}
                className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-6 py-2.5 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50 transition-colors"
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
