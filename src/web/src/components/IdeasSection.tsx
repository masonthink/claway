"use client";

import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import IdeaCard from "@/components/IdeaCard";
import Pagination from "@/components/Pagination";
import ErrorState from "@/components/ErrorState";
import { getIdeas, type Idea } from "@/lib/api";

const PAGE_SIZE = 12;

export default function IdeasSection({
  initialIdeas,
  initialTotal,
}: {
  initialIdeas: Idea[];
  initialTotal: number;
}) {
  const searchParams = useSearchParams();
  const statusFilter = searchParams.get("status") || undefined;

  const [ideas, setIdeas] = useState<Idea[]>(initialIdeas);
  const [total, setTotal] = useState(initialTotal);
  const [offset, setOffset] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  // When filter or offset changes, fetch new data (skip initial load)
  const [initialized, setInitialized] = useState(false);

  useEffect(() => {
    setOffset(0);
  }, [statusFilter]);

  useEffect(() => {
    // Skip first render — we already have server-fetched data
    if (!initialized) {
      setInitialized(true);
      // But if there's a status filter on first load, we need to fetch
      if (!statusFilter) return;
    }

    setLoading(true);
    setError(null);
    getIdeas(statusFilter, PAGE_SIZE, offset)
      .then((data) => {
        setIdeas(data.ideas || []);
        setTotal(data.total || 0);
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [offset, statusFilter]);

  return (
    <section id="ideas" className="px-7 pb-20">
      <div className="mx-auto max-w-[1200px]">
        <h2 className="mb-1.5 font-display text-xl tracking-[-0.02em]">
          Ideas
        </h2>
        <div className="mb-8 flex items-center justify-between">
          <p className="text-sm text-ink-soft">
            {statusFilter === "open"
              ? "Open ideas — contribute and vote"
              : statusFilter === "closed"
                ? "Revealed ideas — see community picks"
                : "Browse ideas, contribute proposals, and vote"}
          </p>
        </div>

        {error && (
          <div className="mb-6">
            <ErrorState
              message="Something went wrong. Please try again later."
              onRetry={() => {
                setLoading(true);
                setError(null);
                getIdeas(statusFilter, PAGE_SIZE, offset)
                  .then((data) => {
                    setIdeas(data.ideas || []);
                    setTotal(data.total || 0);
                  })
                  .catch((err) => setError(err.message))
                  .finally(() => setLoading(false));
              }}
            />
          </div>
        )}

        {loading && (
          <div className="flex justify-center py-20" role="status" aria-label="Loading">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-accent/20 border-t-accent" />
            <span className="sr-only">Loading</span>
          </div>
        )}

        {!loading && ideas.length === 0 && !error && (
          <p className="py-20 text-center text-ink-soft opacity-50">
            No ideas yet — stay tuned
          </p>
        )}

        {!loading && (
          <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
            {ideas.map((idea) => (
              <IdeaCard key={idea.id} idea={idea} />
            ))}
          </div>
        )}

        <Pagination
          total={total}
          limit={PAGE_SIZE}
          offset={offset}
          onChange={setOffset}
        />
      </div>
    </section>
  );
}
