package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetEntityPopulation(entityID string) (*EntityPopulation, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/entity-population/%s", c.HostURL, entityID), nil)
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

	population := EntityPopulation{}
	err = json.Unmarshal(body, &population)
	if err != nil {
		return nil, err
	}

	return &population, nil
}

func (c *Client) FindEntityPopulationByName(sourceName string) (*EntityPopulation, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/entity-population", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("name", sourceName)
	req.URL.RawQuery = q.Encode()

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, nil
	}

	population := EntityPopulation{}
	err = json.Unmarshal(body, &population)
	if err != nil {
		return nil, err
	}

	return &population, nil
}

func (c *Client) CreateEntityPopulation(creationRequest EntityPopulation) (*EntityPopulation, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/entity-population", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateEntityPopulation(entityID string, creationRequest EntityPopulation) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/entity-population/%s", c.HostURL, entityID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteEntityPopulation(entityID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/entity-population/%s", c.HostURL, entityID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
