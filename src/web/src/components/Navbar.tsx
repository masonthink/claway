"use client";

import { useState, useEffect, useRef } from "react";
import Link from "next/link";
import Logo from "./Logo";
import { DIRECT_API_BASE } from "@/lib/api";
import { isLoggedIn, removeToken, getToken } from "@/lib/auth";

// --- SVG Icons (inline, no external deps) ---

function XIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" className={className} fill="currentColor">
      <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z" />
    </svg>
  );
}

function GitHubIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" className={className} fill="currentColor">
      <path d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" />
    </svg>
  );
}

function GoogleIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" className={className}>
      <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 01-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4" />
      <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
      <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05" />
      <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
    </svg>
  );
}

function UserIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" className={className} fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round">
      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
      <circle cx="12" cy="7" r="4" />
    </svg>
  );
}

export default function Navbar() {
  const [loggedIn, setLoggedIn] = useState(false);
  const [menuOpen, setMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    setLoggedIn(isLoggedIn());
  }, []);

  // Close menu on outside click
  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setMenuOpen(false);
      }
    }
    if (menuOpen) document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, [menuOpen]);

  function handleLogout() {
    removeToken();
    setLoggedIn(false);
    setMenuOpen(false);
  }

  return (
    <nav
      className="sticky top-0 z-10 backdrop-blur-[18px]"
      style={{ background: "var(--nav-bg)", borderBottom: "1px solid var(--line)" }}
      aria-label="Main navigation"
    >
      <div className="mx-auto flex max-w-[1200px] items-center gap-6 px-7 py-4">
        <Link href="/" className="flex items-center gap-2.5">
          <Logo className="h-7 w-7" />
          <span className="font-display text-[1.35rem] font-bold tracking-[-0.03em]">
            Claway
          </span>
          <span
            className="rounded px-1 py-px text-[0.5rem] font-medium uppercase leading-none tracking-[0.04em]"
            style={{ color: "var(--accent)", border: "1px solid var(--accent)", opacity: 0.7 }}
          >
            beta
          </span>
        </Link>

        <div className="flex gap-4 text-[0.92rem] text-ink-soft">
          <Link href="/#ideas" className="hover:text-ink">All</Link>
          <Link href="/?status=open#ideas" className="hover:text-ink">Open</Link>
          <Link href="/?status=closed#ideas" className="hover:text-ink">Revealed</Link>
        </div>

        <div className="flex-1" />

        {/* Auth */}
        <div className="relative" ref={menuRef}>
          <button
            onClick={() => setMenuOpen(!menuOpen)}
            className="flex h-9 w-9 items-center justify-center rounded-full transition-colors hover:bg-white/10"
            aria-label={loggedIn ? "Account menu" : "Sign in"}
          >
            <UserIcon className="h-5 w-5" />
          </button>

          {menuOpen && (
            <div
              className="absolute right-0 top-full mt-2 w-56 rounded-xl py-2 shadow-lg"
              style={{ background: "var(--surface)", border: "1px solid var(--line)" }}
            >
              {loggedIn ? (
                <button
                  onClick={handleLogout}
                  className="w-full px-4 py-2.5 text-left text-sm hover:bg-white/5 transition-colors"
                >
                  Sign Out
                </button>
              ) : (
                <>
                  <p className="px-4 py-2 text-xs font-medium text-ink-soft uppercase tracking-wider">
                    Sign in with
                  </p>
                  <a
                    href={`${DIRECT_API_BASE}/auth/x`}
                    className="flex w-full items-center gap-3 px-4 py-2.5 text-sm hover:bg-white/5 transition-colors"
                  >
                    <XIcon className="h-4 w-4" />
                    <span>X (Twitter)</span>
                  </a>
                  <a
                    href={`${DIRECT_API_BASE}/auth/github`}
                    className="flex w-full items-center gap-3 px-4 py-2.5 text-sm hover:bg-white/5 transition-colors"
                  >
                    <GitHubIcon className="h-4 w-4" />
                    <span>GitHub</span>
                  </a>
                  <a
                    href={`${DIRECT_API_BASE}/auth/google`}
                    className="flex w-full items-center gap-3 px-4 py-2.5 text-sm hover:bg-white/5 transition-colors"
                  >
                    <GoogleIcon className="h-4 w-4" />
                    <span>Google</span>
                  </a>
                </>
              )}
            </div>
          )}
        </div>
      </div>
    </nav>
  );
}
