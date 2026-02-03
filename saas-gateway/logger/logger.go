package logger

import (
	"os"
	"time"

	"github.com/evalphobia/logrus_sentry"
	log "github.com/sirupsen/logrus"
	"init.bulbul.tv/bulbul-backend/saas-gateway/config"
)

func Setup(config config.Config) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	logLevel, err := log.ParseLevel("debug")
	if err != nil {
		logLevel = log.DebugLevel
	}
	log.SetLevel(logLevel)

	if config.Env != "local" {
		hook, err := logrus_sentry.NewSentryHook(config.SENTRY_DSN, []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
		})
		if err == nil {
			hook.Timeout = time.Minute * 1
			hook.SetEnvironment(config.Env)
			hook.StacktraceConfiguration.Enable = true
			log.AddHook(hook)
		} else {
			log.Fatal("Error creating sentryhook")
		}
	}
}
