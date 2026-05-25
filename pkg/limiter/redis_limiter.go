package limiter

import (
	"context"
	"strings"

	"github.com/JunLang-7/blog-service/global"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Redis Lua 脚本，在 Redis 端原子执行
var script = redis.NewScript(`
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local fill_interval = tonumber(ARGV[2])
local quantum = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])

-- 获取当前令牌数和上次补充时间戳
local data = redis.call('HMGET', key, 'tokens', 'last_refill_ts')
local tokens = tonumber(data[1])
local last_refill_ts = tonumber(data[2])

-- 获取 Redis 服务器时间 (秒 + 微秒)
local time = redis.call('TIME')
local now_ms = tonumber(time[1]) * 1000 + math.floor(tonumber(time[2]) / 1000)

-- 如果 key 不存在，初始化
if tokens == nil then
    tokens = capacity
    last_refill_ts = now_ms
end

-- 计算经过时间并补充令牌
local elapsed = now_ms - last_refill_ts
if elapsed > 0 then
    local refill_count = math.floor(elapsed / fill_interval) * quantum
    tokens = math.min(capacity, tokens + refill_count)
    last_refill_ts = now_ms
end

-- 计算 TTL：完整补充周期 * 2 作为安全边距，最少 10 秒
local ttl = fill_interval * math.ceil(capacity / quantum) * 2
if ttl < 10000 then ttl = 10000 end

-- 更新令牌和过期时间
redis.call('HSET', key, 'tokens', tokens, 'last_refill_ts', last_refill_ts)
redis.call('PEXPIRE', key, ttl)

-- 判断是否有足够令牌
if tokens >= requested then
    tokens = tokens - requested
    redis.call('HSET', key, 'tokens', tokens)
    return 1
end
return 0
`)

type RedisLimiter struct {
	redisBuckets map[string]LimitBucketRule
}

func NewRedisLimiter() *RedisLimiter {
	return &RedisLimiter{
		redisBuckets: make(map[string]LimitBucketRule),
	}
}

func (r RedisLimiter) Key(c *gin.Context) string {
	uri := c.Request.RequestURI
	if idx := strings.Index(uri, "?"); idx != -1 {
		return uri[:idx]
	}
	return uri
}

func (r RedisLimiter) Take(key string) bool {
	rule, ok := r.redisBuckets[key]
	if !ok {
		return true
	}

	redisKey := "rate_limiter:" + key
	res, err := script.Run(
		context.Background(),
		global.RedisClient,
		[]string{redisKey},
		rule.Capacity,
		rule.FillInterval.Milliseconds(),
		rule.Quantum,
		1,
	).Int()
	if err != nil {
		global.Logger.Errorf(context.Background(), "RedisLimiter.Take err: %v", err)
		return true
	}
	return res > 0
}

func (r RedisLimiter) AddBuckets(rules ...LimitBucketRule) ILimiter {
	for _, rule := range rules {
		if _, ok := r.redisBuckets[rule.Key]; !ok {
			r.redisBuckets[rule.Key] = rule
		}
	}
	return r
}
