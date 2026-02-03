package middleware

import (
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	eventGrpcContext "init.bulbul.tv/bulbul-backend/event-grpc/context"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

func EventGrpcContextToContextMiddleware(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "sinkContext", eventGrpcContext.New(c, config)))
		c.Next()
	}
}

