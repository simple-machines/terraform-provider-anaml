package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetUserGroup(userGroupID string) (*UserGroup, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user-group/%s", c.HostURL, userGroupId), nil)
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

	userGroup := UserGroup{}
	err = json.Unmarshal(body, &userGroup)
	if err != nil {
		return nil, err
	}

	return &userGroup, nil
}

func (c *Client) CreateUserGroup(creationRequest UserGroup) (*UserGroup, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user-group", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateUserGroup(userGroupID string, creationRequest UserGroup) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/user-group/%s", c.HostURL, userGroupId), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteUserGroup(userGroupID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/user-group/%s", c.HostURL, userGroupId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
