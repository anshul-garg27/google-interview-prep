package entry

type AssetInfo struct {
	Id       int     `json:"id,omitempty"`
	Url      *string `json:"url,omitempty"`
	Width    *int    `json:"width,omitempty"`
	Height   *int    `json:"height,omitempty"`
	Duration *int    `json:"duration,omitempty"`
}

type DynamicLink struct {
	Id          int     `json:"id,omitempty"`
	ClientId    *string `json:"clientId,omitempty"`
	Source      *string `json:"source,omitempty"`
	Identifier  *string `json:"identifier,omitempty"`
	DynamicLink *string `json:"dynamicLink,omitempty"`
}
