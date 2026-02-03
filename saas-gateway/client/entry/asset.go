package entry

type AssetWithThumbnail struct {
	Asset     AssetInfo  `json:"asset,omitempty"`
	Thumbnail *AssetInfo `json:"thumbnail,omitempty"`
}

type AssetInfo struct {
	Id       int     `json:"id,omitempty"`
	Url      *string `json:"url,omitempty"`
	Width    *int    `json:"width,omitempty"`
	Height   *int    `json:"height,omitempty"`
	Duration *int    `json:"duration,omitempty"`
}

type DynamicLink struct {
	Id          int     `json:"id,omitempty,omitempty"`
	ClientId    string  `json:"clientId,omitempty,omitempty"`
	Source      *string `json:"source,omitempty,omitempty"`
	Identifier  *string `json:"identifier,omitempty,omitempty"`
	DynamicLink *string `json:"dynamicLink,omitempty,omitempty"`

	Channel  *string `json:"channel,omitempty"`
	Feature  *string `json:"feature,omitempty"`
	Campaign *string `json:"campaign,omitempty"`

	Stage *string `json:"stage,omitempty"`

	Tags []string `json:"tags,omitempty"`

	CustomData map[string]interface{} `json:"customData,omitempty"`

	CanonicalIdentifier *string `json:"canonicalIdentifier,omitempty"`

	OgTitle       *string `json:"ogTitle,omitempty"`
	OgDescription *string `json:"ogDescription,omitempty"`
	OgImageUrl    *string `json:"ogImageUrl,omitempty"`
	OgVideo       *string `json:"ogVideo,omitempty"`
	OgUrl         *string `json:"ogUrl,omitempty"`

	DeeplinkPath        string  `json:"deeplinkPath,omitempty"`
	AndroidDeeplinkPath *string `json:"androidDeeplinkPath,omitempty"`
	AndroidUrl          *string `json:"androidUrl,omitempty"`
	IosDeeplinkPath     *string `json:"iosDeeplinkPath,omitempty"`
	DesktopDeeplinkPath *string `json:"desktopDeeplinkPath,omitempty"`

	MarketingTitle *string `json:"marketingTitle,omitempty"`

	DeeplinkNoAttribution bool `json:"deeplinkNoAttribution,omitempty"`

	LinkClient        string `json:"linkClient,omitempty"`
	PromptAppDownload bool   `json:"promptAppDownload,omitempty"`
	BranchType        int    `json:"branchType,omitempty"`

	DomainPrefix *string `json:"domainPrefix,omitempty"`

	ShortCode *string `json:"shortCode,omitempty"`

	DeeplinkPathWithParams *string `json:"deeplinkPathWithParams,omitempty"`
}
