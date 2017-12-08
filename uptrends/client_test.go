package uptrends

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
)

var (
	u = "random@username.com"
	p = "Cryp7!cP@$$w0rd"
)

type mockReqFailingCreator struct{}

func (m *mockReqFailingCreator) NewRequest(method, path string, body io.Reader) (req *http.Request, err error) {
	return req, errors.New("failed to create a request")
}

type mockReqCreator struct{}

func (m *mockReqCreator) NewRequest(method, path string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, path, body)
	return req, err
}

func TestClientNew(t *testing.T) {
	c, err := New(u, p)
	if err != nil {
		t.Error(err)
	}
	if c.username != u {
		t.Errorf("username %s, expected %s", c.username, u)
	}
	if c.password != p {
		t.Errorf("password %s, expected %s", c.password, p)
	}
	n := url.URL{}
	if c.baseURL == n {
		t.Errorf("baseURL nil, expected net/url.URL")
	}
}

func TestClientNewRequest(t *testing.T) {

	c := &Client{
		url.URL{},
		&http.Client{},
		u,
		p,
		&mockReqFailingCreator{},
	}

	r, err := c.newRequest("GET", "foo", nil)
	if err == nil {
		t.Fatalf("test request construction didn't error as expected: %s", err)
	}

	c = &Client{
		url.URL{},
		&http.Client{},
		u,
		p,
		&mockReqCreator{},
	}

	r, err = c.newRequest("GET", "foo", nil)
	if err != nil {
		t.Fatalf("test request construction errored: %s", err)
	}
	if r == nil {
		t.Fatalf("test request construction failed: %s", err)
	}
	if r.URL.Path != "/v3/foo" {
		t.Fatalf("test request wrong path. expected 'v3/foo', constructed '%s'", r.URL.Path)
	}

	c, err = New(u, p)
	if err != nil {
		t.Fatalf("test client creation errored: %s", err)
	}

	r, err = c.newRequest("GET", "foo", nil)
	if err != nil {
		t.Fatalf("test request construction errored: %s", err)
	}
	if r == nil {
		t.Fatalf("test request construction failed: %s", err)
	}
	if r.URL.Path != "/v3/foo" {
		t.Fatalf("test request wrong path. expected 'v3/foo', constructed '%s'", r.URL.Path)
	}
}
