"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Lightbulb } from "lucide-react";
import { createIdea } from "@/lib/api";
import { isLoggedIn } from "@/lib/auth";
import { useToast } from "@/components/Toast";

export default function NewIdeaPage() {
  const router = useRouter();
  const { toast } = useToast();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [targetUserHint, setTargetUserHint] = useState("");
  const [packageType, setPackageType] = useState<"light" | "standard">("standard");
  const [initiatorCut, setInitiatorCut] = useState(20);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (!isLoggedIn()) {
      router.replace("/");
    }
  }, [router]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !description.trim()) {
      toast("error", "请填写标题和描述");
      return;
    }
    setSubmitting(true);
    try {
      const idea = await createIdea({
        title: title.trim(),
        description: description.trim(),
        target_user_hint: targetUserHint.trim(),
        package_type: packageType,
        initiator_cut_percent: initiatorCut,
      });
      toast("success", "想法创建成功");
      router.push(`/ideas/${idea.id}`);
    } catch (err) {
      toast("error", err instanceof Error ? err.message : "创建失败");
    } finally {
      setSubmitting(false);
    }
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
            <Lightbulb className="h-5 w-5 text-white" />
          </div>
          <div>
            <h1 className="font-display text-xl tracking-[-0.02em]">
              发起想法
            </h1>
            <p className="text-sm text-ink-soft">
              描述你的产品创意，社区将协作完成文档
            </p>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="flex flex-col gap-5">
          {/* Title */}
          <div>
            <label className="mb-2 block text-sm font-medium">
              标题 <span style={{ color: "var(--accent)" }}>*</span>
            </label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="用一句话描述你的产品想法"
              className="w-full rounded-[12px] px-4 py-2.5 text-sm outline-none"
              style={{
                background: "var(--surface-muted)",
                border: "1px solid var(--line)",
                color: "var(--ink)",
              }}
            />
          </div>

          {/* Description */}
          <div>
            <label className="mb-2 block text-sm font-medium">
              详细描述 <span style={{ color: "var(--accent)" }}>*</span>
            </label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="详细描述产品的目标用户、核心功能、解决的问题..."
              rows={5}
              className="w-full resize-y rounded-[12px] px-4 py-3 text-sm leading-relaxed outline-none"
              style={{
                background: "var(--surface-muted)",
                border: "1px solid var(--line)",
                color: "var(--ink)",
              }}
            />
          </div>

          {/* Target user hint */}
          <div>
            <label className="mb-2 block text-sm font-medium">
              目标用户提示
            </label>
            <input
              type="text"
              value={targetUserHint}
              onChange={(e) => setTargetUserHint(e.target.value)}
              placeholder="例如：独立开发者、小型团队、设计师..."
              className="w-full rounded-[12px] px-4 py-2.5 text-sm outline-none"
              style={{
                background: "var(--surface-muted)",
                border: "1px solid var(--line)",
                color: "var(--ink)",
              }}
            />
          </div>

          {/* Package type */}
          <div>
            <label className="mb-2 block text-sm font-medium">套餐类型</label>
            <div className="flex gap-3">
              {(["light", "standard"] as const).map((type) => (
                <button
                  key={type}
                  type="button"
                  onClick={() => setPackageType(type)}
                  className="flex-1 rounded-[12px] px-4 py-3 text-center text-sm font-medium"
                  style={{
                    background:
                      packageType === type
                        ? "rgba(255,107,74,0.1)"
                        : "var(--surface-muted)",
                    border: `1px solid ${
                      packageType === type
                        ? "var(--border-ui-active)"
                        : "var(--line)"
                    }`,
                    color:
                      packageType === type
                        ? "var(--accent-deep)"
                        : "var(--ink)",
                  }}
                >
                  {type === "light" ? "轻量版" : "标准版"}
                  <div className="mt-0.5 text-xs opacity-60">
                    {type === "light" ? "基础文档" : "完整 9 文档"}
                  </div>
                </button>
              ))}
            </div>
          </div>

          {/* Initiator cut */}
          <div>
            <label className="mb-2 flex items-center justify-between text-sm font-medium">
              <span>发起人分成比例</span>
              <span
                className="rounded-[8px] px-2 py-0.5 text-xs font-bold"
                style={{
                  background: "rgba(255,107,74,0.1)",
                  color: "var(--accent-deep)",
                }}
              >
                {initiatorCut}%
              </span>
            </label>
            <input
              type="range"
              min={10}
              max={30}
              step={1}
              value={initiatorCut}
              onChange={(e) => setInitiatorCut(Number(e.target.value))}
              className="w-full accent-[var(--accent)]"
            />
            <div className="mt-1 flex justify-between text-xs text-ink-soft">
              <span>10%</span>
              <span>30%</span>
            </div>
          </div>

          {/* Submit */}
          <button
            type="submit"
            disabled={submitting}
            className="mt-2 inline-flex items-center justify-center gap-2 rounded-[10px] px-6 py-3 text-[0.95rem] font-semibold text-white hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-60 disabled:hover:translate-y-0"
            style={{
              background: "linear-gradient(135deg, var(--accent), var(--accent-deep))",
            }}
          >
            <Lightbulb className="h-4 w-4" />
            {submitting ? "创建中..." : "发起想法"}
          </button>
        </form>
      </div>
    </div>
  );
}
