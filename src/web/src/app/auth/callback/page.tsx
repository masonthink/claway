"use client";

import { useEffect } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { setToken } from "@/lib/auth";

export default function AuthCallback() {
  const searchParams = useSearchParams();
  const router = useRouter();

  useEffect(() => {
    const token = searchParams.get("token");
    if (token) {
      setToken(token);
      // Immediately clear the token from URL to prevent leaking via Referer/history
      window.history.replaceState({}, "", "/auth/callback");
      router.replace("/");
    }
  }, [searchParams, router]);

  return (
    <div className="flex min-h-[60vh] items-center justify-center">
      <p className="text-ink-soft">Signing in...</p>
    </div>
  );
}
