package partnerservice

import (
	"coffee/constants"
	"coffee/core/appcontext"
	"coffee/core/client"
	"context"
	"strconv"
)

type Client struct {
	*client.BaseClient
}

func New(ctx context.Context) *Client {
	return &Client{client.New(ctx, "PARTNER_SERVICE")}
}

func (c *Client) GetPartnerContract(partnerId int64) (*PartnerResponse, error) {
	appCtx := c.Context.Value(constants.AppContextKey).(*appcontext.RequestContext)
	response := &PartnerResponse{}
	_, err := c.GenerateAuthorizedRequest(appCtx.Authorization).
		SetResult(response).
		SetPathParams(map[string]string{"partnerId": strconv.FormatInt(partnerId, 10)}).
		Get(c.BaseUrl + "/product-service/api/partner/contract/partnerId/{partnerId}/byType/SAAS")
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}
