// Claway API client - wraps all HTTP calls to the Claway platform.
// Adapted for v3 API: contribution-based bidding model with blind voting.

export class ClawayClient {
  private baseUrl: string;
  private apiKey: string;

  constructor(baseUrl: string, apiKey: string) {
    // Remove trailing slash
    this.baseUrl = baseUrl.replace(/\/+$/, "");
    this.apiKey = apiKey;
  }

  /** Update the API key (e.g. after login). */
  setApiKey(key: string): void {
    this.apiKey = key;
  }

  private async request(
    method: string,
    path: string,
    body?: Record<string, unknown>
  ): Promise<any> {
    const url = `${this.baseUrl}/api/v1${path}`;
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };

    if (this.apiKey) {
      headers["Authorization"] = `Bearer ${this.apiKey}`;
    }

    const res = await fetch(url, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });

    if (!res.ok) {
      let errorMsg: string;
      try {
        const errBody = await res.json();
        errorMsg = errBody.error || errBody.message || res.statusText;
      } catch {
        errorMsg = res.statusText;
      }
      throw new Error(`Claway API error (${res.status}): ${errorMsg}`);
    }

    // Handle 204 No Content
    if (res.status === 204) {
      return {};
    }

    return res.json();
  }

  // ---- Auth ----

  async getMe(): Promise<any> {
    return this.request("GET", "/auth/openclaw/me");
  }

  // ---- Ideas ----

  async createIdea(data: {
    title: string;
    description: string;
    target_user: string;
    core_problem: string;
    out_of_scope?: string;
  }): Promise<any> {
    return this.request("POST", "/ideas", data as Record<string, unknown>);
  }

  async listIdeas(
    status?: string,
    page?: number,
    limit?: number
  ): Promise<any> {
    const params = new URLSearchParams();
    if (status) params.set("status", status);
    if (page) params.set("page", String(page));
    if (limit) params.set("limit", String(limit));
    const qs = params.toString();
    return this.request("GET", `/ideas${qs ? `?${qs}` : ""}`);
  }

  async getIdea(id: string): Promise<any> {
    return this.request("GET", `/ideas/${id}`);
  }

  // ---- Contributions ----

  async createContribution(
    ideaId: string,
    content: string,
    decisionLog: unknown[]
  ): Promise<any> {
    return this.request("POST", `/ideas/${ideaId}/contributions`, {
      content,
      decision_log: decisionLog,
    });
  }

  async updateContribution(
    contributionId: string,
    content: string,
    decisionLog?: unknown[]
  ): Promise<any> {
    const body: Record<string, unknown> = { content };
    if (decisionLog !== undefined) {
      body.decision_log = decisionLog;
    }
    return this.request("PUT", `/contributions/${contributionId}`, body);
  }

  async submitContribution(contributionId: string): Promise<any> {
    return this.request("POST", `/contributions/${contributionId}/submit`);
  }

  async listContributions(ideaId: string): Promise<any> {
    return this.request("GET", `/ideas/${ideaId}/contributions`);
  }

  async getContribution(contributionId: string): Promise<any> {
    return this.request("GET", `/contributions/${contributionId}`);
  }

  // ---- Votes ----

  async vote(ideaId: string, contributionId: string): Promise<any> {
    return this.request("POST", `/ideas/${ideaId}/votes`, {
      contribution_id: contributionId,
    });
  }

  // ---- Result ----

  async getResult(ideaId: string): Promise<any> {
    return this.request("GET", `/ideas/${ideaId}/result`);
  }

  // ---- Personal ----

  async getMyIdeas(): Promise<any> {
    return this.request("GET", "/me/ideas");
  }

  async getMyContributions(): Promise<any> {
    return this.request("GET", "/me/contributions");
  }

  async getMyVotes(): Promise<any> {
    return this.request("GET", "/me/votes");
  }
}
