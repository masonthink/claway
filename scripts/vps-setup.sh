#!/bin/bash
# Claway VPS Setup Script
# Run on VPS as deploy user
# Usage: ssh vps 'bash -s' < scripts/vps-setup.sh

set -euo pipefail

echo "=== 1. Create deployment directory ==="
sudo mkdir -p /opt/claway
sudo chown -R deploy:deploy /opt/claway

echo "=== 2. Create .env.prod template ==="
if [ ! -f /opt/claway/.env.prod ]; then
  cat > /opt/claway/.env.prod << 'ENVEOF'
# Claway Production Environment
DATABASE_URL=postgresql://claway:CHANGE_ME_CLAWAY_DB_PASSWORD@localhost:5432/claway?sslmode=disable

# Auth
JWT_SECRET=CHANGE_ME_JWT_SECRET

# Server
PORT=8081

# OpenClaw OAuth
OPENCLAW_CLIENT_ID=
OPENCLAW_CLIENT_SECRET=
OPENCLAW_BASE_URL=https://api.openclaw.ai

# LLM Proxy (upstream provider)
UPSTREAM_LLM_BASE_URL=https://api.openai.com
UPSTREAM_LLM_API_KEY=

# Docker image (auto-updated by CI/CD)
BACKEND_IMAGE=ghcr.io/mason2047/claway-backend:latest
ENVEOF
  echo ".env.prod created at /opt/claway/.env.prod — EDIT IT with real values!"
else
  echo ".env.prod already exists, skipping"
fi

echo "=== 3. Login to GHCR ==="
echo "Run manually: echo \$GITHUB_TOKEN | docker login ghcr.io -u mason2047 --password-stdin"

echo ""
echo "=== SETUP COMPLETE ==="
echo ""
echo "Next steps:"
echo "  1. Edit /opt/claway/.env.prod with real passwords"
echo "  2. scp docker-compose.prod.yml to /opt/claway/"
echo "  3. Add DNS records in Cloudflare:"
echo "     - api.claway.cc → VPS IP (proxied)"
echo "     - claway.cc → Vercel (CNAME)"
echo "  4. Set Cloudflare SSL to Flexible"
echo "  5. Set GitHub repo secrets: VPS_HOST, VPS_SSH_KEY"
echo "  6. Push code to trigger first deployment"
