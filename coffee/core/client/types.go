package client

type SFStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Key     string `json:"key"`
	Type    string `json:"type"`
	Count   int    `json:"count"`
}

type SFBaseResponse struct {
	Status SFStatus `json:"status"`
}
