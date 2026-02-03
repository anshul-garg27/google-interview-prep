package amqp

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	amqp "github.com/rabbitmq/amqp091-go"
)

const DefaultMessageUUIDHeaderKey = "_watermill_message_uuid"

// deprecated, please use DefaultMessageUUIDHeaderKey instead
const MessageUUIDHeaderKey = DefaultMessageUUIDHeaderKey

type CustomMarshaler struct {
	// PostprocessPublishing can be used to make some extra processing with amqp.Publishing,
	// for example add CorrelationId and ContentType:
	//
	//  amqp.DefaultMarshaler{
	//		PostprocessPublishing: func(publishing stdAmqp.Publishing) stdAmqp.Publishing {
	//			publishing.CorrelationId = "correlation"
	//			publishing.ContentType = "application/json"
	//
	//			return publishing
	//		},
	//	}
	PostprocessPublishing func(amqp.Publishing) amqp.Publishing

	// When true, DeliveryMode will be not set to Persistent.
	//
	// DeliveryMode Transient means higher throughput, but messages will not be
	// restored on broker restart. The delivery mode of publishings is unrelated
	// to the durability of the queues they reside on. Transient messages will
	// not be restored to durable queues, persistent messages will be restored to
	// durable queues and lost on non-durable queues during server restart.
	NotPersistentDeliveryMode bool

	// Header used to store and read message UUID.
	//
	// If value is empty, DefaultMessageUUIDHeaderKey value is used.
	// If header doesn't exist, empty value is passed as message UUID.
	MessageUUIDHeaderKey string
}

func (d CustomMarshaler) Marshal(msg *message.Message) (amqp.Publishing, error) {
	headers := make(amqp.Table, len(msg.Metadata)+1) // metadata + plus uuid

	for key, value := range msg.Metadata {
		headers[key] = value
	}
	headers[d.computeMessageUUIDHeaderKey()] = msg.UUID

	publishing := amqp.Publishing{
		Body:    msg.Payload,
		Headers: headers,
	}
	if !d.NotPersistentDeliveryMode {
		publishing.DeliveryMode = amqp.Persistent
	}

	if d.PostprocessPublishing != nil {
		publishing = d.PostprocessPublishing(publishing)
	}

	return publishing, nil
}

func (d CustomMarshaler) Unmarshal(amqpMsg amqp.Delivery) (*message.Message, error) {
	msgUUIDStr, err := d.unmarshalMessageUUID(amqpMsg)
	if err != nil {
		return nil, err
	}

	msg := message.NewMessage(msgUUIDStr, amqpMsg.Body)
	msg.Metadata = make(message.Metadata, len(amqpMsg.Headers)-1) // headers - minus uuid

	for key, value := range amqpMsg.Headers {
		if key == d.computeMessageUUIDHeaderKey() {
			continue
		}

		finalValue := ""
		switch t := value.(type) {
		case []interface{}:
			values, _ := value.([]interface{})
			finalValue = values[0].(string)
		case string:
			finalValue, _ = value.(string)
		default:
			log.Printf("Ignoring Unsupported type: %T\n", t)
		}
		log.Println(key, value)
		log.Println(key, finalValue)
		msg.Metadata[key] = finalValue
	}

	return msg, nil
}

func (d CustomMarshaler) unmarshalMessageUUID(amqpMsg amqp.Delivery) (string, error) {
	var msgUUIDStr string

	msgUUID, hasMsgUUID := amqpMsg.Headers[d.computeMessageUUIDHeaderKey()]
	if !hasMsgUUID {
		return "", nil
	}

	msgUUIDStr, hasMsgUUID = msgUUID.(string)
	if !hasMsgUUID {
		return "", errors.Errorf("message UUID is not a string, but: %#v", msgUUID)
	}

	return msgUUIDStr, nil
}

func (d CustomMarshaler) computeMessageUUIDHeaderKey() string {
	if d.MessageUUIDHeaderKey != "" {
		return d.MessageUUIDHeaderKey
	}

	return DefaultMessageUUIDHeaderKey
}
