package anaml

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// HostURL - Default Anaml URL
const HostURL string = "http://localhost:8080"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Auth       *AuthStruct
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
		c.Auth = &AuthStruct{Username: *username, Password: *password}
	} else {
		return nil, errors.New("No username or password set")
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)

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
