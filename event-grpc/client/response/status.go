package response

import "init.bulbul.tv/bulbul-backend/event-grpc/client/entry"

type EmptyResponse struct {
	Status entry.Status `json:"status,omitempty"`
}
