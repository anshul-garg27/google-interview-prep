package middlewares

import (
	"coffee/core/appcontext"
	"github.com/ThreeDotsLabs/watermill/message"
	log "github.com/sirupsen/logrus"
)

func MessageApplicationContext(h message.HandlerFunc) message.HandlerFunc {
	return func(message *message.Message) ([]*message.Message, error) {
		log.Println("Middleware for message header to context enrichment")
		message, _ = appcontext.CreateApplicationContextPubSub(message)
		return h(message)
	}
}
