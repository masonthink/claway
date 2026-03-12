"use client";

import { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { Suspense } from "react";
import { setToken } from "@/lib/auth";

function CallbackHandler() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // New flow: backend redirects here with ?token=JWT
    const token = searchParams.get("token");
    if (token) {
      setToken(token);
      router.push("/dashboard");
      return;
    }

    // Error from OAuth
    const errorMsg = searchParams.get("error");
    if (errorMsg) {
      setError(errorMsg);
      return;
    }

    setError("Missing authentication token");
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
