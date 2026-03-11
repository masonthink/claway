// ClawBeach API client - wraps all HTTP calls to the ClawBeach platform.

export class ClawBeachClient {
  private baseUrl: string;
  private apiKey: string;

  constructor(baseUrl: string, apiKey: string) {
    // Remove trailing slash
    this.baseUrl = baseUrl.replace(/\/+$/, "");
    this.apiKey = apiKey;
  }

  private async request(
    method: string,
    path: string,
    body?: Record<string, unknown>
  ): Promise<any> {
    const url = `${this.baseUrl}/api/v1${path}`;
    const headers: Record<string, string> = {
      Authorization: `Bearer ${this.apiKey}`,
      "Content-Type": "application/json",
    };

    const res = await fetch(url, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });

    if (!res.ok) {
      let errorMsg: string;
      try {
        const errBody = await res.json();
        errorMsg = errBody.error || res.statusText;
      } catch {
        errorMsg = res.statusText;
      }
      throw new Error(`ClawBeach API error (${res.status}): ${errorMsg}`);
    }

    return res.json();
  }

  // ---- Auth ----

  async getMe(): Promise<any> {
    return this.request("GET", "/auth/me");
  }

  // ---- Ideas ----

  async listIdeas(status?: string, limit?: number, offset?: number): Promise<any> {
    const params = new URLSearchParams();
    if (status) params.set("status", status);
    if (limit) params.set("limit", String(limit));
    if (offset) params.set("offset", String(offset));
    const qs = params.toString();
    return this.request("GET", `/ideas${qs ? `?${qs}` : ""}`);
  }

  async getIdea(id: number): Promise<any> {
    return this.request("GET", `/ideas/${id}`);
  }

  async createIdea(data: {
    title: string;
    description: string;
    target_user_hint?: string;
    problem_definition?: string;
    initiator_cut_percent: number;
    package_type: string;
  }): Promise<any> {
    return this.request("POST", "/ideas", data as Record<string, unknown>);
  }

  async getIdeaContext(id: number): Promise<any> {
    return this.request("GET", `/ideas/${id}/context`);
  }

  async getIdeaTasks(id: number): Promise<any> {
    return this.request("GET", `/ideas/${id}/tasks`);
  }

  // ---- Tasks ----

  async getTask(id: number): Promise<any> {
    return this.request("GET", `/tasks/${id}`);
  }

  async claimTask(id: number): Promise<any> {
    return this.request("POST", `/tasks/${id}/claim`);
  }

  async unclaimTask(id: number): Promise<any> {
    return this.request("DELETE", `/tasks/${id}/claim`);
  }

  async submitTask(id: number, content: string, note: string): Promise<any> {
    return this.request("POST", `/tasks/${id}/submit`, { content, note });
  }

  async reviewTask(
    id: number,
    action: string,
    qualityScore?: number,
    rejectReason?: string
  ): Promise<any> {
    return this.request("POST", `/tasks/${id}/review`, {
      action,
      quality_score: qualityScore,
      reject_reason: rejectReason,
    });
  }

  // ---- Documents ----

  async getDocument(taskId: number): Promise<any> {
    return this.request("GET", `/tasks/${taskId}/document`);
  }

  async updateDocument(taskId: number, content: string): Promise<any> {
    return this.request("PUT", `/tasks/${taskId}/document`, { content });
  }

  // ---- Compute ----

  async getMyCompute(): Promise<any> {
    return this.request("GET", "/me/compute");
  }

  // ---- Credits ----

  async getMyCredits(): Promise<any> {
    return this.request("GET", "/me/credits");
  }

  async getMyContributions(): Promise<any> {
    return this.request("GET", "/me/contributions");
  }

  // ---- PRD ----

  async publishPRD(ideaId: number): Promise<any> {
    return this.request("POST", `/ideas/${ideaId}/publish`);
  }

  async purchasePRD(prdId: number): Promise<any> {
    return this.request("POST", `/prd/${prdId}/purchase`);
  }
}
