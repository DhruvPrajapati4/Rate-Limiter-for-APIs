# API Rate Limiter

A production-grade API rate limiter built in **Go** using goroutines and channels, featuring configurable **token bucket / sliding window** algorithms with multi-tiered throttling: API-specific limits, global user quotas, and overall API thresholds. Designed for horizontal scalability with **Redis-backed** distributed state.

## Table of Contents
- [Features](#features)
- [Architecture](#architecture)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Quick Start](#quick-start)
  - [Docker](#docker)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Rate Limiting Algorithms](#rate-limiting-algorithms)
- [Multi-Tier Throttling](#multi-tier-throttling)
- [Testing](#testing)
- [Available Make Commands](#available-make-commands)

## Features

- **Dual algorithms** — Token Bucket and Sliding Window Counter, switchable via config
- **Multi-tier throttling** — API-specific, global user, and global API limits checked concurrently
- **Redis-backed** — Atomic Lua scripts for distributed state with no race conditions
- **Concurrent tier checks** — All 3 tiers evaluated in parallel using goroutines and channels
- **Standard headers** — `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`, `Retry-After`
- **Reusable middleware** — Plug-and-play Gin middleware for any endpoint
- **Flexible user identification** — Path param > `X-User-ID` header > query param > client IP

## Architecture

```
Request → Gin Middleware → Multi-Tier Throttler → Redis (Lua Scripts)
                              ├── API-Specific Limiter   (goroutine)
                              ├── Global User Limiter    (goroutine)
                              └── Global API Limiter     (goroutine)
```

```
cmd/server/main.go              # Entry point
internal/
  config/config.go              # YAML config loader
  limiter/
    limiter.go                  # Limiter interface + factory
    token_bucket.go             # Token bucket algorithm
    sliding_window.go           # Sliding window counter algorithm
  throttle/multi_tier.go        # Concurrent multi-tier composer
  middleware/ratelimit.go       # Gin HTTP middleware
  redis/
    client.go                   # Redis client wrapper
    scripts.go                  # Atomic Lua scripts
```

## Technologies Used

- **Go** 1.23
- **Gin** — HTTP framework
- **Redis** 7 — Distributed state store
- **go-redis/v9** — Redis client
- **Docker Compose** — Local development

## Getting Started

### Prerequisites

- Go 1.23+
- Redis 7+
- Docker (optional)

### Quick Start

1. **Clone the repository:**
   ```bash
   git clone https://github.com/DhruvPrajapati4/Rate-limiter-for-APIs.git
   cd Rate-limiter-for-APIs
   ```

2. **Start Redis:**
   ```bash
   make redis-up
   # or if Redis is already running locally, skip this step
   ```

3. **Configure** (edit `config.yaml` as needed):
   ```yaml
   redis:
     host: localhost
     port: 6379
     db: 0
   ```

4. **Build and run:**
   ```bash
   make run
   ```

5. **Test a request:**
   ```bash
   curl -v http://localhost:8090/api/user1/alpha
   ```

### Docker

Run both Redis and the app with a single command:
```bash
make docker-up
```

## Configuration

All configuration lives in `config.yaml`:

```yaml
server:
  port: 8090

redis:
  host: localhost
  port: 6379
  db: 0
  poolSize: 10
  minIdleConns: 5
  username: redisuser

rate_limit:
  algorithm: "token_bucket"  # or "sliding_window"
  tiers:
    api_specific:
      limit: 100             # Max requests per user per API
      window: 60s
      burst_size: 120        # Token bucket burst capacity
      refill_rate: 100       # Tokens refilled per window
    global_user:
      limit: 500             # Max requests per user across all APIs
      window: 60s
      burst_size: 600
      refill_rate: 500
    global_api:
      limit: 10000           # Max requests for any API across all users
      window: 1s
      burst_size: 12000
      refill_rate: 10000
```

| Field | Description |
|---|---|
| `algorithm` | `"token_bucket"` or `"sliding_window"` |
| `limit` | Maximum allowed requests in the window |
| `window` | Time window duration (e.g., `60s`, `1s`) |
| `burst_size` | Token bucket max capacity (allows short bursts) |
| `refill_rate` | Tokens replenished per window period |

### Environment Variables

| Variable | Description |
|---|---|
| `REDIS_PASSWORD` | Overrides `redis.password` from config. Use this to avoid storing secrets in `config.yaml`. |

**Example:**
```bash
REDIS_PASSWORD=mysecret make run
```

Or in Docker Compose:
```yaml
environment:
  - REDIS_PASSWORD=mysecret
```

## API Endpoints

### Sample Test APIs

| Method | Endpoint | Description |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/api/:userId/alpha` | Sample API Alpha |
| POST | `/api/:userId/beta` | Sample API Beta |
| PUT | `/api/:userId/gamma` | Sample API Gamma |
| DELETE | `/api/:userId/delta` | Sample API Delta |

### Response Headers

Every response includes rate limit headers:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 97
X-RateLimit-Reset: 1740512345
```

When rate limited (HTTP 429):
```
Retry-After: 1
```
```json
{
  "error": "rate limit exceeded",
  "retry_after": 0.6
}
```

## Rate Limiting Algorithms

### Token Bucket

Allows controlled bursts while maintaining an average rate. Each key has a bucket of tokens that refills at a steady rate. Each request consumes one token. When the bucket is empty, requests are denied.

- Allows bursts up to `burst_size`
- Refills at `refill_rate` tokens per window
- Implemented as an atomic Redis Lua script (no race conditions)

### Sliding Window Counter

A hybrid approach that tracks request counts in fixed windows and uses weighted interpolation for accuracy. Avoids the boundary-burst problem of fixed windows.

- Tracks current and previous window counts
- Weighted formula: `prev_count * overlap_ratio + current_count`
- More precise than fixed windows, lower memory than sliding logs

## Multi-Tier Throttling

Every request is checked against three tiers **concurrently** using goroutines and channels:

| Tier | Key Pattern | Purpose |
|---|---|---|
| API-Specific | `api:{name}:user:{id}` | Limit a user's requests to a specific API |
| Global User | `user:{id}` | Limit a user's total requests across all APIs |
| Global API | `global_api:{name}` | Limit total requests to an API across all users |

If **any** tier denies the request, the response is HTTP 429. The most restrictive tier's limits are returned in the headers.

## Testing

**Single request:**
```bash
curl -v http://localhost:8090/api/user1/alpha
```

**Flood test (trigger rate limiting):**
```bash
for i in $(seq 1 110); do
  echo -n "Request $i: "
  curl -s -o /dev/null -w "%{http_code}" http://localhost:8090/api/user1/alpha
  echo
done
```

**Test different users:**
```bash
curl http://localhost:8090/api/user1/alpha
curl http://localhost:8090/api/user2/alpha
```

**Test with header-based user ID:**
```bash
curl -H "X-User-ID: user1" http://localhost:8090/health
```

**Test different APIs (separate API-specific limits):**
```bash
curl http://localhost:8090/api/user1/alpha
curl -X POST http://localhost:8090/api/user1/beta
```

## Available Make Commands

Run `make` to see all commands:

```
Usage: make [target]

Targets:
  help            Show available commands
  build           Build the binary
  run             Build and run the server locally
  test            Run all tests
  vet             Run go vet static analysis
  lint            Run golangci-lint
  clean           Remove binary and build cache
  tidy            Tidy go module dependencies
  redis-up        Start only Redis
  redis-down      Stop only Redis
  docker-build    Build Docker image
  docker-up       Start Redis and app with Docker Compose
  docker-down     Stop and remove Docker containers
```
