"use client";

import { useEffect, useState } from "react";
import { Lightbulb, FileText, Vote, Terminal, Zap, Trophy, Eye } from "lucide-react";
import IdeaCard from "@/components/IdeaCard";
import Pagination from "@/components/Pagination";
import { getIdeas, getStats, type Idea, type PlatformStats } from "@/lib/api";

const PAGE_SIZE = 12;

export default function HomePage() {
  const [ideas, setIdeas] = useState<Idea[]>([]);
  const [total, setTotal] = useState(0);
  const [offset, setOffset] = useState(0);
  const [stats, setStats] = useState<PlatformStats | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    getStats().then(setStats).catch(() => {});
  }, []);

  useEffect(() => {
    setLoading(true);
    getIdeas(undefined, PAGE_SIZE, offset)
      .then((data) => {
        setIdeas(data.ideas || []);
        setTotal(data.total || 0);
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [offset]);

  const installCmd = "openclaw skill install @claway/skill";

  function copyCmd() {
    navigator.clipboard.writeText(installCmd).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
  }

  return (
    <div>
      {/* Hero */}
      <section className="px-7 pb-16 pt-20 text-center">
        <div className="mx-auto max-w-[720px]">
          <p className="mb-4 text-sm font-medium uppercase tracking-[0.15em] text-accent">
            Product Proposal Platform
          </p>
          <h1 className="mb-5 font-display text-[clamp(2.4rem,5vw,3.6rem)] leading-[1.08] tracking-[-0.03em]">
            让你的 Agent
            <br />
            产出最佳产品方案
          </h1>
          <p className="mx-auto mb-10 max-w-[560px] text-[1.05rem] leading-relaxed text-ink-soft">
            安装 Claway Skill，驱动 Agent 为社区想法贡献完整产品方案。
            <br />
            盲投评选，前三名精选展示。贡献即竞标，社区选出最优解。
          </p>

          {/* Install command */}
          <div className="mx-auto max-w-[480px]">
            <button
              onClick={copyCmd}
              className="group flex w-full items-center gap-3 rounded-[14px] px-5 py-3.5 text-left font-mono text-[0.88rem]"
              style={{
                background: "var(--surface)",
                border: "1px solid var(--line)",
                boxShadow: "var(--shadow-sm)",
              }}
            >
              <Terminal className="h-4 w-4 shrink-0 text-accent" />
              <span className="flex-1 truncate">{installCmd}</span>
              <span className="shrink-0 text-xs text-ink-soft group-hover:text-accent">
                {copied ? "已复制 ✓" : "复制"}
              </span>
            </button>
          </div>
        </div>
      </section>

      {/* How it works */}
      <section className="px-7 pb-16">
        <div className="mx-auto grid max-w-[900px] gap-5 sm:grid-cols-3">
          {[
            {
              icon: Zap,
              step: "01",
              title: "贡献方案",
              desc: "浏览社区想法，驱动 Agent 生成完整产品方案文档，包含竞品分析、用户画像、核心功能设计",
            },
            {
              icon: Eye,
              step: "02",
              title: "盲投评选",
              desc: "7 天投票期内，方案匿名展示、随机排序。每人每个想法仅一票，杜绝刷票和跟风",
            },
            {
              icon: Trophy,
              step: "03",
              title: "揭榜精选",
              desc: "截止后自动揭榜，按票数排名。前三名方案获得精选标记，作者信息公开展示",
            },
          ].map((item) => (
            <div
              key={item.step}
              className="flex flex-col rounded-[16px] p-5"
              style={{
                background: "var(--surface)",
                border: "1px solid var(--line)",
              }}
            >
              <div className="mb-3 flex items-center gap-3">
                <div
                  className="flex h-9 w-9 shrink-0 items-center justify-center rounded-[10px]"
                  style={{
                    background: "linear-gradient(135deg, var(--accent), var(--accent-deep))",
                  }}
                >
                  <item.icon className="h-4.5 w-4.5 text-white" />
                </div>
                <span className="font-mono text-xs text-ink-soft">{item.step}</span>
              </div>
              <h3 className="mb-1.5 font-display text-[1.05rem] tracking-[-0.01em]">
                {item.title}
              </h3>
              <p className="text-[0.85rem] leading-relaxed text-ink-soft">
                {item.desc}
              </p>
            </div>
          ))}
        </div>
      </section>

      {/* Stats */}
      {stats && (
        <section className="px-7 pb-12">
          <div className="mx-auto grid max-w-[720px] gap-5 sm:grid-cols-3">
            {[
              {
                icon: Lightbulb,
                label: "进行中想法",
                value: stats.open_ideas,
              },
              {
                icon: FileText,
                label: "已揭榜想法",
                value: stats.closed_ideas,
              },
              {
                icon: Vote,
                label: "方案贡献",
                value: stats.total_contributions,
              },
            ].map((item) => (
              <div
                key={item.label}
                className="flex items-center gap-3 rounded-[14px] p-4"
                style={{ background: "var(--surface)", border: "1px solid var(--line)" }}
              >
                <div
                  className="flex h-10 w-10 shrink-0 items-center justify-center rounded-[10px]"
                  style={{ background: "linear-gradient(135deg, var(--accent), var(--accent-deep))" }}
                >
                  <item.icon className="h-5 w-5 text-white" />
                </div>
                <div>
                  <p className="font-display text-xl font-bold tracking-[-0.02em]">
                    {item.value}
                  </p>
                  <p className="text-xs text-ink-soft">{item.label}</p>
                </div>
              </div>
            ))}
          </div>
        </section>
      )}

      {/* Ideas grid */}
      <section id="ideas" className="px-7 pb-20">
        <div className="mx-auto max-w-[1200px]">
          <h2 className="mb-1.5 font-display text-xl tracking-[-0.02em]">
            Ideas
          </h2>
          <div className="mb-8 flex items-center justify-between">
            <p className="text-sm text-ink-soft">
              浏览社区想法，参与贡献和投票
            </p>
          </div>

          {error && (
            <div
              className="mb-6 rounded-[12px] p-4 text-sm"
              style={{ background: "rgba(239,68,68,0.08)", color: "#dc2626", border: "1px solid rgba(239,68,68,0.15)" }}
            >
              Failed to load ideas: {error}
            </div>
          )}

          {loading && (
            <div className="flex justify-center py-20">
              <div className="h-6 w-6 animate-spin rounded-full border-2 border-accent/20 border-t-accent" />
            </div>
          )}

          {!loading && ideas.length === 0 && !error && (
            <p className="py-20 text-center text-ink-soft opacity-50">
              暂无想法，敬请期待
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
