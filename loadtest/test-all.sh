#!/bin/bash
# Comprehensive load test covering all endpoints

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
DURATION="${DURATION:-30s}"
RATE="${RATE:-100}"

echo "----------------------------"
echo "Running comprehensive load test..."
echo "Base URL: $BASE_URL"
echo "Duration: $DURATION"
echo "Rate: $RATE req/s"
echo ""

echo "⚙️  Setup phase..."

for i in {1..3}; do
  curl -s -X POST "$BASE_URL/team/add" \
    -H "Content-Type: application/json" \
    -d '{
      "team_name": "team'$i'",
      "members": [
        {"user_id": "t'$i'_user1", "username": "User1_T'$i'", "is_active": true},
        {"user_id": "t'$i'_user2", "username": "User2_T'$i'", "is_active": true},
        {"user_id": "t'$i'_user3", "username": "User3_T'$i'", "is_active": false}
      ]
    }' > /dev/null
done

for i in {1..10}; do
  team=$((i % 3 + 1))
  curl -s -X POST "$BASE_URL/pullRequest/create" \
    -H "Content-Type: application/json" \
    -d '{"pull_request_id":"loadtest_pr'$i'","pull_request_name":"Load Test PR '$i'","author_id":"t'$team'_user1"}' > /dev/null
done

echo "Setup complete: 3 teams, 10 PRs created"
echo ""

cat > /tmp/vegeta-targets.txt <<EOF
GET $BASE_URL/health
GET $BASE_URL/health
GET $BASE_URL/team/get?team_name=team1
GET $BASE_URL/team/get?team_name=team2
GET $BASE_URL/team/get?team_name=team3
GET $BASE_URL/users/getReview?user_id=t1_user2
GET $BASE_URL/users/getReview?user_id=t2_user2
GET $BASE_URL/users/getReview?user_id=t3_user2
EOF

echo "Running mixed workload test (read-heavy)..."
vegeta attack -duration=$DURATION -rate=$RATE -targets=/tmp/vegeta-targets.txt | \
  tee /tmp/vegeta-results.bin | \
  vegeta report -type=text

vegeta report -type=json /tmp/vegeta-results.bin | jq '{
  latencies: .latencies,
  success_rate: .success,
  requests: .requests,
  throughput: .throughput,
  duration: .duration,
  status_codes: .status_codes
}'

echo ""
vegeta plot /tmp/vegeta-results.bin > report.html
echo "Load test complete"
echo ""
echo "HTML report generated at report.html"
echo "View with: open report.html"
echo "----------------------------"

# Cleanup
rm -f /tmp/vegeta-targets.txt