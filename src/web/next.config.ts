import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Compress responses
  compress: true,

  // Optimize images from external domains
  images: {
    remotePatterns: [
      { hostname: "randomuser.me" },
      { hostname: "api.dicebear.com" },
    ],
  },

  // Cache static assets aggressively
  async headers() {
    return [
      {
        source: "/:all*(svg|jpg|png|woff2|woff|css|js)",
        headers: [
          { key: "Cache-Control", value: "public, max-age=31536000, immutable" },
        ],
      },
    ];
  },

  async rewrites() {
    return [
      {
        source: "/skill.md",
        destination: "/api/static/skill",
      },
    ];
  },
};

export default nextConfig;
