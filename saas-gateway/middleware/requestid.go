package middleware

import (
	"github.com/gin-gonic/gin"
	"init.bulbul.tv/bulbul-backend/saas-gateway/generator"
)

func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		generator.SetNXRequestIdOnContext(c)
		c.Next()
	}
}
