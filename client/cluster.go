package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetCluster(clusterID string) (*Cluster, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/cluster/%s", c.HostURL, clusterID), nil)
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

	cluster := Cluster{}
	err = json.Unmarshal(body, &cluster)
	if err != nil {
		return nil, err
	}

	return &cluster, nil
}

func (c *Client) CreateCluster(creationRequest Cluster) (*Cluster, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cluster", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateCluster(clusterID string, creationRequest Cluster) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/cluster/%s", c.HostURL, clusterID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteCluster(clusterID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/cluster/%s", c.HostURL, clusterID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
