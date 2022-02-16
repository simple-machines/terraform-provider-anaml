package anaml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetEventStore(EventStoreId string) (*EventStore, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/event-store/%s", c.HostURL, EventStoreId), nil)
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

	entity := EventStore{}
	err = json.Unmarshal(body, &entity)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (c *Client) CreateEventStore(creationRequest EventStore) (*EventStore, error) {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/event-store", c.HostURL), strings.NewReader(string(rb)))
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

func (c *Client) UpdateEventStore(EventStoreId string, creationRequest EventStore) error {
	rb, err := json.Marshal(creationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/event-store/%s", c.HostURL, EventStoreId), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteEventStore(EventStoreId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/event-store/%s", c.HostURL, EventStoreId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) FindEventStoreByName(name string) (*EventStore, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/event-store", c.HostURL), nil)
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

	item := EventStore{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
