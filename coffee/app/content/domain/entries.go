package domain

import coredomain "coffee/core/domain"

type TrendingContentEntry struct {
	ID                int64                    `json:"id,omitempty"`
	Category          string                   `json:"category,omitempty"`
	Platform          string                   `json:"platform,omitempty"`
	Language          string                   `json:"language,omitempty"`
	Views             *int64                   `json:"views,omitempty"`
	Plays             *int64                   `json:"plays,omitempty"`
	Likes             *int64                   `json:"likes,omitempty"`
	Comments          *int64                   `json:"comments,omitempty"`
	Duration          *int64                   `json:"duration,omitempty"`
	PublishedAt       int64                    `json:"publishedAt,omitempty"`
	PostId            string                   `json:"postId,omitempty"`
	PostType          string                   `json:"postType,omitempty"`
	PostLink          string                   `json:"postLink,omitempty"`
	Thumbnail         string                   `json:"thumbnail,omitempty"`
	PlatformAccountId *int64                   `json:"platformAccountId"`
	Month             string                   `json:"month,omitempty"`
	Title             *string                  `json:"title,omitempty"`
	Description       *string                  `json:"description,omitempty"`
	Engagement        *int64                   `json:"engagement,omitempty"`
	EngagementRate    *float64                 `json:"engagementRate,omitempty"`
	AudienceAge       map[string]*float64      `json:"audienceAge,omitempty"`
	AudienceGender    map[string]*float64      `json:"audienceGender,omitempty"`
	AccountInfo       coredomain.SocialAccount `json:"accountInfo"`
}
