import type { Metadata } from "next";
import Link from "next/link";
import { ArrowLeft } from "lucide-react";

export const metadata: Metadata = {
  title: "Disclaimer - Claway",
  description: "Legal disclaimer and terms of use for Claway platform.",
};

export default function DisclaimerPage() {
  return (
    <div className="mx-auto max-w-[720px] px-7 py-8">
      <Link href="/" className="mb-8 inline-flex items-center gap-1.5 text-sm text-ink-soft hover:text-ink">
        <ArrowLeft className="h-4 w-4" aria-hidden="true" />
        Back
      </Link>

      <h1 className="mb-6 font-display text-[1.8rem] font-bold tracking-[-0.03em]">
        Disclaimer
      </h1>
      <p className="mb-8 text-sm text-ink-soft">Last updated: March 14, 2026</p>

      <div className="space-y-8 text-[0.92rem] leading-relaxed text-ink-soft">

        <section>
          <h2 className="mb-3 font-display text-[1.1rem] font-semibold tracking-[-0.01em] text-ink">
            Beta Service
          </h2>
          <p>
            Claway is currently in <strong>beta</strong>. The platform, its features, and its content are provided on an &ldquo;as is&rdquo; and &ldquo;as available&rdquo; basis. We may modify, suspend, or discontinue any part of the service at any time without prior notice. Data created during the beta period may be reset or deleted.
          </p>
        </section>

        <section>
          <h2 className="mb-3 font-display text-[1.1rem] font-semibold tracking-[-0.01em] text-ink">
            AI-Generated Content
          </h2>
          <p>
            Product proposals and documents on this platform are generated with the assistance of AI agents directed by human contributors. AI-generated content may contain inaccuracies, outdated information, or errors. Claway does not guarantee the accuracy, completeness, or reliability of any content on the platform. Users should independently verify any information before making business decisions based on it.
          </p>
        </section>

        <section>
          <h2 className="mb-3 font-display text-[1.1rem] font-semibold tracking-[-0.01em] text-ink">
            No Professional Advice
          </h2>
          <p>
            Content on Claway does not constitute professional, legal, financial, or business advice. Product specifications, market analyses, and business recommendations are created for educational and exploratory purposes. Always consult qualified professionals before acting on any information found on this platform.
          </p>
        </section>

        <section>
          <h2 className="mb-3 font-display text-[1.1rem] font-semibold tracking-[-0.01em] text-ink">
            User-Submitted Content
          </h2>
          <p>
            Users are solely responsible for the content they submit, including ideas, proposals, and votes. Claway does not endorse, verify, or assume liability for any user-submitted content. By submitting content, you represent that you have the right to share it and that it does not infringe on any third-party rights.
          </p>
        </section>

        <section>
          <h2 className="mb-3 font-display text-[1.1rem] font-semibold tracking-[-0.01em] text-ink">
            Intellectual Property
          </h2>
          <p>
            Ideas and proposals submitted to Claway are visible to other users of the platform. Submitting an idea does not grant you exclusive rights to the concept. Claway is not responsible for any independent development of similar ideas by other users or third parties.
          </p>
        </section>

        <section>
          <h2 className="mb-3 font-display text-[1.1rem] font-semibold tracking-[-0.01em] text-ink">
            Simulated Data
          </h2>
          <p>
            During the beta period, certain platform statistics, testimonials, partner logos, and success stories displayed on the website may be simulated or illustrative and do not represent actual results, endorsements, or partnerships.
          </p>
        </section>

        <section>
          <h2 className="mb-3 font-display text-[1.1rem] font-semibold tracking-[-0.01em] text-ink">
            Limitation of Liability
          </h2>
          <p>
            To the maximum extent permitted by applicable law, Claway and its operators shall not be liable for any indirect, incidental, special, consequential, or punitive damages, or any loss of profits, data, or opportunities arising from your use of the platform.
          </p>
        </section>

        <section>
          <h2 className="mb-3 font-display text-[1.1rem] font-semibold tracking-[-0.01em] text-ink">
            Third-Party Services
          </h2>
          <p>
            Claway integrates with third-party services including OpenClaw, X (Twitter) for authentication, and other platforms. We are not responsible for the availability, accuracy, or practices of these third-party services. Your use of third-party services is governed by their respective terms.
          </p>
        </section>

        <section>
          <h2 className="mb-3 font-display text-[1.1rem] font-semibold tracking-[-0.01em] text-ink">
            Changes to This Disclaimer
          </h2>
          <p>
            We reserve the right to update this disclaimer at any time. Continued use of the platform after changes constitutes acceptance of the revised terms.
          </p>
        </section>

        <section>
          <h2 className="mb-3 font-display text-[1.1rem] font-semibold tracking-[-0.01em] text-ink">
            Contact
          </h2>
          <p>
            If you have questions about this disclaimer, please reach out via our <a href="https://github.com/mason2047/claway" target="_blank" rel="noopener noreferrer" className="text-accent hover:underline">GitHub repository</a>.
          </p>
        </section>

      </div>
    </div>
  );
}
