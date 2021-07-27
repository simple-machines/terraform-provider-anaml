package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetWebhook(WebhookId string) (*Webhook, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/webhook/%s", c.HostURL, WebhookId), nil)
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

	Webhook := Webhook{}
	err = json.Unmarshal(body, &Webhook)
	if err != nil {
		return nil, err
	}

	return &Webhook, nil
}

func (c *Client) CreateWebhook(creationRequest Webhook) (*Webhook, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/webhook", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateWebhook(WebhookId string, creationRequest Webhook) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/webhook/%s", c.HostURL, WebhookId), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteWebhook(WebhookId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/webhook/%s", c.HostURL, WebhookId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
