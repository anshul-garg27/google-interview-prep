package dao

type GenreOverviewEntity struct {
	ID                int64 `gorm:"primaryKey;column:id"`
	Category          string
	Month             string
	Platform          string
	Country           string
	Language          string
	ProfileType       string
	Creators          *int64
	Uploads           *int64
	Views             *int64
	Plays             *int64
	Followers         *int64
	Likes             *int64
	Comments          *int64
	Engagement        *int64
	EngagementRate    *float64
	AudienceAgeGender *string
	AudienceGender    *string
	Enabled           bool
}

func (GenreOverviewEntity) TableName() string {
	return "genre_overview"
}
