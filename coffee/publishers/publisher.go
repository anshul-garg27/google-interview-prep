package publishers

import (
	amqp2 "coffee/amqp"
	"sync"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/spf13/viper"
)

var singletonAMPQ *amqp.Publisher
var ampqInit sync.Once
var err error

func AMQP() (*amqp.Publisher, error) {
	ampqInit.Do(func() {
		amqpURI := viper.Get("AMQP_URI").(string)
		amqpConfig := amqp2.GetConfig(amqpURI)
		singletonAMPQ, err = amqp.NewPublisher(amqpConfig, watermill.NewStdLogger(false, false))
	})
	return singletonAMPQ, err
}
func PublishMessage(jsonBytes []byte, topic string) error {
	publisher, _ := AMQP()
	err = publisher.Publish(topic, &message.Message{
		Payload: jsonBytes,
	})
	return err
}
