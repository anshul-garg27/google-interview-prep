package safego

import (
	"init.bulbul.tv/bulbul-backend/event-grpc/context"
	"log"
	"runtime/debug"
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
