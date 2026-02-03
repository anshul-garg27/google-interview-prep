package listeners

import (
	amqp2 "coffee/amqp"
	campaignprofileListners "coffee/app/campaignprofiles/listeners"
	collectionanalytics "coffee/app/collectionanalytics/listeners"
	socialdiscoveryListners "coffee/app/discovery/listeners"
	keywordcollection "coffee/app/keywordcollection/listeners"
	"coffee/app/leaderboard"
	postcollectionlisteners "coffee/app/postcollection/listeners"
	profilecollectionlisteners "coffee/app/profilecollection/listeners"
	"coffee/listeners/middlewares"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/spf13/viper"
)

func SetupListeners() {
	amqpURI := viper.Get("AMQP_URI").(string)
	amqpConfig := amqp2.GetConfig(amqpURI)
	logger := watermill.NewStdLogger(false, false)
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}
	router.AddMiddleware(
		middlewares.MessageApplicationContext,
		middlewares.TransactionSessionHandler,
		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
			Logger:          logger,
			OnRetryHook:     OnRetryHook,
		}.Middleware,

		// Recoverer handles panics from handlers.
		// In this case, it passes them as errors to the Retry middleware.
		middlewares.Recoverer,
	)

	subscriber, err := amqp.NewSubscriber(
		amqpConfig,
		logger,
	)
	if err != nil {
		panic(err)
	}
	socialdiscoveryListners.SetupListeners(router, subscriber)
	profilecollectionlisteners.SetupListeners(router, subscriber)
	postcollectionlisteners.SetupListeners(router, subscriber)
	leaderboard.SetupListeners(router, subscriber)
	collectionanalytics.SetupListeners(router, subscriber)
	keywordcollection.SetupListeners(router, subscriber)
	campaignprofileListners.SetupListeners(router, subscriber)

	go func() {
		err := router.Run(context.Background())
		if err != nil {
			log.Error("Unable to link listeners", err)
			//panic(err)
		}
	}()
}

func OnRetryHook(retryNum int, delay time.Duration) {
	fmt.Println("")
}
