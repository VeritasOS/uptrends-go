package uptapi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

// CheckpointServer is an Uptrends.com checkpoint server
type CheckpointServer struct {
	ServerID       int64  `json:"ServerID"`
	CheckpointID   int64  `json:"CheckPointID"`
	CheckpointName string `json:"CheckPointName"`
	IPAddress      string `json:"IPAddress"`
	IPv6Address    string `json:"IPv6Address,omitempty"`
}

// CheckpointServerIPs returns a list of Uptrends.com server IPs
func (c *Client) CheckpointServerIPs() ([]string, error) {
	req, err := c.newRequest("GET", "checkpointservers", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if os.Getenv("TF_LOG") != "" {
		log.Println("Response body: ", string(data))
	}

	s := []CheckpointServer{}
	err = json.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}

	l := make([]string, 0)
	for _, x := range s {
		l = append(l, x.IPAddress)
	}

	return l, err
}
