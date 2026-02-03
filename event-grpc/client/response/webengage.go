package response

type WebengageResponse struct {
	Response WebengageResponseResponse `json:"response,omitempty"`
}

type WebengageResponseResponse struct {
	Status string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}