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

	feature := Source{}
	err = json.Unmarshal(body, &feature)
	if err != nil {
		return nil, err
	}

	return &feature, nil
}
