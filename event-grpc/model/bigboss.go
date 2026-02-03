package model

import "time"

type BigBossVoteLog struct {
	Meta            JSONB     `sql:"type:jsonb"`
	Country         string    `json:"country"`
	Key             string    `json:"key"`
	Vendorcode      string    `json:"vendorcode"`
	Contestantid    string    `json:"contestantid"`
	SdcTableVersion int64     `json:"_sdc_table_version"`
	Createdat       time.Time `json:"createdat"`
	Updatedat       time.Time `json:"updatedat"`
	Identifier      string    `json:"identifier"`
	Languages       string    `json:"languages"`
	SdcReceivedAt   string    `json:"_sdc_received_at"`
	SdcSequence     int64     `json:"_sdc_sequence"`
	ID              string    `json:"_id"`
	SdcBatchedAt    string    `json:"_sdc_batched_at"`
	SdcExtractedAt  string    `json:"_sdc_extracted_at"`
}

func (BigBossVoteLog) TableName() string {
	return "bigboss_vote_log"
}
