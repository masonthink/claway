"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { LayoutDashboard, LogIn, LogOut } from "lucide-react";
import { isLoggedIn, removeToken } from "@/lib/auth";

export default function Navbar() {
  const [loggedIn, setLoggedIn] = useState(false);

  useEffect(() => {
    setLoggedIn(isLoggedIn());
  }, []);

  const handleLogout = () => {
    removeToken();
    setLoggedIn(false);
    window.location.href = "/";
  };

  return (
    <nav
      className="sticky top-0 z-10 backdrop-blur-[18px]"
      style={{ background: "var(--nav-bg)", borderBottom: "1px solid var(--line)" }}
    >
      <div className="mx-auto flex max-w-[1200px] items-center gap-6 px-7 py-4">
        <Link href="/" className="font-display text-[1.35rem] font-bold tracking-[-0.03em]">
          Claway
        </Link>

        <div className="flex gap-4 text-[0.92rem] text-ink-soft">
          <Link href="/" className="hover:text-ink">Ideas</Link>
          {loggedIn && (
            <Link href="/dashboard" className="flex items-center gap-1.5 hover:text-ink">
              <LayoutDashboard className="h-3.5 w-3.5" />
              Dashboard
            </Link>
          )}
        </div>

        <div className="flex-1" />

        {loggedIn ? (
          <button
            onClick={handleLogout}
            className="inline-flex items-center gap-2 rounded-[10px] px-3.5 py-2 text-sm font-medium text-ink-soft hover:text-ink"
            style={{ border: "1px solid var(--line)" }}
          >
            <LogOut className="h-3.5 w-3.5" />
            Logout
          </button>
        ) : (
          <a
            href={`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081/api/v1"}/auth/x`}
            className="inline-flex items-center gap-2 rounded-[10px] px-4 py-2 text-sm font-semibold text-white hover:-translate-y-0.5"
            style={{ background: "linear-gradient(135deg, var(--accent), var(--accent-deep))" }}
          >
            <LogIn className="h-3.5 w-3.5" />
            Sign in with X
          </a>
        )}
      </div>
    </nav>
  );
}
