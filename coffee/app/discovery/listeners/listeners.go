package listeners

import (
	campaignprofiledao "coffee/app/campaignprofiles/dao"
	manager2 "coffee/app/campaignprofiles/domain"
	cpmanager "coffee/app/campaignprofiles/manager"
	"coffee/app/discovery/dao"
	"coffee/app/discovery/manager"
	"coffee/constants"
	"coffee/core/domain"
	"context"
	"reflect"
	"strconv"

	upserttracker "coffee/client/upserttracker"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/imdario/mergo"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func SetupListeners(router *message.Router, subscriber *amqp.Subscriber) {
	router.AddNoPublisherHandler(
		"upsert_instagram_account",
		"upserttracker.dx___upsert_instagram_account_q",
		subscriber,
		PerformUpsertInstagramAccount,
	)

	router.AddNoPublisherHandler(
		"upsert_youtube_account",
		"upserttracker.dx___upsert_youtube_account_q",
		subscriber,
		PerformUpsertYoutubeAccount,
	)
}

func PerformUpsertInstagramAccount(message *message.Message) error {
	var instaAccountFromDB *dao.InstagramAccountEntity
	var totalPhrase, instaLabel, searchPhraseAdmin *string
	var location, locationList, keywords []string
	var enabledForSaas bool
	var keywordsAdmin *[]string
	var err error

	ctx := message.Context()
	instaAccount := &upserttracker.UpsertInstagramEntry{}
	json.Unmarshal(message.Payload, instaAccount)

	if instaAccount.Handle == nil || instaAccount.IgID == nil {
		return nil
	}

	igIdTrimmed := strings.TrimSpace(*instaAccount.IgID)
	if igIdTrimmed == "" {
		return nil
	}

	instagramAccount := formIAFromInput(instaAccount)
	instaDao := dao.CreateInstagramAccountDao(ctx)
	instaAccountFromDB, err = instaDao.FindByIgId(ctx, *instagramAccount.IgID)
	if err != nil && err.Error() != "record not found" {
		log.Error("error while finding IA using IG-ID")
		return nil
	}
	if instaAccountFromDB == nil && err.Error() == "record not found" {
		instaAccountFromDB, err = instaDao.Create(ctx, instagramAccount)
		if err != nil {
			log.Error("error creating instagram account: ", err, instagramAccount)
			return nil
		}
		profileId := fmt.Sprintf("IA_%d", instaAccountFromDB.ID)
		instagramAccount.ProfileId = &profileId
		instagramAccount.GccProfileId = &profileId
		instaAccountFromDB, err = instaDao.Update(ctx, instagramAccount.ID, instagramAccount)
		if err != nil {
			log.Error("error setting profile_id of IA: ", err)
			return nil
		}
	}

	cpDao := campaignprofiledao.CreateCampaignProfileDao(ctx)
	cpManager := cpmanager.CreateManager(ctx)
	platformAccountId := strconv.FormatInt(instaAccountFromDB.ID, 10)
	campaignEntity, err := cpDao.FindByPlatformAndPlatformId(ctx, string(constants.InstagramPlatform), platformAccountId)
	if err == nil && campaignEntity == nil {
		instaProfile := manager.GetProfileFromInstagramEntity(ctx, *instagramAccount, string(constants.GCCPlatform))
		campaignEntity, err = cpManager.CreateCampaignProfileForNewHandles(ctx, instaProfile)
		if err != nil {
			log.Error("error while create CP (call from listener) ", err)
			return nil
		}
	}

	instagramAccount = mergingEventAndDBDataForInstagram(instaAccountFromDB, instagramAccount)
	campaignProfileEntry, _ := manager2.ToCampaignEntry(ctx, campaignEntity, nil)
	keywords, totalPhrase, enabledForSaas, location, locationList, instaLabel, keywordsAdmin, searchPhraseAdmin = updateParametersForInstagram(ctx, instagramAccount, campaignProfileEntry)
	instagramAccount = updateInstagramEntity(instagramAccount, keywords, totalPhrase, enabledForSaas, location, locationList, instaLabel, keywordsAdmin, searchPhraseAdmin)
	_, err = instaDao.Update(ctx, instaAccountFromDB.ID, instagramAccount)
	if err != nil {
		log.Error("error updating IA: ", err)
		return nil
	}
	return nil
}

func formIAFromInput(instaAccount *upserttracker.UpsertInstagramEntry) *dao.InstagramAccountEntity {
	var thumbnail, businessId, igId, handle, name, bio, city, state, country, accountType *string
	var followers, following, postcount, imagereach, storyreach *int64
	var avglikes, avgcomments, avgreach, avgreelsplaycount, avgviews, ffratio *float64
	var engagementRate float64
	var isPrivate *bool

	if instaAccount.ProfilePic != nil {
		thumbnail = instaAccount.ProfilePic
	}
	if instaAccount.BusinessID != nil {
		businessId = instaAccount.BusinessID
	}
	if instaAccount.IgID != nil {
		igId = instaAccount.IgID
	}
	if instaAccount.Handle != nil {
		handle = instaAccount.Handle
	}
	if instaAccount.Name != nil {
		name = instaAccount.Name
	}
	if instaAccount.Bio != nil {
		bio = instaAccount.Bio
	}
	if instaAccount.Followers != nil {
		followers = instaAccount.Followers
	}
	if instaAccount.Following != nil {
		following = instaAccount.Following
	}
	if instaAccount.PostCount != nil {
		postcount = instaAccount.PostCount
	}
	if instaAccount.ImageReach != nil {
		imageReach := int64(math.Round(*instaAccount.ImageReach))
		imagereach = &imageReach
	}
	if instaAccount.StoryReach != nil {
		storyReach := int64(math.Round(*instaAccount.StoryReach))
		storyreach = &storyReach
	}
	if instaAccount.AvgLikes != nil {
		avglikes = instaAccount.AvgLikes
	}
	if instaAccount.AvgComments != nil {
		avgcomments = instaAccount.AvgComments
	}
	if instaAccount.AvgReach != nil {
		avgreach = instaAccount.AvgReach
	}
	if instaAccount.AvgReelsPlayCount != nil {
		avgreelsplaycount = instaAccount.AvgReelsPlayCount
	}
	if instaAccount.AvgViews != nil {
		avgviews = instaAccount.AvgViews
	}
	if instaAccount.Ffratio != nil {
		ffratio = instaAccount.Ffratio
	}
	if instaAccount.EngagementRate != nil {
		engagementRate = *instaAccount.EngagementRate
	}
	if instaAccount.City != nil {
		cityValue := strings.Title(*instaAccount.City)
		city = &cityValue
	}
	if instaAccount.State != nil {
		stateValue := strings.Title(*instaAccount.State)
		state = &stateValue
	}
	if instaAccount.Country != nil {
		countryValue := strings.Title(*instaAccount.Country)
		country = &countryValue
	}
	if instaAccount.AccountType != nil {
		accountType = instaAccount.AccountType
	}

	if instaAccount.IsPrivate != nil {
		isPrivate = instaAccount.IsPrivate
	}

	return &dao.InstagramAccountEntity{
		Name:              name,
		Handle:            handle,
		BusinessID:        businessId,
		IgID:              igId,
		Thumbnail:         thumbnail,
		Bio:               bio,
		Followers:         followers,
		Following:         following,
		PostCount:         postcount,
		Ffratio:           ffratio,
		AvgViews:          avgviews,
		AvgLikes:          avglikes,
		AvgReach:          avgreach,
		AvgReelsPlayCount: avgreelsplaycount,
		EngagementRate:    engagementRate,
		StoryReach:        imagereach,
		ImageReach:        storyreach,
		AvgComments:       avgcomments,
		EnabledForSaas:    false,
		City:              city,
		State:             state,
		Country:           country,
		AccountType:       accountType,
		IsPrivate:         isPrivate,
	}
}

func mergingEventAndDBDataForInstagram(instaAccountFromDB *dao.InstagramAccountEntity, instagramAccount *dao.InstagramAccountEntity) *dao.InstagramAccountEntity {
	mergo.Merge(instagramAccount, instaAccountFromDB, mergo.WithOverride, mergo.WithTransformers(blankTransformer{}))
	return instagramAccount
}

func updateParametersForInstagram(ctx context.Context, instagramAccount *dao.InstagramAccountEntity, campaignProfileEntry *domain.CampaignProfileEntry) ([]string, *string, bool, []string, []string, *string, *[]string, *string) {
	var keywords, location, locationList, keywordsAdmin []string
	var totalPhrase, city, state, countryCode, instaLabel *string
	var searchPhraseAdmin, search_phrase string
	if instagramAccount.Name != nil && *instagramAccount.Name != "" {
		keywords = append(keywords, strings.Fields(strings.ToLower(*instagramAccount.Name))...)
		search_phrase += strings.ToLower(*instagramAccount.Name) + " "
	}
	if instagramAccount.Handle != nil && *instagramAccount.Handle != "" {
		keywords = append(keywords, strings.Fields(strings.ToLower(*instagramAccount.Handle))...)
		search_phrase += strings.ToLower(*instagramAccount.Handle) + " "
	}
	search_phrase = strings.TrimSpace(search_phrase)
	totalPhrase = &search_phrase
	enabledForSaas := manager.EnabledForSaasInstagram(instagramAccount)

	if instagramAccount.City != nil && *instagramAccount.City != "" {
		city = instagramAccount.City
		locationList = append(locationList, *city)
		cityPrefix := "city_" + (*city)
		location = append(location, cityPrefix)
	}
	if instagramAccount.State != nil && *instagramAccount.State != "" {
		state = instagramAccount.State
		locationList = append(locationList, *state)
		statePrefix := "state_" + (*state)
		location = append(location, statePrefix)
	}
	if instagramAccount.Country != nil && *instagramAccount.Country != "" {
		countryCode = instagramAccount.Country
		locationList = append(locationList, manager.ConvertCountryToCode(*countryCode))
		countryPrefix := "country_" + *countryCode
		location = append(location, countryPrefix)
	}
	if instagramAccount.Followers != nil {
		followers := instagramAccount.Followers
		var label string
		if *followers < 1000 {
			label = "Nano"
		} else if *followers < 75000 {
			label = "Micro"
		} else if *followers < 1000000 {
			label = "Macro"
		} else {
			label = "Mega"
		}
		instaLabel = &label
	}

	if campaignProfileEntry != nil {
		instaProfile := manager.GetProfileFromInstagramEntity(ctx, *instagramAccount, string(constants.GCCPlatform))
		searchPhraseAdmin, keywordsAdmin = manager2.CreateKeywordsAndSearchPhraseAdmin(campaignProfileEntry.AdminDetails, campaignProfileEntry.UserDetails, string(constants.InstagramPlatform), instaProfile)
		if searchPhraseAdmin == "" && len(keywordsAdmin) == 0 {
			return keywords, totalPhrase, enabledForSaas, location, locationList, instaLabel, nil, nil
		}
		return keywords, totalPhrase, enabledForSaas, location, locationList, instaLabel, &keywordsAdmin, &searchPhraseAdmin
	}

	return keywords, totalPhrase, enabledForSaas, location, locationList, instaLabel, &keywords, totalPhrase
}

func updateInstagramEntity(instagramAccount *dao.InstagramAccountEntity, keywords []string, totalPhrase *string, enabledForSaas bool, location []string, locationList []string, instaLabel *string, keywordsAdmin *[]string, searchPhraseAdmin *string) *dao.InstagramAccountEntity {
	if len(keywords) != 0 {
		instagramAccount.Keywords = pq.StringArray(keywords)
		instagramAccount.SearchPhrase = totalPhrase
	}
	if keywordsAdmin != nil {
		instagramAccount.KeywordsAdmin = pq.StringArray(*keywordsAdmin)
	}
	if searchPhraseAdmin != nil {
		instagramAccount.SearchPhraseAdmin = searchPhraseAdmin
	}
	instagramAccount.EnabledForSaas = enabledForSaas
	instagramAccount.Location = location
	instagramAccount.LocationList = locationList
	instagramAccount.Label = instaLabel
	instagramAccount.Enabled = true
	flag := false
	instagramAccount.Deleted = &flag
	return instagramAccount
}

func PerformUpsertYoutubeAccount(message *message.Message) error {
	var totalPhrase, youtubeLabel, searchPhraseAdmin *string
	var location, locationList, keywords []string
	var keywordsAdmin *[]string
	var enabledForSaas bool
	var err error
	var ytAccountFromDB *dao.YoutubeAccountEntity

	ctx := message.Context()
	ytAccount := &upserttracker.UpsertYoutubeEntry{}
	json.Unmarshal(message.Payload, ytAccount)

	if ytAccount.ChannelId == nil {
		return nil
	}

	youtubeAccount := formYAFromInput(ytAccount)
	ytDao := dao.CreateYoutubeAccountDao(ctx)
	ytAccountFromDB, err = ytDao.FindByHandle(ctx, *youtubeAccount.ChannelId, false)
	if err != nil && err.Error() != "record not found" {
		log.Error("error while finding IA using IG-ID")
		return nil
	}
	if ytAccountFromDB == nil && err.Error() == "record not found" {
		ytAccountFromDB, err = ytDao.Create(ctx, youtubeAccount)
		if err != nil {
			log.Error("error creating youtube account: ", err, youtubeAccount)
			return nil
		}
		profileId := fmt.Sprintf("YA_%d", ytAccountFromDB.ID)
		youtubeAccount.ProfileId = &profileId
		youtubeAccount.GccProfileId = youtubeAccount.ProfileId
		ytAccountFromDB, err = ytDao.Update(ctx, youtubeAccount.ID, youtubeAccount)
		if err != nil {
			log.Error("error setting profile_id of YA: ", err)
			return nil
		}
	}

	cpDao := campaignprofiledao.CreateCampaignProfileDao(ctx)
	cpManager := cpmanager.CreateManager(ctx)
	platformAccountId := strconv.FormatInt(ytAccountFromDB.ID, 10)
	campaignEntity, err := cpDao.FindByPlatformAndPlatformId(ctx, string(constants.YoutubePlatform), platformAccountId)
	if err == nil && campaignEntity == nil {
		instaProfile := manager.GetProfileFromYoutubeEntity(ctx, *youtubeAccount, string(constants.GCCPlatform))
		campaignEntity, err = cpManager.CreateCampaignProfileForNewHandles(ctx, instaProfile)
		if err != nil {
			log.Error("error while create CP (call from listener) ", err)
			return nil
		}
	}

	youtubeAccount = mergingEventAndDBDataForYoutube(ytAccountFromDB, youtubeAccount)
	campaignProfileEntry, _ := manager2.ToCampaignEntry(ctx, campaignEntity, nil)
	keywords, totalPhrase, enabledForSaas, location, locationList, youtubeLabel, keywordsAdmin, searchPhraseAdmin = updateParametersForYoutube(ctx, youtubeAccount, campaignProfileEntry)
	youtubeAccount = updateYoutubeEntity(youtubeAccount, keywords, totalPhrase, enabledForSaas, location, locationList, youtubeLabel, keywordsAdmin, searchPhraseAdmin)
	_, err = ytDao.Update(ctx, ytAccountFromDB.ID, youtubeAccount)
	if err != nil {
		log.Error("error updating YA: ", err)
		return nil
	}
	return nil
}

func formYAFromInput(ytAccount *upserttracker.UpsertYoutubeEntry) *dao.YoutubeAccountEntity {
	var channelId, title, thumbnail, city, state, country *string
	var followers, viewscount, uploadscount *int64

	if ytAccount.Thumbnail != nil {
		thumbnail = ytAccount.Thumbnail
	}
	if ytAccount.ChannelId != nil {
		channelId = ytAccount.ChannelId
	}
	if ytAccount.Title != nil {
		title = ytAccount.Title
	}
	if ytAccount.Followers != nil {
		followers = ytAccount.Followers
	}
	if ytAccount.ViewsCount != nil {
		viewscount = ytAccount.ViewsCount
	}
	if ytAccount.UploadsCount != nil {
		uploadscount = ytAccount.UploadsCount
	}
	if ytAccount.City != nil {
		cityValue := strings.Title(*ytAccount.City)
		city = &cityValue
	}
	if ytAccount.State != nil {
		stateValue := strings.Title(*ytAccount.State)
		state = &stateValue
	}
	if ytAccount.Country != nil {
		countryValue := strings.Title(*ytAccount.Country)
		country = &countryValue
	}

	return &dao.YoutubeAccountEntity{
		ChannelId:      channelId,
		Followers:      followers,
		ViewsCount:     viewscount,
		UploadsCount:   uploadscount,
		Title:          title,
		Thumbnail:      thumbnail,
		EnabledForSaas: false,
		City:           city,
		State:          state,
		Country:        country,
	}
}

func mergingEventAndDBDataForYoutube(ytAccountFromDB *dao.YoutubeAccountEntity, youtubeAccount *dao.YoutubeAccountEntity) *dao.YoutubeAccountEntity {
	mergo.Merge(youtubeAccount, ytAccountFromDB, mergo.WithOverride, mergo.WithTransformers(blankTransformer{}))
	return youtubeAccount
}

func updateParametersForYoutube(ctx context.Context, youtubeAccount *dao.YoutubeAccountEntity, campaignProfileEntry *domain.CampaignProfileEntry) ([]string, *string, bool, []string, []string, *string, *[]string, *string) {
	var keywords, location, locationList, keywordsAdmin []string
	var totalPhrase, city, state, countryCode, youtubeLabel *string
	var searchPhraseAdmin, search_phrase string
	if youtubeAccount.Title != nil && *youtubeAccount.Title != "" {
		keywords = append(keywords, strings.Fields(strings.ToLower(*youtubeAccount.Title))...)
		search_phrase += strings.ToLower(*youtubeAccount.Title) + " "
	}
	if youtubeAccount.Username != nil && *youtubeAccount.Username != "" {
		keywords = append(keywords, strings.Fields(strings.ToLower(*youtubeAccount.Username))...)
		search_phrase += strings.ToLower(*youtubeAccount.Username) + " "
	}
	search_phrase = strings.TrimSpace(search_phrase)
	totalPhrase = &search_phrase
	enabledForSaas := manager.EnabledForSaasYoutube(youtubeAccount)
	if youtubeAccount.City != nil && *youtubeAccount.City != "" {
		city = youtubeAccount.City
		locationList = append(locationList, *city)
		cityPrefix := "city_" + (*city)
		location = append(location, cityPrefix)
	}
	if youtubeAccount.State != nil && *youtubeAccount.State != "" {
		state = youtubeAccount.State
		locationList = append(locationList, *state)
		statePrefix := "state_" + (*state)
		location = append(location, statePrefix)
	}
	if youtubeAccount.Country != nil && *youtubeAccount.Country != "" {
		countryCode = youtubeAccount.Country
		locationList = append(locationList, manager.ConvertCountryToCode(*countryCode))
		countryPrefix := "country_" + *countryCode
		location = append(location, countryPrefix)
	}
	if youtubeAccount.Followers != nil {
		followers := youtubeAccount.Followers
		var label string
		if *followers < 1000 {
			label = "Nano"
		} else if *followers < 75000 {
			label = "Micro"
		} else if *followers < 1000000 {
			label = "Macro"
		} else {
			label = "Mega"
		}
		youtubeLabel = &label
	}
	if campaignProfileEntry != nil {
		ytProfile := manager.GetProfileFromYoutubeEntity(ctx, *youtubeAccount, string(constants.GCCPlatform))
		searchPhraseAdmin, keywordsAdmin = manager2.CreateKeywordsAndSearchPhraseAdmin(campaignProfileEntry.AdminDetails, campaignProfileEntry.UserDetails, string(constants.YoutubePlatform), ytProfile)
		if searchPhraseAdmin == "" && len(keywordsAdmin) == 0 {
			return keywords, totalPhrase, enabledForSaas, location, locationList, youtubeLabel, nil, nil
		}
		return keywords, totalPhrase, enabledForSaas, location, locationList, youtubeLabel, &keywordsAdmin, &searchPhraseAdmin
	}

	return keywords, totalPhrase, enabledForSaas, location, locationList, youtubeLabel, &keywords, totalPhrase
}

func updateYoutubeEntity(youtubeAccount *dao.YoutubeAccountEntity, keywords []string, totalPhrase *string, enabledForSaas bool, location []string, locationList []string, youtubeLabel *string, keywordsAdmin *[]string, searchPhraseAdmin *string) *dao.YoutubeAccountEntity {
	if len(keywords) != 0 {
		youtubeAccount.Keywords = pq.StringArray(keywords)
		youtubeAccount.SearchPhrase = totalPhrase
	}
	if keywordsAdmin != nil {
		youtubeAccount.KeywordsAdmin = pq.StringArray(*keywordsAdmin)
	}
	if searchPhraseAdmin != nil {
		youtubeAccount.SearchPhraseAdmin = searchPhraseAdmin
	}
	youtubeAccount.EnabledForSaas = enabledForSaas
	youtubeAccount.Location = location
	youtubeAccount.LocationList = locationList
	youtubeAccount.Label = youtubeLabel
	youtubeAccount.Enabled = true
	flag := false
	youtubeAccount.Deleted = &flag
	return youtubeAccount
}

type blankTransformer struct{}

func (t blankTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	switch typ.Kind() {
	case reflect.Ptr:
		return func(dst, src reflect.Value) error {
			if src.Kind() == reflect.Ptr && dst.Kind() == reflect.Ptr && dst.IsNil() && !src.IsNil() {
				dst.Elem().SetString(src.Elem().String())
			} else if src.Kind() == reflect.Ptr && dst.Kind() == reflect.Ptr && dst.IsNil() && src.IsNil() {
				return nil
			} else {
				if dst.Elem().Kind() == reflect.Slice || dst.Elem().Kind() == reflect.Bool || dst.Elem().Kind() == reflect.Int64 || dst.Elem().Kind() == reflect.Float64 {
					dst.Elem().Set(dst.Elem())
				} else {
					dst.Elem().SetString(dst.Elem().String())
				}
			}
			return nil
		}
	case reflect.Float64:
		return func(dst, src reflect.Value) error {
			dst.SetFloat(dst.Float())
			return nil
		}
	default:
		return nil
	}
}
