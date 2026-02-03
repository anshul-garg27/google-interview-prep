package identity

import (
	"init.bulbul.tv/bulbul-backend/event-grpc/client"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/input"
	"init.bulbul.tv/bulbul-backend/event-grpc/client/response"
	"init.bulbul.tv/bulbul-backend/event-grpc/context"
)

type Client struct {
	*client.BaseClient
}

func New(ctx *context.Context) *Client {
	return &Client{client.New(ctx)}
}

func (c *Client) SearchUserClientAccount(input *input.SearchQuery) (*response.UserClientAccountResponse, error) {
	response := &response.UserClientAccountResponse{}
	_, err := c.GenerateRequest(c.Context.ClientAuthorization).
		SetBody(input).
		SetResult(response).
		Post(c.Context.Config.IDENTITY_URL + "/userclientaccount/search")

	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}