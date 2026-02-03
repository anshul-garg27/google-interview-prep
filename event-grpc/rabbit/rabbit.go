package rabbit

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"init.bulbul.tv/bulbul-backend/event-grpc/context"
	"log"
	"net/http/httptest"
	"strings"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/header"
	"init.bulbul.tv/bulbul-backend/event-grpc/safego"
	"init.bulbul.tv/bulbul-backend/event-grpc/transformer"
)

type RabbitConnection struct {
	connection *amqp.Connection
}

var (
	singletonRabbit  *RabbitConnection
	rabbitCloseError chan *amqp.Error
	rabbitInit       sync.Once
)

func Rabbit(config config.Config) *RabbitConnection {
	rabbitInit.Do(func() {
		rabbitConnected := make(chan bool)
		safego.GoNoCtx(func() {
			rabbitConnector(config, rabbitConnected)
		})
		select {
		case <-rabbitConnected:
		case <-time.After(5 * time.Second):
		}
	})
	return singletonRabbit
}

func connectToRabbit(config config.Config) *RabbitConnection {
	for {
		rabbituser := config.RABBIT_CONNECTION_CONFIG.RABBIT_USER
		rabbitpassword := config.RABBIT_CONNECTION_CONFIG.RABBIT_PASSWORD
		rabbithost := config.RABBIT_CONNECTION_CONFIG.RABBIT_HOST
		rabbitport := config.RABBIT_CONNECTION_CONFIG.RABBIT_PORT
		rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbituser, rabbitpassword, rabbithost, rabbitport)

		conn, err := amqp.DialConfig(rabbitURL, amqp.Config{
			Vhost:      config.RABBIT_CONNECTION_CONFIG.RABBIT_VHOST,
			ChannelMax: 200,
			Heartbeat:  time.Duration(config.RABBIT_CONNECTION_CONFIG.RABBIT_HEARTBEAT) * time.Millisecond,
		})

		if err == nil {
			return &RabbitConnection{
				connection: conn,
			}
		}

		log.Printf("Some error with RabbitConnection connection: %v", err)
		log.Printf("Trying to reconnect to RabbitMQ at %s\n", rabbitURL)
		time.Sleep(time.Second)
	}
}

func rabbitConnector(config config.Config, connected chan bool) {
	for {
		log.Printf("Connecting to rabbit\n")

		singletonRabbit = connectToRabbit(config)
		connected <- true
		rabbitCloseError = make(chan *amqp.Error)
		singletonRabbit.connection.NotifyClose(rabbitCloseError)
		<-rabbitCloseError
	}
}

func (rabbit *RabbitConnection) Publish(exchange string, routingKey string, data []byte, headers map[string]interface{}) error {
	if rabbit.connection == nil || rabbit.connection.IsClosed() {
		return errors.New("Active rabbit connection not found")
	} else {
		if channel, err := rabbit.connection.Channel(); err != nil || channel == nil {
			return errors.New("Active rabbit connection channel not found")
		} else {
			defer func() {
				channel.Close()
			}()
			return channel.Publish(
				exchange,
				routingKey,
				false,
				false,
				amqp.Publishing{
					Headers:      headers,
					ContentType:  "application/json",
					Body:         data,
					DeliveryMode: amqp.Transient,
				})
		}
	}
}

func HeaderTableForRabbitFromXHeader(xheader map[string][]string) map[string]interface{} {
	headerTable := make(map[string]interface{})
	for k, _ := range xheader {
		if xheader[k] != nil && len(xheader[k]) > 0 {
			headerTable[k] = xheader[k][0]
		}
	}
	return headerTable
}

type BufferedConsumerFunc func(delivery amqp.Delivery, c chan interface{}) bool
type ConsumerFunc func(delivery amqp.Delivery) bool

type RabbitConsumerConfig struct {
	QueueName            string
	Exchange             string
	RoutingKey           string
	RetryOnError         bool
	ErrorExchange        *string
	ErrorRoutingKey      *string
	ConsumerCount        int
	ExitCh               <-chan bool
	ConsumerFunc         ConsumerFunc
	BufferedConsumerFunc BufferedConsumerFunc
	BufferChan           chan interface{}
}

func (rabbit *RabbitConnection) InitConsumer(consumerConfig RabbitConsumerConfig) {
	for i := 0; i < consumerConfig.ConsumerCount; i++ {
		safego.GoNoCtx(func() {
			log.Println("Started consumer")
			err := rabbit.Consume(consumerConfig)

			if err != nil {
				log.Printf("Error: %v", err)
			}

			log.Println("Exited consumer")
		})
	}
}

func (rabbit *RabbitConnection) Consume(rabbitConsumerConfig RabbitConsumerConfig) error {
	if rabbit.connection == nil || rabbit.connection.IsClosed() {
		return errors.New("Active rabbit connection not found")
	}

	if channel, err := rabbit.connection.Channel(); err != nil || channel == nil {
		return errors.New("Active rabbit connection channel not found")
	} else {
		defer func() {
			channel.Close()
		}()
		// create the queue if it doesn't already exist
		_, err := channel.QueueDeclare(rabbitConsumerConfig.QueueName, true, false, false, false, nil)
		if err != nil {
			return err
		}

		// bind the queue to the routing key
		err = channel.QueueBind(rabbitConsumerConfig.QueueName, rabbitConsumerConfig.RoutingKey, rabbitConsumerConfig.Exchange, false, nil)
		if err != nil {
			return err
		}

		err = channel.Qos(1, 0, false)
		if err != nil {
			return err
		}

		msgs, err := channel.Consume(
			rabbitConsumerConfig.QueueName, // queue
			"",                             // consumer
			false,                          // auto-ack
			false,                          // exclusive
			false,                          // no-local
			false,                          // no-wait
			nil,                            // args
		)

		if err != nil {
			log.Println("Error consuming message from rabbit")
		}

		for msg := range msgs {
			log.Printf("Message received from queue - %v", rabbitConsumerConfig.QueueName)
			buffering := rabbitConsumerConfig.BufferChan != nil
			listenerResponse := false
			if buffering {
				listenerResponse = rabbitConsumerConfig.BufferedConsumerFunc(msg, rabbitConsumerConfig.BufferChan)
			} else {
				listenerResponse = rabbitConsumerConfig.ConsumerFunc(msg)
			}
			if listenerResponse {
				err = msg.Ack(false)
			} else if !rabbitConsumerConfig.RetryOnError {
				if rabbitConsumerConfig.ErrorExchange != nil && rabbitConsumerConfig.ErrorRoutingKey != nil {
					err := rabbit.Publish(*rabbitConsumerConfig.ErrorExchange, *rabbitConsumerConfig.ErrorRoutingKey, msg.Body, msg.Headers)
					if err != nil {
						log.Printf("Some error publishing message to error q: %v", err)
					}
				}
				err = msg.Ack(false)
			} else {
				log.Printf("Message processing failed from queue - %v", rabbitConsumerConfig.QueueName)

				requeue := true
				retryCount := 1
				if retryCountInterface, ok := msg.Headers[header.XRetryCount]; ok {
					retryCount = transformer.InterfaceToInt(retryCountInterface)
					if retryCount >= 2 {
						requeue = false

						if rabbitConsumerConfig.ErrorExchange != nil && rabbitConsumerConfig.ErrorRoutingKey != nil {
							err := rabbit.Publish(*rabbitConsumerConfig.ErrorExchange, *rabbitConsumerConfig.ErrorRoutingKey, msg.Body, msg.Headers)
							if err != nil {
								log.Printf("Some error publishing message to error q: %v", err)
							}
						}
					} else {
						retryCount++
					}
				}

				if msg.Headers == nil {
					msg.Headers = map[string]interface{}{}
				}

				msg.Headers[header.XRetryCount] = retryCount

				if requeue {
					err := rabbit.Publish(rabbitConsumerConfig.Exchange, rabbitConsumerConfig.RoutingKey, msg.Body, msg.Headers)
					if err != nil {
						log.Printf("Some error publishing message: %v", err)
						err = msg.Nack(false, true)
						if err != nil {
							log.Println("Some error nacking: %v", err)
						}
					} else {
						log.Printf("RePublished so ack: %v", err)
						err = msg.Ack(false)
						if err != nil {
							log.Println("Some error acking: %v", err)
						}
					}
				} else {
					err = msg.Ack(false)
					if err != nil {
						log.Println("Some error aacking: %v", err)
					}
				}
			}
		}

		log.Println("Done with all messages waiting for exit")

		<-rabbitConsumerConfig.ExitCh
		return err
	}
}

func validHeader(header string) bool {
	for _, v := range context.XHEADERS {
		if strings.EqualFold(v, header) {
			return true
		}
	}
	return false
}
func PopulateContextFromRabbitHeaderTable(headers amqp.Table) (*context.Context, error) {
	ginContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	gc := context.New(ginContext, config.New())

	for k, v := range headers {
		if vStr, ok := v.(string); ok && validHeader(k) {
			gc.XHeader.Set(k, vStr)
		} else if vInts, ok := v.([]interface{}); ok && validHeader(k) && vInts != nil && len(vInts) > 0 {
			if vStr, ok := vInts[0].(string); ok {
				gc.XHeader.Set(k, vStr)
			}
		}
	}

	return gc, nil
}
