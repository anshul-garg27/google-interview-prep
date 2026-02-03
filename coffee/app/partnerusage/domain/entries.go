package domain

import "time"

type ContractConfiguration struct {
	Key       string   `json:"key,omitempty"`
	Values    []string `json:"values,omitempty"`
	Type      string   `json:"type,omitempty"`
	Value     int64    `json:"value,omitempty"`
	Consumed  int64    `json:"consumed"`
	Frequency string   `json:"frequency,omitempty"`
}

type ContractNamespace struct {
	Namespace      string                  `json:"namespace,omitempty"`
	Enabled        *bool                   `json:"enabled,omitempty"`
	ModulePlanType *string                 `json:"planType,omitempty"`
	Configuration  []ContractConfiguration `json:"configuration,omitempty"`
}

type Contract struct {
	ID               int                 `json:"id,omitempty"`
	PartnerID        int                 `json:"partnerId,omitempty"`
	ContractType     string              `json:"contractType,omitempty"`
	Status           string              `json:"status,omitempty"`
	OnboardingStatus string              `json:"onboardingStatus,omitempty"`
	Plan             string              `json:"plan,omitempty"`
	StartTime        int64               `json:"startTime,omitempty"`
	EndTime          int64               `json:"endTime,omitempty"`
	Configurations   []ContractNamespace `json:"configurations,omitempty"`
}

type PartnerInput struct {
	ID                   int        `json:"id,omitempty"`
	Name                 string     `json:"name,omitempty"`
	Status               string     `json:"status,omitempty"`
	InventorySyncEnabled bool       `json:"inventorySyncEnabled,omitempty"`
	Contracts            []Contract `json:"contracts,omitempty"`
}

type PartnerUsageEntry struct {
	ID          int64     `json:"id,omitempty"`
	PartnerId   int64     `json:"partnerId,omitempty"`
	Namespace   string    `json:"namespace,omitempty"`
	Key         string    `json:"key,omitempty"`
	Limit       int64     `json:"limit,omitempty"`
	Consumed    *int64    `json:"consumed,omitempty"`
	NextResetOn time.Time `json:"nextResetOn,omitempty"`
	Frequency   string    `json:"frequency,omitempty"`
	StartDate   time.Time `json:"startDate,omitempty"`
	EndDate     time.Time `json:"endDate,omitempty"`
	Enabled     *bool     `json:"enabled,omitempty"`
}

type PartnerProfileTrackEntry struct {
	ID        int64  `json:"id,omitempty"`
	PartnerId int64  `json:"partnerId,omitempty"`
	Key       string `json:"key,omitempty"`
	Enabled   *bool  `json:"endDate,omitempty"`
}

type ActivityInput struct {
	Id        int64       `json:"id,omitempty"`
	AccountId int64       `json:"accountId,omitempty"`
	PartnerId int64       `json:"partnerId,omitempty"`
	Type      string      `json:"type,omitempty"`
	Meta      interface{} `json:"meta,omitempty"`
	UpdatedAt int64       `json:"updatedAt,omitempty"`
}
