package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetTableMonitoring(TableMonitoringId string) (*TableMonitoring, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/table-monitoring/%s", c.HostURL, TableMonitoringId), nil)
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

	TableMonitoringJob := TableMonitoring{}
	err = json.Unmarshal(body, &TableMonitoringJob)
	if err != nil {
		return nil, err
	}

	return &TableMonitoringJob, nil
}

func (c *Client) CreateTableMonitoring(creationRequest TableMonitoring) (*TableMonitoring, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/table-monitoring", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateTableMonitoring(TableMonitoringId string, creationRequest TableMonitoring) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/table-monitoring/%s", c.HostURL, TableMonitoringId), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteTableMonitoring(TableMonitoringId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/table-monitoring/%s", c.HostURL, TableMonitoringId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
