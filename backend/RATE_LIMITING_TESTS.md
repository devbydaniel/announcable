# Rate Limiting Integration Tests

This document describes the integration tests for the rate limiting implementation.

## Overview

The rate limiting implementation includes:
- **IP-based rate limiting**: Limits requests per IP address
- **Org-based rate limiting**: Limits requests per organization ID
- **Rate limit headers**: Includes X-RateLimit-Limit, X-RateLimit-Remaining, and X-RateLimit-Retry-After headers

## Running the Integration Tests

### Prerequisites

1. Ensure backend services are running:
   ```bash
   cd backend
   make dev-services  # Start postgres, mail, minio
   ```

2. Start the backend server:
   ```bash
   cd backend
   make dev-air
   ```

3. Wait for the server to start (should be available at http://localhost:8080)

### Running the Test Script

```bash
cd backend
./test_rate_limiting.sh
```

### What the Tests Verify

The integration test script (`test_rate_limiting.sh`) verifies:

1. **Server is running**: Ensures the API endpoint is accessible
2. **Initial request succeeds**: Verifies the endpoint returns 200 or 404
3. **Rate limit headers present**: Confirms X-RateLimit-Limit and X-RateLimit-Remaining headers are included
4. **IP-based rate limiting**: Makes 101 requests to trigger rate limiting, expects 429 responses
5. **Retry-after header**: Verifies 429 responses include X-RateLimit-Retry-After header
6. **Rate limit recovery**: Waits for rate limit to reset and confirms requests succeed again
7. **Org-based rate limiting**: Tests that different organizations have separate rate limits

## Manual Testing

You can also manually test rate limiting using curl:

### Test IP-based rate limiting

```bash
# Make multiple requests to the same endpoint
for i in {1..101}; do
  curl -s -o /dev/null -w '%{http_code}\n' \
    http://localhost:8080/api/widget-config/00000000-0000-0000-0000-000000000000
done | tail -10
```

Expected: Last few requests should return 429

### Test rate limit headers

```bash
# Check headers in response
curl -I http://localhost:8080/api/widget-config/00000000-0000-0000-0000-000000000000
```

Expected headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
```

### Test 429 response with retry-after

```bash
# Make enough requests to trigger rate limit
for i in {1..101}; do
  curl -s -o /dev/null http://localhost:8080/api/widget-config/00000000-0000-0000-0000-000000000000
done

# Check headers on next request (should be rate limited)
curl -I http://localhost:8080/api/widget-config/00000000-0000-0000-0000-000000000000
```

Expected headers on 429 response:
```
HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Retry-After: 60
```

### Test org-based rate limiting

```bash
# Make requests to two different orgs from the same IP
# Org 1
for i in {1..60}; do
  curl -s -o /dev/null -w '%{http_code}\n' \
    http://localhost:8080/api/widget-config/00000000-0000-0000-0000-000000000000
done

# Org 2 (should have separate limit)
for i in {1..60}; do
  curl -s -o /dev/null -w '%{http_code}\n' \
    http://localhost:8080/api/widget-config/11111111-1111-1111-1111-111111111111
done
```

Expected: Second organization should not be rate limited even though same IP made 60 requests to first org

## Configuration

Rate limits are configured via environment variables (see `.env.example`):

```bash
# Public API rate limiting
RATE_LIMIT_PUBLIC_REQUESTS_PER_INTERVAL=100
RATE_LIMIT_PUBLIC_MAX_TOKENS=100
RATE_LIMIT_PUBLIC_REFILL_INTERVAL_SECONDS=60

# Authenticated API rate limiting
RATE_LIMIT_AUTHENTICATED_REQUESTS_PER_INTERVAL=1000
RATE_LIMIT_AUTHENTICATED_MAX_TOKENS=1000
RATE_LIMIT_AUTHENTICATED_REFILL_INTERVAL_SECONDS=60
```

## Acceptance Criteria

✅ Public endpoints return HTTP 429 when rate limit is exceeded
✅ Rate limits are configurable via environment variables
✅ Rate limiting uses IP-based and org-based strategies
✅ Rate limit headers (X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Retry-After) are included in responses

## Troubleshooting

### Server not responding
- Check if backend is running: `ps aux | grep air`
- Check server logs in terminal running `make dev-air`
- Verify services are running: `docker ps`

### Rate limit not triggering
- Check environment variables are set correctly
- Verify rate limit configuration in backend logs
- Check if rate limit interval is too long or max tokens too high

### Headers missing
- Ensure middleware is properly applied to routes
- Check middleware initialization in `main.go`
- Verify rate limiter is initialized with config values
