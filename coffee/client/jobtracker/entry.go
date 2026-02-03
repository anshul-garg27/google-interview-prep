package jobtracker

import "coffee/core/domain"

type JobEntry struct {
	Id                         *int64            `json:"id,omitempty"`
	JobType                    *string           `json:"jobType,omitempty"`
	JobName                    *string           `json:"jobName,omitempty"`
	Status                     *string           `json:"status,omitempty"`
	TotalStepCount             *int              `json:"totalStepCount,omitempty"`
	CompletedStepCount         *int              `json:"completedStepCount,omitempty"`
	Remark                     *string           `json:"remark,omitempty"`
	Input                      map[string]string `json:"input,omitempty"`
	TerminationRoutingKey      *string           `json:"terminationRoutingKey,omitempty"`
	JobIdentifier              *string           `json:"jobIdentifier,omitempty"`
	InputAssetInformation      *domain.AssetInfo `json:"inputAssetInformation,omitempty"`
	OutputFileAssetInformation *domain.AssetInfo `json:"outputFileAssetInformation,omitempty"`
}

type JobBatchEntry struct {
	Id            int64          `json:"id,omitempty"`
	Job           JobEntry       `json:"jobEntry,omitempty"`
	Status        string         `json:"status,omitempty"`
	Remarks       string         `json:"remarks,omitempty"`
	JobBatchItems []JobBatchItem `json:"jobBatchItems,omitempty"`
}

type JobBatchItem struct {
	Id         int64             `json:"id,omitempty"`
	Status     string            `json:"status,omitempty"`
	InputNode  map[string]string `json:"inputNode,omitempty"`
	OutputNode map[string]string `json:"outputNode,omitempty"`
}
