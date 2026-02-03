package middlewares

import (
	"coffee/constants"
	"coffee/core/appcontext"
	"coffee/core/persistence/postgres"
	"github.com/ThreeDotsLabs/watermill/message"
	log "github.com/sirupsen/logrus"
)

func TransactionSessionHandler(h message.HandlerFunc) message.HandlerFunc {
	return func(message *message.Message) ([]*message.Message, error) {
		ctx := message.Context()
		log.Println("Middleware for listener message txn session")
		session := postgres.InitializePgSession(ctx)
		reqContext := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
		reqContext.SetSession(ctx, session)
		msgs, err := h(message)
		if err != nil {
			session.Rollback(ctx)
		} else {
			session.Commit(ctx)
		}
		return msgs, err
	}
}
