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

package h1

// Imports
import (
	"github.com/google/go-querystring/query"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	libraryVersion = "1"
	defaultBaseURL = "https://api.hackerone.com/v1/"
	userAgent      = "hackeroni/v" + libraryVersion
)

// A Client manages communication with the H1 API.
type Client struct {
	client *http.Client // HTTP client used to communicate with the API.

	// Base URL for API requests. Defaults to the public H1 API. BaseURL should always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the H1 API.
	UserAgent string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the H1 API.
	Report  *ReportService
	Program *ProgramService
}

type service struct {
	client *Client
}

// ListOptions specifies the optional parameters to various List methods that support pagination.
type ListOptions struct {
	// For paginated results which page to retrieve.
	Page uint64 `url:"page[number],omitempty"`

	// For paginated results the size of pages to retrieve
	PageSize uint64 `url:"page[size],omitempty"`

	// For lists the index to sort by
	Sort string `url:"sort,omitempty"`
}

// addOptions adds the parameters in opt as URL query parameters to s
func addOptions(s string, opt interface{}, listOpts *ListOptions) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	lqs, _ := query.Values(listOpts)
	if lqs != nil {
		for k, v := range lqs {
			qs[k] = v
		}
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// NewClient returns a new H1 API client. If a nil httpClient is provided, http.DefaultClient will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.common.client = c
	c.Report = (*ReportService)(&c.common)
	c.Program = (*ProgramService)(&c.common)

	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, c.BaseURL.ResolveReference(rel).String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)

	return req, nil
}

// ResponseLinks represents a JSONAPI ResponseLinks object
type ResponseLinks struct {
	First string `json:"first"`
	Prev  string `json:"prev"`
	Self  string `json:"self"`
	Next  string `json:"next"`
	Last  string `json:"last"`
}

// FirstPageNumber extracts the page number from the ResponseLinks structure
func (l *ResponseLinks) FirstPageNumber() uint64 {
	parsedURL, _ := url.Parse(l.First)
	page, _ := strconv.ParseUint(parsedURL.Query().Get("page[number]"), 10, 64)
	return page
}

// PrevPageNumber extracts the page number from the ResponseLinks structure
func (l *ResponseLinks) PrevPageNumber() uint64 {
	parsedURL, _ := url.Parse(l.Prev)
	page, _ := strconv.ParseUint(parsedURL.Query().Get("page[number]"), 10, 64)
	return page
}

// SelfPageNumber extracts the page number from the ResponseLinks structure
func (l *ResponseLinks) SelfPageNumber() uint64 {
	parsedURL, _ := url.Parse(l.Self)
	page, _ := strconv.ParseUint(parsedURL.Query().Get("page[number]"), 10, 64)
	return page
}

// NextPageNumber extracts the page number from the ResponseLinks structure
func (l *ResponseLinks) NextPageNumber() uint64 {
	parsedURL, _ := url.Parse(l.Next)
	page, _ := strconv.ParseUint(parsedURL.Query().Get("page[number]"), 10, 64)
	return page
}

// LastPageNumber extracts the page number from the ResponseLinks structure
func (l *ResponseLinks) LastPageNumber() uint64 {
	parsedURL, _ := url.Parse(l.Last)
	page, _ := strconv.ParseUint(parsedURL.Query().Get("page[number]"), 10, 64)
	return page
}

// Response is a H1 API response.  This wraps the standard http.Response and provides convenience fields for pagination
type Response struct {
	*http.Response

	// Links relating to the response
	Links ResponseLinks `json:"links"`
}

// ErrorSource represents an ErrorSource from the JSONAPI specification.
type ErrorSource struct {
	Parameter string `json:"parameter"`
}

// Error represents an Error from the JSONAPI specification.
type Error struct {
	Status uint        `json:"status"`
	Title  string      `json:"title"`
	Detail string      `json:"detail"`
	Source ErrorSource `json:"source"`
}

// ErrorResponse wraps a http.Response and is returned when the API returns an error.
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Errors   []Error        `json:"errors"` // The individual errors that occured
}

// ErrorResponse needs to implement Error to be a valid error type.
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %+v", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Errors)
}

// Used to parse responses
type responseWrapper struct {
	*Response
	Data interface{} `json:"data"`
}

// CheckResponse determines if the given http.Response was an error and converts it to a h1.ErrorResponse if so
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		// Ignore errors here so we always pass out an ErrorResponse
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

// Do sends an API request and returns the API response.
func (c *Client) Do(req *http.Request, resource interface{}) (*Response, error) {
	// Actually do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Make a response object
	response := &Response{Response: resp}

	// If API returned an error, return the response and err back to user to inspect
	if err := CheckResponse(resp); err != nil {
		return response, err
	}

	// Wrap the response object so we can get data as well
	wrapper := &responseWrapper{
		Response: response,
		Data:     resource,
	}
	if err := json.NewDecoder(resp.Body).Decode(wrapper); err != nil {
		return response, err
	}

	// Return success
	return response, nil
}
