package app

import (
	campaignprofilesapi "coffee/app/campaignprofiles/api"
	collectionanalyticsapi "coffee/app/collectionanalytics/api"
	collectiongroupapi "coffee/app/collectiongroup/api"
	contentapi "coffee/app/content/api"
	discoveryapi "coffee/app/discovery/api"
	genreinsightsapi "coffee/app/genreinsights/api"
	keywordcollectionapi "coffee/app/keywordcollection/api"
	"coffee/app/leaderboard"
	postcollectionapi "coffee/app/postcollection/api"
	profilecollectionapi "coffee/app/profilecollection/api"

	partnerusageapi "coffee/app/partnerusage/api"
	"coffee/core/rest"
	"context"
)

func SetupContainer() rest.ApplicationContainer {
	container := rest.NewApplicationContainer()
	profileCollectionService := profilecollectionapi.NewProfileCollectionService()
	profileCollectionApi := profilecollectionapi.NewProfileCollectionApi(profileCollectionService)
	container.AddService(profileCollectionService, profileCollectionApi)

	postCollectionService := postcollectionapi.CreatePostCollectionService(context.Background())
	postCollectionApi := postcollectionapi.NewPostCollectionApi(postCollectionService)
	container.AddService(postCollectionService, postCollectionApi)

	socialDiscoveryService := discoveryapi.NewSocialDiscoveryService()
	socialDiscoveryAPI := discoveryapi.NewSocialDiscoveryApi(socialDiscoveryService)
	container.AddService(socialDiscoveryService, socialDiscoveryAPI)

	leaderboardService := leaderboard.NewLeaderboardService()
	leaderboardAPI := leaderboard.NewLeaderboardApi(leaderboardService)
	container.AddService(leaderboardService, leaderboardAPI)

	collectionAnalyticsService := collectionanalyticsapi.CreateCollectionAnalyticsService(context.Background())
	collectionAnalyticsApi := collectionanalyticsapi.NewPostCollectionApi(collectionAnalyticsService)
	container.AddService(collectionAnalyticsService, collectionAnalyticsApi)

	genreInsightsService := genreinsightsapi.NewGenreInsightsService(context.Background())
	genreInsightsApi := genreinsightsapi.NewGenreInsightsImplApi(genreInsightsService)
	container.AddService(genreInsightsService, genreInsightsApi)

	contentService := contentapi.NewContentService(context.Background())
	contentApi := contentapi.NewContentModuleImplApi(contentService)
	container.AddService(contentService, contentApi)

	campaignProfileService := campaignprofilesapi.NewCampaignService(context.Background())
	campaignProfileApi := campaignprofilesapi.NewCampaignApi(campaignProfileService)
	container.AddService(campaignProfileService, campaignProfileApi)

	collectionGroupService := collectiongroupapi.NewCollectionGroupService(context.Background())
	collectionGroupApi := collectiongroupapi.NewCollectionGroupApi(collectionGroupService)
	container.AddService(collectionGroupService, collectionGroupApi)

	contractService := partnerusageapi.NewPartnerUsageService()
	contractApi := partnerusageapi.NewPartnerUsageApi(contractService)
	container.AddService(contractService, contractApi)

	keywordCollectionService := keywordcollectionapi.CreateKeywordCollectionService(context.Background())
	keywordCollectionApi := keywordcollectionapi.NewKeywordCollectionApi(keywordCollectionService)
	container.AddService(keywordCollectionService, keywordCollectionApi)

	return container
}
