package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetEntityMapping(entityID string) (*EntityMapping, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/entity-mapping/%s", c.HostURL, entityID), nil)
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

	entity := EntityMapping{}
	err = json.Unmarshal(body, &entity)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (c *Client) CreateEntityMapping(creationRequest EntityMapping) (*EntityMapping, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/entity-mapping", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateEntityMapping(entityID string, creationRequest EntityMapping) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/entity-mapping/%s", c.HostURL, entityID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteEntityMapping(entityID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/entity-mapping/%s", c.HostURL, entityID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
