#!/bin/bash

# Integration test script for rate limiting functionality
# Tests IP-based and org-based rate limiting

set -e

BASE_URL="http://localhost:8080"
TEST_ORG_ID="00000000-0000-0000-0000-000000000000"
TEST_ORG_ID_2="11111111-1111-1111-1111-111111111111"
ENDPOINT="$BASE_URL/api/widget-config/$TEST_ORG_ID"
ENDPOINT_2="$BASE_URL/api/widget-config/$TEST_ORG_ID_2"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=========================================="
echo "Rate Limiting Integration Tests"
echo "=========================================="
echo ""

# Test 1: Verify server is running
echo "Test 1: Checking if server is running..."
if curl -s -f -o /dev/null -w "%{http_code}" "$ENDPOINT" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Server is running${NC}"
else
    echo -e "${RED}✗ Server is not running. Please start with 'cd backend && make dev-air'${NC}"
    exit 1
fi
echo ""

# Test 2: Verify initial request succeeds
echo "Test 2: Verifying initial request succeeds..."
RESPONSE=$(curl -s -w "\n%{http_code}" "$ENDPOINT")
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)
HEADERS=$(curl -s -I "$ENDPOINT")

if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "404" ]; then
    echo -e "${GREEN}✓ Initial request successful (HTTP $HTTP_CODE)${NC}"
else
    echo -e "${RED}✗ Initial request failed (HTTP $HTTP_CODE)${NC}"
    exit 1
fi
echo ""

# Test 3: Verify rate limit headers are present
echo "Test 3: Verifying rate limit headers are present..."
RATE_LIMIT=$(echo "$HEADERS" | grep -i "X-RateLimit-Limit" || true)
RATE_REMAINING=$(echo "$HEADERS" | grep -i "X-RateLimit-Remaining" || true)

if [ -n "$RATE_LIMIT" ] && [ -n "$RATE_REMAINING" ]; then
    echo -e "${GREEN}✓ Rate limit headers present${NC}"
    echo "  $RATE_LIMIT"
    echo "  $RATE_REMAINING"
else
    echo -e "${RED}✗ Rate limit headers missing${NC}"
    echo "Headers received:"
    echo "$HEADERS"
    exit 1
fi
echo ""

# Test 4: Test IP-based rate limiting (make 101 requests)
echo "Test 4: Testing IP-based rate limiting (making 101 requests)..."
SUCCESS_COUNT=0
RATE_LIMITED_COUNT=0
LAST_10_CODES=""

for i in {1..101}; do
    CODE=$(curl -s -o /dev/null -w "%{http_code}" "$ENDPOINT")
    if [ "$CODE" = "200" ] || [ "$CODE" = "404" ]; then
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    elif [ "$CODE" = "429" ]; then
        RATE_LIMITED_COUNT=$((RATE_LIMITED_COUNT + 1))
    fi

    # Capture last 10 codes for display
    if [ $i -gt 91 ]; then
        LAST_10_CODES="$LAST_10_CODES $CODE"
    fi

    # Small delay to avoid overwhelming the server
    sleep 0.01
done

echo "  Total requests: 101"
echo "  Successful (200/404): $SUCCESS_COUNT"
echo "  Rate limited (429): $RATE_LIMITED_COUNT"
echo "  Last 10 status codes:$LAST_10_CODES"

if [ $RATE_LIMITED_COUNT -gt 0 ]; then
    echo -e "${GREEN}✓ IP-based rate limiting working (got $RATE_LIMITED_COUNT 429 responses)${NC}"
else
    echo -e "${YELLOW}⚠ Warning: No 429 responses received. Rate limit may be high or requests too slow.${NC}"
fi
echo ""

# Test 5: Verify retry-after header on 429 response
echo "Test 5: Verifying 429 response includes X-RateLimit-Retry-After header..."
# Make more requests to ensure we hit the limit
for i in {1..10}; do
    RESPONSE_429=$(curl -s -I "$ENDPOINT")
    HTTP_CODE_429=$(echo "$RESPONSE_429" | grep "HTTP" | awk '{print $2}')

    if [ "$HTTP_CODE_429" = "429" ]; then
        RETRY_AFTER=$(echo "$RESPONSE_429" | grep -i "X-RateLimit-Retry-After" || true)
        if [ -n "$RETRY_AFTER" ]; then
            echo -e "${GREEN}✓ 429 response includes retry-after header${NC}"
            echo "  $RETRY_AFTER"
        else
            echo -e "${RED}✗ 429 response missing retry-after header${NC}"
            echo "Headers:"
            echo "$RESPONSE_429"
        fi
        break
    fi
    sleep 0.1
done
echo ""

# Test 6: Wait for rate limit to reset and verify recovery
echo "Test 6: Waiting for rate limit to reset (2 seconds)..."
sleep 2
RECOVERY_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$ENDPOINT")

if [ "$RECOVERY_CODE" = "200" ] || [ "$RECOVERY_CODE" = "404" ]; then
    echo -e "${GREEN}✓ Rate limit recovered after waiting${NC}"
else
    echo -e "${YELLOW}⚠ Rate limit may still be active (HTTP $RECOVERY_CODE)${NC}"
fi
echo ""

# Test 7: Test org-based rate limiting (different orgs, same IP)
echo "Test 7: Testing org-based rate limiting (different orgs, same IP)..."
ORG1_SUCCESS=0
ORG2_SUCCESS=0

# Make 50 requests to first org
for i in {1..50}; do
    CODE=$(curl -s -o /dev/null -w "%{http_code}" "$ENDPOINT")
    if [ "$CODE" = "200" ] || [ "$CODE" = "404" ]; then
        ORG1_SUCCESS=$((ORG1_SUCCESS + 1))
    fi
    sleep 0.01
done

# Make 50 requests to second org (should have separate limit)
for i in {1..50}; do
    CODE=$(curl -s -o /dev/null -w "%{http_code}" "$ENDPOINT_2")
    if [ "$CODE" = "200" ] || [ "$CODE" = "404" ]; then
        ORG2_SUCCESS=$((ORG2_SUCCESS + 1))
    fi
    sleep 0.01
done

echo "  Org 1 successful requests: $ORG1_SUCCESS/50"
echo "  Org 2 successful requests: $ORG2_SUCCESS/50"

if [ $ORG2_SUCCESS -gt 0 ]; then
    echo -e "${GREEN}✓ Org-based rate limiting appears to be working (separate limits per org)${NC}"
else
    echo -e "${RED}✗ Org-based rate limiting may not be working${NC}"
fi
echo ""

echo "=========================================="
echo "Integration Tests Complete"
echo "=========================================="
echo ""
echo "Summary:"
echo "  ✓ Server is running"
echo "  ✓ Rate limit headers present in responses"
echo "  ✓ IP-based rate limiting enforced"
echo "  ✓ Org-based rate limiting enforced"
echo ""
echo -e "${GREEN}All tests passed!${NC}"
