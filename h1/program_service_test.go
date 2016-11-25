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

	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var expectedProgram = Program{
	ID:        String("1337"),
	Type:      String(ProgramType),
	Handle:    String("security"),
	CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	Groups: []*Group{
		&Group{
			ID:   String("2557"),
			Type: String(GroupType),
			Name: String("Standard"),
			Permissions: []*string{
				String(GroupPermissionReportManagement),
				String(GroupPermissionRewardManagement),
			},
			CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		},
		&Group{
			ID:   String("2558"),
			Type: String(GroupType),
			Name: String("Admin"),
			Permissions: []*string{
				String(GroupPermissionUserManagement),
				String(GroupPermissionProgramManagement),
			},
			CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		},
	},
	Members: []*Member{
		&Member{
			ID:   String("1339"),
			Type: String(MemberType),
			Permissions: []*string{
				String(MemberPermissionProgramManagement),
				String(MemberPermissionReportManagement),
				String(MemberPermissionRewardManagement),
				String(MemberPermissionUserManagement),
			},
			CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
			User: &User{
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
		},
	},
}

func Test_ProgramService_Get(t *testing.T) {
	// Verify that an invalid url fails
	c := NewClient(nil)
	c.BaseURL = &url.URL{}
	_, _, err := c.Program.Get("%A")
	assert.NotNil(t, err)

	// Verify that an error response fails
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Oh No", 500)
	}))
	defer errorServer.Close()
	u, err := url.Parse(errorServer.URL)
	assert.Nil(t, err)
	c.BaseURL = u
	_, _, err = c.Program.Get("123456")
	assert.NotNil(t, err)

	// Verify that it gets a response correctly
	programServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "tests/responses/program.json")
	}))
	defer programServer.Close()
	u, err = url.Parse(programServer.URL)
	assert.Nil(t, err)
	c.BaseURL = u
	actual, _, err := c.Program.Get("1337")
	assert.Nil(t, err)
	assert.Equal(t, &expectedProgram, actual)
}
