package gotinyclient

import (
	"fmt"
	"io/ioutil"
	"net/http"

	gotiny "github.com/chrisvdg/gotiny/backend"
	"github.com/pkg/errors"
)

var (
	// ErrEntryNotFound represents an error where a tiny url entry could not be found
	ErrEntryNotFound = errors.New("tiny URL entry not found")
	// ErrReadUnauthorized represents an error where a request with read permissions failed to authorized
	ErrReadUnauthorized = errors.New("unauthorized, read access token may be invalid")
	// ErrWriteUnauthorized represents an error where a request with write permissions failed to authorized
	ErrWriteUnauthorized = errors.New("unauthorized, write access token may be invalid")
	// ErrInvalidRequest represents an error where an invalid request was made
	ErrInvalidRequest = errors.New("invalid request")
)

const (
	authHeader   = "Authorization"
	bearerPrefix = "bearer"
)

var (
	successCodes = []int{
		http.StatusOK,
		http.StatusCreated,
	}
)

// TinyURL represents a tiny url entry
type TinyURL gotiny.TinyURL

// JSONTime represents a json parsable time struct
type JSONTime gotiny.JSONTime

// New returns a new gotiny client
func New(baseURL, readToken, writeToken string) (*Client, error) {
	return &Client{
		baseURL:    baseURL,
		readToken:  readToken,
		writeToken: writeToken,
	}, nil
}

// Client represents a gotiny client
type Client struct {
	baseURL    string
	readToken  string
	writeToken string
	http       *http.Client
}

func (c *Client) do(req *http.Request) (*response, error) {
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request to gotiny server")
	}
	defer resp.Body.Close()

	r := &response{}
	r.body, err = ioutil.ReadAll(resp.Request.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}
	r.httpResp = resp

	return r, nil
}

func (c *Client) addAuthHeader(auth authType, req *http.Request) {
	switch auth {
	case authRead:
		if c.readToken != "" {
			req.Header.Add(authHeader, fmt.Sprintf("%s %s", bearerPrefix, c.readToken))
		}
	case authWrite:
		if c.writeToken != "" {
			req.Header.Add(authHeader, fmt.Sprintf("%s %s", bearerPrefix, c.writeToken))
		}
	}
}

type authType int

const (
	authRead authType = iota + 1
	authWrite
)

type response struct {
	httpResp *http.Response
	body     []byte
}
