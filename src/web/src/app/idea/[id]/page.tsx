"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Clock, Users, Vote, Eye, Trophy } from "lucide-react";
import StatusBadge from "@/components/StatusBadge";
import MarkdownRenderer from "@/components/MarkdownRenderer";
import {
  getIdea,
  getContributions,
  castVote,
  type Idea,
  type Contribution,
} from "@/lib/api";
import { isLoggedIn } from "@/lib/auth";

function timeLeft(deadline: string): string {
  const diff = new Date(deadline).getTime() - Date.now();
  if (diff <= 0) return "已截止";
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
  if (days > 0) return `${days}天${hours}小时`;
  const mins = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
  if (hours > 0) return `${hours}小时${mins}分`;
  return `${mins}分钟`;
}

export default function IdeaDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [idea, setIdea] = useState<Idea | null>(null);
  const [contributions, setContributions] = useState<Contribution[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [voting, setVoting] = useState<number | null>(null);
  const [voteError, setVoteError] = useState<string | null>(null);
  const [voteSuccess, setVoteSuccess] = useState(false);

  const loadData = () => {
    if (!id) return;
    Promise.all([
      getIdea(id).then(setIdea),
      getContributions(id).then(setContributions),
    ]).catch((err) => setError(err.message));
  };

  useEffect(() => {
    loadData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  const handleVote = async (contributionId: number) => {
    if (!id) return;
    setVoting(contributionId);
    setVoteError(null);
    try {
      await castVote(id, contributionId);
      setVoteSuccess(true);
      loadData();
    } catch (err) {
      setVoteError(err instanceof Error ? err.message : "投票失败");
    } finally {
      setVoting(null);
    }
  };

  if (error) {
    return (
      <div className="mx-auto max-w-[860px] px-7 py-12">
        <div
          className="rounded-[12px] p-4 text-sm"
          style={{ background: "rgba(239,68,68,0.08)", color: "#dc2626", border: "1px solid rgba(239,68,68,0.15)" }}
        >
          Failed to load: {error}
        </div>
      </div>
    );
  }

  if (!idea) {
    return (
      <div className="mx-auto max-w-[860px] px-7 py-12 text-center text-ink-soft">
        Loading...
      </div>
    );
  }

  const isOpen = idea.status === "open";
  const isClosed = idea.status === "closed";

  return (
    <div className="mx-auto max-w-[860px] px-7 py-8">
      <Link href="/" className="mb-6 inline-flex items-center gap-1.5 text-sm text-ink-soft hover:text-ink">
        <ArrowLeft className="h-4 w-4" />
        返回列表
      </Link>

      {/* Header */}
      <div
        className="mb-8 rounded-[20px] p-6"
        style={{ background: "var(--surface)", border: "1px solid var(--line)", boxShadow: "var(--shadow-sm)" }}
      >
        <div className="mb-3 flex items-start justify-between gap-3">
          <div>
            <h1 className="font-display text-2xl tracking-[-0.02em]">{idea.title}</h1>
            <p className="mt-1.5 flex flex-wrap items-center gap-2 text-sm text-ink-soft">
              <span>
                by{" "}
                <Link href={`/user/${idea.initiator_username}`} className="text-accent hover:underline">
                  @{idea.initiator_username}
                </Link>
              </span>
            </p>
          </div>
          <StatusBadge status={idea.status} />
        </div>

        <p className="mb-4 text-[0.95rem] leading-relaxed text-ink-soft">{idea.description}</p>

        {/* Idea details */}
        <div className="mb-4 space-y-3">
          <DetailItem label="目标用户" value={idea.target_user} />
          <DetailItem label="核心问题" value={idea.core_problem} />
          {idea.out_of_scope && <DetailItem label="范围外" value={idea.out_of_scope} />}
        </div>

        {/* Meta */}
        <div className="flex flex-wrap items-center gap-4 text-sm text-ink-soft">
          <span className="flex items-center gap-1">
            <Users className="h-4 w-4" />
            {idea.contribution_count} 贡献
          </span>
          <span className="flex items-center gap-1">
            <Vote className="h-4 w-4" />
            {idea.voter_count} 投票人
          </span>
          {isOpen && (
            <span className="flex items-center gap-1 text-accent">
              <Clock className="h-4 w-4" />
              {timeLeft(idea.deadline)}
            </span>
          )}
          {isClosed && (
            <Link
              href={`/idea/${idea.id}/result`}
              className="ml-auto inline-flex items-center gap-1.5 rounded-[10px] px-4 py-2 text-sm font-semibold text-white hover:-translate-y-0.5"
              style={{ background: "linear-gradient(135deg, var(--accent), var(--accent-deep))" }}
            >
              <Trophy className="h-4 w-4" />
              查看揭榜结果
            </Link>
          )}
        </div>
      </div>

      {/* Vote feedback */}
      {voteError && (
        <div
          className="mb-6 rounded-[12px] p-4 text-sm"
          style={{ background: "rgba(239,68,68,0.08)", color: "#dc2626", border: "1px solid rgba(239,68,68,0.15)" }}
        >
          {voteError}
        </div>
      )}
      {voteSuccess && (
        <div
          className="mb-6 rounded-[12px] p-4 text-sm"
          style={{ background: "rgba(43,198,164,0.1)", color: "rgb(26,107,91)", border: "1px solid rgba(43,198,164,0.2)" }}
        >
          投票成功！
        </div>
      )}

      {/* Contributions */}
      <div>
        <h2 className="mb-4 font-display text-lg tracking-[-0.02em]">
          贡献 ({contributions.length})
        </h2>

        {contributions.length === 0 && (
          <p className="py-8 text-center text-ink-soft opacity-50">
            暂无贡献
          </p>
        )}

        <div className="space-y-4">
          {contributions.map((contrib, idx) => (
            <div
              key={contrib.id}
              className="rounded-[16px] p-5"
              style={{ background: "var(--surface)", border: "1px solid var(--line)" }}
            >
              <div className="mb-3 flex items-start justify-between gap-3">
                <div className="flex items-center gap-2">
                  {isClosed && contrib.author_id ? (
                    <Link
                      href={`/user/${contrib.author_name || ""}`}
                      className="font-mono text-sm text-accent hover:underline"
                    >
                      @{contrib.author_name || `User #${contrib.author_id}`}
                    </Link>
                  ) : (
                    <span className="font-mono text-sm text-ink-soft">
                      匿名贡献者 #{idx + 1}
                    </span>
                  )}
                  <StatusBadge status={contrib.status} />
                </div>
                <span className="flex items-center gap-1 text-xs text-ink-soft">
                  <Eye className="h-3 w-3" />
                  {contrib.view_count}
                </span>
              </div>

              {/* Content: preview for open, full for closed */}
              {isOpen && contrib.preview && (
                <p className="text-sm leading-relaxed text-ink-soft">
                  {contrib.preview}
                </p>
              )}
              {isClosed && contrib.content && (
                <MarkdownRenderer content={contrib.content} />
              )}

              {/* Vote button for open ideas */}
              {isOpen && isLoggedIn() && !voteSuccess && (
                <button
                  onClick={() => handleVote(contrib.id)}
                  disabled={voting !== null}
                  className="mt-3 inline-flex items-center gap-1.5 rounded-[10px] px-4 py-2 text-sm font-medium text-ink-soft hover:text-ink disabled:opacity-50"
                  style={{ border: "1px solid var(--line)" }}
                >
                  <Vote className="h-4 w-4" />
                  {voting === contrib.id ? "投票中..." : "投票"}
                </button>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

function DetailItem({ label, value }: { label: string; value: string }) {
  return (
    <div
      className="rounded-[10px] p-3"
      style={{ background: "var(--surface-muted)" }}
    >
      <p className="mb-0.5 text-xs font-semibold uppercase tracking-wider text-ink-soft">
        {label}
      </p>
      <p className="text-sm leading-relaxed">{value}</p>
    </div>
  );
}
