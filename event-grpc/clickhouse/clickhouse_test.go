package clickhouse

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mssola/useragent"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/model"
	"log"
	"testing"
	"time"
)

func TestClickhouse(t *testing.T) {
	godotenv.Load("../.env")

	clickhouse := Clickhouse(config.New(), nil)

	if clickhouse == nil {
		log.Fatal("Nil clickhouse connection")
	} else {
		log.Printf("Clickhouse %v", clickhouse)
	}
}

func TestClickhouse3(t *testing.T) {

	godotenv.Load("../.env")

	//var launchReferEvent model.LaunchReferEvent
	//Clickhouse(config.New(), nil).Find(&launchReferEvent, "device_id = ?", "D4A8FACC-BAD1-4A22-BBC4-5F5068A7C2F1")
	//log.Printf("Clickhouse %v", launchReferEvent)

	dbName := "dbt"
	var count int64
	//var stgClickEvents []model.StgClickEvent
	db := Clickhouse(config.New(), &dbName)
	//db.Find(&stgClickEvents, "client_ip = ? and date(event_timestamp) = ?", "157.35.58.121", time.Now().Format("2006-01-02")).Count(&count)

	db.Raw("select count() from stg_click_events where client_ip = ? and date(event_timestamp) = ?", "152.58.159.154", time.Now().Format("2006-01-02")).Scan(&count)
	//db.Raw("select date(event_timestamp) d, count() from stg_click_events where event_timestamp >= '2023-05-27' and client_ip='157.35.58.121' and bot=0 group by d having count() > 200", nil).Count(&count)

	log.Printf("Clickhouse %d", count)
}

func TestClickhouse4(t *testing.T) {

	godotenv.Load("../.env")

	var launchReferEvents []model.LaunchReferEvent
	Clickhouse(config.New(), nil).
		Where("device_id = ? OR bb_device_id = ?", "D4A8FACC-BAD1-4A22-BBC4-5F5068A7C2F1", "32954653").
		First(&launchReferEvents)

	log.Printf("Clickhouse %v", launchReferEvents)
}

func TestClickhouse2(t *testing.T) {
	godotenv.Load("../.env")

	clickhouse := Clickhouse(config.New(), nil)

	if clickhouse == nil {
		log.Fatal("Nil clickhouse connection")
	} else {
		log.Printf("Clickhouse %v", clickhouse)

		//event1 := model.Event{
		//	EventId:              "hdfbvhjdbvhjdbvjhdbf",
		//	EventName:            "vhjbfvjdbvhjdbjh",
		//	EventTimestamp:       time.Now(),
		//	InsertTimestamp:      time.Now(),
		//	PpId: "fdss",
		//	EventParams: map[string]interface{}{
		//		"kjbdvjh": "vkhjdbvjhd",
		//		"fkjvndfkjvd": "jfhvjh",
		//	},
		//}
		//
		//event2 := model.Event{
		//	EventId:              "skgbhv",
		//	EventName:            "vhjbfvjdbvhjdbjh",
		//	EventTimestamp:       time.Now(),
		//	InsertTimestamp:      time.Now(),
		//	EventParams: map[string]interface{}{
		//		"kjbdvjh": "vkhjdbvjhd",
		//		"fkjvndfkjvd": "jfhvjh",
		//	},
		//}
		//
		//events := []model.Event{event1, event2}
		//result := clickhouse.Create(events)

		//abLog := model.AbLog{
		//	Name:       "jhghgkhg",
		//	Timestamp:  time.Now(),
		//	Salt:       "jhvfjghsfd",
		//	Event:      "dkshfgkhgsfvsjkhgb",
		//	Checksum:   "dfjhhsvfsjhbvd",
		//	Params: map[string]interface{}{
		//			"kjbdvjh": "vkhjdbvjhd",
		//			"fkjvndfkjvd": "jfhvjh",
		//	},
		//}
		//
		//result := clickhouse.Create(&abLog)

		//em := model.EntityMetric{
		//	EntityType: "kdgfhskghdfv",
		//	EntityId:   "kiguhfvdfvgfk",
		//	Metric:     "jhgksdhsfkghs",
		//	Count:      10,
		//	Timestamp:  time.Now(),
		//}
		//result := clickhouse.Create(&em)

		//aba := model.AbAssignment{
		//	Timestamp:   time.Now(),
		//	UserId:      0,
		//	BBDeviceId:  "12313",
		//	Checksum:    234234234234,
		//	Assignments: map[string]interface{}{
		//			"kjbdvjh": "vkhjdbvjhd",
		//			"fkjvndfkjvd": "jfhvjh",
		//	},
		//}

		//graphyEvent := model.GraphyEvent{
		//	Id:              uuid.New().String(),
		//	InsertTimestamp: time.Now(),
		//	EventName:       "hshshs",
		//	EventData: map[string]interface{}{
		//		"abc": "xyz",
		//	},
		//	Headers:         map[string]interface{}{
		//		"abc": "xyz",
		//	},
		//}
		//
		//result := clickhouse.Create(&graphyEvent)

		//input := "{\"hostName\":\"t1-backend-2\",\"serviceName\":\"identity\",\"timestamp\":{\"nano\":935000000,\"epochSecond\":1662808346},\"method\":\"GET\",\"uri\":\"https://acl.mgapis.com/ekyc-ms/v1/payoutDetails/60f3e9de55daad00144f3e2f?vendorCode=gcc\",\"headers\":{\"Accept\":[\"application/json\"],\"Content-Type\":[\"application/json\"],\"apikey\":[\"75d7b062bfbe795fb484530fdab9030e\"],\"Content-Length\":[\"0\"]},\"remoteAddress\":\"acl.mgapis.com\",\"responseStatus\":200,\"requestBody\":\"\",\"responseBody\":\"{\\\"name\\\":\\\"OK\\\",\\\"status\\\":true,\\\"code\\\":200,\\\"message\\\":\\\"Payout details not found\\\",\\\"data\\\":{\\\"message\\\":\\\"Payout details not found\\\"}}\",\"onlyRequest\":false}\n"
		//traceLogEvent := parser.GetTraceLogEvent([]byte(input))
		//
		//result := clickhouse.Create(&traceLogEvent)

		//webengageEvent := model.WebengageEvent{
		//	Id:          "fihdfhidfsh",
		//	EventName:   "bhgfvh",
		//	AccountType: "hbdfvbhfdv",
		//	Timestamp:   time.Time{},
		//	UserId:      "",
		//	AnonymousId: "",
		//	Data:        nil,
		//}
		//result := clickhouse.Create(&webengageEvent)

		abAssignment := model.AbAssignment{
			Timestamp:  time.Time{},
			UserId:     123,
			BBDeviceId: "djsvd",
			NewUser:    1,
			Checksum:   "jhdfbvjdhfb",
			Ver:        0,
			Assignments: map[string]interface{}{
				"true": true,
			},
		}
		result := clickhouse.Create(&abAssignment)

		log.Println(result)
	}
}

func TestClickhouse5(t *testing.T) {
	if ua := useragent.New("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"); ua != nil {
		fmt.Println(ua)
	}
}
