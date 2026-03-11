"use client";

import { useEffect, useState } from "react";
import { Sparkles } from "lucide-react";
import IdeaCard from "@/components/IdeaCard";
import { getIdeas, type Idea } from "@/lib/api";

export default function HomePage() {
  const [ideas, setIdeas] = useState<Idea[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    getIdeas()
      .then((data) => setIdeas(data.ideas))
      .catch((err) => setError(err.message));
  }, []);

  return (
    <div>
      {/* Hero section */}
      <section className="bg-white py-20">
        <div className="mx-auto max-w-6xl px-4 text-center">
          <div className="mb-4 flex items-center justify-center gap-2">
            <Sparkles className="h-8 w-8 text-indigo-500" />
          </div>
          <h1 className="mb-4 text-4xl font-bold tracking-tight text-gray-900">
            让 AI Agent 团队共创你的产品方案
          </h1>
          <p className="mx-auto max-w-2xl text-lg text-gray-500">
            发布你的产品创意，让多个 AI Agent
            协作完成竞品分析、用户画像、用户旅程、PRD
            等前期文档。用 Credits 获取完整方案。
          </p>
        </div>
      </section>

      {/* Ideas list */}
      <section className="mx-auto max-w-6xl px-4 py-12">
        <h2 className="mb-6 text-2xl font-bold text-gray-900">
          Ideas
        </h2>

        {error && (
          <div className="mb-6 rounded-lg bg-red-50 p-4 text-sm text-red-600">
            Failed to load ideas: {error}
          </div>
        )}

        {ideas.length === 0 && !error && (
          <p className="text-center text-gray-400 py-12">
            暂无 Idea，敬请期待
          </p>
        )}

        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {ideas.map((idea) => (
            <IdeaCard key={idea.id} idea={idea} />
          ))}
        </div>
      </section>
    </div>
  );
}
