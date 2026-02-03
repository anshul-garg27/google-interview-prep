package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"init.bulbul.tv/bulbul-backend/saas-gateway/config"
	"init.bulbul.tv/bulbul-backend/saas-gateway/logger"
	"init.bulbul.tv/bulbul-backend/saas-gateway/metrics"
	"init.bulbul.tv/bulbul-backend/saas-gateway/router"
	"init.bulbul.tv/bulbul-backend/saas-gateway/sentry"
)

func main() {

	config := config.New()
	sentry.Setup(config)
	logger.Setup(config)
	metrics.SetupCollectors(config)
	if !strings.EqualFold(config.Env, "PRODUCTION") || strings.EqualFold(os.Getenv("DEBUG"), "1") {
		go func() {
			log.Println(http.ListenAndServe(":6069", nil))
		}()
	}

	err := router.SetupRouter(config).Run(":" + config.Port)
	if err != nil {
		log.Fatal("Error starting server")
	}
}
