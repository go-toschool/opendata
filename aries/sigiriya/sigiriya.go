package sigiriya

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var baseURL = url.URL{
	Scheme: "https",
	Host:   "api.finciero.com",
	Path:   "services",
}

type clientErr struct {
	msg string
	err error
}

func newError(msg string, err error) *clientErr {
	return &clientErr{msg, err}
}

func (ce *clientErr) Error() string {
	res := "sigiriya: " + ce.msg
	return res
}

// Client represents a sigiriya client, this client
// let us post information to sigiriya service.
type Client struct {
	client http.Client
	apiURL string
	token  string
}

// Config represent basic information to setup a new
// connection to sigiriya
type Config struct {
	Token string
}

// NewClient returns a new sigiriya.Client that allow to connect
// to sigiriya service.
func NewClient(c *Config) *Client {
	return &Client{
		client: http.Client{},
		apiURL: baseURL.String(),
		token:  c.Token,
	}
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, newError("request failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, newError("record not found", errors.New("record not found"))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, newError("failed to read response body", err)
	}

	return body, nil
}

// Get make a GET request.
func (c *Client) Get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.url(path), nil)
	if err != nil {
		return nil, newError("failed to create request", err)
	}
	req.Header.Set("Content-Type", "application/json")

	payloadBytes, err := c.do(req)
	if err != nil {
		return nil, err
	}

	return payloadBytes, nil
}

func (c *Client) url(path string) string { return c.apiURL + path }
