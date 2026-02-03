package parser

import (
	"encoding/json"
	"init.bulbul.tv/bulbul-backend/event-grpc/model"
	"log"
	"time"
)

func GetAffiliateOrderEvent(input []byte) *model.AffiliateOrderEvent {
	rawAffiliateOrdersEvent := map[string]interface{}{}
	err := json.Unmarshal(input, &rawAffiliateOrdersEvent)
	if err == nil {

		if rawAffiliateOrdersEvent == nil {
			log.Printf("Incomplete event data: %v", rawAffiliateOrdersEvent)
			return nil
		}
		created_on, _ := time.Parse("2006-01-02 03:04:05", rawAffiliateOrdersEvent["created_on"].(string))
		last_modified_on, _ := time.Parse("2006-01-02 03:04:05", rawAffiliateOrdersEvent["last_modified_on"].(string))
		var delivered_on time.Time
		if rawAffiliateOrdersEvent["delivered_on"] != nil {
			delivered_on, _ = time.Parse("2006-01-02 03:04:05", rawAffiliateOrdersEvent["delivered_on"].(string))
		}

		event := &model.AffiliateOrderEvent{
			OrderID:        getString(rawAffiliateOrdersEvent["order_id"]),
			Platform:       getString(rawAffiliateOrdersEvent["platform"]),
			ReferralCode:   getString(rawAffiliateOrdersEvent["referral_code"]),
			Meta:           rawAffiliateOrdersEvent["meta"].(map[string]interface{}),
			CreatedOn:      created_on,
			LastModifiedOn: last_modified_on,
			Amount:         int64(getFloat64(rawAffiliateOrdersEvent["amount"])),
			Status:         getString(rawAffiliateOrdersEvent["status"]),
			DeliveredOn:    delivered_on,
		}

		return event
	}
	return nil
}
