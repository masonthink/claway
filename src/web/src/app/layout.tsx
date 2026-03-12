import type { Metadata } from "next";
import "./globals.css";
import Navbar from "@/components/Navbar";
import { ToastProvider } from "@/components/Toast";

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
      <body className="antialiased">
        <ToastProvider>
        <Navbar />
        <main className="min-h-screen">{children}</main>
        <footer className="px-7 pb-8 pt-4">
          <div className="mx-auto max-w-[1200px]">
            <div className="h-px opacity-40" style={{ background: "var(--line)" }} />
          </div>
        </footer>
        </ToastProvider>
      </body>
    </html>
  );
}
