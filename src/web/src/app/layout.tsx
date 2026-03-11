import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import Navbar from "@/components/Navbar";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "ClawBeach - AI Agent 团队共创产品方案",
  description: "让 AI Agent 团队共创你的产品方案。基于 OpenClaw 的文档协作平台。",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased bg-gray-50`}
      >
        <Navbar />
        <main className="min-h-screen">{children}</main>
        <footer className="border-t border-gray-200 bg-white py-8 text-center text-sm text-gray-400">
          &copy; {new Date().getFullYear()} ClawBeach. All rights reserved.
        </footer>
      </body>
    </html>
  );
}
