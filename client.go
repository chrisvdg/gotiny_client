package gotinyclient

import (
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
}

func (c *Client) getReq() (*response, error) {

	return nil, nil
}

type response struct {
	httpResp *http.Response
	body     []byte
}
