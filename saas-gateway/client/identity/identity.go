package identity

import (
	"strconv"
	"time"

	"init.bulbul.tv/bulbul-backend/saas-gateway/client"
	"init.bulbul.tv/bulbul-backend/saas-gateway/client/input"
	"init.bulbul.tv/bulbul-backend/saas-gateway/client/response"
	"init.bulbul.tv/bulbul-backend/saas-gateway/context"
)

type Client struct {
	*client.BaseClient
}

func New(ctx *context.Context) *Client {
	c := &Client{client.New(ctx)}
	c.SetTimeout(time.Duration(1 * time.Minute))
	return c
}

func (c *Client) UserAuth(userAuthInput *input.UserAuth) (*response.UserAuthResponse, error) {
	response := &response.UserAuthResponse{}
	_, err := c.GenerateRequest(c.Context.ClientAuthorization).
		SetBody(userAuthInput).
		SetResult(response).
		Post(c.Context.Config.IDENTITY_URL + "/auth/login")
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func (c *Client) VerifyToken(input *input.Token, skipTokenCache bool) (userAuthResponse *response.UserAuthResponse, error error) {

	response := &response.UserAuthResponse{}
	_, err := c.GenerateRequest(c.Context.ClientAuthorization).
		SetBody(input).
		SetResult(response).
		SetQueryParam("skipTokenCache", strconv.FormatBool(skipTokenCache)).
		Post(c.Context.Config.IDENTITY_URL + "/auth/verify/token?loadPreferences=true")

	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}
