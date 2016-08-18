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
	"fmt"
)

// TeamService handles communication with the report related methods of the H1 API.
type TeamService service

// TeamProfile represents a H1 team
type TeamProfile struct {
	Name          *string `json:"name"`
	TwitterHandle *string `json:"twitter_handle"`
	Website       *string `json:"website"`
	About         string  `json:"about"`
}

// Team represents a H1 team
type Team struct {
	ID                                uint64             `json:"id"`
	Handle                            *string            `json:"handle"`
	URL                               *string            `json:"url"`
	Profile                           TeamProfile        `json:"profile"`
	Policy                            *string            `json:"policy"`
	Scopes                            []string           `json:"scopes"`
	CoverColor                        *string            `json:"cover_color"`
	TwitterHandle                     *string            `json:"twitter_handle"`
	IBB                               *bool              `json:"ibb"`
	HasCoverPhoto                     *bool              `json:"has_cover_photo"`
	ProfilePictureURLs                ProfilePictureURLs `json:"profile_picture_urls"`
	ExternalURL                       *string            `json:"external_url"`
	RejectingSubmissions              *bool              `json:"rejecting_submissions"`
	OffersSwag                        *bool              `json:"offers_swag"`
	OffersBounties                    *bool              `json:"offers_bounties"`
	BountiesPaid                      *float64           `json:"bounties_paid"`
	ResearcherCount                   *uint64            `json:"researcher_count"`
	BugCount                          *uint64            `json:"bug_count"`
	BaseBounty                        *uint64            `json:"base_bounty"`
	ShowTotalBountiesPaid             *bool              `json:"show_total_bounties_paid"`
	ShowAverageBounty                 *bool              `json:"show_average_bounty"`
	ShowTopBounties                   *bool              `json:"show_top_bounties"`
	ShowMeanBountyTime                *bool              `json:"show_mean_bounty_time"`
	ShowMeanFirstResponseTime         *bool              `json:"show_mean_first_response_time"`
	ShowMeanResolutionTime            *bool              `json:"show_mean_resolution_time"`
	TargetSignal                      *int64             `json:"target_signal"`
	CurrentUserReachedAbuseLimit      *bool              `json:"current_user_reached_abuse_limit"`
	CurrentUserReachedTeamSignalLimit *bool              `json:"current_user_reached_team_signal_limit"`
	TeamsUploadCoverPhotoEnabled      *bool              `json:"teams_upload_cover_photo_enabled"`
	CanViewThanks                     *bool              `json:"can_view_thanks"`
	CanViewPolicyVersions             *bool              `json:"can_view_policy_versions"`
	LastPolicyChangeAt                *Timestamp         `json:"last_policy_change_at"`
	CanManageTeamMemberGroups         *bool              `json:"can_manage_team_member_groups"`
	CanInviteTeamMember               *bool              `json:"can_invite_team_member"`
	HackbotGeniusEnabled              *bool              `json:"hackbot_genius_enabled"`
	Permissions                       []string           `json:"permissions"`
}

// GetByHandle a team by handle
func (s *TeamService) Get(handle string) (*Team, *Response, error) {
	req, err := s.client.NewRequest("GET", handle, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	team := new(Team)
	resp, err := s.client.Do(req, team)
	if err != nil {
		return nil, resp, err
	}

	return team, resp, err
}

// TeamProfileMetrics represents a H1 team's profile metrics
type TeamProfileMetrics struct {
	MeanTimeToFirstResponse *float64 `json:"mean_time_to_first_response"`
	MeanTimeToResolution    *float64 `json:"mean_time_to_resolution"`
	MeanTimeToBounty        *float64 `json:"mean_time_to_bounty"`
	TotalBountiesPaidPrefix *string  `json:"total_bounties_paid_prefix"` // TODO: Is this the correct type?
	TotalBountiesPaid       *float64 `json:"total_bounties_paid"`
	AverageBountyLowerRange *float64 `json:"average_bounty_lower_range"`
	AverageBountyUpperRange *float64 `json:"average_bounty_upper_range"`
	TopBountyLowerRange     *float64 `json:"top_bounty_lower_range"`
	TopBountyUpperRange     *float64 `json:"top_bounty_upper_range"`
}

// GetProfileMetrics returns the profile metrics of a team by handle
func (s *TeamService) GetProfileMetrics(handle string) (*TeamProfileMetrics, *Response, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("%s/reporters", handle), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	profileMetrics := new(TeamProfileMetrics)
	resp, err := s.client.Do(req, profileMetrics)
	if err != nil {
		return nil, resp, err
	}

	return profileMetrics, resp, err
}

// ListReporters returns the reporters to a team by handle
func (s *TeamService) ListReporters(handle string) ([]User, *Response, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("%s/reporters", handle), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	var users []User
	resp, err := s.client.Do(req, &users)
	if err != nil {
		return nil, resp, err
	}

	return users, resp, err
}

// CommonResponse represents a common ersponse object
type CommonResponse struct {
	ID      *uint64 `json:"id"`
	Title   *string `json:"title"`
	Message *string `json:"message"`
}

// ListCommonResponses returns the common responses for a team by handle
func (s *TeamService) ListCommonResponses(handle string) ([]CommonResponse, *Response, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("%s/common_responses.json", handle), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	var commonResponses []CommonResponse
	resp, err := s.client.Do(req, &commonResponses)
	if err != nil {
		return nil, resp, err
	}

	return commonResponses, resp, err
}
