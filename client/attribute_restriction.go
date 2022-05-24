package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetAttributeRestriction(attributeID string) (*AttributeRestriction, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/attribute/%s", c.HostURL, attributeID), nil)
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

	attribute := AttributeRestriction{}
	err = json.Unmarshal(body, &attribute)
	if err != nil {
		return nil, err
	}

	return &attribute, nil
}

func (c *Client) CreateAttributeRestriction(creationRequest AttributeRestriction) (*AttributeRestriction, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/attribute", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateAttributeRestriction(attributeID string, creationRequest AttributeRestriction) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/attribute/%s", c.HostURL, attributeID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteAttributeRestriction(attributeID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/attribute/%s", c.HostURL, attributeID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
