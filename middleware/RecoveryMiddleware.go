package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"xmarket_gin/response"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				response.Fail(c, nil, fmt.Sprint(err))
			}
		}()
		c.Next()
	}
}
