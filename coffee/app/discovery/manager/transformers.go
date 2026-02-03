package manager

import (
	cpdomain "coffee/app/campaignprofiles/domain"
	"coffee/app/discovery/dao"
	"coffee/app/discovery/domain"
	beatservice "coffee/client/beat"
	"coffee/constants"
	"coffee/core/appcontext"
	coredomain "coffee/core/domain"
	"context"

	"coffee/helpers"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func ToSocialProfileHashtagsEntry(entity dao.SocialProfileHashatagsEntity) (*domain.SocialProfileHashatagsEntry, error) {
	hashtagCloud := make(map[string]int64)

	for i := 0; i < len(entity.Hashtags); i++ {
		hashtagCloud[entity.Hashtags[i]] = entity.HashtagsCounts[i]
	}
	return &domain.SocialProfileHashatagsEntry{
		PlatformProfileId: entity.PlatformProfileId,
		Platform:          entity.Platform,
		Hashtags:          hashtagCloud,
	}, nil
}

func ToLocationsEntry(entity dao.LocationsEntity) (*domain.LocationsEntry, error) {
	return &domain.LocationsEntry{
		Name:     entity.Name,
		FullName: entity.FullName,
		Type:     entity.Type,
	}, nil
}

func ToSocialProfileAudienceInfoEntry(entity dao.SocialProfileAudienceInfoEntity) (*domain.SocialProfileAudienceInfoEntry, error) {
	var audienceLanguageMap map[string]*float64
	if entity.AudienceLanguage != nil {
		json.Unmarshal([]byte(*entity.AudienceLanguage), &audienceLanguageMap)
	}

	var qualityScoreBreakupMap map[string]*string
	if entity.QualityScoreBreakup != nil {
		json.Unmarshal([]byte(*entity.QualityScoreBreakup), &qualityScoreBreakupMap)
	}
	var audienceGenderMap map[string]*float64
	if entity.AudienceGenderSplit != nil {
		json.Unmarshal([]byte(*entity.AudienceGenderSplit), &audienceGenderMap)
	}

	audienceLocation := domain.AudienceLocation{}
	if entity.AudienceLocationSplit != nil {
		json.Unmarshal([]byte(*entity.AudienceLocationSplit), &audienceLocation)
	}
	audienceAgeGender := domain.AudienceAgeGender{}
	if entity.AudienceAgeGenderSplit != nil {
		json.Unmarshal([]byte(*entity.AudienceAgeGenderSplit), &audienceAgeGender)
	}

	topAge := make(map[string]*float64)
	if entity.AudienceAgeGenderSplit != nil {
		maxValue := -1.00
		maxAge := ""
		for k, v := range audienceAgeGender.Female {
			if val, found := audienceAgeGender.Male[k]; found {
				if (*v)+(*val) > maxValue {
					maxValue = (*v) + (*val)
					maxAge = k
				}
			}
		}
		topAge[maxAge] = &maxValue
	}
	delete(audienceLocation.City, "Others")
	return &domain.SocialProfileAudienceInfoEntry{
		PlatformProfileId:              entity.PlatformProfileId,
		Platform:                       entity.Platform,
		CityWiseAudienceLocation:       audienceLocation.City,
		CountryWiseAudienceLocation:    audienceLocation.Country,
		AudienceGenderSplit:            audienceGenderMap,
		AudienceAgeGenderSplit:         audienceAgeGender,
		TopAge:                         topAge,
		AudienceLanguage:               audienceLanguageMap,
		AudienceReachabilityPercentage: entity.AudienceReachabilityPercentage,
		AudienceAuthenticityPercentage: entity.AudienceAuthenticityPercentage,
		CommentRatePercentage:          entity.CommentRatePercentage,
		QualityAudiencePercentage:      entity.QualityAudiencePercentage,
		QualityAudienceScore:           entity.QualityAudienceScore,
		QualityScoreGrade:              entity.QualityScoreGrade,
		QualityScoreBreakup:            qualityScoreBreakupMap,
	}, nil
}

func convertStringToArray(str string) []float64 {
	var arr []float64
	str = strings.Trim(str, "[]{}\"") // remove curly braces and double quotes from the string
	if str != "" {
		strArr := strings.Split(str, ",")
		//arr := make([]float64, len(strArr))
		for _, s := range strArr {
			f, err := strconv.ParseFloat(s, 64) // convert each substring to a float64
			if err != nil {
				arr = append(arr, 0)
			}
			arr = append(arr, f) // add the float to the float slice
		}
	}

	return arr
}

func ToGroupMetricsEntry(entity dao.GroupMetricsEntity) (*coredomain.GroupMetricsEntry, error) {
	var binStartArray, binEndArray, binHeightrray []float64
	if entity.BinStart != nil {
		binStartArray = convertStringToArray(*entity.BinStart)
		binEndArray = convertStringToArray(*entity.BinEnd)
		binHeightrray = convertStringToArray(*entity.BinHeight)
	}
	var total float64
	var erGraph []coredomain.ERHistogram
	for i := 0; i < len(binStartArray); i++ {
		var erEntry coredomain.ERHistogram
		erEntry.Start = math.Ceil((binStartArray[i]*100)*100) / 100
		erEntry.End = math.Ceil((binEndArray[i]*100)*100) / 100
		erEntry.Value = binHeightrray[i]
		total += binHeightrray[i]
		erGraph = append(erGraph, erEntry)
	}
	for i := range erGraph {
		erGraph[i].Value = erGraph[i].Value / total * 100
	}
	var engagementRatePercentage float64
	if entity.GroupAvgEngagementRate != nil {
		engagementRatePercentage = *entity.GroupAvgEngagementRate * 100
	}
	return &coredomain.GroupMetricsEntry{
		GroupKey:                     entity.GroupKey,
		Profiles:                     entity.Profiles,
		GroupAvgLikes:                entity.GroupAvgLikes,
		GroupAvgComments:             entity.GroupAvgComments,
		GroupAvgCommentsRate:         entity.GroupAvgCommentsRate,
		GroupAvgFollowers:            entity.GroupAvgFollowers,
		GroupAvgEngagementRate:       helpers.ToFloat64(engagementRatePercentage),
		GroupAvgLikesToCommentRatio:  entity.GroupAvgLikesToCommentRatio,
		GroupAvgFollowersGrowth7d:    entity.GroupAvgFollowersGrowth7d,
		GroupAvgFollowersGrowth30d:   entity.GroupAvgFollowersGrowth30d,
		GroupAvgFollowersGrowth90d:   entity.GroupAvgFollowersGrowth90d,
		GroupAvgAudienceReachability: entity.GroupAvgAudienceReachability,
		GroupAvgAudienceAuthencity:   entity.GroupAvgAudienceAuthencity,
		GroupAvgAudienceQuality:      entity.GroupAvgAudienceQuality,
		GroupAvgPostCount:            entity.GroupAvgPostCount,
		GroupAvgReactionRate:         entity.GroupAvgReactionRate,
		GroupAvgLikesSpread:          entity.GroupAvgLikesSpread,
		GroupAvgFollowersGrowth1y:    entity.GroupAvgFollowersGrowth1y,
		GroupAvgPostsPerWeek:         entity.GroupAvgPostsPerWeek,
		GroupAvgReelsReach:           entity.GroupAvgReelsReach,
		GroupAvgImageReach:           entity.GroupAvgImageReach,
		GroupAvgStoryReach:           entity.GroupAvgStoryReach,
		GroupAvgVideoReach:           entity.GroupAvgVideoReach,
		GroupAvgShortsReach:          entity.GroupAvgShortsReach,
		ErGraph:                      erGraph,
	}, nil
}
func makeInstagramlinkedSocials(entity dao.InstagramAccountEntity, source string) *coredomain.SocialAccount {
	audienceGenderAvailable, audienceCityAvailable := checkAvailabilityOfAudienceCityAndGenderInstagram(entity)
	return &coredomain.SocialAccount{
		Code:                     entity.ID,
		ProfileCode:              entity.ProfileId,
		GCCProfileCode:           entity.GccProfileId,
		Name:                     entity.Name,
		Platform:                 string(constants.InstagramPlatform),
		Handle:                   entity.Handle,
		Username:                 entity.Handle,
		IgId:                     entity.IgID,
		PlatformThumbnail:        entity.Thumbnail,
		EngagementRatePercenatge: entity.EngagementRate * 100,
		Followers:                entity.Followers,
		Following:                entity.Following,
		AvgLikes:                 entity.AvgLikes,
		AvgComments:              entity.AvgComments,
		Uploads:                  entity.PostCount,
		AvgViews:                 entity.AvgViews,
		AvgReelsPlay30d:          entity.AvgReelsPlay30d,
		AvgReach:                 entity.AvgReach,
		StoryReach:               entity.StoryReach,
		ImageReach:               entity.ImageReach,
		ReelsReach:               entity.ReelsReach,
		IsVerified:               helpers.ToBool(entity.IsVerified),
		Source:                   source,
		AccountType:              entity.AccountType,
		AvgReelsPlayCount:        entity.AvgReelsPlayCount,
		AudienceGenderAvailable:  &audienceGenderAvailable,
		AudienceCityAvailable:    &audienceCityAvailable,
		IsPrivate:                entity.IsPrivate,
	}
}
func GetProfileFromInstagramEntity(ctx context.Context, entity dao.InstagramAccountEntity, linkedSource string) *coredomain.Profile {
	var categoryMap map[string]*string
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)

	if entity.CategoryRank != nil {
		json.Unmarshal([]byte(*entity.CategoryRank), &categoryMap)
	}

	var estdPriceMap map[string]*float64
	if entity.EstPostPrice != nil {
		json.Unmarshal([]byte(*entity.EstPostPrice), &estdPriceMap)
	}

	var similarGroupDataMap map[string]*float64
	if entity.SimilarProfileGroupData != nil {
		json.Unmarshal([]byte(*entity.SimilarProfileGroupData), &similarGroupDataMap)
	}
	if entity.Name == nil {
		entity.Name = entity.Handle
	}
	var linkedSocials []*coredomain.SocialAccount

	if linkedSource == string(constants.GccLinkedProfile) {
		gccLinkedSocials := makeInstagramlinkedSocials(entity, string(constants.GccLinkedProfile))
		linkedSocials = append(linkedSocials, gccLinkedSocials)
	} else {
		saasLinkedSocials := makeInstagramlinkedSocials(entity, string(constants.SaasLinkedProfile))
		linkedSocials = append(linkedSocials, saasLinkedSocials)
	}
	if appCtx.PlanType != nil && *appCtx.PlanType == constants.FreePlan {
		var encodedEmail, encodedPhone string
		if entity.Email != nil {
			encodedEmail = helpers.EncodeEmail(*entity.Email)
			entity.Email = &encodedEmail
		}
		if entity.Phone != nil {
			encodedPhone = helpers.EncodePhone(*entity.Phone)
			entity.Phone = &encodedPhone
		}
	}
	var estImpressions coredomain.InstagramImpressions
	if entity.EstImpressions != nil {
		json.Unmarshal([]byte(*entity.EstImpressions), &estImpressions)
	}
	profile := coredomain.Profile{
		Code:               entity.ProfileId,
		GccCode:            entity.GccProfileId,
		PlatformCode:       strconv.FormatInt(entity.ID, 10),
		Name:               entity.Name,
		Handle:             entity.Handle,
		Username:           entity.Handle,
		Phone:              entity.Phone,
		Email:              entity.Email,
		IgId:               entity.IgID,
		Description:        entity.Bio,
		Platform:           string(constants.InstagramPlatform),
		ProfileType:        entity.Label,
		LinkedSocials:      linkedSocials,
		LinkedChannelId:    entity.LinkedChannelId,
		GccLinkedChannelId: entity.GccLinkedChannelId,
		Thumbnail:          entity.Thumbnail,
		CategoryRank:       categoryMap,
		EstPostPrice:       estdPriceMap,
		Categories:         (*[]string)(&entity.Categories),
		Gender:             entity.Gender,
		LocationList:       (*[]string)(&entity.LocationList),
		Languages:          (*[]string)(&entity.Languages),
		Metrics: coredomain.ProfileMetrics{
			Followers:                entity.Followers,
			EngagementRatePercenatge: entity.EngagementRate * 100,
			FollowersGrowth7d:        entity.FollowersGrowth7d,
			AuthenticEngagement:      entity.AuthenticEngagement,
			CommentRatePercentage:    entity.CommentRatePercentage,
			LikesSpreadPercentage:    entity.LikesSpreadPercentage,
			CountryRank:              entity.CountryRank,
			StoryReach:               entity.StoryReach,
			ImageReach:               entity.ImageReach,
			ReelsReach:               entity.ReelsReach,
			AvgReach:                 entity.AvgReach,
			AvgReelsPlayCount:        entity.AvgReelsPlayCount,
			AvgLikes:                 entity.AvgLikes,
			AvgComments:              entity.AvgComments,
			AvgViews:                 entity.AvgViews,
			AvgReelsPlay30d:          entity.AvgReelsPlay30d,
			Following:                entity.Following,
			Uploads:                  entity.PostCount,
			Ffratio:                  entity.Ffratio,
			Plays30d:                 entity.Plays30d,
			Uploads30d:               entity.Uploads30d,
			InstagramAccountType:     entity.AccountType,
			Impressions30d:           entity.Impressions30d,
			ReelsImpressions:         estImpressions.ReelPosts,
			ImageImpressions:         estImpressions.ImagePosts,
			FollowersGrowth30d:       entity.FollowersGrowth30d,
		},
		Grades: &coredomain.Grades{
			ErGrade:                   entity.ErGrade,
			AvgCommentsGrade:          entity.AvgCommentsGrade,
			AvgLikesGrade:             entity.AvgLikesGrade,
			CommentsRateGrade:         entity.CommentsRateGrade,
			FollowersGrade:            entity.FollowersGrade,
			EngagementRateGrade:       entity.EngagementRateGrade,
			ReelsReachGrade:           entity.ReelsReachGrade,
			StoryReachGrade:           entity.StoryReachGrade,
			ImageReachGrade:           entity.ImageReachGrade,
			LikesToCommentRatioGrade:  entity.LikesToCommentRatioGrade,
			FollowersGrowth7dGrade:    entity.FollowersGrowth7dGrade,
			FollowersGrowth30dGrade:   entity.FollowersGrowth30dGrade,
			FollowersGrowth90dGrade:   entity.FollowersGrowth90dGrade,
			AudienceReachabilityGrade: entity.AudienceReachabilityGrade,
			AudienceAuthencityGrade:   entity.AudienceAuthencityGrade,
			AudienceQualityGrade:      entity.AudienceQualityGrade,
			PostCountGrade:            entity.PostCountGrade,
			LikesSpreadGrade:          entity.LikesSpreadGrade,
			ImageImpressionsGrade:     entity.ImageImpressionsGrade,
			ReelsImpressionsGrade:     entity.ReelsImpressionsGrade,
		},
		UpdatedAt:               entity.UpdatedAt.Unix(),
		IsVerified:              helpers.ToBool(entity.IsVerified),
		GroupData:               similarGroupDataMap,
		ProfileCollectionItemId: 0,
	}
	if entity.CampaignProfile != nil {
		profile.CampaignProfile, _ = cpdomain.ToCampaignEntry(ctx, entity.CampaignProfile, &profile)
	}
	return &profile
}

func makeSaasYoutubelinkedSocials(entity dao.YoutubeAccountEntity, source string) *coredomain.SocialAccount {
	audienceGenderAvailable, audienceCityAvailable := checkAvailabilityOfAudienceCityAndGenderYoutube(entity)
	return &coredomain.SocialAccount{
		Code:                    entity.ID,
		ProfileCode:             entity.ProfileId,
		Name:                    entity.Title,
		Platform:                string(constants.YoutubePlatform),
		Handle:                  entity.ChannelId,
		Username:                entity.Username,
		PlatformThumbnail:       entity.Thumbnail,
		Followers:               entity.Followers,
		Uploads:                 entity.UploadsCount,
		AvgVideoViews30d:        entity.AvgVideoViews30d,
		AvgReach:                helpers.ToFloat64(0),
		Views:                   entity.ViewsCount,
		AvgViews:                entity.AvgViews,
		Source:                  source,
		GCCProfileCode:          entity.GccProfileId,
		AudienceCityAvailable:   &audienceCityAvailable,
		AudienceGenderAvailable: &audienceGenderAvailable,
	}
}
func GetProfileFromYoutubeEntity(ctx context.Context, entity dao.YoutubeAccountEntity, linkedSource string) *coredomain.Profile {
	var categoryMap map[string]*string
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if entity.CategoryRank != nil {
		json.Unmarshal([]byte(*entity.CategoryRank), &categoryMap)
	}
	var estdPriceMap map[string]*float64
	if entity.EstPostPrice != nil {
		json.Unmarshal([]byte(*entity.EstPostPrice), &estdPriceMap)
	}
	var title *string
	if entity.Title != nil {
		title = entity.Title
	} else {
		title = entity.Username
	}
	var publishTimestamp int64
	if entity.LatestVideoPublishTime != nil {
		publishTime, err := time.Parse("2006-01-02T15:04:05Z", *entity.LatestVideoPublishTime)
		if err == nil {
			publishTimestamp = publishTime.Unix()
		}
	}
	var similarGroupDataMap map[string]*float64
	if entity.SimilarProfileGroupData != nil {
		json.Unmarshal([]byte(*entity.SimilarProfileGroupData), &similarGroupDataMap)
	}
	var linkedSocials []*coredomain.SocialAccount

	if linkedSource == string(constants.GccLinkedProfile) {
		gccLinkedSocials := makeSaasYoutubelinkedSocials(entity, string(constants.GccLinkedProfile))
		linkedSocials = append(linkedSocials, gccLinkedSocials)
	} else {
		saasLinkedSocials := makeSaasYoutubelinkedSocials(entity, string(constants.SaasLinkedProfile))
		linkedSocials = append(linkedSocials, saasLinkedSocials)
	}

	if appCtx.PlanType != nil && *appCtx.PlanType == constants.FreePlan {
		var encodedEmail, encodedPhone string
		if entity.Email != nil {
			encodedEmail = helpers.EncodeEmail(*entity.Email)
			entity.Email = &encodedEmail
		}
		if entity.Phone != nil {
			encodedPhone = helpers.EncodePhone(*entity.Phone)
			entity.Phone = &encodedPhone
		}
	}

	profile := coredomain.Profile{
		Code:            entity.ProfileId,
		GccCode:         entity.GccProfileId,
		PlatformCode:    strconv.FormatInt(entity.ID, 10),
		Name:            title,
		Handle:          entity.ChannelId,
		Username:        entity.Username,
		Phone:           entity.Phone,
		Email:           entity.Email,
		IgId:            new(string),
		Description:     entity.Description,
		Platform:        string(constants.YoutubePlatform),
		ProfileType:     entity.Label,
		LinkedSocials:   linkedSocials,
		LinkedHandle:    entity.LinkedHandle,
		GccLinkedHandle: entity.GccLinkedInstagramHandle,
		Thumbnail:       entity.Thumbnail,
		CategoryRank:    categoryMap,
		EstPostPrice:    estdPriceMap,
		Categories:      (*[]string)(&entity.Categories),
		Gender:          entity.Gender,
		LocationList:    (*[]string)(&entity.LocationList),
		Languages:       (*[]string)(&entity.Languages),
		Metrics: coredomain.ProfileMetrics{
			Followers:              entity.Followers,
			Views:                  entity.ViewsCount,
			FollowersGrowth7d:      entity.FollowersGrowth7d,
			AuthenticEngagement:    entity.AuthenticEngagement,
			CommentRatePercentage:  entity.CommentRatePercentage,
			CountryRank:            entity.CountryRank,
			VideoViewsLast30:       entity.VideoViewsLast30,
			AvgViews:               entity.AvgViews,
			VideoReach:             entity.VideoReach,
			ShortsReach:            entity.ShortsReach,
			ReactionRate:           entity.ReactionRate,
			CommentsRate:           entity.CommentsRate,
			Cpm:                    entity.Cpm,
			AvgPostsPerWeek:        entity.AvgPostsPerWeek,
			FollowersGrowth1y:      entity.FollowersGrowth1y,
			AvgShortsViews30d:      entity.AvgShortsViews30d,
			AvgVideoViews30d:       entity.AvgVideoViews30d,
			LatestVideoPublishTime: publishTimestamp,
			Uploads:                entity.UploadsCount,
			Views30d:               entity.Views30d,
			Uploads30d:             entity.Uploads30d,
		},
		Grades: &coredomain.Grades{
			CommentsRateGrade:      entity.CommentsRateGrade,
			FollowersGrade:         entity.FollowersGrade,
			FollowersGrowth7dGrade: entity.FollowersGrowth7dGrade,
			AvgPostsPerWeekGrade:   entity.AvgPostsPerWeekGrade,
			FollowersGrowth1yGrade: entity.FollowersGrowth1yGrade,
			Views30dGrade:          entity.Views30dGrade,
		},
		UpdatedAt:               entity.UpdatedAt.Unix(),
		GroupData:               similarGroupDataMap,
		ProfileCollectionItemId: 0,
	}
	if entity.CampaignProfile != nil {
		profile.CampaignProfile, _ = cpdomain.ToCampaignEntry(ctx, entity.CampaignProfile, &profile)
	}
	return &profile
}

func GetGrowthFromSocialProfileTimeSeriesEntity(entities []dao.SocialProfileTimeSeriesEntity, monthly_stats bool) *domain.Growth {
	var platformCode string
	var platform string
	var followerGraph, followingGraph, viewGraph, playsGraph, UploadsGraph, followersChangeGraph, engagementRateGraph, totalViewGraph, totalPlaysGraph []domain.Graph

	for _, entity := range entities {
		platformCode = strconv.FormatInt(entity.PlatformProfileId, 10)
		platform = entity.Platform
		var label string
		if monthly_stats {
			label = entity.Date.Format("Jan 2006")
		} else {
			currentDate := entity.Date
			formattedDay := helpers.GetDaySuffix(currentDate.Day())
			formattedDate := fmt.Sprintf("%s %s", formattedDay, currentDate.Month().String()[:3])

			prevDate := entity.PrevDate
			prevDate = prevDate.AddDate(0, 0, 1)

			formattedPrevDay := helpers.GetDaySuffix(prevDate.Day())
			formattedPrevDate := fmt.Sprintf("%s %s", formattedPrevDay, prevDate.Month().String()[:3])

			label = formattedPrevDate + " - " + formattedDate
		}
		if entity.Followers != nil {
			followerGraph = append(followerGraph, domain.Graph{
				Label: label,
				Value: float64(*entity.Followers),
			})
		}
		if entity.Following != nil {
			followingGraph = append(followingGraph, domain.Graph{
				Label: label,
				Value: float64(*entity.Following),
			})
		}
		if entity.EngagementRate != nil {
			engagementRateGraph = append(engagementRateGraph, domain.Graph{
				Label: label,
				Value: *entity.EngagementRate,
			})
		}
		if entity.Views != nil {
			viewGraph = append(viewGraph, domain.Graph{
				Label: label,
				Value: float64(*entity.Views),
			})
		}
		if entity.Plays != nil {
			playsGraph = append(playsGraph, domain.Graph{
				Label: label,
				Value: float64(*entity.Plays),
			})
		}
		if entity.Uploads != nil {
			UploadsGraph = append(UploadsGraph, domain.Graph{
				Label: label,
				Value: float64(*entity.Uploads),
			})
		}
		if entity.FollowersChange != nil {
			followersChangeGraph = append(followersChangeGraph, domain.Graph{
				Label: label,
				Value: float64(*entity.FollowersChange),
			})
		}
		if entity.ViewsTotal != nil {
			totalViewGraph = append(totalViewGraph, domain.Graph{
				Label: label,
				Value: float64(*entity.ViewsTotal),
			})
		}
		if entity.PlaysTotal != nil {
			totalPlaysGraph = append(totalPlaysGraph, domain.Graph{
				Label: label,
				Value: float64(*entity.PlaysTotal),
			})
		}
	}

	return &domain.Growth{
		PlatformCode:         platformCode,
		Platform:             platform,
		FollowerGraph:        followerGraph,
		FollowingGraph:       followingGraph,
		ViewGraph:            viewGraph,
		PlaysGraph:           playsGraph,
		UploadsGraph:         UploadsGraph,
		EngagementGraph:      engagementRateGraph,
		FollowersChangeGraph: followersChangeGraph,
		TotalViewGraph:       totalViewGraph,
		TotalPlaysGraph:      totalPlaysGraph,
	}
}

func MakeYoutubeLinkedSocials(entity dao.YoutubeAccountEntity) coredomain.SocialAccount {
	audienceGenderAvailable, audienceCityAvailable := checkAvailabilityOfAudienceCityAndGenderYoutube(entity)
	return coredomain.SocialAccount{
		Code:                    entity.ID,
		ProfileCode:             entity.ProfileId,
		GCCProfileCode:          entity.GccProfileId,
		Name:                    entity.Title,
		Platform:                string(constants.YoutubePlatform),
		Handle:                  entity.ChannelId,
		Username:                entity.Username,
		PlatformThumbnail:       entity.Thumbnail,
		Followers:               entity.Followers,
		Uploads:                 entity.UploadsCount,
		AvgVideoViews30d:        entity.AvgVideoViews30d,
		AvgReach:                helpers.ToFloat64(0),
		Views:                   entity.ViewsCount,
		AudienceGenderAvailable: &audienceGenderAvailable,
		AudienceCityAvailable:   &audienceCityAvailable,
	}
}

func MakeInstagramLinkedSocials(entity dao.InstagramAccountEntity) coredomain.SocialAccount {
	audienceGenderAvailable, audienceCityAvailable := checkAvailabilityOfAudienceCityAndGenderInstagram(entity)
	return coredomain.SocialAccount{
		Code:                     entity.ID,
		ProfileCode:              entity.ProfileId,
		Name:                     entity.Name,
		Platform:                 string(constants.InstagramPlatform),
		Handle:                   entity.Handle,
		Username:                 entity.Handle,
		PlatformThumbnail:        entity.Thumbnail,
		EngagementRatePercenatge: entity.EngagementRate * 100,
		Followers:                entity.Followers,
		Following:                entity.Following,
		AvgLikes:                 entity.AvgLikes,
		AvgComments:              entity.AvgComments,
		Uploads:                  entity.PostCount,
		AvgReelsPlay30d:          entity.AvgReelsPlay30d,
		AvgReach:                 entity.AvgReach,
		StoryReach:               entity.StoryReach,
		ImageReach:               entity.ImageReach,
		ReelsReach:               entity.ReelsReach,
		IsVerified:               helpers.ToBool(entity.IsVerified),
		AudienceGenderAvailable:  &audienceGenderAvailable,
		AudienceCityAvailable:    &audienceCityAvailable,
		IsPrivate:                entity.IsPrivate,
	}
}

func EnrichProfileWithSaasLinkedSocials(profiles []coredomain.Profile, linkedMap map[string]coredomain.SocialAccount) []coredomain.Profile {
	for i := range profiles {
		if profiles[i].Platform == string(constants.InstagramPlatform) {
			channelId := profiles[i].LinkedChannelId
			if channelId != nil {
				if _, isPresent := linkedMap[*channelId]; isPresent {
					linkedDetails := linkedMap[*channelId]
					linkedDetails.Source = string(constants.SaasLinkedProfile)
					profiles[i].LinkedSocials = append(profiles[i].LinkedSocials, &linkedDetails)
				}
			}

		} else if profiles[i].Platform == string(constants.YoutubePlatform) {
			handle := profiles[i].LinkedHandle
			if handle != nil {
				if _, isPresent := linkedMap[*handle]; isPresent {
					linkedDetails := linkedMap[*handle]
					linkedDetails.Source = string(constants.SaasLinkedProfile)
					profiles[i].LinkedSocials = append(profiles[i].LinkedSocials, &linkedDetails)
				}
			}
		}
	}
	return profiles
}

func EnrichProfileWithGccLinkedSocials(profiles []coredomain.Profile, linkedMap map[string]coredomain.SocialAccount) []coredomain.Profile {
	for i := range profiles {
		if profiles[i].Platform == string(constants.InstagramPlatform) {
			channelId := profiles[i].GccLinkedChannelId
			if channelId != nil {
				if _, isPresent := linkedMap[*channelId]; isPresent {
					linkedDetails := linkedMap[*channelId]
					linkedDetails.Source = string(constants.GCCPlatform)
					profiles[i].LinkedSocials = append(profiles[i].LinkedSocials, &linkedDetails)
				}
			}

		} else if profiles[i].Platform == string(constants.YoutubePlatform) {
			handle := profiles[i].GccLinkedHandle
			if handle != nil {
				if _, isPresent := linkedMap[*handle]; isPresent {
					linkedDetails := linkedMap[*handle]
					linkedDetails.Source = string(constants.GCCPlatform)
					profiles[i].LinkedSocials = append(profiles[i].LinkedSocials, &linkedDetails)
				}
			}
		}
	}
	return profiles
}

func TransformDataToSocialProfilePosts(platformProfileId int64, data beatservice.Data, platform string, handle string) []coredomain.SocialProfilePostsEntry {

	var socialProfilePostsEntry coredomain.SocialProfilePostsEntry
	var posts []coredomain.SocialProfilePostsEntry
	for _, post := range data.RecentPosts {

		var publishTime time.Time
		var publishTimestamp int64
		publishTime, err := time.Parse("2006-01-02T15:04:05", post.PublishTime)
		if err != nil {
			fmt.Println(err)
		}
		publishTimestamp = publishTime.Unix()
		socialProfilePostsEntry = coredomain.SocialProfilePostsEntry{
			Platform:       platform,
			Handle:         handle,
			PostId:         post.PostId,
			PostLink:       post.PostUrl,
			PostType:       post.PostType,
			Thumbnail:      post.Dimensions.ThumbnailUrl,
			LikesCount:     post.Metrics.Likes,
			CommentsCount:  post.Dimensions.Comments,
			ViewsCount:     post.Metrics.Reach,
			EngagementRate: post.Metrics.EngagementRate,
			PublishedAt:    publishTimestamp,
		}
		posts = append(posts, socialProfilePostsEntry)
	}
	return posts
}

func ConvertCountryToCode(country string) string {
	if value, ok := constants.CountryCodes[country]; ok {
		return value
	} else {
		return country
	}
}

func checkAvailabilityOfAudienceCityAndGenderInstagram(entity dao.InstagramAccountEntity) (bool, bool) {
	audienceGenderAvailable := false
	if entity.AudienceGender != nil {
		audienceGenderAvailable = true
	}
	audienceCityAvailable := false
	if entity.AudienceLocation != nil {
		var data map[string]map[string]float64
		json.Unmarshal([]byte(*entity.AudienceLocation), &data)
		if _, ok := data["city"]; ok && len(data["city"]) > 0 {
			audienceCityAvailable = true
		}
	}
	return audienceGenderAvailable, audienceCityAvailable
}

func checkAvailabilityOfAudienceCityAndGenderYoutube(entity dao.YoutubeAccountEntity) (bool, bool) {
	audienceGenderAvailable := false
	if entity.AudienceGender != nil {
		audienceGenderAvailable = true
	}
	audienceCityAvailable := false
	if entity.AudienceLocation != nil {
		var data map[string]map[string]float64
		json.Unmarshal([]byte(*entity.AudienceLocation), &data)
		if _, ok := data["city"]; ok && len(data["city"]) > 0 {
			audienceCityAvailable = true
		}
	}
	return audienceGenderAvailable, audienceCityAvailable
}
