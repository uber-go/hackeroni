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
	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_Session_GetCurrentUser(t *testing.T) {
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "tests/current_user.json")
	}))
	defer errorServer.Close()
	u, err := url.Parse(errorServer.URL)
	assert.Nil(t, err)

	client := NewClient(nil)
	client.BaseURL = u
	assert.Nil(t, err)

	user, _, err := client.Session.GetCurrentUser()
	assert.Nil(t, err)
	assert.Equal(t, &SessionUser{
		CSRFToken: String("jJz31wvEtnbmhsoa+Iv4eDJknSHbccnAr0Qss10rp5QJQIwh+lPuDa3WPMV0DgHiIvwDpRCxX5XKkoQ+bw4vOw=="),
		SignedIn:  Bool(false),
	}, user)
}
