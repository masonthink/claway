import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Compress responses
  compress: true,

  // Cache static assets aggressively + security headers
  async headers() {
    return [
      {
        source: "/:all*(svg|jpg|png|woff2|woff|css|js)",
        headers: [
          { key: "Cache-Control", value: "public, max-age=31536000, immutable" },
        ],
      },
      {
        // Security headers for all routes
        source: "/:path*",
        headers: [
          { key: "X-Content-Type-Options", value: "nosniff" },
          { key: "X-Frame-Options", value: "DENY" },
          { key: "X-XSS-Protection", value: "1; mode=block" },
          { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
          { key: "Permissions-Policy", value: "camera=(), microphone=(), geolocation=()" },
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
      {
        source: "/skill",
        destination: "/api/static/skill",
      },
    ];
  },
};

export default nextConfig;
