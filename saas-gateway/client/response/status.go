package response

import "init.bulbul.tv/bulbul-backend/saas-gateway/client/entry"

type EmptyResponse struct {
	Status entry.Status `json:"status,omitempty"`
}
