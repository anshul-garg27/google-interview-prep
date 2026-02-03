package jobtracker

type UpsertInstagramEntry struct {
	ProfilePic        *string  `json:"profile_pic,omitempty"`
	BusinessID        *string  `json:"business_id,omitempty"`
	IgID              *string  `json:"ig_id,omitempty"`
	Handle            *string  `json:"handle,omitempty"`
	Name              *string  `json:"name,omitempty"`
	Bio               *string  `json:"bio,omitempty"`
	Followers         *int64   `json:"followers,omitempty"`
	Following         *int64   `json:"following,omitempty"`
	PostCount         *int64   `json:"num_of_posts,omitempty"`
	EngagementRate    *float64 `json:"avg_engagement,omitempty"`
	AvgLikes          *float64 `json:"avg_likes,omitempty"`
	AvgComments       *float64 `json:"avg_comments,omitempty"`
	AvgReach          *float64 `json:"avg_reach,omitempty"`
	AvgReelsPlayCount *float64 `json:"avg_play_count,omitempty"`
	AvgViews          *float64 `json:"avg_views_count,omitempty"`
	ImageReach        *float64 `json:"avg_static_reach,omitempty"`
	Ffratio           *float64 `json:"ffratio,omitempty"`
	StoryReach        *float64 `json:"story_reach,omitempty"`
	City              *string  `json:"city,omitempty"`
	State             *string  `json:"state,omitempty"`
	Country           *string  `json:"country,omitempty"`
	AccountType       *string  `json:"accountType,omitempty"`
	IsPrivate         *bool    `json:"is_private,omitempty"`
}

type UpsertYoutubeEntry struct {
	ChannelId    *string `json:"handle,omitempty"`
	Followers    *int64  `json:"subscribers,omitempty"`
	ViewsCount   *int64  `json:"views,omitempty"`
	UploadsCount *int64  `json:"uploads,omitempty"`
	Title        *string `json:"title,omitempty"`
	Thumbnail    *string `json:"profile_pic,omitempty"`
	City         *string `json:"city,omitempty"`
	State        *string `json:"state,omitempty"`
	Country      *string `json:"country,omitempty"`
}
