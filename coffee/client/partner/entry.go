package partnerservice

import "coffee/app/partnerusage/domain"

type PartnerResponse struct {
	Status     Status                `json:"status"`
	ServerTime int64                 `json:"serverTime"`
	Partners   []domain.PartnerInput `json:"partners"`
	Success    bool                  `json:"success"`
}
type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Key     string `json:"key"`
	Type    string `json:"type"`
	Count   int    `json:"count"`
}
