package socialstream

import (
	"init.bulbul.tv/bulbul-backend/event-grpc/client"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/input"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/response"
	"init.bulbul.tv/bulbul-backend/event-grpc/context"
	"strconv"
)

type Client struct {
	*client.BaseClient
}

func New(ctx *context.Context) *Client {
	return &Client{client.New(ctx)}
}

func (c *Client) SearchInfluencerCampaigns(input input.SearchQuery, page int, size int, fetchDetails bool) (*response.InfluencerProgramCampaignResponse, error) {
	response := &response.InfluencerProgramCampaignResponse{}
	_, err := c.GenerateRequest(c.Context.DeviceAuthorization).
		SetBody(input).
		SetQueryParam("page", strconv.Itoa(page)).
		SetQueryParam("size", strconv.Itoa(size)).
		SetQueryParam("fetchDetails", strconv.FormatBool(fetchDetails)).
		SetResult(response).
		Post(c.Context.Config.SOCIAL_STREAM_URL + "/influencerprogramcampaign/searchV2")

	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}