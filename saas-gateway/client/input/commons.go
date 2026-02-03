package input

type SearchQuery struct {
	Filters        []SearchFilter `json:"filters,omitempty"`
	OrderBy        *string        `json:"orderBy,omitempty"`
	OrderDirection *string        `json:"orderDirection,omitempty"`
}

type SearchFilter struct {
	FilterType      string  `json:"filterType,omitempty"`
	Field           string  `json:"field,omitempty"`
	Value           string  `json:"value,omitempty"`
	Expression      *string `json:"expression"`
	ExpressionValue *string `json:"expressionValue"`
}
