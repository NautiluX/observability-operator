package framework

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/common/model"
)

// PrometheusClient is an HTTP-based client for querying a Prometheus server
type PrometheusClient struct {
	baseURL string
	client  *http.Client
}

// PrometheusResponse is used to contain prometheus query results
type PrometheusResponse struct {
	Status string                 `json:"status"`
	Error  string                 `json:"error"`
	Data   prometheusResponseData `json:"data"`
}

type prometheusResponseData struct {
	ResultType string       `json:"resultType"`
	Result     model.Vector `json:"result"`
}

func NewPrometheusClient(url string) *PrometheusClient {
	return &PrometheusClient{
		baseURL: url,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *PrometheusClient) Query(query string) (*PrometheusResponse, error) {
	url := fmt.Sprintf("%s/api/v1/query?query=%s", c.baseURL, query)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to query Prometheus: %v", err)
	}

	var result PrometheusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("unable to parse query response: %v", err)
	}

	return &result, nil
}
