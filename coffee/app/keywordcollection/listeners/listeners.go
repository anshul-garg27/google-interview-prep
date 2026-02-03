package listeners

import (
	manager2 "coffee/app/keywordcollection/manager"
	"encoding/json"
	"strings"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	log "github.com/sirupsen/logrus"
)

func SetupListeners(router *message.Router, subscriber *amqp.Subscriber) {
	router.AddNoPublisherHandler(
		"update_keyword_collection_report",
		"coffee.dx___keyword_collection_out_q",
		subscriber,
		UpdateKeywordCollectionReport,
	)
}

type BeatReportPayload struct {
	JobID                  string      `json:"job_id"`
	ReportAsset            ReportAsset `json:"report_asset"`
	DemographicReportAsset ReportAsset `json:"demographic_report_asset"`
}

type ReportAsset struct {
	Path   string `json:"path"`
	Bucket string `json:"bucket"`
}

func UpdateKeywordCollectionReport(message *message.Message) error {
	payload := &BeatReportPayload{}
	if err := json.Unmarshal(message.Payload, payload); err == nil {
		if payload == nil || (*payload).JobID == "" || (*payload).ReportAsset.Bucket == "" || (*payload).ReportAsset.Path == "" || (*payload).DemographicReportAsset.Path == "" {
			log.Errorf("Invalid Collection Report Job Payload - %v", payload)
			return nil // Invalid Payload
		}
		parts := strings.Split(payload.JobID, "_")
		if len(parts) != 3 {
			log.Errorf("Invalid Collection Report Job ID - %s", payload.JobID)
			return nil // Trash message
		}
		ctx := message.Context()
		collectionId := parts[2]
		manager := manager2.CreateKeywordCollectionManager(ctx)
		collection, err := manager.Manager.FindById(ctx, collectionId)
		if err != nil || collection == nil {
			if err != nil {
				log.Errorf("Unable to process message - %s", err)
			} else {
				log.Errorf("Unable to find collection - %s", collectionId)
			}
			return nil // Collection Missing
		}
		collection.ReportBucket = &((*payload).ReportAsset.Bucket)
		collection.ReportPath = &((*payload).ReportAsset.Path)
		collection.DemographyReportPath = &((*payload).DemographicReportAsset.Path)
		_, err = manager.Update(ctx, collectionId, collection)
		if err != nil {
			log.Errorf("Unable to update collection report details - %s", err)
		}
		return nil
	}
	return nil
}
