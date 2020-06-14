package gotinyclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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

// CreateEntry creates a new tiny URL entry
// Returns the created TinyUrl entry
// When called without an ID, the created ID can be found there
func (c *Client) CreateEntry(entryID, entryURL string) (TinyURL, error) {
	emptyResult := TinyURL{}
	reqURL := c.getReqURL(resourceTiny)
	form := url.Values{}
	form.Add("id", entryID)
	form.Add("url", entryURL)
	req, err := http.NewRequest("POST", reqURL, strings.NewReader(form.Encode()))
	if err != nil {
		return emptyResult, errors.Wrap(err, "failed to create request to create new entry")
	}

	c.addAuthHeader(authWrite, req)
	res, err := c.do(req)
	if err != nil {
		return emptyResult, errors.Wrap(err, "request create new entry failed")
	}

	err = c.checkDefaultErrors(res, req, nil, false)
	if err != nil {
		return emptyResult, err
	}

	return formatEntry(res.body)
}

// UpdateEntry updates a tiny URL entry
func (c *Client) UpdateEntry(entryID, entryURL string) error {
	reqURL := c.getReqURLWithID(resourceTiny, entryID)
	form := url.Values{}
	form.Add("url", entryURL)
	req, err := http.NewRequest("POST", reqURL, strings.NewReader(form.Encode()))
	if err != nil {
		return errors.Wrap(err, "failed to create request to create new entry")
	}

	c.addAuthHeader(authWrite, req)
	res, err := c.do(req)
	if err != nil {
		return errors.Wrap(err, "request create new entry failed")
	}

	return c.checkDefaultErrors(res, req, nil, true)
}

// GetEntry fetches the entry information of provided ID
func (c *Client) GetEntry(entryID string) (TinyURL, error) {
	emptyResult := TinyURL{}
	url := c.getReqURLWithID(resourceTiny, entryID)
	url = fmt.Sprintf("%s/expand", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return emptyResult, errors.Wrap(err, "failed to create request to fetch entry")
	}
	c.addAuthHeader(authRead, req)
	res, err := c.do(req)
	if err != nil {
		return emptyResult, errors.Wrap(err, "request getting entry failed")
	}

	err = c.checkDefaultErrors(res, req, nil, true)
	if err != nil {
		return emptyResult, err
	}

	return formatEntry(res.body)
}

// DeleteEntry deletes an tiny URL entry
func (c *Client) DeleteEntry(entryID string) error {
	url := c.getReqURLWithID(resourceTiny, entryID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request to fetch entry")
	}
	c.addAuthHeader(authWrite, req)
	res, err := c.do(req)
	if err != nil {
		return errors.Wrap(err, "request getting entry failed")
	}

	return c.checkDefaultErrors(res, req, nil, true)
}

func formatEntryList(data []byte) ([]TinyURL, error) {
	result := []TinyURL{}
	err := unmarshal(data, &result)
	return result, err
}

func formatEntry(data []byte) (TinyURL, error) {
	result := TinyURL{}
	err := unmarshal(data, &result)
	return result, err
}

func unmarshal(data []byte, object interface{}) error {
	if string(data) != "" && string(data) != "[]" && string(data) != "{}" {
		err := json.Unmarshal(data, &object)
		return errors.Wrap(err, "failed to unmarshal response")
	}
	return nil
}
