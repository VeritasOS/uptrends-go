package uptapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// ProbeGroup is a unique group of Uptrends.com probes
type ProbeGroup struct {
	GUID string `json:"Guid,omitempty"`
	Name string `json:"name"`
}

// NewProbeGroup creates a new Uptrends.com probe group
// g includes a name for the new group
func (c *Client) NewProbeGroup(g *ProbeGroup) (string, error) {
	data, err := json.Marshal(g)
	if err != nil {
		return "", err
	}

	req, err := c.newRequest("POST", "probegroups", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 201 {
		return "", errors.New(resp.Status)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if os.Getenv("TF_LOG") != "" {
		log.Println("Response body: ", string(data))
	}

	result := struct {
		GUID string `json:"Guid"`
	}{}
	err = json.Unmarshal(data, &result)

	return result.GUID, nil
}

// ProbeGroup returns an existing Uptrends.com probe group
// id is the Uptrends.com Probe Group GUID
func (c *Client) ProbeGroup(id string) (*ProbeGroup, error) {
	path := fmt.Sprintf("probegroups/%s", id)
	req, err := c.newRequest("GET", path, nil)
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

	// When no content is returned, assume the group doesn't exist
	//  so return empty GUID to set resource ID blank
	if resp.ContentLength == 0 {
		return &ProbeGroup{
			GUID: "",
			Name: "TERRAFORM: GROUP NOT FOUND, NO CONTENT RETURNED",
		}, nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if os.Getenv("TF_LOG") != "" {
		log.Println("Response body: ", string(data))
	}

	result := &ProbeGroup{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// UpdateProbeGroup updates the name of an existing Uptrends.com probe group
// g includes a new name and unique GUID for an existing group
func (c *Client) UpdateProbeGroup(g *ProbeGroup) error {
	data, err := json.Marshal(g)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("probegroups/%s", g.GUID)
	req, err := c.newRequest("PUT", path, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	return nil
}

// DeleteProbeGroup deletes an existing Uptrends.com probe group
// id is the Uptrends.com Probe Group GUID
func (c *Client) DeleteProbeGroup(id string) error {
	path := fmt.Sprintf("probegroups/%s", id)
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 400 {
		// if delete fails with BadRequest, assume deletion occurred previously
		return nil
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	return nil
}
