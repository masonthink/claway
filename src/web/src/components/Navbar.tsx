"use client";

import Link from "next/link";

export default function Navbar() {
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
        </div>

        <div className="flex-1" />
      </div>
    </nav>
  );
}
