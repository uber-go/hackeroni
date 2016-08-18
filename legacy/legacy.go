// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package legacy

// Imports
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	libraryVersion = "1"
	defaultBaseURL = "https://hackerone.com/"
	userAgent      = "uber/h1/v" + libraryVersion
)

// A Client manages communication with the H1.
type Client struct {
	client *http.Client // HTTP client used to communicate.

	// Base URL for requests. Defaults to the public H1. BaseURL should always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with H1.
	UserAgent string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of H1.
	Session *SessionService
	User    *UserService
	Team    *TeamService
	Report  *ReportService
}

type service struct {
	client *Client
}

// NewClient returns a new client. If a nil httpClient is provided, http.DefaultClient will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.common.client = c
	c.Session = (*SessionService)(&c.common)
	c.User = (*UserService)(&c.common)
	c.Team = (*TeamService)(&c.common)
	c.Report = (*ReportService)(&c.common)

	return c
}

// NewRequest creates an request. A relative URL can be provided in urlStr
func (c *Client) NewRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, c.BaseURL.ResolveReference(rel).String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", c.UserAgent)

	return req, nil
}

// Response is a H1 response.  This wraps the standard http.Response and provides convenience fields for pagination
type Response struct {
	*http.Response
}

// ErrorResponse wraps a http.Response and is returned when the API returns an error.
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
}

// ErrorResponse needs to implement Error to be a valid error type.
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode)
}

// CheckResponse determines if the given http.Response was an error and converts it to a h1.ErrorResponse if so
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	return errorResponse
}

// Do sends an request and returns the response.
func (c *Client) Do(req *http.Request, resource interface{}) (*Response, error) {
	// Actually do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Make a response object
	response := &Response{Response: resp}

	// If request returned an error, return the response and err back to user to inspect
	if err := CheckResponse(resp); err != nil {
		return response, err
	}

	// Decode if a resource was provided
	if resource != nil {
		if err := json.NewDecoder(resp.Body).Decode(resource); err != nil {
			return response, err
		}
	}

	// Return success
	return response, nil
}
