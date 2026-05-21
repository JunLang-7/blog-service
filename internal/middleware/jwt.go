package middleware

import (
	"errors"

	"github.com/JunLang-7/blog-service/pkg/app"
	"github.com/JunLang-7/blog-service/pkg/errcode"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			tokenString string
			ecode       = errcode.Success
		)
		if s, exists := c.GetQuery("token"); exists {
			tokenString = s
		} else {
			tokenString = c.GetHeader("token")
		}
		if tokenString == "" {
			ecode = errcode.InvalidParams
		} else {
			_, err := app.ParseToken(tokenString)
			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					ecode = errcode.UnauthorizedTokenTimeout
				} else {
					ecode = errcode.UnauthorizedTokenError
				}
			}
		}

		if !errors.Is(ecode, errcode.Success) {
			response := app.NewResponse(c)
			response.ToErrorResponse(ecode)
			c.Abort()
			return
		}

		c.Next()
	}
}
