#!/bin/bash
# Simple load test for PR endpoints using vegeta

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
DURATION="${DURATION:-10s}"
RATE="${RATE:-50}"

echo "----------------------------"
echo "Load testing PR endpoints..."
echo "Base URL: $BASE_URL"
echo "Duration: $DURATION"
echo "Rate: $RATE req/s"
echo ""

echo "Setup: Creating team and initial PRs..."
curl -s -X POST "$BASE_URL/team/add" \
  -H "Content-Type: application/json" \
  -d @loadtest/team-add.json > /dev/null

for i in {1..5}; do
  curl -s -X POST "$BASE_URL/pullRequest/create" \
    -H "Content-Type: application/json" \
    -d '{"pull_request_id":"pr'$i'","pull_request_name":"Test PR '$i'","author_id":"user1"}' > /dev/null
done

echo "Setup complete"
echo ""

echo "Testing GET /users/getReview"
echo "GET $BASE_URL/users/getReview?user_id=user2" | \
  vegeta attack -duration=$DURATION -rate=$RATE | \
  vegeta report -type=text

echo ""
echo "PR endpoints test complete"
echo "----------------------------"
