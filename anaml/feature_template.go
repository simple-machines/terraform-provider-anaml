package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetFeatureTemplate(featureID string) (*FeatureTemplate, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/feature-template/%s", c.HostURL, featureID), nil)
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

	feature := FeatureTemplate{}
	err = json.Unmarshal(body, &feature)
	if err != nil {
		return nil, err
	}

	return &feature, nil
}

func (c *Client) CreateFeatureTemplate(creationRequest FeatureTemplate) (*FeatureTemplate, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/feature-template", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateFeatureTemplate(templateID string, creationRequest FeatureTemplate) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/feature-template/%s", c.HostURL, templateID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteFeatureTemplate(templateID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/feature-template/%s", c.HostURL, templateID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
