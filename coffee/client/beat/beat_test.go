package beatservice

import (
	"coffee/constants"
	"coffee/core/appcontext"
	"context"
	"fmt"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestClient_RecentPostsByPlatformAndPlatformId(t *testing.T) {
	ctx, err := setupTest()
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping beat tests, set environment variable INTEGRATION")
	}
	response, err := New(ctx).RecentPostsByPlatformAndPlatformId("instagram", "2094200507")
	if err != nil {
		log.Error(err)
		log.Fatal(fmt.Sprintf("failed: %+v", response))
	} else {
		log.Info(fmt.Sprintf("success: %+v", response))
		expectedProfileId := "2094200507"
		assert.Equal(t, expectedProfileId, response.Data.ProfileId, "Profile ID should be same")
	}
}

func TestClient_FindYoutubeChannelIdByHandle(t *testing.T) {
	ctx, err := setupTest()
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping beat tests, set environment variable INTEGRATION")
	}
	channelResponse, err := New(ctx).FindYoutubeChannelIdByHandle("@tseries")
	if err != nil {
		log.Error(err)
		log.Fatal(fmt.Sprintf("failed: %+v", channelResponse))
	} else {
		log.Info(fmt.Sprintf("success: %+v", channelResponse))
		expectedChannelId := "UCq-Fj5jknLsUf-MWSy4_brA"
		assert.Equal(t, expectedChannelId, channelResponse.YoutubeChannel.ChannelId, "Channel ID of T-Series should be correct.")
	}
}

func TestClient_FindProfileByHandle(t *testing.T) {
	ctx, err := setupTest()
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping beat tests, set environment variable INTEGRATION")
	}
	profileResponse, err := New(ctx).FindProfileByHandle("INSTAGRAM", "virat.kohli")
	if err != nil {
		log.Error(err)
		log.Fatal(fmt.Sprintf("failed: %+v", profileResponse))
	} else {
		log.Info(fmt.Sprintf("success: %+v", profileResponse))
		expectedHandle := "virat.kohli"
		assert.Equal(t, expectedHandle, profileResponse.Profile.Handle, "Handle in response should be same.")
	}
}

func TestClient_FindCreatorInsightsByHandle(t *testing.T) {
	ctx, err := setupTest()
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping beat tests, set environment variable INTEGRATION")
	}
	profileResponse, err := New(ctx).FindCreatorInsightsByHandle("little_princess_maira", "EAAEYVETYlrgBAGYn6yn3ZAZCYZBenMDeP3DTZAg2h35mvNHSmBZAIPjTEdHYJPIOqZAS8l755Lj7SyMBwUci3ur1mkkZAS5nDstJHVIIptVmDGC8NHmXNtpI5t1Fb7JOAqiv2ZANEUaPeX0JsLUKzofyFIsnOrZAFmqIM9CPZBP5sMdkNkH8m1zVtLJyXZAR9TocZBGXGZBzZBrQ1PaQZDZD", "17841448048489396")
	if err != nil {
		log.Fatal(fmt.Sprintf("failed: %+v", profileResponse))
	} else {
		log.Info(fmt.Sprintf("success: %+v", profileResponse))
		expectedHandle := "little_princess_maira"
		assert.Equal(t, expectedHandle, profileResponse.CompleteProfile.Handle, "Handle in response should be same.")
	}
}

func TestClient_FindCreatorAudienceInsightsByHandle(t *testing.T) {
	ctx, err := setupTest()
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping beat tests, set environment variable INTEGRATION")
	}
	profileResponse, err := New(ctx).FindCreatorAudienceInsightsByHandle("little_princess_maira", "EAAEYVETYlrgBAGYn6yn3ZAZCYZBenMDeP3DTZAg2h35mvNHSmBZAIPjTEdHYJPIOqZAS8l755Lj7SyMBwUci3ur1mkkZAS5nDstJHVIIptVmDGC8NHmXNtpI5t1Fb7JOAqiv2ZANEUaPeX0JsLUKzofyFIsnOrZAFmqIM9CPZBP5sMdkNkH8m1zVtLJyXZAR9TocZBGXGZBzZBrQ1PaQZDZD", "17841448048489396")
	if err != nil {
		log.Fatal(fmt.Sprintf("failed: %+v", profileResponse))
	} else {
		log.Info(fmt.Sprintf("success: %+v", profileResponse))
		expectedHandle := "little_princess_maira"
		assert.Equal(t, expectedHandle, profileResponse.Profile.Handle, "Handle in response should be same.")
	}
}

func setupTest() (context.Context, error) {
	ctx := context.Background()
	viper.SetConfigFile("../../.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err)
		log.Fatal(err)
	}
	appCtx := &appcontext.RequestContext{}
	ctx = context.WithValue(ctx, constants.AppContextKey, appCtx)
	return ctx, err
}
