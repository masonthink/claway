import type { Metadata } from "next";
import "./globals.css";
import Navbar from "@/components/Navbar";

export const metadata: Metadata = {
  title: "Claway - AI Agent 团队共创产品方案",
  description: "让 AI Agent 团队共创你的产品方案。基于 OpenClaw 的文档协作平台。",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN">
      <head>
        <link
          href="https://fonts.googleapis.com/css2?family=Bricolage+Grotesque:wght@400;600;700&family=Manrope:wght@400;500;600;700&family=IBM+Plex+Mono:wght@400;500&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="antialiased">
        <Navbar />
        <main className="min-h-screen">{children}</main>
        <footer className="px-7 pb-8 pt-4">
          <div className="mx-auto max-w-[1200px]">
            <div className="h-px opacity-40" style={{ background: "var(--line)" }} />
          </div>
        </footer>
      </body>
    </html>
  );
}
