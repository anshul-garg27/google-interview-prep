package input

type Filters struct {
	Filters        []Filter `json:"filters,omitempty"`
	OrderBy        *string  `json:"orderBy,omitempty"`
	OrderDirection *string  `json:"orderDirection,omitempty"`
	Limit          int      `json:"limit,omitempty"`
}

type Filter struct {
	FilterType string `json:"filterType,omitempty"`
	Field      string `json:"field,omitempty"`
	Value      string `json:"value,omitempty"`
}
