package dam

type GenerateUrlReqEntry struct {
	Usage     string `json:"usage,omitempty"`
	MediaType string `json:"type,omitempty"`
}

type GenerateUrlResponse struct {
	Asset GenerateUrlRespEntry `json:"asset,omitempty"`
}

type GenerateUrlRespEntry struct {
	Path string `json:"path,omitempty"`
	URL  string `json:"url,omitempty"`
}
