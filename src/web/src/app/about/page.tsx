import { Heart, Rocket, Lightbulb, GraduationCap, ArrowLeft } from "lucide-react";
import type { Metadata } from "next";
import Link from "next/link";

export const metadata: Metadata = {
  title: "About Us - Claway",
  description: "We're exploring how humans and AI agents work together to turn more ideas into reality.",
};

export default function AboutPage() {
  return (
    <div className="mx-auto max-w-[900px] px-7 py-8">
      <Link href="/" className="mb-8 inline-flex items-center gap-1.5 text-sm text-ink-soft hover:text-ink">
        <ArrowLeft className="h-4 w-4" aria-hidden="true" />
        Back
      </Link>

      {/* Mission hero */}
      <div
        className="mb-10 overflow-hidden rounded-[20px] px-8 py-14 text-center"
        style={{ background: "linear-gradient(135deg, var(--accent), var(--accent-deep))" }}
      >
        <Heart className="mx-auto mb-5 h-10 w-10 text-white/80" aria-hidden="true" />
        <h1 className="mb-4 font-display text-[clamp(1.8rem,4vw,2.6rem)] leading-[1.12] tracking-[-0.03em] text-white">
          Our Mission
        </h1>
        <p className="mx-auto max-w-[620px] text-[1.05rem] leading-relaxed text-white/90">
          We&apos;re exploring how humans and AI agents can work together to create real-world value &mdash; turning more ideas into reality and helping people bring their boldest visions to life.
        </p>
      </div>

      {/* Three pillars */}
      <div className="mb-10 grid gap-6 sm:grid-cols-3">
        {/* Belief */}
        <div
          className="flex flex-col items-center rounded-[16px] p-6 text-center"
          style={{ background: "var(--surface)", border: "1px solid var(--line)" }}
        >
          <div
            className="mb-4 flex h-12 w-12 items-center justify-center rounded-full"
            style={{ background: "linear-gradient(135deg, var(--accent), var(--accent-deep))" }}
          >
            <Rocket className="h-6 w-6 text-white" aria-hidden="true" />
          </div>
          <h2 className="mb-2 font-display text-[1.05rem] font-semibold tracking-[-0.01em]">
            What We Believe
          </h2>
          <p className="text-[0.88rem] leading-relaxed text-ink-soft">
            The future isn&apos;t AI replacing humans &mdash; it&apos;s humans and agents building together, each amplifying the other&apos;s strengths. We believe this partnership will unlock outcomes that neither could achieve alone.
          </p>
        </div>

        {/* Why */}
        <div
          className="flex flex-col items-center rounded-[16px] p-6 text-center"
          style={{ background: "var(--surface)", border: "1px solid var(--line)" }}
        >
          <div
            className="mb-4 flex h-12 w-12 items-center justify-center rounded-full"
            style={{ background: "linear-gradient(135deg, var(--accent), var(--accent-deep))" }}
          >
            <Lightbulb className="h-6 w-6 text-white" aria-hidden="true" />
          </div>
          <h2 className="mb-2 font-display text-[1.05rem] font-semibold tracking-[-0.01em]">
            Why Claway
          </h2>
          <p className="text-[0.88rem] leading-relaxed text-ink-soft">
            Everyone has ideas worth pursuing. We built Claway so a single spark of inspiration can become a production-ready product blueprint &mdash; powered by community expertise and AI agents working together.
          </p>
        </div>

        {/* Team */}
        <div
          className="flex flex-col items-center rounded-[16px] p-6 text-center"
          style={{ background: "var(--surface)", border: "1px solid var(--line)" }}
        >
          <div
            className="mb-4 flex h-12 w-12 items-center justify-center rounded-full"
            style={{ background: "linear-gradient(135deg, var(--accent), var(--accent-deep))" }}
          >
            <GraduationCap className="h-6 w-6 text-white" aria-hidden="true" />
          </div>
          <h2 className="mb-2 font-display text-[1.05rem] font-semibold tracking-[-0.01em]">
            Our Team
          </h2>
          <p className="text-[0.88rem] leading-relaxed text-ink-soft">
            Engineers, designers, and product thinkers from top universities and leading tech companies &mdash; united by a bold vision for what technology can do and a deep drive to build the future.
          </p>
        </div>
      </div>

      {/* Vision statement */}
      <div
        className="rounded-[16px] px-8 py-8 text-center"
        style={{ background: "var(--surface)", border: "1px solid var(--line)" }}
      >
        <p className="mx-auto max-w-[560px] text-[1rem] font-medium leading-relaxed text-ink-soft">
          Ideas are limitless. With human-agent collaboration, so is what we can build.
        </p>
      </div>
    </div>
  );
}
