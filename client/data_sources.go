package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) FindSource(sourceName string) (*Source, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/source", c.HostURL), nil)
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

	source := Source{}
	err = json.Unmarshal(body, &source)
	if err != nil {
		return nil, err
	}

	return &source, nil
}

func (c *Client) FindDestination(sourceName string) (*Destination, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/destination", c.HostURL), nil)
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

	destination := Destination{}
	err = json.Unmarshal(body, &destination)
	if err != nil {
		return nil, err
	}

	return &destination, nil
}

func (c *Client) FindCluster(sourceName string) (*Cluster, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/cluster", c.HostURL), nil)
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

	cluster := Cluster{}
	err = json.Unmarshal(body, &cluster)
	if err != nil {
		return nil, err
	}

	return &cluster, nil
}
