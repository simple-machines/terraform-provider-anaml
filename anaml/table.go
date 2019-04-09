package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetTable(tableID string) (*Table, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/table/%s", c.HostURL, tableID), nil)
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

	table := Table{}
	err = json.Unmarshal(body, &table)
	if err != nil {
		return nil, err
	}

	return &table, nil
}

func (c *Client) CreateTable(creationRequest Table) (*Table, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/table", c.HostURL), strings.NewReader(string(rb)))
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

	creationRequest.Id = V
	return &creationRequest, nil
}

func (c *Client) UpdateTable(tableID string, creationRequest Table) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/table/%s", c.HostURL, tableID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteTable(tableId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/table/%s", c.HostURL, tableId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
