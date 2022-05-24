package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetLabelRestriction(labelID string) (*LabelRestriction, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/allowed-label/%s", c.HostURL, labelID), nil)
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

	label := LabelRestriction{}
	err = json.Unmarshal(body, &label)
	if err != nil {
		return nil, err
	}

	return &label, nil
}

func (c *Client) CreateLabelRestriction(creationRequest LabelRestriction) (*LabelRestriction, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/allowed-label", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateLabelRestriction(labelID string, creationRequest LabelRestriction) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/allowed-label/%s", c.HostURL, labelID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteLabelRestriction(labelID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/allowed-label/%s", c.HostURL, labelID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
