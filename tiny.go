package gotinyclient

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

const (
	resourceTiny = "tiny"
)

// ListEntries lists all tiny URL entries
func (c *Client) ListEntries() ([]TinyURL, error) {
	url := c.getReqURL(resourceTiny)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request to list entries")
	}
	c.addAuthHeader(authRead, req)
	res, err := c.do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request listing entries failed")
	}

	err = c.checkDefaultErrors(res, req, nil, false)
	if err != nil {
		return nil, err
	}

	return formatEntryList(res.body)
}

func formatEntryList(data []byte) ([]TinyURL, error) {
	result := []TinyURL{}
	if string(data) != "" && string(data) != "[]" {
		err := json.Unmarshal(data, &result)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal response to entry list")
		}
	}

	return result, nil
}

// CreateEntry creates a new tiny URL entry
// Returns the created TinyUrl entry
// When called without an ID, the created ID can be found there
func (c *Client) CreateEntry(id, url string) (TinyURL, error) {

	return TinyURL{}, nil
}

// UpdateEntry updates a tiny URL entry
func (c *Client) UpdateEntry(id, url string) error {

	return nil
}

// GetEntry fetches the entry information of provided ID
func (c *Client) GetEntry(id string) (TinyURL, error) {

	return TinyURL{}, nil
}

// DeleteEntry deletes an tiny URL entry
func (c *Client) DeleteEntry(id string) error {

	return nil
}
