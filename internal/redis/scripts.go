package redis

import "github.com/redis/go-redis/v9"

// TokenBucketScript is a Lua script that atomically checks and updates a token bucket.
// KEYS[1]: the bucket key
// ARGV[1]: burst size (max tokens after refill)
// ARGV[2]: refill rate (tokens per second)
// ARGV[3]: current time (unix seconds, float)
// ARGV[4]: TTL in seconds
// ARGV[5]: initial token count (limit)
// Returns: {allowed (0/1), remaining tokens}
var TokenBucketScript = redis.NewScript(`
local key = KEYS[1]
local burst = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local ttl = tonumber(ARGV[4])
local initial = tonumber(ARGV[5])

local data = redis.call('HMGET', key, 'tokens', 'last_refill')
local tokens = tonumber(data[1])
local last_refill = tonumber(data[2])

if tokens == nil then
    tokens = initial
    last_refill = now
end

local elapsed = math.max(0, now - last_refill)
tokens = math.min(burst, tokens + elapsed * rate)

local allowed = 0
local remaining = math.floor(tokens)

if tokens >= 1 then
    tokens = tokens - 1
    allowed = 1
    remaining = math.floor(tokens)
end

redis.call('HMSET', key, 'tokens', tostring(tokens), 'last_refill', tostring(now))
redis.call('EXPIRE', key, ttl)

return {allowed, remaining}
`)

// SlidingWindowScript is a Lua script that atomically checks and updates a sliding window counter.
// KEYS[1]: current window key
// KEYS[2]: previous window key
// ARGV[1]: limit
// ARGV[2]: window size in seconds
// ARGV[3]: current time (unix seconds, float)
// ARGV[4]: TTL in seconds
// Returns: {allowed (0/1), remaining, weighted count}
var SlidingWindowScript = redis.NewScript(`
local curr_key = KEYS[1]
local prev_key = KEYS[2]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local ttl = tonumber(ARGV[4])

local curr_start = math.floor(now / window) * window
local elapsed = now - curr_start
local weight = (window - elapsed) / window

local prev_count = tonumber(redis.call('GET', prev_key) or '0') or 0
local curr_count = tonumber(redis.call('GET', curr_key) or '0') or 0

local weighted = math.floor(prev_count * weight + curr_count)

if weighted >= limit then
    local remaining = math.max(0, limit - weighted)
    return {0, remaining, weighted}
end

curr_count = redis.call('INCR', curr_key)
redis.call('EXPIRE', curr_key, ttl)
redis.call('EXPIRE', prev_key, ttl)

weighted = math.floor(prev_count * weight + curr_count)
local remaining = math.max(0, limit - weighted)

return {1, remaining, weighted}
`)
