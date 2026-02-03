package beatservice

type Status struct {
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}
type Data struct {
	Platform    string  `json:"platform,omitempty"`
	ProfileId   string  `json:"profile_id,omitempty"`
	RecentPosts []Posts `json:"recent_posts,omitempty"`
}
type Posts struct {
	ShortCode   string     `json:"shortcode,omitempty"`
	PostId      string     `json:"post_id,omitempty"`
	PostType    string     `json:"post_type,omitempty"`
	PostUrl     string     `json:"post_url,omitempty"`
	PublishTime string     `json:"publish_time,omitempty"`
	Metrics     Metrics    `json:"metrics,omitempty"`
	Dimensions  Dimensions `json:"dimensions,omitempty"`
}
type Metrics struct {
	Reach          int64   `json:"reach,omitempty"`
	Views          int64   `json:"views,omitempty"`
	EngagementRate float64 `json:"engagement_rate,omitempty"`
	Likes          int64   `json:"likes,omitempty"`
	PlayCount      int64   `json:"play_count,omitempty"`
}

type Dimensions struct {
	Caption      string `json:"caption,omitempty"`
	Comments     int64  `json:"comments,omitempty"`
	ThumbnailUrl string `json:"thumbnail_url,omitempty"`
	DisplayUrl   string `json:"display_url,omitempty"`
}

type RecentPostsResponse struct {
	Status Status `json:"status"`
	Data   Data   `json:"data"`
}

type Profile struct {
	Platform           string  `json:"platform,omitempty"`
	Handle             string  `json:"handle,omitempty"`
	FbId               string  `json:"fbid,omitempty"`
	ProfilePicUrl      string  `json:"profile_pic_url,omitempty"`
	ProfileId          string  `json:"profile_id,omitempty"`
	Biography          string  `json:"biography,omitempty"`
	Following          int64   `json:"following,omitempty"`
	Followers          int64   `json:"followers,omitempty"`
	FullName           string  `json:"full_name,omitempty"`
	IsPrivate          *bool   `json:"is_private,omitempty"` // changing bool to *bool because if its boolean, then in case of null, it is storing false value.
	RecentPosts        []Posts `json:"recent_posts,omitempty"`
	PostMetrics        Metrics `json:"post_metrics,omitempty"`
	Uploads            int64   `json:"media_count,omitempty"`
	UpdatedAt          string  `json:"updated_at,omitempty"`
	YoutubeChannelId   string  `json:"channel_id,omitempty"`
	YoutubeTitle       string  `json:"title,omitempty"`
	YoutubeSubscribers int64   `json:"subscribers,omitempty"`
	YoutubeUploads     int64   `json:"uploads,omitempty"`
	YoutubeViews       int64   `json:"views,omitempty"`
	YoutubeThumbnail   string  `json:"thumbnail,omitempty"`
}

type ProfileResponse struct {
	Status  Status   `json:"status"`
	Profile *Profile `json:"profile"`
}

type YoutubeChannelResponse struct {
	Status         Status         `json:"status"`
	YoutubeChannel YoutubeChannel `json:"response"`
}

type YoutubeChannel struct {
	Handle    string `json:"handle,omitempty"`
	ChannelId string `json:"channel_id,omitempty"`
}

type CompleteProfileResponse struct {
	Status          Status          `json:"status"`
	CompleteProfile CompleteProfile `json:"profile"`
}
type CompleteProfile struct {
	Platform           string                     `json:"platform,omitempty"`
	Handle             string                     `json:"handle,omitempty"`
	FbId               string                     `json:"fbid,omitempty"`
	ProfilePicUrl      string                     `json:"profile_pic_url,omitempty"`
	ProfileId          string                     `json:"profile_id,omitempty"`
	Biography          string                     `json:"biography,omitempty"`
	Following          int64                      `json:"following,omitempty"`
	Followers          int64                      `json:"followers,omitempty"`
	FullName           string                     `json:"full_name,omitempty"`
	IsPrivate          *bool                      `json:"is_private,omitempty"`
	RecentPosts        []Posts                    `json:"recent_posts,omitempty"`
	PostMetrics        CompleteProfilePostMetrics `json:"post_metrics,omitempty"`
	UpdatedAt          string                     `json:"updated_at,omitempty"`
	YoutubeChannelId   string                     `json:"channel_id,omitempty"`
	YoutubeTitle       string                     `json:"title,omitempty"`
	YoutubeSubscribers int64                      `json:"subscribers,omitempty"`
	YoutubeUploads     int64                      `json:"uploads,omitempty"`
	YoutubeViews       int64                      `json:"views,omitempty"`
	YoutubeThumbnail   string                     `json:"thumbnail,omitempty"`
	Uploads            *int64                     `json:"media_count,omitempty"`
}
type CompleteProfilePostMetrics struct {
	Static PostMetrics `json:"static,omitempty"`
	Reels  PostMetrics `json:"reels,omitempty"`
	Story  PostMetrics `json:"story,omitempty"`
	All    PostMetrics `json:"all,omitempty"`
}
type PostMetrics struct {
	AverageEngagementRate *float64 `json:"avg_engagement,omitempty"`
	AveragePlays          *float64 `json:"avg_plays,omitempty"`
	AverageLikes          *float64 `json:"avg_likes,omitempty"`
	AverageComments       *float64 `json:"avg_comments,omitempty"`
	AverageReach          *float64 `json:"avg_reach,omitempty"`
}

type CreatorInsight struct {
	Status  Status                 `json:"status"`
	Profile CreatorInsightAudience `json:"profile"`
}

type CreatorInsightAudience struct {
	Handle       string       `json:"handle,omitempty"`
	AudienceData AudienceData `json:"audience_data,omitempty"`
}

type AudienceData struct {
	City      map[string]int64       `json:"city,omitempty"`
	GenderAge map[string]interface{} `json:"gender_age,omitempty"`
}
