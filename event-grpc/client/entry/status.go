package entry

type Status struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Key     string `json:"key,omitempty"`
	Type    string `json:"type,omitempty"`
	Count   int    `json:"count,omitempty"`
}
