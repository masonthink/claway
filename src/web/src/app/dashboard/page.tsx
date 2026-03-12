"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Cpu, FileText, Coins } from "lucide-react";
import CreditBalance from "@/components/CreditBalance";
import StatusBadge from "@/components/StatusBadge";
import { isLoggedIn } from "@/lib/auth";
import {
  getMe,
  getMyCredits,
  getMyContributions,
  getMyCompute,
  type User,
  type CreditInfo,
  type Contribution,
  type ComputeUsage,
} from "@/lib/api";

export default function DashboardPage() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [credits, setCredits] = useState<CreditInfo | null>(null);
  const [contributions, setContributions] = useState<Contribution[]>([]);
  const [compute, setCompute] = useState<ComputeUsage | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isLoggedIn()) {
      router.push("/");
      return;
    }
    Promise.all([
      getMe().then(setUser),
      getMyCredits().then(setCredits),
      getMyContributions().then((d) => setContributions(d.contributions || [])),
      getMyCompute().then(setCompute),
    ]).catch((err) => setError(err.message));
  }, [router]);

  if (!isLoggedIn()) return null;

  if (error) {
    return (
      <div className="mx-auto max-w-[860px] px-7 py-12">
        <div className="rounded-[12px] p-4 text-sm" style={{ background: "rgba(239,68,68,0.08)", color: "#dc2626", border: "1px solid rgba(239,68,68,0.15)" }}>
          {error}
        </div>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-[860px] px-7 py-8">
      <div className="mb-8 flex items-center justify-between">
        <h1 className="font-display text-2xl tracking-[-0.02em]">
          {user ? `${user.username} 的 Dashboard` : "Dashboard"}
        </h1>
        {credits && <CreditBalance balance={credits.balance} />}
      </div>

      {/* Transactions */}
      {credits && credits.transactions && credits.transactions.length > 0 && (
        <section className="mb-8">
          <h2 className="mb-4 flex items-center gap-2 font-display text-lg tracking-[-0.02em]">
            <Coins className="h-4.5 w-4.5 text-gold" />
            最近交易
          </h2>
          <div className="overflow-hidden rounded-[16px]" style={{ border: "1px solid var(--line)", background: "var(--surface)" }}>
            {credits.transactions.slice(0, 10).map((tx, i) => (
              <div
                key={tx.id}
                className="flex items-center justify-between px-5 py-3"
                style={{ borderBottom: i === Math.min(credits.transactions.length, 10) - 1 ? "none" : "1px solid var(--line)" }}
              >
                <div>
                  <p className="text-sm font-medium">{tx.description}</p>
                  <p className="text-xs text-ink-soft">{new Date(tx.created_at).toLocaleDateString("zh-CN")}</p>
                </div>
                <span className={`text-sm font-semibold ${tx.amount > 0 ? "text-seafoam" : "text-accent"}`}>
                  {tx.amount > 0 ? "+" : ""}{tx.amount}
                </span>
              </div>
            ))}
          </div>
        </section>
      )}

      {/* Contributions */}
      {contributions.length > 0 && (
        <section className="mb-8">
          <h2 className="mb-4 flex items-center gap-2 font-display text-lg tracking-[-0.02em]">
            <FileText className="h-4.5 w-4.5 text-accent" />
            我的贡献
          </h2>
          <div className="overflow-hidden rounded-[16px]" style={{ border: "1px solid var(--line)", background: "var(--surface)" }}>
            {contributions.map((c, i) => (
              <Link
                key={`${c.idea_id}-${c.task_code}`}
                href={`/ideas/${c.idea_id}`}
                className="flex items-center justify-between px-5 py-3 hover:bg-[rgba(255,107,74,0.05)]"
                style={{ borderBottom: i === contributions.length - 1 ? "none" : "1px solid var(--line)" }}
              >
                <div className="min-w-0">
                  <p className="truncate text-sm font-medium">{c.idea_title}</p>
                  <p className="text-xs text-ink-soft">{c.task_code} - {c.task_name}</p>
                </div>
                <div className="flex shrink-0 items-center gap-2.5">
                  <span className="text-xs text-ink-soft">{c.token_cost.toFixed(2)} tokens</span>
                  <StatusBadge status={c.status} />
                </div>
              </Link>
            ))}
          </div>
        </section>
      )}

      {/* Compute */}
      {compute && (
        <section className="mb-8">
          <h2 className="mb-4 flex items-center gap-2 font-display text-lg tracking-[-0.02em]">
            <Cpu className="h-4.5 w-4.5 text-accent" />
            Compute 使用量
          </h2>
          <div className="mb-4 rounded-[12px] px-5 py-3" style={{ background: "rgba(255,107,74,0.08)" }}>
            <span className="text-sm text-ink-soft">总消耗: </span>
            <span className="text-lg font-bold text-accent-deep">{(compute.total_cost || 0).toFixed(2)} tokens</span>
          </div>
          {compute.breakdown && compute.breakdown.length > 0 && (
            <div className="overflow-hidden rounded-[16px]" style={{ border: "1px solid var(--line)", background: "var(--surface)" }}>
              {compute.breakdown.map((item, i) => (
                <Link
                  key={item.idea_id}
                  href={`/ideas/${item.idea_id}`}
                  className="flex items-center justify-between px-5 py-3 hover:bg-[rgba(255,107,74,0.05)]"
                  style={{ borderBottom: i === compute.breakdown.length - 1 ? "none" : "1px solid var(--line)" }}
                >
                  <span className="text-sm font-medium">{item.idea_title}</span>
                  <span className="text-sm text-ink-soft">{item.cost.toFixed(2)} tokens</span>
                </Link>
              ))}
            </div>
          )}
        </section>
      )}
    </div>
  );
}
