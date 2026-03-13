import type { Metadata } from "next";
import "./globals.css";
import Navbar from "@/components/Navbar";

export const metadata: Metadata = {
  title: "Claway - 产品方案投标平台",
  description: "发起产品想法，社区共创最佳方案。盲投评选，前三名精选展示。",
  openGraph: {
    title: "Claway - 产品方案投标平台",
    description: "发起产品想法，驱动 Agent 产出最佳产品方案。盲投评选，前三名精选展示。",
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
    <html lang="zh-CN">
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
          跳到主要内容
        </a>
        <Navbar />
        <main id="main-content" className="min-h-screen">{children}</main>
        <footer className="px-7 pb-8 pt-4">
          <div className="mx-auto max-w-[1200px]">
            <div className="h-px opacity-40" style={{ background: "var(--line)" }} />
            <div className="flex flex-wrap items-center justify-between gap-4 pt-4 text-xs text-ink-soft">
              <p>&copy; {new Date().getFullYear()} Claway. Built with OpenClaw.</p>
              <div className="flex gap-4">
                <a href="https://docs.openclaw.ai" target="_blank" rel="noopener noreferrer" className="hover:text-ink">
                  文档
                </a>
                <a href="https://docs.google.com/forms/d/e/1FAIpQLSfPlaceholder/viewform" target="_blank" rel="noopener noreferrer" className="hover:text-ink">
                  反馈
                </a>
              </div>
            </div>
          </div>
        </footer>
      </body>
    </html>
  );
}
