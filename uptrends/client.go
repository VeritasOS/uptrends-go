package uptrends

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type clientor interface {
	Do(*http.Request) (*http.Response, error)
}

type requestCreator interface {
	NewRequest(method, path string, body io.Reader) (*http.Request, error)
}

type httpReqCreator struct{}

func (h *httpReqCreator) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, path, body)
}

// Client is an Uptrends.com API Client
type Client struct {
	baseURL url.URL
	clientor
	username string
	password string
	requestCreator
}

// New returns an Uptrends.com API Client
//  username is the Uptrends.com login user
//  password is the corresponding password for username
func New(username string, password string) (*Client, error) {
	u, err := url.Parse("https://api.uptrends.com")
	if err != nil {
		return nil, err
	}

	return &Client{
		*u,
		&http.Client{
			Timeout: time.Second * 10,
		},
		username,
		password,
		&httpReqCreator{},
	}, nil
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	url := c.baseURL
	url.Path = strings.Join([]string{"", "v3", path}, "/")
	req, err := c.NewRequest(method, url.String(), body)
	if err != nil {
		return req, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json, text/json, */*; q=0.1")

	req.SetBasicAuth(c.username, c.password)

	return req, err
}
