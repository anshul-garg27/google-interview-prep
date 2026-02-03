package amqp

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"strings"
)

func GetConfig(amqpURI string) amqp.Config {
	amqpConfig := amqp.NewDurableQueueConfig(amqpURI)
	routingKeyGen := func(topic string) string {
		return strings.Split(topic, "___")[1]
	}
	queueBindRKGen := func(queueName string) string {

		bindRoutingKey := strings.ReplaceAll(queueName, "_q", "")
		fmt.Println("Binding Queue To Routing Key", queueName, bindRoutingKey)
		return bindRoutingKey
	}
	queueGen := func(topic string) string {
		queueName := strings.Split(topic, "___")[1]
		fmt.Println("Binding Queue", queueName)
		return queueName
	}
	exchangeGen := func(topic string) string {
		exchangeName := strings.Split(topic, "___")[0]
		fmt.Println("Binding With Exchange", exchangeName)
		return exchangeName
	}
	amqpConfig.Exchange = amqp.ExchangeConfig{
		GenerateName: exchangeGen,
		Type:         "direct",
		Durable:      true,
		AutoDeleted:  false,
		Internal:     false,
		NoWait:       false,
		Arguments:    nil,
	}
	amqpConfig.Publish = amqp.PublishConfig{
		GenerateRoutingKey: routingKeyGen,
		Mandatory:          false,
		Immediate:          false,
		Transactional:      false,
		ChannelPoolSize:    0,
		ConfirmDelivery:    false,
	}
	amqpConfig.Queue = amqp.QueueConfig{
		GenerateName: queueGen,
		Durable:      true,
		AutoDelete:   false,
		Exclusive:    false,
		NoWait:       false,
		Arguments:    nil,
	}
	amqpConfig.QueueBind = amqp.QueueBindConfig{
		GenerateRoutingKey: queueBindRKGen,
		NoWait:             false,
		Arguments:          nil,
	}
	amqpConfig.Marshaler = CustomMarshaler{}
	return amqpConfig
}
