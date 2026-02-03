package pg

import (
	"github.com/joho/godotenv"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/model"
	"log"
	"testing"
	"time"
)

func TestPG(t *testing.T) {
	godotenv.Load("../.env")

	pg := PG(config.New())

	if pg == nil {
		log.Fatal("Nil pg connection")
	} else {
		log.Printf("PG %v", pg)
	}
}

func TestPG2(t *testing.T) {
	godotenv.Load("../.env")

	pg := PG(config.New())

	if pg == nil {
		log.Fatal("Nil pg connection")
	} else {
		log.Printf("PG %v", pg)

		event := model.Event{
			EventId:              "hdfbvhjdbvhjdbvjhdbf",
			EventName:            "vhjbfvjdbvhjdbjh",
			EventTimestamp:       time.Now(),
			InsertTimestamp:      time.Now(),
			EventParams: map[string]interface{}{
				"kjbdvjh": "vkhjdbvjhd",
				"fkjvndfkjvd": "jfhvjh",
			},
		}

		result := pg.Create(&event)

		log.Println(result)
	}
}
