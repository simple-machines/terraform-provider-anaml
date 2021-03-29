package anaml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// HostURL - Default Anaml URL
const HostURL string = "http://localhost:8080"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	Branch     *string
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	Token string `json:"token"`
}

// NewClient -
func NewClient(host, username, password, branch *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	if branch != nil && *branch != "" {
		c.Branch = branch
	}

	if (username != nil) && (password != nil) {
		// form request body
		rb, err := json.Marshal(AuthStruct{
			Username: *username,
			Password: *password,
		})
		if err != nil {
			return nil, err
		}

		// authenticate
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/login", c.HostURL), strings.NewReader(string(rb)))
		if err != nil {
			return nil, err
		}

		body, err := c.doRequest(req)
		if err != nil {
			return nil, err
		}

		// parse response body
		ar := AuthResponse{}
		err = json.Unmarshal(body, &ar)
		if err != nil {
			return nil, err
		}

		c.Token = ar.Token
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	if c.Branch != nil {
		q := req.URL.Query()
		q.Add("branch", *c.Branch)
		req.URL.RawQuery = q.Encode()
	}

	log.Printf("[DEBUG] Request: %v\n", req)

	if req.Body != nil {
		requestBody, err := ioutil.ReadAll(req.Body)
		if err == nil {
			reader0 := ioutil.NopCloser(bytes.NewBuffer(requestBody))
			reader1 := ioutil.NopCloser(bytes.NewBuffer(requestBody))
			log.Printf("[DEBUG] Request body: %q", reader0)
			req.Body = reader1
		}
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Response: %v\n", res)

	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	reader := ioutil.NopCloser(bytes.NewBuffer(responseBody))
	log.Printf("[DEBUG] Request body: %q", reader)

	if res.StatusCode == 404 {
		return nil, nil
	}

	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, responseBody)
	}

	return responseBody, err
}
