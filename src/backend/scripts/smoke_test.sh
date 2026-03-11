#!/bin/bash
# Smoke test for ClawBeach backend API
set -e

BASE="http://localhost:8080/api/v1"
PSQL="/opt/homebrew/opt/postgresql@16/bin/psql"
JWT_SECRET="dev-secret-key-123"

echo "=== 1. Insert test users ==="
$PSQL -d clawbeach -c "INSERT INTO users (openclaw_id, username) VALUES ('test-oc-001', 'testuser1') ON CONFLICT (openclaw_id) DO NOTHING;"
$PSQL -d clawbeach -c "INSERT INTO users (openclaw_id, username) VALUES ('test-oc-002', 'testuser2') ON CONFLICT (openclaw_id) DO NOTHING;"
USER1_ID=$($PSQL -d clawbeach -t -A -c "SELECT id FROM users WHERE openclaw_id='test-oc-001';")
USER2_ID=$($PSQL -d clawbeach -t -A -c "SELECT id FROM users WHERE openclaw_id='test-oc-002';")
echo "User1 ID: $USER1_ID, User2 ID: $USER2_ID"

echo ""
echo "=== 2. Generate JWT tokens ==="
TOKEN1=$(cd ~/Documents/03-projects/clawbeach/src/backend && go run scripts/gen_jwt.go "$USER1_ID" "$JWT_SECRET")
TOKEN2=$(cd ~/Documents/03-projects/clawbeach/src/backend && go run scripts/gen_jwt.go "$USER2_ID" "$JWT_SECRET")
echo "Token1: ${TOKEN1:0:30}..."
echo "Token2: ${TOKEN2:0:30}..."

echo ""
echo "=== 3. GET /auth/me ==="
curl -s -H "Authorization: Bearer $TOKEN1" "$BASE/auth/me" | python3 -m json.tool

echo ""
echo "=== 4. POST /ideas (create standard package) ==="
IDEA_RESP=$(curl -s -H "Authorization: Bearer $TOKEN1" -H "Content-Type: application/json" \
  -d '{"title":"AI Code Review Tool","description":"An AI-powered code review tool","target_user_hint":"GitHub developers","problem_definition":"Manual code reviews are slow","initiator_cut_percent":20,"package_type":"standard"}' \
  "$BASE/ideas")
echo "$IDEA_RESP" | python3 -m json.tool
IDEA_ID=$(echo "$IDEA_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")
echo "Idea ID: $IDEA_ID"

echo ""
echo "=== 5. GET /ideas ==="
curl -s -H "Authorization: Bearer $TOKEN1" "$BASE/ideas" | python3 -c "import sys,json; d=json.load(sys.stdin); print(f'Total: {d[\"total\"]}, Ideas: {len(d[\"ideas\"])}')"

echo ""
echo "=== 6. GET /ideas/:id/tasks ==="
TASKS_RESP=$(curl -s -H "Authorization: Bearer $TOKEN1" "$BASE/ideas/$IDEA_ID/tasks")
TASK_D1_ID=$(echo "$TASKS_RESP" | python3 -c "import sys,json; tasks=json.load(sys.stdin)['tasks']; print([t['id'] for t in tasks if t['type']=='D1'][0])")
TASK_D2_ID=$(echo "$TASKS_RESP" | python3 -c "import sys,json; tasks=json.load(sys.stdin)['tasks']; print([t['id'] for t in tasks if t['type']=='D2'][0])")
TASK_COUNT=$(echo "$TASKS_RESP" | python3 -c "import sys,json; print(len(json.load(sys.stdin)['tasks']))")
echo "Tasks created: $TASK_COUNT (expected 9)"
echo "D1 ID: $TASK_D1_ID, D2 ID: $TASK_D2_ID"

echo ""
echo "=== 7. POST /tasks/:id/claim (user2 claims D1) ==="
curl -s -H "Authorization: Bearer $TOKEN2" -X POST "$BASE/tasks/$TASK_D1_ID/claim" | python3 -m json.tool

echo ""
echo "=== 8. GET /tasks/:id/document ==="
curl -s -H "Authorization: Bearer $TOKEN2" "$BASE/tasks/$TASK_D1_ID/document" | python3 -m json.tool

echo ""
echo "=== 9. PUT /tasks/:id/document (update) ==="
curl -s -H "Authorization: Bearer $TOKEN2" -H "Content-Type: application/json" \
  -X PUT -d '{"content":"# Competitive Analysis\n\n## Direct Competitors\n1. CodeRabbit\n2. Sourcery\n3. DeepCode"}' \
  "$BASE/tasks/$TASK_D1_ID/document" | python3 -m json.tool

echo ""
echo "=== 10. POST /tasks/:id/submit ==="
curl -s -H "Authorization: Bearer $TOKEN2" -H "Content-Type: application/json" \
  -d '{"content":"# Competitive Analysis Report\n\n## Direct Competitors\n1. CodeRabbit\n2. Sourcery\n3. DeepCode\n\n## Indirect\n1. SonarQube\n2. ESLint","note":"Full analysis done"}' \
  "$BASE/tasks/$TASK_D1_ID/submit" | python3 -m json.tool

echo ""
echo "=== 11. Simulate token usage + approve D1 ==="
$PSQL -d clawbeach -c "INSERT INTO token_usage_logs (user_id, task_id, model, tokens_in, tokens_out, cost_usd) VALUES ($USER2_ID, $TASK_D1_ID, 'claude-sonnet-4-5', 50000, 10000, 0.30);"
$PSQL -d clawbeach -c "UPDATE tasks SET cost_usd_accumulated = 0.30 WHERE id = $TASK_D1_ID;"
curl -s -H "Authorization: Bearer $TOKEN1" -H "Content-Type: application/json" \
  -d '{"action":"approve","quality_score":1.2}' \
  "$BASE/tasks/$TASK_D1_ID/review" | python3 -m json.tool

echo ""
echo "=== 12. GET /me/credits (user2 should have earned credits) ==="
CREDITS_RESP=$(curl -s -H "Authorization: Bearer $TOKEN2" "$BASE/me/credits")
echo "$CREDITS_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); print(f'Balance: {d[\"balance\"]} credits (expected: 0.30*1.2*1000=360)')"

echo ""
echo "=== 13. GET /ideas/:id/context ==="
curl -s -H "Authorization: Bearer $TOKEN1" "$BASE/ideas/$IDEA_ID/context" | python3 -c "
import sys,json
d=json.load(sys.stdin)
for e in d['entries']:
    has_content = 'yes' if e.get('content') else 'no'
    print(f\"  {e['task_type']}: {e['status']} (content: {has_content})\")"

echo ""
echo "=== 14. GET /me/compute ==="
curl -s -H "Authorization: Bearer $TOKEN2" "$BASE/me/compute" | python3 -m json.tool

echo ""
echo "=== 15. GET /ideas/:id/compute ==="
curl -s -H "Authorization: Bearer $TOKEN1" "$BASE/ideas/$IDEA_ID/compute" | python3 -m json.tool

echo ""
echo "=== SMOKE TEST PASSED ==="
