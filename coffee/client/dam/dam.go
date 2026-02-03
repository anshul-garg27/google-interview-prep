package dam

import (
	"bytes"
	"coffee/core/client"
	"coffee/core/domain"
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type Client struct {
	*client.BaseClient
}

func New(ctx context.Context) *Client {
	return &Client{client.New(ctx, "DAM_URL")}
}

func (c *Client) FindAssetById(assetId int) (*domain.AssetInfoResponse, error) {
	response := &domain.AssetInfoResponse{}
	_, err := c.GenerateRequest().
		SetResult(response).
		SetPathParams(map[string]string{"id": strconv.Itoa(assetId)}).
		Get(c.BaseUrl + "/dam/{id}")
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func (c *Client) InitUpload(input *GenerateUrlReqEntry) (*GenerateUrlResponse, error) {
	response := &GenerateUrlResponse{}
	_, err := c.GenerateRequest().
		SetBody(input).
		SetResult(response).
		Post(c.BaseUrl + "/dam/upload/init")

	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func (c *Client) FinishUpload(path string) (*domain.AssetInfoResponse, error) {
	response := &domain.AssetInfoResponse{}
	_, err := c.GenerateRequest().
		SetBody(&GenerateUrlRespEntry{Path: path}).SetResult(response).
		Post(c.BaseUrl + "/dam/upload/finish")

	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func (c *Client) UploadFileToStorage(url string, fileName string, data []byte) (bool, error) {
	reader := bytes.NewReader(data)
	req, err := http.NewRequest("PUT", url, reader)
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "text/csv")
	req.Header.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fileName))
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *Client) DamUpload(input *GenerateUrlReqEntry, data []byte, fileName string) (*domain.AssetInfoResponse, error) {
	resp, err := c.InitUpload(input)
	if err != nil {
		return nil, err
	}
	done, err := c.UploadFileToStorage(resp.Asset.URL, fileName, data)
	if err != nil || done == false {
		return nil, err
	}
	response, err := c.FinishUpload(resp.Asset.Path)
	if err != nil {
		return nil, err
	}
	return response, nil
}
