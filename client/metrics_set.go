package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetMetricsSet(MetricsSetID string) (*MetricsSet, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/metrics-set/%s", c.HostURL, MetricsSetID), nil)
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

	MetricsSet := MetricsSet{}
	err = json.Unmarshal(body, &MetricsSet)
	if err != nil {
		return nil, err
	}

	return &MetricsSet, nil
}

func (c *Client) CreateMetricsSet(creationRequest MetricsSet) (*MetricsSet, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/metrics-set", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateMetricsSet(MetricsSetID string, creationRequest MetricsSet) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/metrics-set/%s", c.HostURL, MetricsSetID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteMetricsSet(MetricsSetID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/metrics-set/%s", c.HostURL, MetricsSetID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
