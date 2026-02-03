package rabbit

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
)

// fixme tests
func TestRabbit(t *testing.T) {
	godotenv.Load("../.env")

	rabbit := Rabbit(config.New())

	if rabbit.connection == nil {
		log.Fatal("No error but nil rabbit connection")
	} else if rabbit.connection.IsClosed() {
		log.Fatal("Rabbit connection closed")
	} else {
		log.Printf("Rabbit %v", rabbit)
	}
}

func TestRabbitConnection_Publish(t *testing.T) {
	godotenv.Load("../.env")

	rabbit := Rabbit(config.New())

	if rabbit.connection == nil {
		log.Fatal("No error but nil rabbit connection")
	} else if rabbit.connection.IsClosed() {
		log.Fatal("Rabbit connection closed")
	} else {
		log.Printf("Rabbit %v", rabbit)
		err := rabbit.Publish("event.dx", "TOUCH", []byte("hahahhaha"))
		if err != nil {
			log.Fatal("Error publishing")
		} else {
			log.Println("Message published")
		}
	}
}

func TestRabbitConnection_Consume(t *testing.T) {
	// fixme impl
}
