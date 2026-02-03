package response

import "init.bulbul.tv/bulbul-backend/event-grpc/client/entry"

type AssetInfoResponse struct {
	Status entry.Status      `json:"status,omitempty"`
	List   []entry.AssetInfo `json:"list,omitempty"`
}

type DynamicLinkResponse struct {
	Status entry.Status        `json:"status,omitempty"`
	Links  []entry.DynamicLink `json:"links,omitempty"`
}

type GenerateUrlResponse struct {
	Asset GeneratedURL `json:"asset,omitempty"`
}

type GeneratedURL struct {
	Path string `json:"path,omitempty"`
	URL  string `json:"url,omitempty"`
}
