package safego

import (
	"runtime/debug"

	log "github.com/sirupsen/logrus"
	"init.bulbul.tv/bulbul-backend/saas-gateway/context"
)

func Go(gCtx *context.Context, f func()) {
	go func() {
		defer func() {
			if panicMessage := recover(); panicMessage != nil {
				stack := debug.Stack()

				gCtx.Logger.Error().Msgf("RECOVERED FROM UNHANDLED PANIC: %v\nSTACK: %s", panicMessage, stack)
			}
		}()

		f()
	}()
}

func GoNoCtx(f func()) {
	go func() {
		defer func() {
			if panicMessage := recover(); panicMessage != nil {
				stack := debug.Stack()

				log.Printf("RECOVERED FROM UNHANDLED PANIC: %v\nSTACK: %s", panicMessage, stack)
			}
		}()

		f()
	}()
}
