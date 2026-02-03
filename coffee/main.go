package main

import (
	"coffee/config"
	"coffee/logger"
	"coffee/sentry"
	"coffee/server"
	"net/http"
)

func main() {
	config.Setup()
	sentry.Setup()
	logger.Setup()
	go func() {
		http.ListenAndServe(":9292", nil)
	}()
	server.Setup()

}
