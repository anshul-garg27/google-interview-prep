package domain

import (
	coredomain "coffee/core/domain"
	"time"
)

type SocialProfileHashatagsEntry struct {
	ID                int64            `json:"id,omitempty"`
	PlatformProfileId int64            `json:"platformProfileId,omitempty"`
	Platform          string           `json:"platform,omitempty"`
	Hashtags          map[string]int64 `json:"hashtags,omitempty"`
}
type LocationsEntry struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	FullName string `json:"fullName,omitempty"`
	Type     string `json:"type,omitempty"`
}

type ProfileAdminDetails struct {
	ID                 *int64     `json:"ID,omitempty"`
	PlatformProfileId  int64      `json:"platformProfileId,omitempty"`
	Platform           string     `json:"platform,omitempty"`
	Gender             *string    `json:"gender,omitempty"`
	Phone              *string    `json:"phone,omitempty"`
	Name               *string    `json:"name,omitempty"`
	Email              *string    `json:"email,omitempty"`
	Label              *string    `json:"label,omitempty"`
	Languages          *[]string  `json:"languages,omitempty"`
	Categories         *[]string  `json:"categories,omitempty"`
	City               *string    `json:"city,omitempty"`
	State              *string    `json:"state,omitempty"`
	Country            *string    `json:"country,omitempty"`
	WhatsappEnabled    *bool      `json:"whatsappEnabled,omitempty"`
	IsBlacklisted      *bool      `json:"isBlacklisted,omitempty"`
	BlacklistedBy      *string    `json:"blacklistedBy,omitempty"`
	BlacklistingReason *string    `json:"blacklistingReason,omitempty"`
	Dob                *time.Time `json:"dob"`
	CreatedBy          string     `json:"createdBy"`
}

type ProfileUserDetails struct {
	ID                int64     `json:"ID,omitempty"`
	PlatformProfileId int64     `json:"platformProfileId,omitempty"`
	Platform          string    `json:"platform,omitempty"`
	Gender            *string   `json:"gender,omitempty"`
	Phone             *string   `json:"phone,omitempty"`
	Name              *string   `json:"name,omitempty"`
	Email             *string   `json:"email,omitempty"`
	Label             *string   `json:"label,omitempty"`
	Languages         *[]string `json:"languages,omitempty"`
	Categories        *[]string `json:"categories,omitempty"`
	Location          *string   `json:"location,omitempty"`
	WhatsappEnabled   *bool     `json:"whatsappEnabled,omitempty"`
	Dob               time.Time `json:"dob"`
	CreatedBy         string    `json:"createdBy"`
}

type Growth struct {
	PlatformCode         string  `json:"platformCode,omitempty"`
	Platform             string  `json:"platform,omitempty"`
	FollowerGraph        []Graph `json:"followerGraph,omitempty"`
	FollowingGraph       []Graph `json:"followingGraph,omitempty"`
	ViewGraph            []Graph `json:"viewGraph,omitempty"`
	PlaysGraph           []Graph `json:"playsGraph,omitempty"`
	UploadsGraph         []Graph `json:"uploadsGraph,omitempty"`
	EngagementGraph      []Graph `json:"engagementGraph,omitempty"`
	FollowersChangeGraph []Graph `json:"followersChangeGraph,omitempty"`
	TotalViewGraph       []Graph `json:"totalViewGraph,omitempty"`
	TotalPlaysGraph      []Graph `json:"TotalPlaysGraph,omitempty"`
}

type Content struct {
	PlatformCode string                               `json:"platformCode,omitempty"`
	Platform     string                               `json:"platform,omitempty"`
	Posts        []coredomain.SocialProfilePostsEntry `json:"posts"`
}
type AudienceAgeGender struct {
	Male   map[string]*float64 `json:"male,omitempty"`
	Female map[string]*float64 `json:"female,omitempty"`
}
type AudienceLocation struct {
	City    map[string]*float64 `json:"city,omitempty"`
	Country map[string]*float64 `json:"country,omitempty"`
}
type SocialProfileAudienceInfoEntry struct {
	PlatformProfileId              int64                      `json:"platformProfileId,omitempty"`
	Platform                       string                     `json:"platform,omitempty"`
	CityWiseAudienceLocation       map[string]*float64        `json:"cityWiseAudienceLocation"`
	CountryWiseAudienceLocation    map[string]*float64        `json:"countryWiseAudienceLocation,omitempty"`
	AudienceGenderSplit            map[string]*float64        `json:"audienceGenderSplit,omitempty"`
	AudienceAgeGenderSplit         AudienceAgeGender          `json:"audienceAgeGenderSplit,omitempty"`
	TopAge                         map[string]*float64        `json:"topAge,omitempty"`
	AudienceLanguage               map[string]*float64        `json:"audienceLanguage,omitempty"`
	NotableFollowers               []coredomain.SocialAccount `json:"notableFollowers,omitempty"`
	AudienceReachabilityPercentage *float64                   `json:"audienceReachabilityPercentage,omitempty"`
	AudienceAuthenticityPercentage *float64                   `json:"audienceAuthenticityPercentage,omitempty"`
	CommentRatePercentage          *float64                   `json:"commentRatePercentage,omitempty"`
	QualityAudiencePercentage      *float64                   `json:"qualityAudiencePercentage,omitempty"`
	QualityAudienceScore           *float64                   `json:"qualityAudienceScore,omitempty"`
	QualityScoreGrade              *string                    `json:"qualityScoreGrade,omitempty"`
	QualityScoreBreakup            map[string]*string         `json:"qualityScoreBreakup,omitempty"`
}
type InstagramApiStatus struct {
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}

var InstagramSearchFilterAdminMap = map[string]string{
	"keywords":                    "instagram_account.keywords_admin",
	"search_phrase":               "instagram_account.search_phrase_admin",
	"gender":                      "campaign_profiles.gender",
	"phone":                       "campaign_profiles.phone",
	"isBlacklisted":               "campaign_profiles.admin_details.isBlacklisted",
	"location":                    "instagram_account.location",
	"location_list":               "campaign_profiles.location",
	"languages":                   "campaign_profiles.languages",
	"dateOfBirth":                 "campaign_profiles.dob",
	"followers":                   "instagram_account.followers",
	"audience_age.13-17":          "instagram_account.audience_age.13-17",
	"audience_age.18-24":          "instagram_account.audience_age.18-24",
	"audience_age.25-34":          "instagram_account.audience_age.25-34",
	"audience_age.35-44":          "instagram_account.audience_age.35-44",
	"audience_age.45-54":          "instagram_account.audience_age.45-54",
	"audience_age.55-64":          "instagram_account.audience_age.55-64",
	"audience_age.65+":            "instagram_account.audience_age.65+",
	"audience_gender.male_per":    "instagram_account.audience_gender.male_per",
	"audience_gender.female_per":  "instagram_account.audience_gender.female_per",
	"audience_gender.other_per":   "instagram_account.audience_gender.other_per",
	"categories":                  "instagram_account.categories",
	"flag_contact_info_available": "instagram_account.flag_contact_info_available",
	"est_post_price.image_posts":  "instagram_account.est_post_price.image_posts",
	"est_post_price.reel_posts":   "instagram_account.est_post_price.reel_posts",
	"est_post_price.story_posts":  "instagram_account.est_post_price.story_posts",
	"avg_likes":                   "instagram_account.avg_likes",
	"averageLikes":                "instagram_account.avg_likes",
	"est_reach.image_posts":       "instagram_account.est_reach.image_posts",
	"est_reach.reel_posts":        "instagram_account.est_reach.reel_posts",
	"est_reach.story_posts":       "instagram_account.est_reach.story_posts",
	"est_impressions.image_posts": "instagram_account.est_impressions.image_posts",
	"est_impressions.reel_posts":  "instagram_account.est_impressions.reel_posts",
	"est_impressions.story_posts": "instagram_account.est_impressions.story_posts",
	"profileLabel":                "instagram_account.label",
	"following":                   "instagram_account.following",
	"ffratio":                     "instagram_account.ffratio",
	"engagement_rate":             "instagram_account.engagement_rate",
	"averageEngagement":           "instagram_account.engagement_rate",
	"avg_views":                   "instagram_account.avg_views",
	"averageVideoViews":           "instagram_account.avg_views",
	"post_count":                  "instagram_account.post_count",
	"numOfPosts":                  "instagram_account.post_count",
	"id":                          "instagram_account.id",
	"name":                        "instagram_account.name",
	"profile_type":                "instagram_account.profile_type",

	"onGcc":           "campaign_profiles.on_gcc",
	"creatorCohorts":  "campaign_profiles.admin_details.creatorCohorts",
	"creatorPrograms": "campaign_profiles.admin_details.creatorPrograms",
	"categoryIds":     "campaign_profiles.admin_details.campaignCategoryIds",

	"shortlistId":        "profile_collection_item.shortlist_id",
	"shortlistingStatus": "profile_collection_item.shortlisting_status",
	"campaignProfileId":  "profile_collection_item.campaign_profile_id",
	"itemEnabled":        "profile_collection_item.enabled",
	"itemHidden":         "profile_collection_item.hidden",

	"campaignId":        "profile_collection.campaign_id",
	"collectionId":      "profile_collection.id",
	"collectionShareId": "profile_collection.share_id",
}

var InstagramSearchFilterMap = map[string]string{
	"keywords":                    "instagram_account.keywords",
	"search_phrase":               "instagram_account.search_phrase",
	"gender":                      "instagram_account.gender",
	"phone":                       "instagram_account.phone",
	"isBlacklisted":               "instagram_account.is_blacklisted",
	"location":                    "instagram_account.location",
	"location_list":               "instagram_account.location_list",
	"languages":                   "instagram_account.languages",
	"dateOfBirth":                 "instagram_account.dob",
	"followers":                   "instagram_account.followers",
	"audience_age.13-17":          "instagram_account.audience_age.13-17",
	"audience_age.18-24":          "instagram_account.audience_age.18-24",
	"audience_age.25-34":          "instagram_account.audience_age.25-34",
	"audience_age.35-44":          "instagram_account.audience_age.35-44",
	"audience_age.45-54":          "instagram_account.audience_age.45-54",
	"audience_age.55-64":          "instagram_account.audience_age.55-64",
	"audience_age.65+":            "instagram_account.audience_age.65+",
	"audience_gender.male_per":    "instagram_account.audience_gender.male_per",
	"audience_gender.female_per":  "instagram_account.audience_gender.female_per",
	"audience_gender.other_per":   "instagram_account.audience_gender.other_per",
	"categories":                  "instagram_account.categories",
	"flag_contact_info_available": "instagram_account.flag_contact_info_available",
	"est_post_price.image_posts":  "instagram_account.est_post_price.image_posts",
	"est_post_price.reel_posts":   "instagram_account.est_post_price.reel_posts",
	"est_post_price.story_posts":  "instagram_account.est_post_price.story_posts",
	"avg_likes":                   "instagram_account.avg_likes",
	"averageLikes":                "instagram_account.avg_likes",
	"est_reach.image_posts":       "instagram_account.est_reach.image_posts",
	"est_reach.reel_posts":        "instagram_account.est_reach.reel_posts",
	"est_reach.story_posts":       "instagram_account.est_reach.story_posts",
	"est_impressions.image_posts": "instagram_account.est_impressions.image_posts",
	"est_impressions.reel_posts":  "instagram_account.est_impressions.reel_posts",
	"est_impressions.story_posts": "instagram_account.est_impressions.story_posts",
	"profileLabel":                "instagram_account.label",
	"following":                   "instagram_account.following",
	"ffratio":                     "instagram_account.ffratio",
	"engagement_rate":             "instagram_account.engagement_rate",
	"averageEngagement":           "instagram_account.engagement_rate",
	"avg_views":                   "instagram_account.avg_views",
	"averageVideoViews":           "instagram_account.avg_views",
	"post_count":                  "instagram_account.post_count",
	"numOfPosts":                  "instagram_account.post_count",
	"id":                          "instagram_account.id",
	"name":                        "instagram_account.name",
	"profile_type":                "instagram_account.profile_type",

	"onGcc":           "campaign_profiles.on_gcc",
	"creatorCohorts":  "campaign_profiles.admin_details.creatorCohorts",
	"creatorPrograms": "campaign_profiles.admin_details.creatorPrograms",
	"categoryIds":     "campaign_profiles.admin_details.campaignCategoryIds",

	"shortlistId":        "profile_collection_item.shortlist_id",
	"shortlistingStatus": "profile_collection_item.shortlisting_status",
	"campaignProfileId":  "profile_collection_item.campaign_profile_id",
	"itemEnabled":        "profile_collection_item.enabled",
	"itemHidden":         "profile_collection_item.hidden",

	"campaignId":        "profile_collection.campaign_id",
	"collectionId":      "profile_collection.id",
	"collectionShareId": "profile_collection.share_id",
}

var InstagramSortingMap = map[string]bool{
	"followers":                             true,
	"plays":                                 true,
	"avg_likes":                             true,
	"engagement_rate":                       true,
	"avg_views":                             true,
	"followers,reels_reach,engagement_rate": true,
	"engagement_rate,engagement_rate_grade_order,followers": true,
	"engagement_rate_grade_order,engagement_rate,followers": true,
	"image_reach,followers,engagement_rate":                 true,
	"reels_reach,followers,engagement_rate":                 true,
	"story_reach,followers,engagement_rate":                 true,
	"avg_likes,avg_likes_grade_order,followers":             true,
	"avg_likes_grade_order,avg_likes,followers":             true,
	"avg_comments,avg_likes,followers":                      true,
	"profile_collection_item.rank":                          true,
}

var YoutubeSortingMap = map[string]bool{
	"followers":                     true,
	"views_count":                   true,
	"avg_views":                     true,
	"followers,views_30d,avg_views": true,
	"views_30d,views_30d_grade_order,followers": true,
	"views_30d_grade_order,views_30d,followers": true,
	"avg_views,followers,views_30d":             true,
	"avg_posts_per_week,followers,views_30d":    true,
	"profile_collection_item.rank":              true,
}

var YoutubeSearchFilterAdminMap = map[string]string{
	"keywords":                    "youtube_account.keywords_admin",
	"search_phrase":               "youtube_account.search_phrase_admin",
	"gender":                      "campaign_profiles.gender",
	"phone":                       "campaign_profiles.phone",
	"isBlacklisted":               "campaign_profiles.admin_details.isBlacklisted",
	"dateOfBirth":                 "campaign_profiles.dob",
	"location":                    "youtube_account.location",
	"location_list":               "campaign_profiles.location",
	"languages":                   "campaign_profiles.languages",
	"followers":                   "youtube_account.followers",
	"numOfSubscribers":            "youtube_account.followers",
	"audience_age.13-17":          "youtube_account.audience_age.13-17",
	"audience_age.18-24":          "youtube_account.audience_age.18-24",
	"audience_age.25-34":          "youtube_account.audience_age.25-34",
	"audience_age.35-44":          "youtube_account.audience_age.35-44",
	"audience_age.45-54":          "youtube_account.audience_age.45-54",
	"audience_age.55-64":          "youtube_account.audience_age.55-64",
	"audience_age.65+":            "youtube_account.audience_age.65+",
	"audience_gender.male_per":    "youtube_account.audience_gender.male_per",
	"audience_gender.female_per":  "youtube_account.audience_gender.female_per",
	"audience_gender.other_per":   "youtube_account.audience_gender.other_per",
	"categories":                  "youtube_account.categories",
	"flag_contact_info_available": "youtube_account.flag_contact_info_available",
	"est_post_price.video_posts":  "youtube_account.est_post_price.video_posts",
	"est_post_price.short_posts":  "youtube_account.est_post_price.short_posts",
	"est_reach.video_posts":       "youtube_account.est_reach.video_posts",
	"est_reach.short_posts":       "youtube_account.est_reach.short_posts",
	"est_impressions.video_posts": "youtube_account.est_impressions.video_posts",
	"est_impressions.short_posts": "youtube_account.est_impressions.short_posts",
	"profileLabel":                "youtube_account.label",
	"avg_views":                   "youtube_account.avg_views",
	"averageVideoViews":           "youtube_account.avg_views",
	"numOfVideos":                 "youtube_account.uploads_count",
	"id":                          "youtube_account.id",
	"name":                        "youtube_account.title",
	"profile_type":                "youtube_account.profile_type",

	"shortlistId":        "profile_collection_item.shortlist_id",
	"shortlistingStatus": "profile_collection_item.shortlisting_status",
	"campaignProfileId":  "profile_collection_item.campaign_profile_id",
	"itemEnabled":        "profile_collection_item.enabled",
	"itemHidden":         "profile_collection_item.hidden",

	"campaignId":        "profile_collection.campaign_id",
	"collectionId":      "profile_collection.id",
	"collectionShareId": "profile_collection.share_id",

	"creatorPrograms": "campaign_profiles.admin_details.creatorPrograms",
	"creatorCohorts":  "campaign_profiles.admin_details.creatorCohorts",
	"onGcc":           "campaign_profiles.on_gcc",
	"categoryIds":     "campaign_profiles.admin_details.campaignCategoryIds",
}

var YoutubeSearchFilterMap = map[string]string{
	"keywords":                    "youtube_account.keywords",
	"search_phrase":               "youtube_account.search_phrase",
	"gender":                      "youtube_account.gender",
	"phone":                       "youtube_account.phone",
	"isBlacklisted":               "youtube_account.is_blacklisted",
	"dateOfBirth":                 "youtube_account.dob",
	"location":                    "youtube_account.location",
	"location_list":               "youtube_account.location_list",
	"languages":                   "youtube_account.languages",
	"followers":                   "youtube_account.followers",
	"numOfSubscribers":            "youtube_account.followers",
	"audience_age.13-17":          "youtube_account.audience_age.13-17",
	"audience_age.18-24":          "youtube_account.audience_age.18-24",
	"audience_age.25-34":          "youtube_account.audience_age.25-34",
	"audience_age.35-44":          "youtube_account.audience_age.35-44",
	"audience_age.45-54":          "youtube_account.audience_age.45-54",
	"audience_age.55-64":          "youtube_account.audience_age.55-64",
	"audience_age.65+":            "youtube_account.audience_age.65+",
	"audience_gender.male_per":    "youtube_account.audience_gender.male_per",
	"audience_gender.female_per":  "youtube_account.audience_gender.female_per",
	"audience_gender.other_per":   "youtube_account.audience_gender.other_per",
	"categories":                  "youtube_account.categories",
	"flag_contact_info_available": "youtube_account.flag_contact_info_available",
	"est_post_price.video_posts":  "youtube_account.est_post_price.video_posts",
	"est_post_price.short_posts":  "youtube_account.est_post_price.short_posts",
	"est_reach.video_posts":       "youtube_account.est_reach.video_posts",
	"est_reach.short_posts":       "youtube_account.est_reach.short_posts",
	"est_impressions.video_posts": "youtube_account.est_impressions.video_posts",
	"est_impressions.short_posts": "youtube_account.est_impressions.short_posts",
	"profileLabel":                "youtube_account.label",
	"avg_views":                   "youtube_account.avg_views",
	"averageVideoViews":           "youtube_account.avg_views",
	"numOfVideos":                 "youtube_account.uploads_count",
	"id":                          "youtube_account.id",
	"name":                        "youtube_account.title",
	"profile_type":                "youtube_account.profile_type",

	"shortlistId":        "profile_collection_item.shortlist_id",
	"shortlistingStatus": "profile_collection_item.shortlisting_status",
	"campaignProfileId":  "profile_collection_item.campaign_profile_id",
	"itemEnabled":        "profile_collection_item.enabled",
	"itemHidden":         "profile_collection_item.hidden",

	"campaignId":        "profile_collection.campaign_id",
	"collectionId":      "profile_collection.id",
	"collectionShareId": "profile_collection.share_id",

	"onGcc":           "campaign_profiles.on_gcc",
	"creatorCohorts":  "campaign_profiles.admin_details.creatorCohorts",
	"creatorPrograms": "campaign_profiles.admin_details.creatorPrograms",
	"categoryIds":     "campaign_profiles.admin_details.campaignCategoryIds",
}

type Graph struct {
	Label string
	Value float64
}

type DiscoveryActivityMeta struct {
	QueryParams map[string]string      `json:"queryParams,omitempty"`
	BodyParam   coredomain.SearchQuery `json:"bodyParams,omitempty"`
	Result      map[string]interface{} `json:"result,omitempty"`
	Platform    string                 `json:"platform,omitempty"`
}
type BodyParam struct {
	Filters []coredomain.SearchFilter `json:"filters,omitempty"`
}
