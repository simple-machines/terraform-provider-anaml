package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetFeatureStore(FeatureStoreID string) (*FeatureStore, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/feature-store/%s", c.HostURL, FeatureStoreID), nil)
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

	FeatureStore := FeatureStore{}
	err = json.Unmarshal(body, &FeatureStore)
	if err != nil {
		return nil, err
	}

	return &FeatureStore, nil
}

func (c *Client) CreateFeatureStore(creationRequest FeatureStore) (*FeatureStore, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/feature-store", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateFeatureStore(FeatureStoreID string, creationRequest FeatureStore) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/feature-store/%s", c.HostURL, FeatureStoreID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteFeatureStore(FeatureStoreID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/feature-store/%s", c.HostURL, FeatureStoreID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
