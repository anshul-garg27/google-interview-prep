package logger

import (
	"coffee/constants"
	"os"
	"time"

	"github.com/evalphobia/logrus_sentry"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Setup() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	logLevel, err := log.ParseLevel(viper.GetString(constants.LogLevelConfig))
	if err != nil {
		logLevel = log.DebugLevel
	}
	log.SetLevel(logLevel)

	if viper.GetString("ENV") != "local" {
		hook, err := logrus_sentry.NewSentryHook(viper.GetString("SENTRY_DSN"), []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
		})
		if err == nil {
			hook.Timeout = time.Minute * 1
			hook.SetEnvironment(viper.GetString("ENV"))
			hook.StacktraceConfiguration.Enable = true
			log.AddHook(hook)
		} else {
			log.Fatal("Error creating sentryhook")
		}
	}
}
