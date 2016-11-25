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

import (
	"github.com/stretchr/testify/assert"

	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var expectedReport = Report{
	ID:    String("1337"),
	Type:  String(ReportType),
	Title: String("XSS in login form"),
	VulnerabilityInformation: String("..."),
	State:     String("new"),
	CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	Program: &Program{
		ID:        String("1337"),
		Type:      String(ProgramType),
		Handle:    String("security"),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	},
	VulnerabilityTypes: []VulnerabilityType{
		VulnerabilityType{
			ID:          String("1337"),
			Type:        String(VulnerabilityTypeType),
			Name:        String("Cross-Site Scripting (XSS)"),
			Description: String("Failure of a site to validate, filter, or encode user input before returning it to another user's web client."),
			CreatedAt:   NewTimestamp("2016-02-02T04:05:06.000Z"),
		},
	},
	Reporter: &User{
		ID:       String("1337"),
		Type:     String(UserType),
		Disabled: Bool(false),
		Username: String("api-example"),
		Name:     String("API Example"),
		ProfilePicture: UserProfilePicture{
			Size62x62:   String("/assets/avatars/default.png"),
			Size82x82:   String("/assets/avatars/default.png"),
			Size110x110: String("/assets/avatars/default.png"),
			Size260x260: String("/assets/avatars/default.png"),
		},
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	},
	Attachments: []Attachment{},
	Swag:        []Swag{},
	Activities:  []Activity{},
	Bounties:    []Bounty{},
	Summaries:   []ReportSummary{},
}

func Test_ReportService_Get(t *testing.T) {
	// Verify that an invalid url fails
	c := NewClient(nil)
	c.BaseURL = &url.URL{}
	_, _, err := c.Report.Get("%A")
	assert.NotNil(t, err)

	// Verify that an error response fails
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Oh No", 500)
	}))
	defer errorServer.Close()
	u, err := url.Parse(errorServer.URL)
	assert.Nil(t, err)
	c.BaseURL = u
	_, _, err = c.Report.Get("123456")
	assert.NotNil(t, err)

	// Verify that it gets a response correctly
	reportServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "tests/responses/report.json")
	}))
	defer reportServer.Close()
	u, err = url.Parse(reportServer.URL)
	assert.Nil(t, err)
	c.BaseURL = u
	actual, _, err := c.Report.Get("123456")
	assert.Nil(t, err)
	assert.Equal(t, &expectedReport, actual)
}

var expectedCommentRequest = `{"data":{"type":"activity-comment","attributes":{"message":"A fix has been deployed. Can you retest, please?","internal":false}}}`
var expectedComment = Activity{
	ID:        String("1337"),
	Type:      String(ActivityCommentType),
	Message:   String("A fix has been deployed. Can you retest, please?"),
	Internal:  Bool(false),
	CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
}

func Test_ReportService_CreateComment(t *testing.T) {
	// Verify that an invalid url fails
	c := NewClient(nil)
	c.BaseURL = &url.URL{}
	_, _, err := c.Report.CreateComment("%A", "A fix has been deployed. Can you retest, please?", false)
	assert.NotNil(t, err)

	// Verify that an error response fails
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Oh No", 500)
	}))
	defer errorServer.Close()
	u, err := url.Parse(errorServer.URL)
	assert.Nil(t, err)
	c.BaseURL = u
	_, _, err = c.Report.CreateComment("123456", "A fix has been deployed. Can you retest, please?", false)
	assert.NotNil(t, err)

	// Verify that it gets a response correctly and it has the correct request body
	commentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(string(body)) != strings.TrimSpace(expectedCommentRequest) {
			http.Error(w, "Non-matching request!", http.StatusBadRequest)
			return
		}
		http.ServeFile(w, r, "tests/responses/report_create-comment.json")
	}))
	defer commentServer.Close()
	u, err = url.Parse(commentServer.URL)
	assert.Nil(t, err)
	c.BaseURL = u
	actual, _, err := c.Report.CreateComment("123456", "A fix has been deployed. Can you retest, please?", false)
	assert.Nil(t, err)
	actual.RawActor = nil
	actual.rawData = nil
	assert.Equal(t, &expectedComment, actual)
}

var expectedUpdateAssigneeRequest1 = `{"data":{"type":"nobody","attributes":{"message":"@member Please check this out!"}}}`
var expectedUpdateAssigneeRequest2 = `{"data":{"id":"1337","type":"user","attributes":{"message":"@member Please check this out!"}}}`
var expectedUpdateAssigneeRequest3 = `{"data":{"id":"1337","type":"group","attributes":{"message":"@member Please check this out!"}}}`

func Test_ReportService_UpdateAssignee(t *testing.T) {
	// Verify that an invalid url fails
	c := NewClient(nil)
	c.BaseURL = &url.URL{}
	_, _, err := c.Report.UpdateAssignee("%A", "", nil)
	assert.NotNil(t, err)

	// Verify that an error response fails
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Oh No", 500)
	}))
	defer errorServer.Close()
	u, err := url.Parse(errorServer.URL)
	assert.Nil(t, err)
	c.BaseURL = u
	_, _, err = c.Report.UpdateAssignee("123456", "", nil)
	assert.NotNil(t, err)

	// Verify that it gets a response correctly and it has the correct request body
	reportServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(string(body)) != strings.TrimSpace(expectedUpdateAssigneeRequest1) {
			http.Error(w, "Non-matching request!", http.StatusBadRequest)
			return
		}
		http.ServeFile(w, r, "tests/responses/report.json")
	}))
	defer reportServer.Close()
	u, err = url.Parse(reportServer.URL)
	assert.Nil(t, err)
	c.BaseURL = u
	actual, _, err := c.Report.UpdateAssignee("123456", "@member Please check this out!", nil)
	assert.Nil(t, err)
	assert.Equal(t, &expectedReport, actual)

	// Verify that it gets a response correctly and it has the correct request body
	reportServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(string(body)) != strings.TrimSpace(expectedUpdateAssigneeRequest2) {
			http.Error(w, "Non-matching request!", http.StatusBadRequest)
			return
		}
		http.ServeFile(w, r, "tests/responses/report.json")
	}))
	defer reportServer.Close()
	u, err = url.Parse(reportServer.URL)
	assert.Nil(t, err)
	c.BaseURL = u
	actual, _, err = c.Report.UpdateAssignee("123456", "@member Please check this out!", &User{
		ID: String("1337"),
	})
	assert.Nil(t, err)
	assert.Equal(t, &expectedReport, actual)

	// Verify that it gets a response correctly and it has the correct request body
	reportServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(string(body)) != strings.TrimSpace(expectedUpdateAssigneeRequest3) {
			http.Error(w, "Non-matching request!", http.StatusBadRequest)
			return
		}
		http.ServeFile(w, r, "tests/responses/report.json")
	}))
	defer reportServer.Close()
	u, err = url.Parse(reportServer.URL)
	assert.Nil(t, err)
	c.BaseURL = u
	actual, _, err = c.Report.UpdateAssignee("123456", "@member Please check this out!", &Group{
		ID: String("1337"),
	})
	assert.Nil(t, err)
	assert.Equal(t, &expectedReport, actual)
}

var expectedStateChangeRequest = `{"data":{"type":"state-change","attributes":{"message":"This vulnerability has been resolved. Thanks!","state":"resolved"}}}`

func Test_ReportService_UpdateState(t *testing.T) {
	// Verify that an invalid url fails
	c := NewClient(nil)
	c.BaseURL = &url.URL{}
	_, _, err := c.Report.UpdateState("%A", "", "")
	assert.NotNil(t, err)

	// Verify that an error response fails
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Oh No", 500)
	}))
	defer errorServer.Close()
	u, err := url.Parse(errorServer.URL)
	assert.Nil(t, err)
	c.BaseURL = u
	_, _, err = c.Report.UpdateState("123456", "", "")
	assert.NotNil(t, err)

	// Verify that it gets a response correctly and it has the correct request body
	reportServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(string(body)) != strings.TrimSpace(expectedStateChangeRequest) {
			http.Error(w, "Non-matching request!", http.StatusBadRequest)
			return
		}
		http.ServeFile(w, r, "tests/responses/report.json")
	}))
	defer reportServer.Close()
	u, err = url.Parse(reportServer.URL)
	assert.Nil(t, err)
	c.BaseURL = u
	actual, _, err := c.Report.UpdateState("123456", ReportStateResolved, "This vulnerability has been resolved. Thanks!")
	assert.Nil(t, err)
	assert.Equal(t, &expectedReport, actual)
}

func Test_ReportService_List(t *testing.T) {

	// Verify that an invalid url fails
	c := NewClient(nil)
	c.BaseURL = &url.URL{
		Scheme: "http://[fe80::1%en0]/",
	}
	_, _, err := c.Report.List(ReportListFilter{}, nil)
	assert.NotNil(t, err)

	// Verify that an error response fails
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Oh No", 500)
	}))
	defer errorServer.Close()
	u, err := url.Parse(errorServer.URL)
	assert.Nil(t, err)
	c.BaseURL = u
	_, _, err = c.Report.List(ReportListFilter{}, nil)
	assert.NotNil(t, err)

	// Verify that it gets a response correctly
	reportServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "tests/responses/report_list.json")
	}))
	defer reportServer.Close()
	u, err = url.Parse(reportServer.URL)
	assert.Nil(t, err)
	c.BaseURL = u
	actual, _, err := c.Report.List(ReportListFilter{}, nil)
	assert.Nil(t, err)
	assert.Equal(t, expectedReport, actual[0])

}

/*

// List returns all Reports matching the specified criteria
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#reports/query
func (s *ReportService) List(filterOpts ReportListOptions, listOpts *ListOptions) ([]Report, *Response, error) {
	opts := reportListRequest{
		ListOptions: listOpts,
		Filter:      filterOpts,
	}
	u, err := addOptions("reports", &opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	reports := new([]Report)
	resp, err := s.client.Do(req, reports)
	if err != nil {
		return nil, resp, err
	}

	return *reports, resp, err
}
*/
