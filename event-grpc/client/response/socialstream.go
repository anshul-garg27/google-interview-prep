package response

import "init.bulbul.tv/bulbul-backend/event-grpc/client/entry"

type InfluencerProgramCampaignResponse struct {
	Status    entry.Status                      `json:"status,omitempty"`
	Campaigns []entry.InfluencerProgramCampaign `json:"list,omitempty"`
}