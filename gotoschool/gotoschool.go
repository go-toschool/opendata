package gotoschool

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var baseURL = url.URL{
	Scheme: "http",
	Host:   "api.gotoschool.com",
	Path:   "",
}

type clientErr struct {
	msg string
	err error
}

func newError(msg string, err error) *clientErr {
	return &clientErr{msg, err}
}

func (ce *clientErr) Error() string {
	res := "gotoschool: " + ce.msg
	return res
}

// Client represents a gotoschool client, this client
// let us post information to gotoschool service.
type Client struct {
	client    http.Client
	apiURL    string
	token     string
	userToken string
}

// Config represent basic information to setup a new
// connection to gotoschool
type Config struct {
	Token string
}

// NewClient returns a new gotoschool.Client that allow to connect
// to gotoschool service.
func NewClient(c *Config) *Client {
	return &Client{
		client: http.Client{},
		apiURL: baseURL.String(),
		token:  c.Token,
	}
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	req.Header.Set("Accept", "application/graphql")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("X-gotoschool-User-Token", fmt.Sprintf("Bearer %s", c.userToken))

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
	req.Header.Set("Content-Type", "application/graphql")

	payloadBytes, err := c.do(req)
	if err != nil {
		return nil, err
	}

	return payloadBytes, nil
}

// Post make a POST request.
func (c *Client) Post(path, body string) ([]byte, error) {
	b := bytes.NewBufferString(body)

	req, err := http.NewRequest("POST", c.url(path), b)
	if err != nil {
		return nil, newError("failed to create request", err)
	}
	req.Header.Set("Content-Type", "application/graphql")

	payloadBytes, err := c.do(req)
	if err != nil {
		return nil, err
	}

	return payloadBytes, nil
}

// SetUserToken ...
func (c *Client) SetUserToken(token string) {
	c.userToken = token
}

func (c *Client) url(path string) string { return c.apiURL + path }
