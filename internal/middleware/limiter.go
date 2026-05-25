package middleware

import (
	"github.com/JunLang-7/blog-service/pkg/app"
	"github.com/JunLang-7/blog-service/pkg/errcode"
	"github.com/JunLang-7/blog-service/pkg/limiter"
	"github.com/gin-gonic/gin"
)

func RateLimiter(l limiter.ILimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !l.Take(l.Key(c)) {
			response := app.NewResponse(c)
			response.ToErrorResponse(errcode.TooManyRequests)
			c.Abort()
			return
		}
		c.Next()
	}
}
