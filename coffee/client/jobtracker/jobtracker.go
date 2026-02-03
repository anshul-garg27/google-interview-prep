package jobtracker

import (
	"coffee/core/client"
	"context"
	"strconv"
)

type Client struct {
	*client.BaseClient
}

func New(ctx context.Context) *Client {
	return &Client{client.New(ctx, "JOBTRACKER_URL")}
}

func (c *Client) CreateJob(jobEntry JobEntry) (*JobResponse, error) {
	response := &JobResponse{}
	c.SetCloseConnection(true)
	req := c.GenerateRequest().
		SetResult(response).
		SetBody(jobEntry)
	resp, err := req.Post(c.BaseUrl + "/job/")
	if err != nil {
		return nil, err
	} else {
		resp.RawResponse.Body.Close()
		return response, nil
	}
}

func (c *Client) UpdateJobBatch(jobbatch JobBatchEntry) (*client.SFBaseResponse, error) {
	response := &client.SFBaseResponse{}
	c.SetCloseConnection(true)
	req := c.GenerateRequest().
		SetResult(response).
		SetBody(jobbatch)
	resp, err := req.Put(c.BaseUrl + "/batch/" + strconv.FormatInt(jobbatch.Id, 10))
	if err != nil {
		return nil, err
	} else {
		resp.RawResponse.Body.Close()
		return response, nil
	}
}
