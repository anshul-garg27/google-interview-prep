package tests

import (
	"coffee/app/campaignprofiles/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_CreateInstagramEntityUsingBeatResponse(t *testing.T) {
	profile := prepareIaProfileForTest()
	beatResponse := mockIaBeatResponse()
	expectedResponse := mockIaEntityToCheckBeatResponse()

	response := domain.CreateInstagramEntityUsingBeatResponse(profile, beatResponse)
	assert.Equal(t, expectedResponse.Name, response.Name, "name should be same")
	assert.Equal(t, expectedResponse.Handle, response.Handle, "handle should be same")
	assert.Equal(t, expectedResponse.IgID, response.IgID, "igid should be same")
	assert.Equal(t, expectedResponse.Thumbnail, response.Thumbnail, "thumbnail should be same")
	assert.Equal(t, expectedResponse.Bio, response.Bio, "bio should be same")
	assert.Equal(t, expectedResponse.Followers, response.Followers, "followers should be same")
	assert.Equal(t, expectedResponse.Following, response.Following, "following should be same")
	assert.Equal(t, expectedResponse.EngagementRate, response.EngagementRate, "engagementrate should be same")
	assert.Equal(t, expectedResponse.AvgLikes, response.AvgLikes, "avglikes should be same")
	assert.Equal(t, expectedResponse.AvgReach, response.AvgReach, "avgreach should be same")
	assert.Equal(t, expectedResponse.ReelsReach, response.ReelsReach, "reelsreach should be same")
	assert.Equal(t, expectedResponse.StoryReach, response.StoryReach, "StoryReach should be same")
	assert.Equal(t, expectedResponse.ImageReach, response.ImageReach, "ImageReach should be same")
	assert.Equal(t, expectedResponse.AvgComments, response.AvgComments, "AvgComments should be same")
	assert.Equal(t, expectedResponse.AvgReelsPlayCount, response.AvgReelsPlayCount, "AvgReelsPlayCount should be same")
	assert.Equal(t, expectedResponse.Label, response.Label, "Label should be same")
	assert.Equal(t, expectedResponse.ID, response.ID, "ID should be same")
	assert.Equal(t, expectedResponse.ProfileId, response.ProfileId, "Profile ID should be same")
	assert.Equal(t, expectedResponse.AvgViews, response.AvgViews, "Average Views should be same")
}

func TestClient_FillSocialDetailsInFullRefreshResultForIA(t *testing.T) {
	cpEntry := prepareIaCampaignProfileForTest()
	beatResponse := mockIaBeatResponse()
	expectedResponse := mockIaCpEntryToCheckBeatResponse()

	response := domain.FillSocialDetailsInFullRefreshResultForIA(cpEntry, beatResponse)
	assert.Equal(t, expectedResponse.SocialDetails.Code, response.SocialDetails.Code, "Code should be same")
	assert.Equal(t, expectedResponse.SocialDetails.ProfileCode, response.SocialDetails.ProfileCode, "ProfileCode should be same")
	assert.Equal(t, expectedResponse.SocialDetails.GCCProfileCode, response.SocialDetails.GCCProfileCode, "GCCProfileCode should be same")
	assert.Equal(t, expectedResponse.SocialDetails.Name, response.SocialDetails.Name, "Name should be same")
	assert.Equal(t, expectedResponse.SocialDetails.Platform, response.SocialDetails.Platform, "Platform should be same")
	assert.Equal(t, expectedResponse.SocialDetails.Handle, response.SocialDetails.Handle, "Handle should be same")
	assert.Equal(t, expectedResponse.SocialDetails.Username, response.SocialDetails.Username, "Username should be same")
	assert.Equal(t, expectedResponse.SocialDetails.PlatformThumbnail, response.SocialDetails.PlatformThumbnail, "PlatformThumbnail should be same")
	assert.Equal(t, expectedResponse.SocialDetails.EngagementRatePercenatge, response.SocialDetails.EngagementRatePercenatge, "EngagementRatePercentage should be same")
	assert.Equal(t, expectedResponse.SocialDetails.Followers, response.SocialDetails.Followers, "Followers should be same")
	assert.Equal(t, expectedResponse.SocialDetails.Following, response.SocialDetails.Following, "Following should be same")
	assert.Equal(t, expectedResponse.SocialDetails.AvgLikes, response.SocialDetails.AvgLikes, "AvgLikes should be same")
	assert.Equal(t, expectedResponse.SocialDetails.AvgComments, response.SocialDetails.AvgComments, "AvgComments should be same")
	assert.Equal(t, expectedResponse.SocialDetails.Uploads, response.SocialDetails.Uploads, "Uploads should be same")
	assert.Equal(t, expectedResponse.SocialDetails.AvgReelsPlay30d, response.SocialDetails.AvgReelsPlay30d, "AvgReelsPlay30d should be same")
	assert.Equal(t, expectedResponse.SocialDetails.AvgReach, response.SocialDetails.AvgReach, "AvgReach should be same")
	assert.Equal(t, expectedResponse.SocialDetails.StoryReach, response.SocialDetails.StoryReach, "StoryReach should be same")
	assert.Equal(t, expectedResponse.SocialDetails.ImageReach, response.SocialDetails.ImageReach, "ImageReach should be same")
	assert.Equal(t, expectedResponse.SocialDetails.ReelsReach, response.SocialDetails.ReelsReach, "ReelsReach should be same")
	assert.Equal(t, expectedResponse.SocialDetails.IsVerified, response.SocialDetails.IsVerified, "IsVerified should be same")
	assert.Equal(t, expectedResponse.SocialDetails.Source, response.SocialDetails.Source, "Source should be same")
	assert.Equal(t, expectedResponse.SocialDetails.AccountType, response.SocialDetails.AccountType, "AccountType should be same")
	assert.Equal(t, expectedResponse.SocialDetails.AudienceCityAvailable, response.SocialDetails.AudienceCityAvailable, "AudienceCityAvailable should be same")
	assert.Equal(t, expectedResponse.SocialDetails.AudienceGenderAvailable, response.SocialDetails.AudienceGenderAvailable, "AudienceGenderAvailable should be same")
	assert.Equal(t, expectedResponse.SocialDetails.IgId, response.SocialDetails.IgId, "IgId should be same")
}

//---------------------------------------------------------------------------------------------
//---------------------------------------------------------------------------------------------

func TestClient_FillSocialDetailsInFullRefreshResultForYA(t *testing.T) {
	cpEntry := prepareYaCampaignProfileForTest()
	beatResponse := mockYaBeatResponse()
	expectedResponse := mockYaCpEntryToCheckBeatResponse()

	response := domain.FillSocialDetailsInFullRefreshResultForYA(cpEntry, beatResponse)
	assert.Equal(t, expectedResponse.SocialDetails.Code, response.SocialDetails.Code, "Code should be same")
	assert.Equal(t, expectedResponse.SocialDetails.ProfileCode, response.SocialDetails.ProfileCode, "ProfileCode should be same")
	assert.Equal(t, expectedResponse.SocialDetails.GCCProfileCode, response.SocialDetails.GCCProfileCode, "GCCProfileCode should be same")
	assert.Equal(t, expectedResponse.SocialDetails.Name, response.SocialDetails.Name, "Name should be same")
	assert.Equal(t, expectedResponse.SocialDetails.Platform, response.SocialDetails.Platform, "Platform should be same")
	assert.Equal(t, expectedResponse.SocialDetails.Handle, response.SocialDetails.Handle, "Handle should be same")
	assert.Equal(t, expectedResponse.SocialDetails.Username, response.SocialDetails.Username, "Username should be same")
	assert.Equal(t, expectedResponse.SocialDetails.PlatformThumbnail, response.SocialDetails.PlatformThumbnail, "PlatformThumbnail should be same")
	assert.Equal(t, expectedResponse.SocialDetails.EngagementRatePercenatge, response.SocialDetails.EngagementRatePercenatge, "EngagementRatePercentage should be same")
	assert.Equal(t, expectedResponse.SocialDetails.Followers, response.SocialDetails.Followers, "Followers should be same")
	assert.Equal(t, expectedResponse.SocialDetails.Uploads, response.SocialDetails.Uploads, "Uploads should be same")
	assert.Equal(t, expectedResponse.SocialDetails.AvgReach, response.SocialDetails.AvgReach, "AvgReach should be same")
	assert.Equal(t, expectedResponse.SocialDetails.Source, response.SocialDetails.Source, "Source should be same")
	assert.Equal(t, expectedResponse.SocialDetails.AudienceCityAvailable, response.SocialDetails.AudienceCityAvailable, "AudienceCityAvailable should be same")
	assert.Equal(t, expectedResponse.SocialDetails.AudienceGenderAvailable, response.SocialDetails.AudienceGenderAvailable, "AudienceGenderAvailable should be same")
	assert.Equal(t, expectedResponse.SocialDetails.Views, response.SocialDetails.Views, "Views should be same")
}

func TestClient_CreateYoutubeEntityUsingBeatResponse(t *testing.T) {
	profile := prepareYaProfileForTest()
	beatResponse := mockYaBeatResponse()
	expectedResponse := mockYAEntityToCheckBeatResponse()

	response := domain.CreateYoutubeEntityUsingBeatResponse(profile, beatResponse)
	assert.Equal(t, expectedResponse.ID, response.ID, "ID should be same")
	assert.Equal(t, expectedResponse.ProfileId, response.ProfileId, "ProfileId should be same")
	assert.Equal(t, expectedResponse.GccProfileId, response.GccProfileId, "GccProfileId should be same")
	assert.Equal(t, expectedResponse.Username, response.Username, "Username should be same")
	assert.Equal(t, expectedResponse.Description, response.Description, "Description should be same")
	assert.Equal(t, expectedResponse.VideoReach, response.VideoReach, "VideoReach should be same")
	assert.Equal(t, expectedResponse.AvgViews, response.AvgViews, "Description should be same")
	assert.Equal(t, expectedResponse.ChannelId, response.ChannelId, "channelId should be same")
	assert.Equal(t, expectedResponse.Title, response.Title, "title should be same")
	assert.Equal(t, expectedResponse.Thumbnail, response.Thumbnail, "Thumbnail should be same")
	assert.Equal(t, expectedResponse.UploadsCount, response.UploadsCount, "UploadsCount should be same")
	assert.Equal(t, expectedResponse.Followers, response.Followers, "followers should be same")
	assert.Equal(t, expectedResponse.ViewsCount, response.ViewsCount, "ViewsCount should be same")
	assert.Equal(t, expectedResponse.Label, response.Label, "Label should be same")
}
