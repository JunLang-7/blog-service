package limiter

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

type MethodLimiter struct {
	*Limiter
}

func NewMethodLimiter() ILimiter {
	return MethodLimiter{Limiter: &Limiter{
		limiterBuckets: make(map[string]*ratelimit.Bucket),
	}}
}

func (m MethodLimiter) Key(c *gin.Context) string {
	uri := c.Request.RequestURI
	index := strings.Index(uri, "?")
	if index == -1 {
		return uri
	}
	return uri[:index]
}

func (m MethodLimiter) Take(key string) bool {
	bucket, ok := m.Limiter.limiterBuckets[key]
	if !ok {
		return true
	}
	return bucket.TakeAvailable(1) == 1
}

func (m MethodLimiter) AddBuckets(rules ...LimitBucketRule) ILimiter {
	for _, rule := range rules {
		if _, ok := m.limiterBuckets[rule.Key]; !ok {
			m.limiterBuckets[rule.Key] = ratelimit.NewBucketWithQuantum(rule.FillInterval, rule.Capacity, rule.Quantum)
		}
	}
	return m
}
