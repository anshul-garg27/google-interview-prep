package entry

type CoverImageSource string

const (
	AUTO    CoverImageSource = "AUTO"
	MANUAL  CoverImageSource = "MANUAL"
	SOCIAL_NETWORK   CoverImageSource = "SOCIAL_NETWORK"
)

type InfluencerProgramCampaign struct {
	Id                    int                    `json:"id,omitempty"`
	IdempotencyKey        string                 `json:"idempotencyKey,omitempty"`
	HostId                int                    `json:"hostId,omitempty"`
	HostClientId          string                 `json:"hostClientId,omitempty"`
	SocialNetworkType     SocialNetworkType      `json:"socialNetworkType,omitempty"`
	IsExternal            bool                   `json:"isExternal,omitempty"`
	SocialStreamId        *int                   `json:"socialStreamId,omitempty"`
	SocialStream          *SocialStream          `json:"socialStream,omitempty"`
	MarketingCampaignId   *int                   `json:"marketingCampaignId,omitempty"`
	SocialNetworkVideoUrl *string                `json:"socialNetworkVideoUrl,omitempty"`
	Status                string                 `json:"status,omitempty"`
	LiveDate              int64                  `json:"liveDate,omitempty"`
	Meta                  map[string]interface{} `json:"meta,omitempty"`
	BuyOnlineLink         *string                `json:"buyOnlineLink,omitempty"`
	BuyFromWhatsappLink   *string                `json:"buyFromWhatsappLink,omitempty"`
}

type SocialStream struct {
	Id              int                   `json:"id,omitempty"`
	HostId          int                   `json:"hostId,omitempty"`
	HostClientId    string                `json:"hostClientId,omitempty"`
	TargetClientId  string                `json:"targetClientId,omitempty"`
	Code            *string               `json:"code,omitempty"`
	Slug            *string               `json:"slug,omitempty"`
	Title           string                `json:"title,omitempty"`
	Description     *string               `json:"description,omitempty"`
	Status          string                `json:"status,omitempty"`
	StreamingInfo   StreamingInfo         `json:"streamingInfo,omitempty"`
	Products        []SocialStreamProduct `json:"products,omitempty"`
	VideoProductIds []int                 `json:"videoProductIds,omitempty"`

	TargetMonth       *int   `json:"targetMonth,omitempty"`
	TargetWeek        *int   `json:"targetWeek,omitempty"`
	TentativeLiveDate *int64 `json:"tentativeLiveDate,omitempty"`
	ActualLiveDate    *int64 `json:"actualLiveDate,omitempty"`

	ProducerId         *int     `json:"producerId,omitempty"`
	TrackingLabel      *string  `json:"trackingLabel,omitempty"`
	SharedCatalogIds   []int    `json:"sharedCatalogIds,omitempty"`
	SeoTags            []string `json:"seoTags,omitempty"`
	SuggestionTemplate *string  `json:"suggestionTemplate,omitempty"`
	IntuitionCatalogId *int     `json:"intuitionCatalogId,omitempty"`

	CategoryIds           []int    `json:"categoryIds,omitempty"`
	SubcategoryIds        []int    `json:"subcategoryIds,omitempty"`
	ProductTypeIds        []int    `json:"productTypeIds,omitempty"`
	SpokenLanguages       []string `json:"spokenLanguages,omitempty"`
	TargetAudienceGenders []string `json:"targetAudienceGenders,omitempty"`
	DestinationStores     []string `json:"destinationStores,omitempty"`
	DestinationChannels   []string `json:"destinationChannels,omitempty"`

	Label              *string `json:"label,omitempty"`
	AffiliatePartnerId *int    `json:"affiliatePartnerId,omitempty"`

	CustomerAppShareLinkId *int    `json:"customerAppShareLinkId,omitempty"`
	CustomerAppShareLink   *string `json:"customerAppShareLink,omitempty"`

	IsInfluencerProgramVideo   bool                        `json:"isInfluencerProgramVideo,omitempty"`
	InfluencerProgramCampaigns []InfluencerProgramCampaign `json:"influencerProgramCampaigns,omitempty"`

	ContentCategoryL1Ids        []int    `json:"contentCategoryL1Ids,omitempty"`
	ContentCategoryL2Ids        []int    `json:"contentCategoryL2Ids,omitempty"`
}

type SocialStreamProduct struct {
	Id             int        `json:"id,omitempty"`
	ProductId      int        `json:"productId,omitempty"`
	SkuId          int        `json:"skuId,omitempty"`
	OrderLineId    string     `json:"orderLineId,omitempty"`
	OrderId        string     `json:"orderId,omitempty"`
	StoreOrderId   string     `json:"storeOrderId,omitempty"`
	Rank           int        `json:"rank,omitempty"`
	ShipmentStatus *string    `json:"shipmentStatus,omitempty"`
	IsPartOfVideo  bool       `json:"isPartOfVideo,omitempty"`
	ConfirmedOn    *int64     `json:"confirmedOn,omitempty"`
	CancelledOn    *int64     `json:"cancelledOn,omitempty"`
	DeliveredOn    *int64     `json:"deliveredOn,omitempty"`
	SkuThumbnailId *int       `json:"skuThumbnailId,omitempty"`
	SkuThumbnail   *AssetInfo `json:"skuThumbnail,omitempty"`
}

type StreamingInfo struct {
	VideoOrientation string            `json:"videoOrientation,omitempty"`
	StreamingMode    string            `json:"streamingMode,omitempty"`
	Provider         SocialNetworkType `json:"provider,omitempty"`
	CoverImageSource CoverImageSource  `json:"coverImageSource,omitempty"`
	CoverImageId     *int              `json:"coverImageId,omitempty"`
	CoverImage       *AssetInfo        `json:"coverImage,omitempty"`
	VideoId          *int              `json:"videoId,omitempty"`
	Video            *AssetInfo        `json:"video,omitempty"`
	VideoUrl         *string           `json:"videoUrl,omitempty"`
	CoverImageUrl    *string           `json:"coverImageUrl,omitempty"`
	ChannelId        *string           `json:"channelId,omitempty"`
}