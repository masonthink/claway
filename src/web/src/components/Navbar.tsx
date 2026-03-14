"use client";

import Link from "next/link";
import Logo from "./Logo";

export default function Navbar() {
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
      </div>
    </nav>
  );
}
