import type { NextConfig } from "next";

const nextConfig: NextConfig = {
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
