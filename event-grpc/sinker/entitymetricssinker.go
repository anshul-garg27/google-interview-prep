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

type EntityMetric struct {
	EntityType string `json:"entityType,omitempty"`
	EntityId   string `json:"entityId,omitempty"`
	Metric     string `json:"metric,omitempty"`
	Count      int    `json:"count,omitempty"`
	Timestamp  int64  `json:"timestamp,omitempty"`
}

func SinkEntityMetricsToClickhouse(delivery amqp.Delivery) bool {

	entityMetric := &EntityMetric{}
	if err := json.Unmarshal(delivery.Body, entityMetric); err == nil {

		entityModel := &model.EntityMetric{
			EntityType: entityMetric.EntityType,
			EntityId:   entityMetric.EntityId,
			Metric:     entityMetric.Metric,
			Count:      entityMetric.Count,
			Timestamp:  time.Unix(0, entityMetric.Timestamp*int64(time.Millisecond)),
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
