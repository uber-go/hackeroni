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
//	"fmt"
)

// UserService handles communication with the report related methods of the H1 API.
type UserService service

// UserTeamContext appears on some user objects
type UserTeamContext struct {
	NumberOfReportsToSameTeam            *uint64 `json:"number_of_reports_to_same_team"`
	NumberOfResolvedReportsToSameTeam    *uint64 `json:"number_of_resolved_reports_to_same_team"`
	NumberOfBountiesReceivedFromSameTeam *uint64 `json:"number_of_bounties_received_from_same_team"`
	SumBountyAmountReceivedFromSameTeam  *string `json:"sum_bounty_amount_received_from_same_team"`
}

// User represents a H1 user
type User struct {
	ID                 *uint64            `json:"id"`
	Username           *string            `json:"username"`
	Name               *string            `json:"name"`
	Biography          *string            `json:"bio"`
	URL                *string            `json:"url"`
	ProfilePictureURLs ProfilePictureURLs `json:"profile_picture_urls"`
	Disabled           *bool              `json:"disabled"`
	ReportCount        *uint64            `json:"report_count"`
	TargetCount        *uint64            `json:"target_count"`
	Reputation         *uint64            `json:"reputation"`
	Rank               *uint64            `json:"rank"`
	Signal             *float64           `json:"signal"`
	Impact             *float64           `json:"impact"`
	SignalPercentile   *uint              `json:"signal_percentile"`
	ImpactPercentile   *uint              `json:"impact_percentile"`
	TeamContext        *UserTeamContext   `json:"team_context"`
}

/*// Get looks up a user by username
func (s *UserService) Get(username string) (*User, *Response, error) {
	req, err := s.client.NewRequest("GET", username, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, err
}

// ListPublicTeams returns the public teams associated with a username
func (s *UserService) ListPublicTeams(username string) ([]Team, *Response, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("%s/public_teams", username), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	wrapper := &struct {
		Items []Team `json:"items"`
	}{
		Items: []Team{},
	}
	resp, err := s.client.Do(req, wrapper)
	if err != nil {
		return nil, resp, err
	}

	return wrapper.Items, resp, err
}

// ListThanks returns the thanks associated with a usernamme
func (s *UserService) ListThanks(username string) ([]Team, *Response, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("%s/thanks", username), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	wrapper := &struct {
		Items []Team `json:"items"`
	}{
		Items: []Team{},
	}
	resp, err := s.client.Do(req, wrapper)
	if err != nil {
		return nil, resp, err
	}

	return wrapper.Items, resp, err
}

// ListBadges returns the thanks associated with a usernamme
func (s *UserService) ListBadges(username string) ([]Badge, *Response, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("%s/badges", username), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	wrapper := &struct {
		Items []Badge `json:"items"`
	}{
		Items: []Badge{},
	}
	resp, err := s.client.Do(req, wrapper)
	if err != nil {
		return nil, resp, err
	}

	return wrapper.Items, resp, err
}*/
