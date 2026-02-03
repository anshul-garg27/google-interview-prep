package webengage

import (
	"errors"
	"init.bulbul.tv/bulbul-backend/event-grpc/client"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/input"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/response"
	"init.bulbul.tv/bulbul-backend/event-grpc/context"
	"strings"
)

type Client struct {
	*client.BaseClient
}

func New(ctx *context.Context) *Client {
	return &Client{client.New(ctx)}
}

func (c *Client) DispatchWebengageEvent(event *input.WebengageEvent) (*response.WebengageResponse, error) {

	apiKey := c.Context.Config.WEBENGAGE_API_KEY
	postUrl := c.Context.Config.WEBENGAGE_URL
	if event.AppType == nil {
		if strings.EqualFold(event.EventName, "Extension Profile Visited") || strings.HasPrefix(event.EventName, "VS -") {
			appType := input.WEBENGAGE_APP_TYPE_BRAND
			event.AppType = &appType
		} else {
			return nil, errors.New("Invalid app type")
		}
	}

	if *event.AppType == input.WEBENGAGE_APP_TYPE_BRAND {
		apiKey = c.Context.Config.WEBENGAGE_BRAND_API_KEY
		postUrl = c.Context.Config.WEBENGAGE_BRAND_URL
	}

	// make it nil so that it doesnt land to webengage
	event.AppType = nil

	response := &response.WebengageResponse{}
	_, err := c.GenerateRequest("Bearer " + apiKey).
		SetBody(event).
		SetResult(response).
		Post(postUrl)

	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func (c *Client) DispatchWebengageUserEvent(event *input.WebengageUserEvent) (*response.WebengageResponse, error) {

	apiKey := c.Context.Config.WEBENGAGE_API_KEY
	postUrl := c.Context.Config.WEBENGAGE_USER_URL
	if event.AppType == nil {
		return nil, errors.New("Invalid app type")
	} else if *event.AppType == input.WEBENGAGE_APP_TYPE_BRAND {
		apiKey = c.Context.Config.WEBENGAGE_BRAND_API_KEY
		postUrl = c.Context.Config.WEBENGAGE_BRAND_USER_URL
	}

	// make it nil so that it doesnt land to webengage
	event.AppType = nil

	response := &response.WebengageResponse{}
	_, err := c.GenerateRequest("Bearer " + apiKey).
		SetBody(event).
		SetResult(response).
		Post(postUrl)

	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}
