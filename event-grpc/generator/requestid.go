package generator

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"init.bulbul.tv/bulbul-backend/event-grpc/header"
)

func SetNXRequestIdOnContext(c *gin.Context) string {
	if c == nil {
		return uuid.New().String()
	}

	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}

	if c.Keys[header.RequestID] == nil {
		uuid := uuid.New().String()
		c.Keys[header.RequestID] = uuid
		if c.Request != nil && c.Request.Header != nil {
			c.Request.Header.Set(header.RequestID, uuid)
		}

		if c.Writer != nil && c.Writer.Header() != nil {
			c.Writer.Header().Set(header.RequestID, uuid)
		}
		return uuid
	} else {
		return c.Keys[header.RequestID].(string)
	}
}
