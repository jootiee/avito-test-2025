#!/bin/bash
# Simple load test for team endpoints using vegeta

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
DURATION="${DURATION:-10s}"
RATE="${RATE:-50}"

echo "----------------------------"
echo "Load testing team endpoints..."
echo "Base URL: $BASE_URL"
echo "Duration: $DURATION"
echo "Rate: $RATE req/s"
echo ""

echo "Setup: Creating initial team..."
curl -s -X POST "$BASE_URL/team/add" \
  -H "Content-Type: application/json" \
  -d @loadtest/team-add.json > /dev/null

echo ""

echo "Testing GET /team/get"
echo "GET $BASE_URL/team/get?team_name=vegeta_team" | \
  vegeta attack -duration=$DURATION -rate=$RATE | \
  vegeta report -type=text

echo ""
echo "Team endpoints test complete"
echo "----------------------------"
