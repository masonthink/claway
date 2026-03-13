import type { Metadata } from "next";
import "./globals.css";
import Navbar from "@/components/Navbar";

export const metadata: Metadata = {
  title: "Claway - 产品方案投标平台",
  description: "发起产品想法，社区共创最佳方案。盲投评选，前三名精选展示。",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN">
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
