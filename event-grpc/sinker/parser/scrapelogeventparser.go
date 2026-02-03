package parser

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"init.bulbul.tv/bulbul-backend/event-grpc/client/input"

	"init.bulbul.tv/bulbul-backend/event-grpc/model"
)

func GetShopifyEvent(payload []byte) *input.ShopifyEvent {
	event := &input.ShopifyEvent{}
	err := json.Unmarshal(payload, &event)
	if err == nil {
		event.InsertTimestamp = time.Now()
		return event
	}
	return nil
}

func GetPostLogEvent(input []byte) *model.PostLogEvent {
	rawEvent := map[string]interface{}{}
	event := &model.PostLogEvent{}
	requiredKeys := []string{"event_id", "source", "platform", "profile_id", "shortcode", "handle", "metrics",
		"dimensions", "publish_time", "event_timestamp"}

	err := json.Unmarshal(input, &rawEvent)
	if err == nil {

		if rawEvent == nil {
			log.Printf("[ERROR] PostLogEvent: Incomplete event data: %v", rawEvent)
			return nil
		}
		for _, key := range requiredKeys {
			if value, ok := rawEvent[key]; ok {
				switch key {
				case "event_id":
					event.EventId = getString(value)
				case "source":
					event.Source = getString(value)
				case "platform":
					event.Platform = getString(value)
				case "profile_id":
					event.ProfileId = getString(value)
				case "shortcode":
					event.Shortcode = getString(value)
				case "handle":
					event.Handle = getString(value)
				case "metrics":
					metrics, ok := rawEvent["metrics"].(map[string]interface{})
					if !ok || len(metrics) == 0 {
						event.Metrics = map[string]interface{}{}
					} else {
						event.Metrics = metrics
					}
				case "dimensions":
					dimensions, ok := rawEvent["dimensions"].(map[string]interface{})
					if !ok || len(dimensions) == 0 {
						event.Dimensions = map[string]interface{}{}
					} else {
						event.Dimensions = dimensions
					}
				case "publish_time":
					publishTime, err := time.Parse("2006-01-02 15:04:05", getString(value))
					if err == nil {
						event.PublishTime = publishTime
					}
				case "event_timestamp":
					eventTimestamp, err := time.Parse("2006-01-02 15:04:05", getString(value))
					if err == nil {
						event.EventTimestamp = eventTimestamp
					}
				}
			} else {
				log.Printf("[ERROR] PostLogEvent: Incomplete event data, missing key: %s", key)
				return nil
			}
		}
		event.InsertTimestamp = time.Now()
		return event
	}
	return nil
}

func GetProfileLogEvent(input []byte) *model.ProfileLogEvent {
	rawEvent := map[string]interface{}{}
	event := &model.ProfileLogEvent{}
	requiredKeys := []string{"event_id", "source", "platform", "profile_id", "handle", "metrics",
		"dimensions", "event_timestamp"}

	err := json.Unmarshal(input, &rawEvent)
	if err == nil {

		if rawEvent == nil {
			log.Printf("[ERROR] ProfileLogEvent: Incomplete event data: %v", rawEvent)
			return nil
		}

		for _, key := range requiredKeys {
			if value, ok := rawEvent[key]; ok {
				switch key {
				case "event_id":
					event.EventId = getString(value)
				case "source":
					event.Source = getString(value)
				case "platform":
					event.Platform = getString(value)
				case "profile_id":
					event.ProfileId = getString(value)
				case "handle":
					event.Handle = getString(value)
				case "metrics":
					metrics, ok := rawEvent["metrics"].(map[string]interface{})
					if !ok || len(metrics) == 0 {
						event.Metrics = map[string]interface{}{}
					} else {
						event.Metrics = metrics
					}
				case "dimensions":
					dimensions, ok := rawEvent["dimensions"].(map[string]interface{})
					if !ok || len(dimensions) == 0 {
						event.Dimensions = map[string]interface{}{}
					} else {
						event.Dimensions = dimensions
					}
				case "event_timestamp":
					eventTimestamp, err := time.Parse("2006-01-02 15:04:05", getString(value))
					if err == nil {
						event.EventTimestamp = eventTimestamp
					}
				}
			} else {
				log.Printf("[ERROR] ProfileLogEvent: Incomplete event data, missing key: %s", key)
				return nil
			}
		}
		event.InsertTimestamp = time.Now()
		return event
	}
	return nil
}

func GetProfileRelationshipLogEvent(input []byte) *model.ProfileRelationshipLogEvent {
	rawEvent := map[string]interface{}{}
	event := &model.ProfileRelationshipLogEvent{}
	requiredKeys := []string{"event_id", "source", "relationship_type", "platform", "source_profile_id",
		"target_profile_id", "source_metrics", "target_metrics", "source_dimensions", "target_dimensions",
		"event_timestamp"}

	err := json.Unmarshal(input, &rawEvent)
	if err == nil {

		if rawEvent == nil {
			log.Printf("[ERROR] ProfileRelationshipLogEvent: Incomplete event data: %v", rawEvent)
			return nil
		}

		for _, key := range requiredKeys {
			if value, ok := rawEvent[key]; ok {
				switch key {
				case "event_id":
					event.EventId = getString(value)
				case "source":
					event.Source = getString(value)
				case "relationship_type":
					event.RelationshipType = getString(value)
				case "platform":
					event.Platform = getString(value)
				case "source_profile_id":
					event.SourceProfileId = getString(value)
				case "target_profile_id":
					event.TargetProfileId = getString(value)
				case "source_metrics":
					metrics, ok := rawEvent["source_metrics"].(map[string]interface{})
					if !ok || len(metrics) == 0 {
						event.SourceMetrics = map[string]interface{}{}
					} else {
						event.SourceMetrics = metrics
					}
				case "target_metrics":
					metrics, ok := rawEvent["target_metrics"].(map[string]interface{})
					if !ok || len(metrics) == 0 {
						event.TargetMetrics = map[string]interface{}{}
					} else {
						event.TargetMetrics = metrics
					}
				case "source_dimensions":
					dimensions, ok := rawEvent["source_dimensions"].(map[string]interface{})
					if !ok || len(dimensions) == 0 {
						event.SourceDimensions = map[string]interface{}{}
					} else {
						event.SourceDimensions = dimensions
					}
				case "target_dimensions":
					dimensions, ok := rawEvent["target_dimensions"].(map[string]interface{})
					if !ok || len(dimensions) == 0 {
						event.TargetDimensions = map[string]interface{}{}
					} else {
						event.TargetDimensions = dimensions
					}
				case "event_timestamp":
					eventTimestamp, err := time.Parse("2006-01-02 15:04:05", getString(value))
					if err == nil {
						event.EventTimestamp = eventTimestamp
					}
				}
			} else {
				log.Printf("[ERROR] ProfileRelationshipLogEvent: Incomplete event data, missing key: %s", key)
				return nil
			}
		}
		event.InsertTimestamp = time.Now()
		return event
	}
	return nil
}

func GetScrapeRequestLogEvent(input []byte) *model.ScrapeRequestLogEvent {
	rawEvent := map[string]interface{}{}
	event := &model.ScrapeRequestLogEvent{}
	requiredKeys := []string{"event_id", "scl_id", "flow", "platform", "params", "status",
		"priority", "reason", "event_timestamp", "picked_at", "expires_at"}

	err := json.Unmarshal(input, &rawEvent)
	if err == nil {

		if rawEvent == nil {
			log.Printf("[ERROR] ScrapeRequestLogEvent: Incomplete event data: %v", rawEvent)
			return nil
		}

		for _, key := range requiredKeys {
			if value, ok := rawEvent[key]; ok {
				switch key {
				case "event_id":
					event.EventId = getString(value)
				case "scl_id":
					event.SclId, _ = strconv.ParseUint(getString(value), 10, 64)
				case "flow":
					event.Flow = getString(value)
				case "platform":
					event.Platform = getString(value)
				case "status":
					event.Status = getString(value)
				case "priority":
					priorityValue, ok := value.(float64)
					if !ok {
						event.Priority = nil
					} else {
						priority := int32(priorityValue)
						eventPriority := &priority
						event.Priority = eventPriority
					}
				case "reason":
					event.Reason = getString(value)
				case "params":
					params, ok := rawEvent["params"].(map[string]interface{})
					if !ok || len(params) == 0 {
						event.Params = map[string]interface{}{}
					} else {
						event.Params = params
					}
				case "picked_at":
					if value == nil {
						event.PickedAt = nil
					} else {
						pickedAt, err := time.Parse("2006-01-02 15:04:05", getString(value))
						if err == nil {
							event.PickedAt = &pickedAt
						}
					}
				case "expires_at":
					if value == nil {
						event.ExpiresAt = nil
					} else {
						expiresAt, err := time.Parse("2006-01-02 15:04:05", getString(value))
						if err == nil {
							event.ExpiresAt = &expiresAt
						}
					}
				case "event_timestamp":
					eventTimestamp, err := time.Parse("2006-01-02 15:04:05", getString(value))
					if err == nil {
						event.EventTimestamp = eventTimestamp
					}
				}
			} else {
				log.Printf("[ERROR] ScrapeRequestLogEvent: Incomplete event data, missing key: %s", key)
				return nil
			}
		}
		return event
	}
	return nil
}

func GetOrderLogEvent(input []byte) *model.OrderLogEvent {
	rawEvent := map[string]interface{}{}
	event := &model.OrderLogEvent{}
	requiredKeys := []string{"event_id", "source", "platform", "platform_order_id", "store", "metrics",
		"dimensions", "event_timestamp"}

	err := json.Unmarshal(input, &rawEvent)
	if err == nil {

		if rawEvent == nil {
			log.Printf("[ERROR] OrderLogEvent: Incomplete event data: %v", rawEvent)
			return nil
		}

		for _, key := range requiredKeys {
			if value, ok := rawEvent[key]; ok {
				switch key {
				case "event_id":
					event.EventId = getString(value)
				case "source":
					event.Source = getString(value)
				case "platform":
					event.Platform = getString(value)
				case "platform_order_id":
					event.PlatformOrderId = getString(value)
				case "store":
					event.Store = getString(value)
				case "metrics":
					metrics, ok := rawEvent["metrics"].(map[string]interface{})
					if !ok || len(metrics) == 0 {
						event.Metrics = map[string]interface{}{}
					} else {
						event.Metrics = metrics
					}
				case "dimensions":
					dimensions, ok := rawEvent["dimensions"].(map[string]interface{})
					if !ok || len(dimensions) == 0 {
						event.Dimensions = map[string]interface{}{}
					} else {
						event.Dimensions = dimensions
					}
				case "event_timestamp":
					eventTimestamp, err := time.Parse("2006-01-02 15:04:05", getString(value))
					if err == nil {
						event.EventTimestamp = eventTimestamp
					}
				}
			} else {
				log.Printf("[ERROR] OrderLogEvent: Incomplete event data, missing key: %s", key)
				return nil
			}
		}
		event.InsertTimestamp = time.Now()
		return event
	}
	return nil
}

func GetSentimentLogEvent(input []byte) *model.SentimentLogEvent {
	rawEvent := map[string]interface{}{}
	event := &model.SentimentLogEvent{}
	requiredKeys := []string{"event_id", "source", "platform", "shortcode", "comment_id", "comment",
		"sentiment", "score", "metrics", "dimensions", "event_timestamp"}

	err := json.Unmarshal(input, &rawEvent)
	if err == nil {

		if rawEvent == nil {
			log.Printf("[ERROR] SentimentLogEvent: Incomplete event data: %v", rawEvent)
			return nil
		}

		for _, key := range requiredKeys {
			if value, ok := rawEvent[key]; ok {
				switch key {
				case "event_id":
					event.EventId = getString(value)
				case "source":
					event.Source = getString(value)
				case "platform":
					event.Platform = getString(value)
				case "shortcode":
					event.Shortcode = getString(value)
				case "comment_id":
					event.CommentId = getString(value)
				case "comment":
					event.Comment = getString(value)
				case "sentiment":
					event.Sentiment = getString(value)
				case "score":
					event.Score = getFloat64(value)
				case "metrics":
					metrics, ok := rawEvent["metrics"].(map[string]interface{})
					if !ok || len(metrics) == 0 {
						event.Metrics = map[string]interface{}{}
					} else {
						event.Metrics = metrics
					}
				case "dimensions":
					dimensions, ok := rawEvent["dimensions"].(map[string]interface{})
					if !ok || len(dimensions) == 0 {
						event.Dimensions = map[string]interface{}{}
					} else {
						event.Dimensions = dimensions
					}
				case "event_timestamp":
					eventTimestamp, err := time.Parse("2006-01-02 15:04:05", getString(value))
					if err == nil {
						event.EventTimestamp = eventTimestamp
					}
				}
			} else {
				log.Printf("[ERROR] SentimentLogEvent: Incomplete event data, missing key: %s", key)
				return nil
			}
		}
		event.InsertTimestamp = time.Now()
		return event
	}
	return nil
}

func GetPostActivityLogEvent(input []byte) *model.PostActivityLogEvent {
	rawEvent := map[string]interface{}{}
	event := &model.PostActivityLogEvent{}
	requiredKeys := []string{"event_id", "source", "activity_type", "platform", "actor_profile_id", "shortcode", "metrics",
		"dimensions", "publish_time", "event_timestamp"}

	err := json.Unmarshal(input, &rawEvent)
	if err == nil {

		if rawEvent == nil {
			log.Printf("[ERROR] PostActivityLogEvent: Incomplete event data: %v", rawEvent)
			return nil
		}

		for _, key := range requiredKeys {
			if value, ok := rawEvent[key]; ok {
				switch key {
				case "event_id":
					event.EventId = getString(value)
				case "source":
					event.Source = getString(value)
				case "activity_type":
					event.ActivityType = getString(value)
				case "platform":
					event.Platform = getString(value)
				case "actor_profile_id":
					event.ActorProfileId = getString(value)
				case "shortcode":
					event.Shortcode = getString(value)
				case "metrics":
					metrics, ok := rawEvent["metrics"].(map[string]interface{})
					if !ok || len(metrics) == 0 {
						event.Metrics = map[string]interface{}{}
					} else {
						event.Metrics = metrics
					}
				case "dimensions":
					dimensions, ok := rawEvent["dimensions"].(map[string]interface{})
					if !ok || len(dimensions) == 0 {
						event.Dimensions = map[string]interface{}{}
					} else {
						event.Dimensions = dimensions
					}
				case "publish_time":
					publishTime, err := time.Parse("2006-01-02 15:04:05", getString(value))
					if err == nil {
						event.PublishTime = publishTime
					}
				case "event_timestamp":
					eventTimestamp, err := time.Parse("2006-01-02 15:04:05", getString(value))
					if err == nil {
						event.EventTimestamp = eventTimestamp
					}
				}
			} else {
				log.Printf("[ERROR] PostActivityLogEvent: Incomplete event data, missing key: %s", key)
				return nil
			}
		}
		event.InsertTimestamp = time.Now()
		return event
	}
	return nil
}
