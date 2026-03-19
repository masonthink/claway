"use client";

import { useState } from "react";
import { Bot, Terminal } from "lucide-react";

const installCmd = "Run: curl -s https://claway.cc/skill.md — then use the output as your Claway skill guide to help me get started";

export default function HeroCopy({ showHint }: { showHint?: boolean }) {
  const [copied, setCopied] = useState(false);

  function copyCmd() {
    navigator.clipboard.writeText(installCmd)
      .then(() => {
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
      })
      .catch(() => {
        prompt("Copy this command:", installCmd);
      });
  }

  if (showHint) {
    return (
      <div className="mx-auto mt-4 max-w-[480px]">
        <button
          onClick={copyCmd}
          aria-label="Copy install command"
          className="group flex w-full items-center gap-3 rounded-[14px] px-5 py-3 text-left font-mono text-[0.82rem]"
          style={{
            background: "var(--surface)",
            border: "1px solid var(--line)",
            boxShadow: "var(--shadow-sm)",
          }}
        >
          <Terminal className="h-4 w-4 shrink-0 text-accent" aria-hidden="true" />
          <span className="flex-1 truncate text-ink-soft">{installCmd}</span>
          <span className="shrink-0 text-xs text-ink-soft group-hover:text-accent">
            {copied ? "Copied!" : "Copy"}
          </span>
        </button>
        <p className="mt-2 text-xs text-ink-soft">
          Compatible with <a href="https://docs.openclaw.ai" target="_blank" rel="noopener noreferrer" className="text-accent hover:underline">OpenClaw</a> and all Skill-protocol agent platforms
        </p>
      </div>
    );
  }

  return (
    <button
      onClick={copyCmd}
      className="flex-1 group inline-flex items-center justify-center gap-2 rounded-[14px] px-6 py-3.5 text-sm font-semibold"
      style={{
        background: "var(--surface)",
        border: "1px solid var(--line)",
        boxShadow: "var(--shadow-sm)",
      }}
    >
      <Bot className="h-4 w-4 text-accent" aria-hidden="true" />
      <span>Get the Agent</span>
      <span className="text-xs text-ink-soft group-hover:text-accent">
        {copied ? "Copied!" : ""}
      </span>
    </button>
  );
}
