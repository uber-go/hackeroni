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
	"github.com/stretchr/testify/assert"

	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_ResponseLinks(t *testing.T) {
	// Check that the helper page number methods work
	links := ResponseLinks{
		First: "https://api.hackerone.com/v1/reports?page%5Bnumber%5D=1",
		Prev:  "https://api.hackerone.com/v1/reports?page%5Bnumber%5D=2",
		Self:  "https://api.hackerone.com/v1/reports?page%5Bnumber%5D=3",
		Next:  "https://api.hackerone.com/v1/reports?page%5Bnumber%5D=4",
		Last:  "https://api.hackerone.com/v1/reports?page%5Bnumber%5D=5",
	}
	assert.Equal(t, uint64(1), links.FirstPageNumber())
	assert.Equal(t, uint64(2), links.PrevPageNumber())
	assert.Equal(t, uint64(3), links.SelfPageNumber())
	assert.Equal(t, uint64(4), links.NextPageNumber())
	assert.Equal(t, uint64(5), links.LastPageNumber())
}

func Test_addOptions(t *testing.T) {
	// Check that an invalid URL fails
	_, err := addOptions("http://[fe80::1%en0]/", nil)
	assert.NotNil(t, err)

	// Check that invalid parameters fails
	_, err = addOptions("https://hackerone.com/", "Invalid URL Query")
	assert.NotNil(t, err)

	// Check that valid opts marshal correclty
	opts := struct {
		Param string `url:"param"`
	}{
		Param: "value",
	}
	url, err := addOptions("https://hackerone.com/", opts)
	assert.Nil(t, err)
	assert.Equal(t, "https://hackerone.com/?param=value", url)
}

func Test_ErrorResponse(t *testing.T) {
	// Check that .Error results in a correct string
	u, err := url.Parse("https://hackerone.com/admin")
	assert.Nil(t, err)
	errResp := ErrorResponse{
		Response: &http.Response{
			StatusCode: 1337,
			Request: &http.Request{
				Method: "POST",
				URL:    u,
			},
		},
		Errors: []Error{
			Error{
				Status: 400,
				Title:  "Invalid Parameter",
				Detail: "The parameter 'id' is invalid.",
				Source: ErrorSource{
					Parameter: "id",
				},
			},
		},
	}
	assert.Equal(t, "POST https://hackerone.com/admin: 1337 [{Status:400 Title:Invalid Parameter Detail:The parameter 'id' is invalid. Source:{Parameter:id}}]", errResp.Error())
}

func Test_NewClient(t *testing.T) {
	// Check that it uses the DefaultClient
	c := NewClient(nil)
	assert.Equal(t, c.client, http.DefaultClient)
}

func Test_CheckResponse(t *testing.T) {
	// Check that a 200 returns nil
	resp := &http.Response{
		StatusCode: 200,
	}
	err := CheckResponse(resp)
	assert.Nil(t, err)

	// Check that an error response unmarshals correctly
	bodyBytes, err := ioutil.ReadFile("tests/responses/error.json")
	assert.Nil(t, err)
	body := bytes.NewBuffer(bodyBytes)
	resp = &http.Response{
		StatusCode: 400,
		Body:       ioutil.NopCloser(body),
	}
	err = CheckResponse(resp)
	assert.NotNil(t, err)
	expected := &ErrorResponse{
		Response: resp,
		Errors: []Error{
			Error{
				Status: 400,
				Title:  "Invalid Parameter",
				Detail: "The parameter 'id' is invalid.",
				Source: ErrorSource{
					Parameter: "id",
				},
			},
		},
	}
	assert.Equal(t, expected, err)
}

func Test_NewRequest(t *testing.T) {
	// Check that an invalid URL fails
	client := NewClient(nil)
	_, err := client.NewRequest("GET", "http://[fe80::1%en0]/", nil)
	assert.NotNil(t, err)

	// Check that an invalid base URL fails
	badclient := NewClient(nil)
	badclient.BaseURL = &url.URL{
		Scheme: "http://[fe80::1%en0]/",
	}
	_, err = badclient.NewRequest("GET", "/", nil)
	assert.NotNil(t, err)

	// Check that a relative path resolves correctly
	req, err := client.NewRequest("GET", "/relativepath", nil)
	assert.Nil(t, err)
	u, err := url.Parse("https://api.hackerone.com/relativepath")
	assert.Nil(t, err)
	expected := &http.Request{
		Method:     "GET",
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: http.Header{
			"User-Agent": []string{userAgent},
		},
		Host: "api.hackerone.com",
	}
	assert.Equal(t, expected, req)
}

func Test_Client_Do(t *testing.T) {
	// Verify that an invalid request fails
	client := NewClient(nil)
	_, err := client.Do(&http.Request{}, nil)
	assert.NotNil(t, err)

	// Verify that a error response results in an error
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Oh No", 500)
	}))
	defer errorServer.Close()
	u, err := url.Parse(errorServer.URL)
	assert.Nil(t, err)
	_, err = client.Do(&http.Request{
		URL: u,
	}, nil)
	assert.NotNil(t, err)

	// Verify that a invalid response results in an error
	invalidServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Invalid JSON")
	}))
	defer invalidServer.Close()
	u, err = url.Parse(invalidServer.URL)
	assert.Nil(t, err)
	_, err = client.Do(&http.Request{
		URL: u,
	}, nil)
	assert.NotNil(t, err)

	// Verify that a success response results in a struct
	successServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "{}")
	}))
	defer successServer.Close()
	u, err = url.Parse(successServer.URL)
	assert.Nil(t, err)
	successResponse, err := client.Do(&http.Request{
		URL: u,
	}, nil)
	assert.Nil(t, err)
	assert.Equal(t, ResponseLinks{}, successResponse.Links)
}
