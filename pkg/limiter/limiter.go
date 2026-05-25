package limiter

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

type ILimiter interface {
	Key(c *gin.Context) string
	Take(key string) bool
	AddBuckets(rules ...LimitBucketRule) ILimiter
}

type Limiter struct {
	limiterBuckets map[string]*ratelimit.Bucket
}

type LimitBucketRule struct {
	Key          string
	FillInterval time.Duration
	Capacity     int64
	Quantum      int64
}
