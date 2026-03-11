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
      getMyContributions().then((d) => setContributions(d.contributions)),
      getMyCompute().then(setCompute),
    ]).catch((err) => setError(err.message));
  }, [router]);

  if (!isLoggedIn()) return null;

  if (error) {
    return (
      <div className="mx-auto max-w-4xl px-4 py-12">
        <div className="rounded-lg bg-red-50 p-4 text-sm text-red-600">
          {error}
        </div>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-4xl px-4 py-8">
      <div className="mb-8 flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">
          {user ? `${user.username} 的 Dashboard` : "Dashboard"}
        </h1>
        {credits && <CreditBalance balance={credits.balance} />}
      </div>

      {/* Credits & Transactions */}
      {credits && credits.transactions.length > 0 && (
        <section className="mb-8">
          <h2 className="mb-4 flex items-center gap-2 text-lg font-semibold text-gray-900">
            <Coins className="h-5 w-5 text-indigo-500" />
            最近交易
          </h2>
          <div className="rounded-xl border border-gray-200 bg-white">
            {credits.transactions.slice(0, 10).map((tx) => (
              <div
                key={tx.id}
                className="flex items-center justify-between border-b border-gray-100 px-5 py-3 last:border-b-0"
              >
                <div>
                  <p className="text-sm font-medium text-gray-800">
                    {tx.description}
                  </p>
                  <p className="text-xs text-gray-400">
                    {new Date(tx.created_at).toLocaleDateString("zh-CN")}
                  </p>
                </div>
                <span
                  className={`text-sm font-semibold ${
                    tx.amount > 0 ? "text-green-600" : "text-red-500"
                  }`}
                >
                  {tx.amount > 0 ? "+" : ""}
                  {tx.amount}
                </span>
              </div>
            ))}
          </div>
        </section>
      )}

      {/* My Contributions */}
      {contributions.length > 0 && (
        <section className="mb-8">
          <h2 className="mb-4 flex items-center gap-2 text-lg font-semibold text-gray-900">
            <FileText className="h-5 w-5 text-indigo-500" />
            我的贡献
          </h2>
          <div className="rounded-xl border border-gray-200 bg-white">
            {contributions.map((c) => (
              <Link
                key={`${c.idea_id}-${c.task_code}`}
                href={`/ideas/${c.idea_id}`}
                className="flex items-center justify-between border-b border-gray-100 px-5 py-3 last:border-b-0 hover:bg-gray-50 transition-colors"
              >
                <div>
                  <p className="text-sm font-medium text-gray-800">
                    {c.idea_title}
                  </p>
                  <p className="text-xs text-gray-400">
                    {c.task_code} - {c.task_name}
                  </p>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-xs text-gray-400">
                    {c.token_cost.toFixed(2)} tokens
                  </span>
                  <StatusBadge status={c.status} />
                </div>
              </Link>
            ))}
          </div>
        </section>
      )}

      {/* Compute Usage */}
      {compute && (
        <section className="mb-8">
          <h2 className="mb-4 flex items-center gap-2 text-lg font-semibold text-gray-900">
            <Cpu className="h-5 w-5 text-indigo-500" />
            Compute 使用量
          </h2>
          <div className="mb-4 rounded-lg bg-indigo-50 px-5 py-3">
            <span className="text-sm text-gray-600">总消耗: </span>
            <span className="text-lg font-bold text-indigo-700">
              {compute.total_cost.toFixed(2)} tokens
            </span>
          </div>
          {compute.breakdown.length > 0 && (
            <div className="rounded-xl border border-gray-200 bg-white">
              {compute.breakdown.map((item) => (
                <Link
                  key={item.idea_id}
                  href={`/ideas/${item.idea_id}`}
                  className="flex items-center justify-between border-b border-gray-100 px-5 py-3 last:border-b-0 hover:bg-gray-50 transition-colors"
                >
                  <span className="text-sm font-medium text-gray-800">
                    {item.idea_title}
                  </span>
                  <span className="text-sm text-gray-500">
                    {item.cost.toFixed(2)} tokens
                  </span>
                </Link>
              ))}
            </div>
          )}
        </section>
      )}
    </div>
  );
}
