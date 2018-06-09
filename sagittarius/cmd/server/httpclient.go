package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	dialTimeout = flag.Duration("dial_timeout", 5*time.Second, "Timeout for dialing an HTTP connection.")
	reqTimeout  = flag.Duration("req_timeout", 20*time.Second, "Timeout for roundtripping an HTTP request.")
)

type httpRequest http.Request

func newHTTPRequest(method, url, token string) (*httpRequest, error) {
	r, err := http.NewRequest(method, url, nil)
	if token != "" {
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	return (*httpRequest)(r), err
}

func (r *httpRequest) WithBody(body []byte) *httpRequest {
	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	r.ContentLength = int64(len(body))
	return r
}

type httpClient struct {
	client *http.Client
	logger *logrus.Logger
}

// http.ProxyFromEnvironment
func newHTTPClient(proxy string, l *logrus.Logger) *httpClient {
	dialer := &net.Dialer{
		Timeout:   *dialTimeout,
		KeepAlive: *reqTimeout / 2,
	}
	t := http.Transport{
		Dial: dialer.Dial,
		ResponseHeaderTimeout: *reqTimeout / 2,
		TLSHandshakeTimeout:   *reqTimeout / 2,
	}

	if u, err := url.Parse(proxy); err == nil {
		t.Proxy = http.ProxyURL(u)
	}

	c := &http.Client{
		Transport: &t,
		Timeout:   *reqTimeout,
	}

	return &httpClient{
		client: c,
		logger: l,
	}
}

// Post returns true if the there is no error and the response HTTP status is
// 200 or 201.
func (c *httpClient) Post(url, token string, body []byte) bool {
	req, err := newHTTPRequest("POST", url, token)
	if err != nil {
		c.logger.WithFields(logrus.Fields{
			"message": "httpclient: error building request",
			"error":   err.Error(),
		}).Fatal("Critical Error")

		return false
	}

	resp, err := c.do(req.WithBody(body))
	if err != nil {
		c.logger.WithFields(logrus.Fields{
			"message": "httpclient: a transport error occurs",
			"error":   err.Error(),
		}).Warning("Warning")
		return false
	}

	// if resp doesn't have a HTTP status 200 or 201 is considered
	// an error and will be retried.
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return false
	}

	return true
}

func (c *httpClient) do(r *httpRequest) (*http.Response, error) {
	return c.client.Do((*http.Request)(r))
}
