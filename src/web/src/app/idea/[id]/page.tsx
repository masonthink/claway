"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Clock, Users, Vote, Eye, Trophy } from "lucide-react";
import StatusBadge from "@/components/StatusBadge";
import MarkdownRenderer from "@/components/MarkdownRenderer";
import ErrorState from "@/components/ErrorState";
import {
  getIdea,
  getContributions,
  castVote,
  type Idea,
  type Contribution,
} from "@/lib/api";
import { isLoggedIn } from "@/lib/auth";
import { timeLeft } from "@/lib/utils";

export default function IdeaDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [idea, setIdea] = useState<Idea | null>(null);
  const [contributions, setContributions] = useState<Contribution[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [voting, setVoting] = useState<number | null>(null);
  const [voteError, setVoteError] = useState<string | null>(null);
  const [voteSuccess, setVoteSuccess] = useState(false);
  const [votedContribId, setVotedContribId] = useState<number | null>(null);

  const loadData = () => {
    if (!id) return;
    setError(null);
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
    if (!confirm("确认投票？每个想法只能投一票，投票后不可更改。")) return;

    setVoting(contributionId);
    setVoteError(null);
    try {
      await castVote(id, contributionId);
      setVoteSuccess(true);
      setVotedContribId(contributionId);
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
        <ErrorState message={error} onRetry={loadData} />
      </div>
    );
  }

  if (!idea) {
    return null; // loading.tsx handles this
  }

  const isOpen = idea.status === "open";
  const isClosed = idea.status === "closed";

  return (
    <div className="mx-auto max-w-[860px] px-7 py-8">
      <Link href="/" className="mb-6 inline-flex items-center gap-1.5 text-sm text-ink-soft hover:text-ink">
        <ArrowLeft className="h-4 w-4" aria-hidden="true" />
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
            <Users className="h-4 w-4" aria-hidden="true" />
            {idea.contribution_count} 贡献
          </span>
          <span className="flex items-center gap-1">
            <Vote className="h-4 w-4" aria-hidden="true" />
            {idea.voter_count} 投票人
          </span>
          {isOpen && (
            <span className="flex items-center gap-1 text-accent">
              <Clock className="h-4 w-4" aria-hidden="true" />
              {timeLeft(idea.deadline)}
            </span>
          )}
          {isClosed && (
            <Link
              href={`/idea/${idea.id}/result`}
              className="ml-auto inline-flex items-center gap-1.5 rounded-[10px] px-4 py-2 text-sm font-semibold text-white hover:-translate-y-0.5"
              style={{ background: "linear-gradient(135deg, var(--accent), var(--accent-deep))" }}
            >
              <Trophy className="h-4 w-4" aria-hidden="true" />
              查看揭榜结果
            </Link>
          )}
        </div>
      </div>

      {/* Vote feedback */}
      {voteError && (
        <div
          className="mb-6 rounded-[12px] p-4 text-sm"
          role="alert"
          style={{ background: "rgba(239,68,68,0.08)", color: "#dc2626", border: "1px solid rgba(239,68,68,0.15)" }}
        >
          {voteError}
        </div>
      )}
      {voteSuccess && (
        <div
          className="mb-6 rounded-[12px] p-4 text-sm"
          role="status"
          style={{ background: "rgba(43,198,164,0.1)", color: "rgb(26,107,91)", border: "1px solid rgba(43,198,164,0.2)" }}
        >
          投票成功！感谢你的参与。
        </div>
      )}

      {/* Contributions */}
      <div>
        <h2 className="mb-4 font-display text-lg tracking-[-0.02em]">
          贡献 ({contributions.length})
        </h2>

        {/* Blind voting hint */}
        {isOpen && contributions.length > 0 && (
          <p
            className="mb-4 rounded-[10px] p-3 text-xs text-ink-soft"
            style={{ background: "var(--surface-muted)" }}
          >
            盲投期间仅展示方案摘要，完整内容将在揭榜后公开。每个想法只能投一票，请仔细阅读后投票。
          </p>
        )}

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
              style={{
                background: "var(--surface)",
                border: votedContribId === contrib.id
                  ? "2px solid var(--accent)"
                  : "1px solid var(--line)",
              }}
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
                  {votedContribId === contrib.id && (
                    <span className="text-xs font-medium text-accent">已投票</span>
                  )}
                </div>
                <span className="flex items-center gap-1 text-xs text-ink-soft">
                  <Eye className="h-3 w-3" aria-hidden="true" />
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
              {isOpen && !voteSuccess && (
                <button
                  onClick={() => {
                    if (!isLoggedIn()) {
                      setVoteError("请先登录后再投票");
                      return;
                    }
                    handleVote(contrib.id);
                  }}
                  disabled={voting !== null}
                  className="mt-3 inline-flex items-center gap-1.5 rounded-[10px] px-4 py-2 text-sm font-medium text-ink-soft hover:text-ink disabled:opacity-50"
                  style={{ border: "1px solid var(--line)" }}
                >
                  <Vote className="h-4 w-4" aria-hidden="true" />
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
