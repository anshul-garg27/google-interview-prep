package sinker

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"init.bulbul.tv/bulbul-backend/event-grpc/clickhouse"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"init.bulbul.tv/bulbul-backend/event-grpc/model"
	"log"
	"time"
)

type AbAssignment struct {
	Timestamp   int64                  `json:"timestamp,omitempty"`
	UserId      *int                   `json:"userId,omitempty"`
	BBDeviceId  *string                `json:"bbDeviceId,omitempty"`
	NewUser     bool                   `json:"newUser,omitempty"`
	Version     *string                `json:"version,omitempty"`
	Ver         int64                  `json:"ver,omitempty"`
	Checksum    string                 `json:"checksum,omitempty"`
	Assignments map[string]interface{} `json:"assignments,omitempty"`
}

func SinkABAssignmentsToClickhouse(delivery amqp.Delivery) bool {

	abAssignment := &AbAssignment{}
	if err := json.Unmarshal(delivery.Body, abAssignment); err == nil {

		entityModel := &model.AbAssignment{
			Timestamp:   time.Unix(abAssignment.Timestamp, 0),
			Ver:         abAssignment.Ver,
			Checksum:    abAssignment.Checksum,
			Assignments: abAssignment.Assignments,
		}

		if abAssignment.NewUser {
			entityModel.NewUser = 1
		}

		if abAssignment.UserId != nil {
			entityModel.UserId = *abAssignment.UserId
		}

		if abAssignment.BBDeviceId != nil {
			entityModel.BBDeviceId = *abAssignment.BBDeviceId
		}

		if abAssignment.Version != nil {
			entityModel.Assignments["app_version"] = *abAssignment.Version
		}

		if result := clickhouse.Clickhouse(config.New(), nil).Create(entityModel); result != nil && result.Error == nil {
			log.Printf("Created successfully: %v", result.RowsAffected)
			return true
		} else if result == nil {
			log.Printf("Creation failed: %v", result)
		} else {
			log.Printf("Creation failed: %v", result.Error)
		}
	}

	return false
}
