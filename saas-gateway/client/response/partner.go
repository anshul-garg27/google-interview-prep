package response

import (
	"init.bulbul.tv/bulbul-backend/saas-gateway/client/entry"
)

type PartnerResponse struct {
	Status   entry.Status `json:"status,omitempty"`
	Partners []Partner    `json:"partners"`
}

type Partner struct {
	ID                   int64      `json:"id"`
	Name                 string     `json:"name"`
	Status               string     `json:"status"`
	InventorySyncEnabled bool       `json:"inventorySyncEnabled"`
	Contracts            []Contract `json:"contracts"`
}

type Contract struct {
	ID               int64  `json:"id"`
	PartnerID        int64  `json:"partnerId"`
	ContractType     string `json:"contractType"`
	Status           string `json:"status"`
	OnboardingStatus string `json:"onboardingStatus"`
	Plan             string `json:"plan"`
	StartTime        int64  `json:"startTime"`
	EndTime          int64  `json:"endTime"`
}
