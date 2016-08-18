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

import (
	"bytes"
	"encoding/json"
	"net/url"
)

// SessionService handles communication with the session related methods of H1.
type SessionService service

// SessionUser describes the current user for a session
type SessionUser struct {
	CSRFToken *string `json:"csrf_token"`
	SignedIn  *bool   `json:"signed_in?"`
}

// GetCurrentUser returns information about the logged in user for the session (including the current CSRF token)
func (s *SessionService) GetCurrentUser() (*SessionUser, *Response, error) {
	req, err := s.client.NewRequest("GET", "current_user", nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, resp, err
	}

	var user SessionUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, resp, err
	}

	return &user, resp, nil
}

// GetTeams returns the current user's teams
/*func (s *SessionService) GetTeams() ([]Team, *Response, error) {
	req, err := s.client.NewRequest("GET", "teams", nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	var teams []Team
	resp, err := s.client.Do(req, &teams)
	if err != nil {
		return nil, resp, err
	}

	return teams, resp, err
}*/

// Acquire attempts to authenticate with the provided credentials
func (s *SessionService) Acquire(email string, password string) (*Response, error) {
	user, resp, err := s.GetCurrentUser()
	if err != nil {
		return resp, err
	}

	body := url.Values{
		"authenticity_token": []string{*user.CSRFToken},
		"user[email]":        []string{email},
		"user[password]":     []string{password},
	}

	req, err := s.client.NewRequest("POST", "users/sign_in", bytes.NewBufferString(body.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "text/html")

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
