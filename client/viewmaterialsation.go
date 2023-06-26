package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetViewMaterialisationJob(id string) (*ViewMaterialisationJob, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/view-materialisation/%s", c.HostURL, id), nil)
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

	MaterialisedView := ViewMaterialisationJob{}
	err = json.Unmarshal(body, &MaterialisedView)
	if err != nil {
		return nil, err
	}

	return &MaterialisedView, nil
}

func (c *Client) CreateViewMaterialisationJob(creationRequest ViewMaterialisationJob) (*ViewMaterialisationJob, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/view-materialisation", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateViewMaterialisationJob(id string, creationRequest ViewMaterialisationJob) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/view-materialisation/%s", c.HostURL, id), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteViewMaterialisationJob(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/view-materialisation/%s", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) FindViewMaterialisationJobByName(name string) (*ViewMaterialisationJob, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/view-materialisation", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, nil
	}

	item := ViewMaterialisationJob{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
