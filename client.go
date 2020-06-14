package gotinyclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	gotiny "github.com/chrisvdg/gotiny/backend"
	"github.com/pkg/errors"
)

var (
	// ErrEntryNotFound represents an error where a tiny url entry could not be found
	ErrEntryNotFound = errors.New("tiny URL entry not found")
	// ErrUnauthorized represents an error where a request failed to authorize
	ErrUnauthorized = errors.New("request unauthorized")
	// ErrReadUnauthorized represents an error where a request with read permissions failed to authorize
	ErrReadUnauthorized = errors.New("unauthorized, read access token may be invalid")
	// ErrWriteUnauthorized represents an error where a request with write permissions failed to authorize
	ErrWriteUnauthorized = errors.New("unauthorized, write access token may be invalid")
	// ErrBadRequest represents an error where a request was made with bad parameters
	ErrBadRequest = errors.New("bad request")
	// UnauthorizedErrors is a list of errors that represent unauthorzied errors
	UnauthorizedErrors = []error{ErrUnauthorized, ErrReadUnauthorized, ErrWriteUnauthorized}
	// ErrUnexpectedCode represents an error where a request returned an unexpected error
	ErrUnexpectedCode = errors.New("recieved unexpected http code")
)

const (
	authHeader   = "Authorization"
	bearerPrefix = "bearer"
	apiPath      = "api"
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
		http: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}, nil
}

// Client represents a gotiny client
type Client struct {
	baseURL    string
	readToken  string
	writeToken string
	http       *http.Client
}

// getReqURL returns the full API call address using the base URL
// and appending the provided path
func (c *Client) getReqURL(resourcePath string) string {
	return fmt.Sprintf("%s/%s", c.baseURL, path.Join(apiPath, resourcePath))
}

// do executed the provided http request and returns the response with read body
func (c *Client) do(req *http.Request) (*response, error) {
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request to gotiny server")
	}
	defer resp.Body.Close()

	r := &response{}
	r.body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}
	r.httpResp = resp

	return r, nil
}

// checkDefaultErrors default response error check
func (c *Client) checkDefaultErrors(res *response, req *http.Request, expectedCodes []int, isResource bool) error {
	switch res.httpResp.StatusCode {
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusNotFound:
		if isResource {
			return ErrEntryNotFound
		}
		return errors.Errorf("Path %s not found", req.URL.String())
	case http.StatusBadRequest:
		return errors.Wrap(ErrBadRequest, string(res.body))
	default:
		if len(expectedCodes) == 0 {
			expectedCodes = successCodes
		}
		if !isExpectedCode(expectedCodes, res.httpResp.StatusCode) {
			return errors.Wrap(ErrUnexpectedCode, string(res.body))
		}
	}

	return nil
}

// addAuthHeader adds authorization header if configured
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

// IsUnauthorized checks wether provided error is an authorized error
func IsUnauthorized(err error) bool {
	for _, r := range UnauthorizedErrors {
		if errors.Cause(err) == r {
			return true
		}
	}

	return false
}

func isExpectedCode(expCodes []int, code int) bool {
	for _, c := range expCodes {
		if c == code {
			return true
		}
	}

	return false
}
