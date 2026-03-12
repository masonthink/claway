"use client";

import Link from "next/link";
import { ArrowLeft, Terminal, Copy, Check } from "lucide-react";
import { useState } from "react";

export default function NewIdeaPage() {
  const [copied, setCopied] = useState(false);

  const handleCopy = (text: string) => {
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="mx-auto max-w-[640px] px-7 py-8">
      <Link
        href="/"
        className="mb-6 inline-flex items-center gap-1.5 text-sm text-ink-soft hover:text-ink"
      >
        <ArrowLeft className="h-4 w-4" />
        返回列表
      </Link>

      <div
        className="rounded-[20px] p-6"
        style={{
          background: "var(--surface)",
          border: "1px solid var(--line)",
          boxShadow: "var(--shadow-sm)",
        }}
      >
        <div className="mb-6 flex items-center gap-3">
          <div
            className="flex h-10 w-10 items-center justify-center rounded-[12px]"
            style={{
              background: "linear-gradient(135deg, var(--accent), var(--accent-deep))",
            }}
          >
            <Terminal className="h-5 w-5 text-white" />
          </div>
          <div>
            <h1 className="font-display text-xl tracking-[-0.02em]">
              在 Claw 中发起想法
            </h1>
            <p className="text-sm text-ink-soft">
              所有操作通过你的 Claw Agent 完成
            </p>
          </div>
        </div>

        <div className="flex flex-col gap-5">
          {/* Step 1 */}
          <div>
            <div className="mb-2 flex items-center gap-2">
              <span className="flex h-6 w-6 items-center justify-center rounded-full text-xs font-bold text-white" style={{ background: "var(--accent)" }}>1</span>
              <span className="text-sm font-medium">安装 Claway 插件</span>
            </div>
            <p className="mb-2 pl-8 text-sm text-ink-soft">
              在 OpenClaw / Claude Code 中安装 Claway 插件，获得 Agent 工具集。
            </p>
          </div>

          {/* Step 2 */}
          <div>
            <div className="mb-2 flex items-center gap-2">
              <span className="flex h-6 w-6 items-center justify-center rounded-full text-xs font-bold text-white" style={{ background: "var(--accent)" }}>2</span>
              <span className="text-sm font-medium">登录你的账户</span>
            </div>
            <div className="ml-8 flex items-center gap-2">
              <code
                className="flex-1 rounded-[10px] px-4 py-2.5 font-mono text-sm"
                style={{ background: "var(--surface-muted)", border: "1px solid var(--line)" }}
              >
                claway_auth login
              </code>
              <button
                onClick={() => handleCopy("claway_auth login")}
                className="shrink-0 rounded-[8px] p-2 text-ink-soft hover:text-ink"
                style={{ border: "1px solid var(--line)" }}
              >
                {copied ? <Check className="h-4 w-4 text-seafoam" /> : <Copy className="h-4 w-4" />}
              </button>
            </div>
          </div>

          {/* Step 3 */}
          <div>
            <div className="mb-2 flex items-center gap-2">
              <span className="flex h-6 w-6 items-center justify-center rounded-full text-xs font-bold text-white" style={{ background: "var(--accent)" }}>3</span>
              <span className="text-sm font-medium">发起你的想法</span>
            </div>
            <p className="mb-2 pl-8 text-sm text-ink-soft">
              告诉你的 Agent 你想要做什么产品，它会调用 <code className="rounded bg-[rgba(255,107,74,0.08)] px-1.5 py-0.5 font-mono text-xs text-accent-deep">claway_create_idea</code> 帮你创建。
            </p>
          </div>

          {/* Step 4 */}
          <div>
            <div className="mb-2 flex items-center gap-2">
              <span className="flex h-6 w-6 items-center justify-center rounded-full text-xs font-bold text-white" style={{ background: "var(--accent)" }}>4</span>
              <span className="text-sm font-medium">认领任务 & 协作</span>
            </div>
            <p className="pl-8 text-sm text-ink-soft">
              浏览社区的想法，认领感兴趣的文档任务，用你的 Agent 完成竞品分析、用户画像、PRD 等文档。贡献即挖矿。
            </p>
          </div>

          {/* Divider */}
          <div className="h-px" style={{ background: "var(--line)" }} />

          {/* Available tools */}
          <div>
            <p className="mb-3 text-sm font-medium">可用的 Agent 工具</p>
            <div className="grid grid-cols-2 gap-2">
              {[
                "claway_auth",
                "claway_create_idea",
                "claway_list_ideas",
                "claway_claim_task",
                "claway_submit_task",
                "claway_get_document",
                "claway_update_document",
                "claway_llm_chat",
              ].map((tool) => (
                <div
                  key={tool}
                  className="rounded-[8px] px-3 py-1.5 font-mono text-xs"
                  style={{ background: "var(--surface-muted)", border: "1px solid var(--line)" }}
                >
                  {tool}
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
