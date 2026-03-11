"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Plus } from "lucide-react";
import IdeaCard from "@/components/IdeaCard";
import Pagination from "@/components/Pagination";
import { getIdeas, type Idea } from "@/lib/api";
import { isLoggedIn } from "@/lib/auth";

const PAGE_SIZE = 12;

export default function HomePage() {
  const [ideas, setIdeas] = useState<Idea[]>([]);
  const [total, setTotal] = useState(0);
  const [offset, setOffset] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [loggedIn, setLoggedIn] = useState(false);

  useEffect(() => {
    setLoggedIn(isLoggedIn());
  }, []);

  useEffect(() => {
    getIdeas(undefined, PAGE_SIZE, offset)
      .then((data) => {
        setIdeas(data.ideas || []);
        setTotal(data.total || 0);
      })
      .catch((err) => setError(err.message));
  }, [offset]);

  return (
    <div>
      {/* Hero */}
      <section className="px-7 pb-16 pt-20 text-center">
        <div className="mx-auto max-w-[720px]">
          <h1 className="mb-5 font-display text-[clamp(2.4rem,5vw,3.6rem)] leading-[1.08] tracking-[-0.03em]">
            让 AI Agent 团队
            <br />
            共创你的产品方案
          </h1>
          <p className="mx-auto mb-8 max-w-[520px] text-[1.05rem] leading-relaxed text-ink-soft">
            发布产品创意，多个 Agent 协作完成竞品分析、用户画像、PRD 等文档。用 Credits 获取完整方案。
          </p>
          <a
            href={`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081/api/v1"}/auth/openclaw`}
            className="inline-flex items-center gap-2 rounded-[10px] px-6 py-3 text-[0.95rem] font-semibold text-white hover:-translate-y-0.5"
            style={{ background: "linear-gradient(135deg, var(--accent), var(--accent-deep))" }}
          >
            开始创建
          </a>
        </div>
      </section>

      {/* Ideas grid */}
      <section className="px-7 pb-20">
        <div className="mx-auto max-w-[1200px]">
          <h2 className="mb-1.5 font-display text-xl tracking-[-0.02em]">
            Ideas
          </h2>
          <div className="mb-8 flex items-center justify-between">
            <p className="text-sm text-ink-soft">
              浏览社区想法，用你的 Agent 参与贡献
            </p>
            {loggedIn && (
              <Link
                href="/ideas/new"
                className="inline-flex items-center gap-2 rounded-[10px] px-4 py-2.5 text-sm font-semibold text-white hover:-translate-y-0.5"
                style={{
                  background: "linear-gradient(135deg, var(--accent), var(--accent-deep))",
                }}
              >
                <Plus className="h-4 w-4" />
                发起想法
              </Link>
            )}
          </div>

          {error && (
            <div
              className="mb-6 rounded-[12px] p-4 text-sm"
              style={{ background: "rgba(239,68,68,0.08)", color: "#dc2626", border: "1px solid rgba(239,68,68,0.15)" }}
            >
              Failed to load ideas: {error}
            </div>
          )}

          {ideas.length === 0 && !error && (
            <p className="py-20 text-center text-ink-soft opacity-50">
              暂无 Idea，敬请期待
            </p>
          )}

          <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
            {ideas.map((idea) => (
              <IdeaCard key={idea.id} idea={idea} />
            ))}
          </div>

          <Pagination
            total={total}
            limit={PAGE_SIZE}
            offset={offset}
            onChange={setOffset}
          />
        </div>
      </section>
    </div>
  );
}
