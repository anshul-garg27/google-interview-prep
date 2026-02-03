package sentry

import (
	"coffee/helpers"
	"io"
	"os"
	"runtime/debug"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Setup() {
	if viper.GetString("ENV") != "local" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              viper.GetString("SENTRY_DSN"),
			Debug:            true,
			Environment:      viper.GetString("ENV"),
			AttachStacktrace: true,
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
	}
}

// for ability to test the PrintPrettyStack function
var recovererErrorWriter io.Writer = os.Stderr

func PrintToSentry(err error) {
	sentry.CaptureException(err)
}

func PrintToSentryWithInterface(rvr interface{}) {
	debugStack := debug.Stack()
	s := helpers.PrettyStack{}
	out, err := s.Parse(debugStack, rvr)
	if err == nil {
		recovererErrorWriter.Write(out)
		sentry.CaptureMessage(string(out))
	} else {
		os.Stderr.Write(debugStack)
	}
}
