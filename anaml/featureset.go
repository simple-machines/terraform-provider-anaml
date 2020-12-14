package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetFeatureSet(FeatureSetID string) (*FeatureSet, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/feature-set/%s", c.HostURL, FeatureSetID), nil)
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

	FeatureSet := FeatureSet{}
	err = json.Unmarshal(body, &FeatureSet)
	if err != nil {
		return nil, err
	}

	return &FeatureSet, nil
}

func (c *Client) CreateFeatureSet(creationRequest FeatureSet) (*FeatureSet, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/feature-set", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateFeatureSet(FeatureSetID string, creationRequest FeatureSet) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/feature-set/%s", c.HostURL, FeatureSetID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteFeatureSet(FeatureSetID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/feature-set/%s", c.HostURL, FeatureSetID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
