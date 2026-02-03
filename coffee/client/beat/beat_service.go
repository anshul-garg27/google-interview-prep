package beatservice

import (
	"coffee/core/client"
	"context"
	"strings"
	"time"
)

type Client struct {
	*client.BaseClient
}

func New(ctx context.Context) *Client {
	return &Client{
		client.New(ctx, "BEAT_SERVICE")}
}

func (c *Client) RecentPostsByPlatformAndPlatformId(platform string, platformId string) (*RecentPostsResponse, error) {
	response := &RecentPostsResponse{}
	_, err := c.GenerateRequest().
		SetResult(response).
		SetPathParams(map[string]string{"platform": strings.ToLower(platform), "id": platformId}).
		Get(c.BaseUrl + "/recent/posts/{platform}/byprofileid/{id}")
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func (c *Client) FindProfileByHandle(platform string, handle string) (*ProfileResponse, error) {
	response := &ProfileResponse{}
	c.SetTimeout(1 * time.Minute)
	_, err := c.GenerateRequest().
		SetResult(response).
		SetPathParams(map[string]string{"platform": strings.ToLower(platform), "handle": handle}).
		Get(c.BaseUrl + "/profiles/{platform}/byhandle/{handle}")
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func (c *Client) FindYoutubeChannelIdByHandle(handle string) (*YoutubeChannelResponse, error) {
	response := &YoutubeChannelResponse{}
	c.SetTimeout(30 * time.Second)
	_, err := c.GenerateRequest().
		SetResult(response).
		SetPathParams(map[string]string{"handle": handle}).
		Get(c.BaseUrl + "/youtube/channel/byhandle/{handle}")
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func (c *Client) FetchCompleteProfileByHandle(platform string, handle string) (*CompleteProfileResponse, error) {
	response := &CompleteProfileResponse{}
	c.SetTimeout(1 * time.Minute)
	_, err := c.GenerateRequest().
		SetResult(response).
		SetPathParams(map[string]string{"platform": strings.ToLower(platform), "handle": handle}).
		Get(c.BaseUrl + "/profiles/{platform}/byhandle/{handle}?full_refresh=1")
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func (c *Client) FindSocialProfileByProfileId(platform string, profileId string) (*CompleteProfileResponse, error) {
	response := &CompleteProfileResponse{}
	c.SetTimeout(1 * time.Minute)
	_, err := c.GenerateRequest().
		SetResult(response).
		SetPathParams(map[string]string{"platform": strings.ToLower(platform), "profileId": profileId}).
		Get(c.BaseUrl + "/profiles/{platform}/byid/{profileId}")
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func (c *Client) FindCreatorInsightsByHandle(handle string, token string, userId string) (*CompleteProfileResponse, error) {
	response := &CompleteProfileResponse{}
	c.SetTimeout(1 * time.Minute)
	_, err := c.GenerateRequest().
		SetResult(response).
		SetPathParams(map[string]string{"handle": handle, "token": token, "userId": userId}).
		Get(c.BaseUrl + "/profiles/INSTAGRAM/byhandle/{handle}/insights?token={token}&user_id={userId}")
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func (c *Client) FindCreatorAudienceInsightsByHandle(handle string, token string, userId string) (*CreatorInsight, error) {
	response := &CreatorInsight{}
	c.SetTimeout(1 * time.Minute)
	_, err := c.GenerateRequest().
		SetResult(response).
		SetPathParams(map[string]string{"handle": handle, "token": token, "userId": userId}).
		Get(c.BaseUrl + "/profiles/INSTAGRAM/byhandle/{handle}/audienceinsights?token={token}&user_id={userId}")
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}
