"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Trophy, Medal } from "lucide-react";
import MarkdownRenderer from "@/components/MarkdownRenderer";
import ErrorState from "@/components/ErrorState";
import {
  getIdea,
  getRevealResult,
  getContributions,
  type Idea,
  type RevealResult,
  type Contribution,
} from "@/lib/api";

const RANK_STYLES: Record<number, { bg: string; color: string; icon: string }> = {
  1: { bg: "rgba(231,187,103,0.2)", color: "#92700a", icon: "gold" },
  2: { bg: "rgba(192,192,192,0.2)", color: "#6b7280", icon: "silver" },
  3: { bg: "rgba(205,127,50,0.2)", color: "#92400e", icon: "bronze" },
};

export default function RevealResultPage() {
  const { id } = useParams<{ id: string }>();
  const [idea, setIdea] = useState<Idea | null>(null);
  const [result, setResult] = useState<RevealResult | null>(null);
  const [contribContents, setContribContents] = useState<Record<number, string>>({});
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;

    Promise.all([
      getIdea(id).then(setIdea),
      getRevealResult(id).then(setResult),
      getContributions(id).then((contribs) => {
        const contents: Record<number, string> = {};
        for (const c of contribs) {
          if (c.content) contents[c.id] = c.content;
        }
        setContribContents(contents);
      }),
    ]).catch((err) => setError(err.message));
  }, [id]);

  const reload = () => {
    if (!id) return;
    setError(null);
    Promise.all([
      getIdea(id).then(setIdea),
      getRevealResult(id).then(setResult),
      getContributions(id).then((contribs) => {
        const contents: Record<number, string> = {};
        for (const c of contribs) {
          if (c.content) contents[c.id] = c.content;
        }
        setContribContents(contents);
      }),
    ]).catch((err) => setError(err.message));
  };

  if (error) {
    return (
      <div className="mx-auto max-w-[860px] px-7 py-12">
        <ErrorState message={error} onRetry={reload} />
      </div>
    );
  }

  if (!idea || !result) {
    return (
      <div className="mx-auto max-w-[860px] px-7 py-12 text-center text-ink-soft">
        Loading...
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-[860px] px-7 py-8">
      <Link href={`/idea/${id}`} className="mb-6 inline-flex items-center gap-1.5 text-sm text-ink-soft hover:text-ink">
        <ArrowLeft className="h-4 w-4" />
        Back to idea
      </Link>

      {/* Header */}
      <div
        className="mb-8 rounded-[20px] p-6"
        style={{ background: "var(--surface)", border: "1px solid var(--line)", boxShadow: "var(--shadow-sm)" }}
      >
        <div className="mb-2 flex items-center gap-2">
          <Trophy className="h-5 w-5 text-gold" />
          <h1 className="font-display text-2xl tracking-[-0.02em]">Results</h1>
        </div>
        <p className="mb-3 text-[0.95rem] text-ink-soft">{idea.title}</p>
        <div className="flex flex-wrap gap-4 text-sm text-ink-soft">
          <span>Total votes: {result.total_votes}</span>
          <span>Revealed: {new Date(result.revealed_at).toLocaleDateString("en-US")}</span>
        </div>
      </div>

      {/* Rankings */}
      <div className="space-y-6">
        {result.results.map((entry) => {
          const style = RANK_STYLES[entry.rank] || { bg: "var(--surface)", color: "var(--ink-soft)", icon: "" };
          const content = contribContents[entry.contribution_id];

          return (
            <div
              key={entry.contribution_id}
              className="overflow-hidden rounded-[20px]"
              style={{ border: "1px solid var(--line)", background: "var(--surface)" }}
            >
              {/* Rank header */}
              <div
                className="flex items-center justify-between px-6 py-4"
                style={{ background: style.bg }}
              >
                <div className="flex items-center gap-3">
                  <span
                    className="flex h-8 w-8 items-center justify-center rounded-[10px] text-sm font-bold"
                    style={{ background: style.bg, color: style.color }}
                  >
                    {entry.rank <= 3 ? <Medal className="h-5 w-5" /> : `#${entry.rank}`}
                  </span>
                  <div>
                    <span className="text-sm font-semibold" style={{ color: style.color }}>
                      #{entry.rank}
                    </span>
                    {entry.is_featured && (
                      <span className="ml-2 text-xs font-medium text-accent">
                        Featured
                      </span>
                    )}
                  </div>
                </div>
                <div className="text-right">
                  <p className="text-sm font-semibold">{entry.vote_count} votes</p>
                  <Link
                    href={`/user/${entry.author_username}`}
                    className="text-xs text-accent hover:underline"
                  >
                    @{entry.author_username}
                  </Link>
                </div>
              </div>

              {/* Content */}
              {content && (
                <div className="p-6">
                  <MarkdownRenderer content={content} />
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
