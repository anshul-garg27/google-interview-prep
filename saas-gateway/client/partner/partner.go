package partner

import (
	"time"

	"init.bulbul.tv/bulbul-backend/saas-gateway/client"
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

func (c *Client) FindPartnerById(partnerId string) (*response.PartnerResponse, error) {
	response := &response.PartnerResponse{}
	_, err := c.GenerateRequest(c.Context.ClientAuthorization).
		SetPathParams(map[string]string{"partnerId": partnerId}).
		SetResult(response).
		Get(c.Context.Config.PARTNER_URL + "/partner/{partnerId}")

	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}
