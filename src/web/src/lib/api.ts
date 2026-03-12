// API client — routes through Next.js proxy to avoid slow direct connections

import { getToken } from "./auth";

// Proxy base: same-origin /api/proxy/* → backend
const PROXY_BASE = "/api/proxy";
// Direct backend URL (for OAuth redirects only)
export const DIRECT_API_BASE =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081/api/v1";

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

  const res = await fetch(`${PROXY_BASE}${path}`, {
    ...options,
    headers,
  });

  if (!res.ok) {
    throw new Error(`API error: ${res.status} ${res.statusText}`);
  }

  return res.json();
}

// Types

export interface Idea {
  id: string;
  title: string;
  description: string;
  status: "active" | "completed";
  package_type: string;
  initiator: string;
  tasks_completed: number;
  tasks_total: number;
  total_compute_cost: number;
  created_at: string;
}

export interface Task {
  id: string;
  idea_id: string;
  name: string;
  code: string; // D1-D9
  status: "open" | "claimed" | "submitted" | "approved" | "rejected";
  claimed_by?: string;
  token_cost: number;
}

export interface PRD {
  id: string;
  idea_id: string;
  title: string;
  content: string;
  preview: string;
  price: number;
  purchased: boolean;
}

export interface User {
  id: string;
  username: string;
  avatar_url?: string;
}

export interface CreditInfo {
  balance: number;
  transactions: {
    id: string;
    amount: number;
    type: string;
    description: string;
    created_at: string;
  }[];
}

export interface Contribution {
  idea_id: string;
  idea_title: string;
  task_code: string;
  task_name: string;
  status: string;
  token_cost: number;
}

export interface ComputeUsage {
  total_cost: number;
  breakdown: {
    idea_id: string;
    idea_title: string;
    cost: number;
  }[];
}

export interface ComputeLeaderboard {
  entries: {
    username: string;
    total_cost: number;
  }[];
}

// API functions

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
  return request(`/ideas${qs ? `?${qs}` : ""}`);
}

export function getIdea(id: string): Promise<Idea> {
  return request(`/ideas/${id}`);
}

export function getIdeaTasks(ideaId: string): Promise<{ tasks: Task[] }> {
  return request(`/ideas/${ideaId}/tasks`);
}

export function getIdeaContext(ideaId: string): Promise<{ context: string }> {
  return request(`/ideas/${ideaId}/context`);
}

export function getIdeaCompute(
  ideaId: string
): Promise<ComputeLeaderboard> {
  return request(`/ideas/${ideaId}/compute`);
}

export function getPRD(id: string): Promise<PRD> {
  return request(`/prd/${id}`);
}

export function purchasePRD(id: string): Promise<{ success: boolean }> {
  return request(`/prd/${id}/purchase`, { method: "POST" });
}

export function getMe(): Promise<User> {
  return request("/me");
}

export function getMyCredits(): Promise<CreditInfo> {
  return request("/me/credits");
}

export function getMyContributions(): Promise<{
  contributions: Contribution[];
}> {
  return request("/me/contributions");
}

export function getMyCompute(): Promise<ComputeUsage> {
  return request("/me/compute");
}

export function createIdea(data: {
  title: string;
  description: string;
  target_user_hint: string;
  package_type: string;
  initiator_cut_percent: number;
}): Promise<Idea> {
  return request("/ideas", { method: "POST", body: JSON.stringify(data) });
}

export function claimTask(taskId: string): Promise<{ success: boolean }> {
  return request(`/tasks/${taskId}/claim`, { method: "POST" });
}

export function unclaimTask(taskId: string): Promise<{ success: boolean }> {
  return request(`/tasks/${taskId}/claim`, { method: "DELETE" });
}

export function submitTask(
  taskId: string,
  data: { output_content: string; output_note: string }
): Promise<{ success: boolean }> {
  return request(`/tasks/${taskId}/submit`, {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export function reviewTask(
  taskId: string,
  data: { quality_score: number; reject_reason?: string }
): Promise<{ success: boolean }> {
  return request(`/tasks/${taskId}/review`, {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export function getTaskDocument(
  taskId: string
): Promise<{ content: string }> {
  return request(`/tasks/${taskId}/document`);
}

export function updateTaskDocument(
  taskId: string,
  content: string
): Promise<{ success: boolean }> {
  return request(`/tasks/${taskId}/document`, {
    method: "PUT",
    body: JSON.stringify({ content }),
  });
}

export function publishIdea(ideaId: string): Promise<{ success: boolean }> {
  return request(`/ideas/${ideaId}/publish`, { method: "POST" });
}
