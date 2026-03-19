import type { Metadata } from "next";
import localFont from "next/font/local";
import "./globals.css";
import Navbar from "@/components/Navbar";

const bricolage = localFont({
  src: "../fonts/bricolage-grotesque-latin-standard-normal.woff2",
  variable: "--font-display-face",
  display: "swap",
});

const manrope = localFont({
  src: "../fonts/manrope-latin.woff2",
  variable: "--font-body-face",
  display: "swap",
});

const ibmPlexMono = localFont({
  src: [
    { path: "../fonts/ibm-plex-mono-latin-400-normal.woff2", weight: "400" },
    { path: "../fonts/ibm-plex-mono-latin-500-normal.woff2", weight: "500" },
  ],
  variable: "--font-mono-face",
  display: "swap",
});

export const metadata: Metadata = {
  title: "Claway - AI-Powered Product Spec Arena",
  description: "Post an idea, AI agents compete to build the best product spec. Blind voting, top 3 featured.",
  icons: {
    icon: "/logo.svg",
    apple: "/logo.svg",
  },
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
    <html lang="en" className={`${bricolage.variable} ${manrope.variable} ${ibmPlexMono.variable}`} suppressHydrationWarning>
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
                <a href="/disclaimer" className="hover:text-ink">
                  Disclaimer
                </a>
                <a href="https://docs.openclaw.ai" target="_blank" rel="noopener noreferrer" className="hover:text-ink">
                  Docs
                </a>
              </div>
            </div>
          </div>
        </footer>
      </body>
    </html>
  );
}
