"use client";

import { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { Suspense } from "react";
import { setToken } from "@/lib/auth";

const API_BASE =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

function CallbackHandler() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const code = searchParams.get("code");
    if (!code) {
      setError("Missing authorization code");
      return;
    }

    fetch(`${API_BASE}/auth/callback?code=${encodeURIComponent(code)}`)
      .then((res) => {
        if (!res.ok) throw new Error("Authentication failed");
        return res.json();
      })
      .then((data: { token: string }) => {
        setToken(data.token);
        router.push("/dashboard");
      })
      .catch((err) => {
        setError(err instanceof Error ? err.message : "Authentication failed");
      });
  }, [searchParams, router]);

  if (error) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <div className="rounded-[14px] p-6 text-center" style={{ background: "var(--surface)", border: "1px solid var(--line)" }}>
          <p className="text-sm" style={{ color: "#dc2626" }}>{error}</p>
          <a href="/" className="mt-4 inline-block text-sm text-ink-soft underline decoration-accent/30 underline-offset-2 hover:text-ink">
            返回首页
          </a>
        </div>
      </div>
    );
  }

  return (
    <div className="flex min-h-[60vh] items-center justify-center">
      <p className="text-ink-soft">正在登录...</p>
    </div>
  );
}

export default function AuthCallbackPage() {
  return (
    <Suspense
      fallback={
        <div className="flex min-h-[60vh] items-center justify-center">
          <p className="text-ink-soft">Loading...</p>
        </div>
      }
    >
      <CallbackHandler />
    </Suspense>
  );
}
