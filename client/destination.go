package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetDestination(destinationID string) (*Destination, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/destination/%s", c.HostURL, destinationID), nil)
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

	destination := Destination{}
	err = json.Unmarshal(body, &destination)
	if err != nil {
		return nil, err
	}

	return &destination, nil
}

func (c *Client) CreateDestination(creationRequest Destination) (*Destination, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/destination", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateDestination(destinationID string, creationRequest Destination) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/destination/%s", c.HostURL, destinationID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteDestination(destinationID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/destination/%s", c.HostURL, destinationID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
