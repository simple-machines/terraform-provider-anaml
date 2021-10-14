package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetBranchProtection(branchProtectionId string) (*BranchProtection, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/branch-protection/%s", c.HostURL, branchProtectionId), nil)
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

	branchProtection := BranchProtection{}
	err = json.Unmarshal(body, &branchProtection)
	if err != nil {
		return nil, err
	}

	return &branchProtection, nil
}

func (c *Client) CreateBranchProtection(creationRequest BranchProtection) (*BranchProtection, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/branch-protection", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateBranchProtection(branchProtectionId string, creationRequest BranchProtection) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/branch-protection/%s", c.HostURL, branchProtectionId), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteBranchProtection(branchProtectionId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/branch-protection/%s", c.HostURL, branchProtectionId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
