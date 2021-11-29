package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (c *Client) GetAccessToken(owner int, tokenId string) (*AccessToken, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user/%s/access-token/%s", c.HostURL, strconv.Itoa(owner), tokenId), nil)
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

	token := AccessToken{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (c *Client) CreateAccessToken(owner int, creationRequest AccessToken) (*AccessToken, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user/%s/access-token", c.HostURL, strconv.Itoa(owner)), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var token AccessToken
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (c *Client) DeleteAccessToken(owner int, tokenId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/user/%s/access-token/%s", c.HostURL, strconv.Itoa(owner), tokenId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
