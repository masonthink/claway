// API client for Claway v3

import { getToken } from "./auth";

// Proxy base: same-origin /api/proxy/* -> backend
const PROXY_BASE = "/api/proxy";
// Direct backend URL (for OAuth redirects only)
export const DIRECT_API_BASE =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081/api/v1";

// Simple in-memory cache for GET requests (stale-while-revalidate pattern)
const cache = new Map<string, { data: unknown; ts: number }>();
const CACHE_TTL = 30_000; // 30 seconds

function getCached<T>(key: string): T | null {
  const entry = cache.get(key);
  if (entry && Date.now() - entry.ts < CACHE_TTL) {
    return entry.data as T;
  }
  return null;
}

function setCache(key: string, data: unknown) {
  cache.set(key, { data, ts: Date.now() });
  // Evict old entries if cache grows too large
  if (cache.size > 100) {
    const oldest = cache.keys().next().value;
    if (oldest) cache.delete(oldest);
  }
}

async function request<T>(
  path: string,
  options?: RequestInit
): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options?.headers as Record<string, string>),
  };
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const method = options?.method || "GET";

  // Use cache for GET requests without auth
  if (method === "GET" && !token) {
    const cached = getCached<T>(path);
    if (cached) return cached;
  }

  const res = await fetch(`${PROXY_BASE}${path}`, {
    ...options,
    headers,
  });

  if (!res.ok) {
    const body = await res.json().catch(() => null);
    const msg = body?.error || `${res.status} ${res.statusText}`;
    throw new Error(msg);
  }

  const data: T = await res.json();

  // Cache GET responses
  if (method === "GET") {
    setCache(path, data);
  }

  return data;
}

// --- Types ---

export interface Idea {
  id: number;
  initiator_id: number;
  title: string;
  description: string;
  target_user: string;
  core_problem: string;
  out_of_scope: string | null;
  status: "open" | "closed" | "cancelled";
  deadline: string;
  revealed_at: string | null;
  created_at: string;
  // enriched fields
  contribution_count: number;
  voter_count: number;
  initiator_username: string;
}

export interface Contribution {
  id: number;
  idea_id: number;
  author_id?: number;
  author_name?: string;
  content: string;
  preview?: string;
  decision_log?: unknown[];
  status: "draft" | "submitted";
  view_count: number;
  created_at: string;
  updated_at: string;
  submitted_at?: string;
  preview_url?: string;
}

export interface RevealResultEntry {
  contribution_id: number;
  author_id: number;
  author_username: string;
  vote_count: number;
  rank: number;
  is_featured: boolean;
}

export interface RevealResult {
  idea_id: number;
  total_votes: number;
  revealed_at: string;
  results: RevealResultEntry[];
}

export interface PlatformStats {
  open_ideas: number;
  closed_ideas: number;
  total_contributions: number;
}

export interface User {
  id: number;
  username: string;
  display_name: string;
  avatar_url: string;
  created_at: string;
}

export interface UserProfile {
  user: User;
  idea_count: number;
  contribution_count: number;
  featured_count: number;
}

// --- Public API ---

export function getStats(): Promise<PlatformStats> {
  return request("/public/stats");
}

export function getIdeas(
  status?: string,
  limit?: number,
  offset?: number
): Promise<{ ideas: Idea[]; total: number }> {
  const params = new URLSearchParams();
  if (status) params.set("status", status);
  if (limit) params.set("limit", String(limit));
  if (offset) params.set("offset", String(offset));
  const qs = params.toString();
  return request(`/public/ideas${qs ? `?${qs}` : ""}`);
}

export function getIdea(id: string): Promise<Idea> {
  return request(`/public/ideas/${id}`);
}

export function getContributions(ideaId: string): Promise<Contribution[]> {
  return request(`/public/ideas/${ideaId}/contributions`);
}

export function getRevealResult(ideaId: string): Promise<RevealResult> {
  return request(`/public/ideas/${ideaId}/result`);
}

export function getUserProfile(username: string): Promise<UserProfile> {
  return request(`/public/users/${username}`);
}

// --- Auth API (read-only — all write operations go through OpenClaw Skill) ---

export function getMe(): Promise<User> {
  return request("/me");
}

export interface MyVote {
  id: number;
  idea_id: number;
  contribution_id: number;
  voted_at: string;
}

export function getMyVoteForIdea(ideaId: string): Promise<MyVote> {
  return request(`/me/votes/${ideaId}`);
}

export function getDraftPreview(contributionId: string): Promise<Contribution> {
  return request(`/draft/${contributionId}`);
}

export function getContribution(id: string): Promise<Contribution> {
  return request(`/contributions/${id}`);
}
