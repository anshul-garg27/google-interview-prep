package parser

import (
	"encoding/json"
	"init.bulbul.tv/bulbul-backend/event-grpc/model"
	"log"
	"time"
)

func GetBigBossVoteLog(input []byte) *model.BigBossVoteLog {
	bigbossVote := map[string]interface{}{}
	err := json.Unmarshal(input, &bigbossVote)
	if err == nil {

		if bigbossVote == nil {
			log.Printf("Incomplete event data: %v", bigbossVote)
			return nil
		}
		created_at, _ := time.Parse("2006-01-02 03:04:05.000000-07:00", bigbossVote["createdat"].(string))
		updated_at, _ := time.Parse("2006-01-02 03:04:05.000000-07:00", bigbossVote["updatedat"].(string))

		event := &model.BigBossVoteLog{
			Meta:            bigbossVote["meta"].(map[string]interface{}),
			Country:         getString(bigbossVote["country"]),
			Key:             getString(bigbossVote["key"]),
			Vendorcode:      getString(bigbossVote["vendorcode"]),
			Contestantid:    getString(bigbossVote["contestantid"]),
			SdcTableVersion: int64(getFloat64(bigbossVote["_sdc_table_version"])),
			Createdat:       created_at,
			Updatedat:       updated_at,
			Identifier:      getString(bigbossVote["identifier"]),
			Languages:       getString(bigbossVote["languages"]),
			SdcReceivedAt:   getString(bigbossVote["_sdc_received_at"]),
			SdcSequence:     int64(getFloat64(bigbossVote["_sdc_sequence"])),
			ID:              getString(bigbossVote["_id"]),
			SdcBatchedAt:    getString(bigbossVote["_sdc_batched_at"]),
			SdcExtractedAt:  getString(bigbossVote["_sdc_extracted_at"]),
		}

		return event
	}
	return nil
}
