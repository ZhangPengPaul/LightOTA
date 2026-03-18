package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ZhangPengPaul/LightOTA/internal/model"
)

type ThirdPartyClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewThirdPartyClient(baseURL, apiKey string) *ThirdPartyClient {
	return &ThirdPartyClient{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *ThirdPartyClient) GetDevice(deviceID string) (*model.ThirdPartyDevice, error) {
	url := fmt.Sprintf("%s/api/v1/devices/%s", c.baseURL, deviceID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result model.ThirdPartyResponse
	result.Data = &model.ThirdPartyDevice{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("third party api returned error code %d", result.Code)
	}

	device, ok := result.Data.(*model.ThirdPartyDevice)
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	return device, nil
}

func (c *ThirdPartyClient) QueryDevices(req *model.ThirdPartyDeviceQueryRequest) ([]model.ThirdPartyDevice, int, error) {
	url := fmt.Sprintf("%s/api/v1/devices/query", c.baseURL)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, 0, err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var result model.ThirdPartyDeviceQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	if result.Code != 0 {
		return nil, 0, fmt.Errorf("third party api returned error code %d", result.Code)
	}

	return result.Data.Selected, result.Data.Total, nil
}
