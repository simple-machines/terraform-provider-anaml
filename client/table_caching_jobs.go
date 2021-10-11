package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetTableCaching(TableCachingId string) (*TableCaching, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/table-caching/%s", c.HostURL, TableCachingId), nil)
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

	TableCachingJob := TableCaching{}
	err = json.Unmarshal(body, &TableCachingJob)
	if err != nil {
		return nil, err
	}

	return &TableCachingJob, nil
}

func (c *Client) CreateTableCaching(creationRequest TableCaching) (*TableCaching, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/table-caching", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateTableCaching(TableCachingId string, creationRequest TableCaching) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/table-caching/%s", c.HostURL, TableCachingId), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteTableCaching(TableCachingId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/table-caching/%s", c.HostURL, TableCachingId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
