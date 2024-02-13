package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetMetricsJob(MetricsJobID string) (*MetricsJob, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/metrics-job/%s", c.HostURL, MetricsJobID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	if body == nil {
		return nil, nil
	}

	MetricsJob := MetricsJob{}
	err = json.Unmarshal(body, &MetricsJob)
	if err != nil {
		return nil, err
	}

	return &MetricsJob, nil
}

func (c *Client) CreateMetricsJob(creationRequest MetricsJob) (*MetricsJob, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/metrics-job", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var V int
	err = json.Unmarshal(body, &V)
	if err != nil {
		return nil, err
	}

	creationRequest.ID = V
	return &creationRequest, nil
}

func (c *Client) UpdateMetricsJob(MetricsJobID string, creationRequest MetricsJob) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/metrics-job/%s", c.HostURL, MetricsJobID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteMetricsJob(MetricsJobID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/metrics-job/%s", c.HostURL, MetricsJobID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
