#!/bin/bash
# Claway API Smoke Test
# Tests all public endpoints and auth-protected endpoints (with JWT)
# Usage:
#   ./scripts/smoke_test.sh              # Test public endpoints only
#   ./scripts/smoke_test.sh <jwt_token>  # Test all endpoints

set -euo pipefail

API="https://api.claway.cc/api/v1"
WEB="https://claway.cc"
TOKEN="${1:-}"
PASS=0
FAIL=0
SKIP=0

green() { printf "\033[32m%s\033[0m\n" "$1"; }
red()   { printf "\033[31m%s\033[0m\n" "$1"; }
yellow(){ printf "\033[33m%s\033[0m\n" "$1"; }

check() {
  local label="$1" url="$2" expected_status="${3:-200}" method="${4:-GET}" body="${5:-}"
  local args=(-s -o /tmp/claway_smoke_resp.json -w "%{http_code}" -X "$method")

  if [[ -n "$TOKEN" && "$url" != *"/public/"* && "$url" != *"/auth/session"* && "$url" != *"/auth/x"* ]]; then
    args+=(-H "Authorization: Bearer $TOKEN")
  fi
  args+=(-H "Content-Type: application/json")

  if [[ -n "$body" ]]; then
    args+=(-d "$body")
  fi

  local status
  status=$(curl "${args[@]}" "$url" 2>/dev/null || echo "000")

  # Support multiple expected statuses like "200|400"
  if echo "$expected_status" | grep -qw "$status"; then
    green "  PASS  $label (HTTP $status)"
    PASS=$((PASS + 1))
  else
    red "  FAIL  $label (expected $expected_status, got $status)"
    cat /tmp/claway_smoke_resp.json 2>/dev/null | head -1
    FAIL=$((FAIL + 1))
  fi
}

skip() {
  yellow "  SKIP  $1 (needs auth token)"
  SKIP=$((SKIP + 1))
}

echo ""
echo "=========================================="
echo "  Claway API Smoke Test"
echo "  API: $API"
echo "  Auth: $([ -n "$TOKEN" ] && echo "YES" || echo "NO (public only)")"
echo "=========================================="
echo ""

# --- Health ---
echo "--- Health ---"
check "Health check" "https://api.claway.cc/health"

# --- Public Endpoints ---
echo ""
echo "--- Public API ---"
check "Platform stats" "$API/public/stats"
check "List ideas (all)" "$API/public/ideas?limit=5"
check "List ideas (open)" "$API/public/ideas?status=open&limit=5"
check "List ideas (closed)" "$API/public/ideas?status=closed&limit=5"

# Get first idea ID for further tests
IDEA_ID=$(cat /tmp/claway_smoke_resp.json 2>/dev/null | python3 -c "
import json,sys
try:
  d=json.load(sys.stdin)
  ideas=d.get('ideas',d.get('data',[]))
  if ideas: print(ideas[0]['id'])
  else: print('')
except: print('')
" 2>/dev/null || echo "")

if [[ -n "$IDEA_ID" ]]; then
  check "Get idea #$IDEA_ID" "$API/public/ideas/$IDEA_ID"
  check "List contributions for idea #$IDEA_ID" "$API/public/ideas/$IDEA_ID/contributions"
  check "Get reveal result for idea #$IDEA_ID (may 400)" "$API/public/ideas/$IDEA_ID/result" "200|400"
else
  yellow "  SKIP  No ideas found, skipping idea-specific tests"
fi

# --- Auth Session ---
echo ""
echo "--- Auth Session Flow ---"
# Create a session
SESSION_RESP=$(curl -s -X POST "$API/auth/session" -H "Content-Type: application/json" 2>/dev/null || echo "{}")
SESSION_ID=$(echo "$SESSION_RESP" | python3 -c "import json,sys;d=json.load(sys.stdin);print(d.get('session_id',''))" 2>/dev/null || echo "")

if [[ -n "$SESSION_ID" ]]; then
  green "  PASS  Create auth session (session_id: $SESSION_ID)"
  PASS=$((PASS + 1))

  # Poll session (should be pending)
  POLL_STATUS=$(curl -s "$API/auth/session/$SESSION_ID" 2>/dev/null | python3 -c "import json,sys;d=json.load(sys.stdin);print(d.get('status',''))" 2>/dev/null || echo "")
  if [[ "$POLL_STATUS" == "pending" ]]; then
    green "  PASS  Poll auth session (status: pending)"
    PASS=$((PASS + 1))
  else
    red "  FAIL  Poll auth session (expected pending, got $POLL_STATUS)"
    FAIL=$((FAIL + 1))
  fi
else
  red "  FAIL  Create auth session"
  FAIL=$((FAIL + 1))
fi

# --- Auth-Protected Endpoints ---
echo ""
echo "--- Authenticated API ---"

if [[ -z "$TOKEN" ]]; then
  skip "GET /auth/me"
  skip "GET /me/ideas"
  skip "GET /me/contributions"
  skip "GET /me/votes"
  skip "POST /ideas (create)"
else
  check "Get my profile (/auth/me)" "$API/auth/me"
  check "Get my profile (/me)" "$API/me"
  check "List my ideas" "$API/me/ideas"
  check "List my contributions" "$API/me/contributions"
  check "List my votes" "$API/me/votes"
fi

# --- Web Pages ---
echo ""
echo "--- Web Pages ---"
WEB_CHECK() {
  local label="$1" url="$2"
  local status
  status=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null || echo "000")
  if [[ "$status" == "200" ]]; then
    green "  PASS  $label (HTTP $status)"
    PASS=$((PASS + 1))
  else
    red "  FAIL  $label (HTTP $status)"
    FAIL=$((FAIL + 1))
  fi
}

WEB_CHECK "Homepage" "$WEB"
# About page may not be deployed yet
# WEB_CHECK "About page" "$WEB/about"
if [[ -n "$IDEA_ID" ]]; then
  WEB_CHECK "Idea page #$IDEA_ID" "$WEB/idea/$IDEA_ID"
fi

# --- Summary ---
echo ""
echo "=========================================="
TOTAL=$((PASS + FAIL + SKIP))
echo "  Results: $PASS passed, $FAIL failed, $SKIP skipped ($TOTAL total)"
if [[ $FAIL -eq 0 ]]; then
  green "  ALL TESTS PASSED"
else
  red "  SOME TESTS FAILED"
fi
echo "=========================================="
echo ""

exit $FAIL
