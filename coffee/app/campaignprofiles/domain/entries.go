package domain

import "coffee/app/discovery/domain"

var CampaignProfileSearchFilterMap = map[string]string{

	"search_phrase": "campaign_profiles.admin_details.name,campaign_profiles.user_details.name",

	"creatorCohorts":  "campaign_profiles.admin_details.creatorCohorts",
	"creatorPrograms": "campaign_profiles.admin_details.creatorPrograms",
	"categoryIds":     "campaign_profiles.admin_details.campaignCategoryIds",
	"isBlacklisted":   "campaign_profiles.admin_details.isBlacklisted",

	"shortlistId":        "profile_collection_item.shortlist_id",
	"shortlistingStatus": "profile_collection_item.shortlisting_status",
	"campaignProfileId":  "profile_collection_item.campaign_profile_id",
	"platform":           "profile_collection_item.platform",

	"itemEnabled": "profile_collection_item.enabled",
	"itemHidden":  "profile_collection_item.hidden",

	"campaignId":        "profile_collection.campaign_id",
	"collectionId":      "profile_collection.id",
	"collectionShareId": "profile_collection.share_id",
}

type CampaignProfileAudienceInsightsEntry struct {
	Id                int64                                 `json:"id"`
	Platform          string                                `json:"platform"`
	PlatformAccountId int64                                 `json:"platformAccountId"`
	Handle            *string                               `json:"handle,omitempty"`
	UpdatedAt         string                                `json:"updatedAt,omitempty"`
	AudienceData      domain.SocialProfileAudienceInfoEntry `json:"audienceData,omitempty"`
}
