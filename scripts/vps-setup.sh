#!/bin/bash
# Claway VPS Setup Script
# Run on VPS as root, then switch deploy user ownership
# Usage: ssh -i ~/.ssh/dtc_deploy_vps deploy@45.32.57.146 'bash -s' < scripts/vps-setup.sh

set -euo pipefail

echo "=== 1. Create deployment directory ==="
sudo mkdir -p /opt/claway
sudo chown -R deploy:deploy /opt/claway

echo "=== 2. Create PostgreSQL database ==="
# Claway shares the DTC PostgreSQL instance
# Create a dedicated database and user
docker exec dtc-postgres psql -U dtc -d postgres -c "
  DO \$\$
  BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'claway') THEN
      CREATE ROLE claway WITH LOGIN PASSWORD 'CHANGE_ME_CLAWAY_DB_PASSWORD';
    END IF;
  END
  \$\$;
  SELECT 'CREATE DATABASE claway OWNER claway'
  WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'claway')
  \gexec
  GRANT ALL PRIVILEGES ON DATABASE claway TO claway;
"

echo "=== 3. Copy docker-compose.prod.yml ==="
# This file should be scp'd separately or committed
# scp -i ~/.ssh/dtc_deploy_vps docker-compose.prod.yml deploy@45.32.57.146:/opt/claway/

echo "=== 4. Create .env.prod template ==="
if [ ! -f /opt/claway/.env.prod ]; then
  cat > /opt/claway/.env.prod << 'ENVEOF'
# Claway Production Environment
# Database (shares DTC PostgreSQL via host network)
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

echo "=== 5. Update Caddy config ==="
# Add Claway site block to existing DTC Caddyfile
CADDYFILE="/opt/digital-twin/deploy/Caddyfile"
if grep -q "claway.concors.ai" "$CADDYFILE" 2>/dev/null; then
  echo "Caddy config for claway.concors.ai already exists, skipping"
else
  cat >> "$CADDYFILE" << 'CADDYEOF'

# ── Claway API ──────────────────────────────────────────────────────────
claway.concors.ai {

    reverse_proxy localhost:8081 {
        health_uri      /health
        health_interval 10s
        health_timeout  5s

        header_up X-Real-IP        {http.request.header.CF-Connecting-IP}
        header_up X-Forwarded-For  {http.request.header.CF-Connecting-IP}
    }

    log {
        output stdout
        format json
    }

    header {
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
        X-Content-Type-Options    "nosniff"
        X-Frame-Options           "DENY"
        Referrer-Policy           "strict-origin-when-cross-origin"
        -Server
    }
}
CADDYEOF
  echo "Added claway.concors.ai to Caddyfile"

  # Reload Caddy
  docker exec dtc-caddy caddy reload --config /etc/caddy/Caddyfile 2>/dev/null || \
    docker restart dtc-caddy
  echo "Caddy reloaded"
fi

echo "=== 6. Login to GHCR ==="
echo "Run manually: echo \$GITHUB_TOKEN | docker login ghcr.io -u mason2047 --password-stdin"

echo ""
echo "=== SETUP COMPLETE ==="
echo ""
echo "Next steps:"
echo "  1. Edit /opt/claway/.env.prod with real passwords"
echo "  2. scp docker-compose.prod.yml to /opt/claway/"
echo "  3. Add DNS record: claway.concors.ai → 45.32.57.146 (Cloudflare, proxied)"
echo "  4. Set Cloudflare SSL to Flexible for this subdomain"
echo "  5. Set GitHub repo secrets: VPS_HOST, VPS_SSH_KEY"
echo "  6. Push code to trigger first deployment"
