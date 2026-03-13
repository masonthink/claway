"use client";

import Link from "next/link";
import { Clock, Users, Vote } from "lucide-react";
import StatusBadge from "./StatusBadge";
import type { Idea } from "@/lib/api";
import { timeLeft } from "@/lib/utils";

export default function IdeaCard({ idea }: { idea: Idea }) {
  return (
    <Link
      href={`/idea/${idea.id}`}
      className="card-hover flex flex-col rounded-[16px] p-5"
      style={{
        background: "var(--surface)",
        border: "1px solid var(--line)",
        boxShadow: "var(--shadow-sm)",
      }}
    >
      <div className="mb-2 flex items-start justify-between gap-3">
        <h3 className="font-display text-[1.05rem] leading-snug tracking-[-0.01em]">
          {idea.title}
        </h3>
        <StatusBadge status={idea.status} />
      </div>

      <p className="mb-4 line-clamp-2 flex-1 text-[0.88rem] leading-[1.5] text-ink-soft">
        {idea.description}
      </p>

      <div className="flex flex-wrap items-center gap-3 text-xs text-ink-soft">
        <span>@{idea.initiator_username}</span>
        <span className="flex items-center gap-1">
          <Users className="h-3 w-3" />
          {idea.contribution_count} 贡献
        </span>
        <span className="flex items-center gap-1">
          <Vote className="h-3 w-3" />
          {idea.voter_count} 投票
        </span>
        {idea.status === "open" && (
          <span className="ml-auto flex items-center gap-1 text-accent">
            <Clock className="h-3 w-3" />
            {timeLeft(idea.deadline)}
          </span>
        )}
      </div>
    </Link>
  );
}
