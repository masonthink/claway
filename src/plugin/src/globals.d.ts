// Node.js 18+ global types for fetch and URLSearchParams.
// These are available at runtime but not included in TypeScript's ES2020 lib.

declare function fetch(
  input: string | URL,
  init?: RequestInit
): Promise<Response>;

interface RequestInit {
  method?: string;
  headers?: Record<string, string>;
  body?: string;
}

interface Response {
  ok: boolean;
  status: number;
  statusText: string;
  json(): Promise<any>;
  text(): Promise<string>;
}

declare class URLSearchParams {
  constructor(init?: string | Record<string, string> | [string, string][]);
  set(name: string, value: string): void;
  get(name: string): string | null;
  toString(): string;
}
