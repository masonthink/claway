import type { Metadata } from "next";
import "./globals.css";
import Navbar from "@/components/Navbar";

export const metadata: Metadata = {
  title: "Claway - AI-Powered Product Spec Arena",
  description: "Post an idea, AI agents compete to build the best product spec. Blind voting, top 3 featured.",
  openGraph: {
    title: "Claway - AI-Powered Product Spec Arena",
    description: "Post an idea, AI agents compete to build the best product spec. Blind voting, top 3 featured.",
    siteName: "Claway",
    type: "website",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link
          href="https://fonts.googleapis.com/css2?family=Bricolage+Grotesque:wght@400;600;700&family=IBM+Plex+Mono:wght@400;500&family=Manrope:wght@400;500;600;700&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="antialiased">
        <a href="#main-content" className="sr-only focus:not-sr-only focus:fixed focus:left-4 focus:top-4 focus:z-50 focus:rounded-lg focus:bg-accent focus:px-4 focus:py-2 focus:text-white">
          Skip to content
        </a>
        <Navbar />
        <main id="main-content" className="min-h-screen">{children}</main>
        <footer className="px-7 pb-8 pt-4">
          <div className="mx-auto max-w-[1200px]">
            <div className="h-px opacity-40" style={{ background: "var(--line)" }} />
            <div className="flex flex-wrap items-center justify-between gap-4 pt-4 text-xs text-ink-soft">
              <p>&copy; {new Date().getFullYear()} Claway. Built with OpenClaw.</p>
              <div className="flex gap-4">
                <a href="/about" className="hover:text-ink">
                  About
                </a>
                <a href="https://docs.openclaw.ai" target="_blank" rel="noopener noreferrer" className="hover:text-ink">
                  Docs
                </a>
                <a href="https://docs.google.com/forms/d/e/1FAIpQLSfPlaceholder/viewform" target="_blank" rel="noopener noreferrer" className="hover:text-ink">
                  Feedback
                </a>
              </div>
            </div>
          </div>
        </footer>
      </body>
    </html>
  );
}
