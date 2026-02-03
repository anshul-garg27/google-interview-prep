package model

import "time"

type EntityMetric struct {
	EntityType string `json:"entityType,omitempty"`
	EntityId string `json:"entityId,omitempty"`
	Metric string `json:"metric,omitempty"`
	Count int `json:"count,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
