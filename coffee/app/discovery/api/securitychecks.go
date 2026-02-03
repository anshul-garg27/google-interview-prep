package api

import (
	"coffee/app/discovery/domain"
	"coffee/constants"
	"coffee/core/appcontext"
	coredomain "coffee/core/domain"
	"strings"

	"context"

	partnermanager "coffee/app/partnerusage/manager"

	"github.com/pkg/errors"
)

func validateDiscoverySearch(ctx context.Context, page int, size int, searchQuery coredomain.SearchQuery, blockFreeUsers bool, blockSaasUsers bool) error {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if ((appCtx.PlanType == nil || *appCtx.PlanType == constants.FreePlan) && (appCtx.PartnerId == nil || *appCtx.PartnerId != -1)) && blockFreeUsers {
		err := errors.New("invalid Request")
		return err
	} else if ((appCtx.PlanType == nil || *appCtx.PlanType == constants.SaasPlan) && (appCtx.PartnerId == nil || *appCtx.PartnerId != -1)) && blockSaasUsers {
		err := errors.New("invalid Request")
		return err
	}
	return nil
}

func checkTimeSeriesValidityForFreeUsers(ctx context.Context, profilePageTrackKey string) error {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)

	if appCtx.PlanType == nil || *appCtx.PlanType == constants.FreePlan {
		err := errors.New("invalid Request")
		return err
	} else if appCtx.PartnerId != nil && *appCtx.PartnerId != -1 {
		partnerManager := partnermanager.CreatePartnerUsageManager(ctx)
		trackEntry, err := partnerManager.FindProfilePageByKey(ctx, *appCtx.PartnerId, profilePageTrackKey)
		if err != nil && trackEntry == nil {
			activityMeta := ""
			_, err := partnerManager.LogPaidActivity(ctx, constants.ProfilePage, constants.ProfilePageTimeSeriesActivity, int64(1), *appCtx.PartnerId, *appCtx.AccountId, "", activityMeta, profilePageTrackKey, "PAID")
			if err != nil {
				return errors.New("invalid request")
			}
		}
	}

	return nil
}
func nullifyAudienceFields(audience domain.SocialProfileAudienceInfoEntry) domain.SocialProfileAudienceInfoEntry {

	audience.CityWiseAudienceLocation = nil
	audience.AudienceAgeGenderSplit.Female = nil
	audience.AudienceAgeGenderSplit.Male = nil
	audience.AudienceReachabilityPercentage = nil
	audience.AudienceAuthenticityPercentage = nil
	audience.QualityAudienceScore = nil
	audience.QualityScoreBreakup = nil

	return audience
}

func blockAudienceDataForFreeUsers(ctx context.Context, audienceData []domain.SocialProfileAudienceInfoEntry, profilePageTrackKey string) []domain.SocialProfileAudienceInfoEntry {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PlanType == nil || *appCtx.PlanType == constants.FreePlan {
		for key := range audienceData {
			audienceData[key] = nullifyAudienceFields(audienceData[key])
		}
	} else if appCtx.PartnerId != nil && *appCtx.PartnerId != -1 {
		partnerManager := partnermanager.CreatePartnerUsageManager(ctx)
		for key := range audienceData {
			_, err := partnerManager.FindProfilePageByKey(ctx, *appCtx.PartnerId, profilePageTrackKey)
			if err != nil {
				audienceData[key] = nullifyAudienceFields(audienceData[key])
			}
		}
	}
	return audienceData
}

func nullifyOverviewFields(profile coredomain.Profile) coredomain.Profile {

	profile.CategoryRank = nil
	profile.EstPostPrice = nil
	profile.Metrics.CountryRank = nil
	profile.Metrics.StoryReach = nil
	profile.Metrics.ImageReach = nil
	profile.Metrics.ReelsReach = nil
	profile.CategoryRank = nil
	profile.EstPostPrice = nil
	profile.Metrics.CountryRank = nil
	profile.Metrics.StoryReach = nil
	profile.Metrics.ImageReach = nil
	profile.Metrics.ReelsReach = nil
	profile.Metrics.LikesSpreadPercentage = nil
	if profile.SimilarProfileGroupData != nil {
		profile.SimilarProfileGroupData.ErGraph = nil
		profile.SimilarProfileGroupData.GroupAvgLikesSpread = nil
	}
	if profile.Grades != nil {
		profile.Grades.LikesSpreadGrade = nil
		profile.Grades.AudienceQualityGrade = nil
		profile.Grades.AudienceReachabilityGrade = nil
		profile.Grades.AudienceAuthencityGrade = nil
	}

	return profile
}

func blockOverviewDataForFreeUsers(ctx context.Context, profiles []coredomain.Profile) []coredomain.Profile {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PlanType == nil || *appCtx.PlanType == constants.FreePlan {
		for key := range profiles {
			profiles[key] = nullifyOverviewFields(profiles[key])
		}
	} else if appCtx.PartnerId != nil && *appCtx.PartnerId != -1 {
		partnerManager := partnermanager.CreatePartnerUsageManager(ctx)
		for key := range profiles {
			profilePageTrackKey := "profile-" + profiles[key].Platform + "-" + profiles[key].PlatformCode
			_, err := partnerManager.FindProfilePageByKey(ctx, *appCtx.PartnerId, profilePageTrackKey)

			if err != nil {
				profiles[key] = nullifyOverviewFields(profiles[key])
			}
		}
	}
	return profiles
}

func nullifyProfileDataFields(profile coredomain.Profile) coredomain.Profile {
	profile.CategoryRank = nil
	profile.EstPostPrice = nil
	profile.Grades.LikesSpreadGrade = nil
	profile.Grades.AudienceQualityGrade = nil
	profile.Grades.AudienceReachabilityGrade = nil
	profile.Grades.AudienceAuthencityGrade = nil
	return profile
}

func blockProfileDataForFreeUsers(ctx context.Context, profiles []coredomain.Profile) []coredomain.Profile {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if appCtx.PlanType == nil || *appCtx.PlanType == constants.FreePlan {
		for key := range profiles {
			profiles[key] = nullifyProfileDataFields(profiles[key])
		}
	} else if appCtx.PartnerId != nil && *appCtx.PartnerId != -1 {

		partnerManager := partnermanager.CreatePartnerUsageManager(ctx)
		for key := range profiles {
			profilePageTrackKey := "profile-" + profiles[key].Platform + "-" + profiles[key].PlatformCode
			_, err := partnerManager.FindProfilePageByKey(ctx, *appCtx.PartnerId, profilePageTrackKey)
			if err != nil {
				profiles[key] = nullifyProfileDataFields(profiles[key])
			}
		}
	}
	return profiles

}

func checkForPaidRequestParams(searchQuery coredomain.SearchQuery, page int, size int) (bool, bool) {
	blockedFiltersForFreeUsers := map[string]bool{
		"audience_location":           true, // no
		"audience_age":                true, //no
		"audience_gender":             true, //no
		"languages":                   true, //no
		"flag_contact_info_available": true,
	}
	blockedFiltersForSaaSUsers := map[string]bool{
		"audience_location": true, // no
		"audience_age":      true, //no
		"audience_gender":   true, //no
		"languages":         true, //no
	}
	blockFreeUser := false
	blockSaasUser := false
	for _, filter := range searchQuery.Filters {
		val := strings.Split(filter.Field, ".")
		if blockedFiltersForFreeUsers[val[0]] {
			blockFreeUser = true
		}
	}
	for _, filter := range searchQuery.Filters {
		val := strings.Split(filter.Field, ".")
		if blockedFiltersForSaaSUsers[val[0]] {
			blockSaasUser = true
		}
	}

	if page > 2 {
		blockFreeUser = true
	}

	return blockFreeUser, blockSaasUser
}
