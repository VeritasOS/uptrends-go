package uptrends

import (
	"net/http"
	"net/url"
	"testing"
)

type mockRespBody struct {
	text string
}

func (m *mockRespBody) Close() (err error) { return err }
func (m *mockRespBody) Read(p []byte) (n int, err error) {
	max := cap(p)
	for k, v := range []byte(m.text) {
		if k > max {
			break
		}
		p[k] = byte(v)
		n = k
	}

	return n, err
}

type mockClientor struct {
	body string
}

func (m *mockClientor) Do(req *http.Request) (resp *http.Response, err error) {
	resp.Body = &mockRespBody{m.body}

	return resp, err
}

func TestCheckpointServerIPs(t *testing.T) {
	type testData struct {
		responseBody string
		ipList       []string
	}

	tester := func(l []testData) {
		for _, v := range l {
			c := &Client{
				url.URL{},
				&mockClientor{body: v.responseBody},
				"u",
				"p",
				&mockReqCreator{},
			}

			list, err := c.CheckpointServerIPs()
			if err != nil {
				t.Errorf("CheckpointServerIPs() errored %s:", err)
			}

			elen := len(v.ipList)
			olen := len(list)
			if olen != elen {
				t.Errorf("CheckpointServerIPs() didn't return expected number of IPs. Expected %d, Observed %d:", elen, olen)
			}
		}
	}

	td := []testData{testData{"foo", []string{"bar"}}}
	tester(td)
}
