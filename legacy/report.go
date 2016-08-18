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
//	"encoding/json"
	"fmt"
	"net/url"
)

// ReportAbilities dictates what can be done with a report
type ReportAbilities struct {
	CanManage                   *bool   `json:"can_manage?"`
	CanExport                   *bool   `json:"can_export?"`
	CanAddComment               *bool   `json:"can_add_comment?"`
	CanChangeState              *bool   `json:"can_change_state?"`
	CanReopen                   *bool   `json:"can_reopen?"`
	CanAwardBounty              *bool   `json:"can_award_bounty?"`
	CanAwardSwag                *bool   `json:"can_award_swag?"`
	CanSuggestBountyAmount      *bool   `json:"can_suggest_bounty_amount?"`
	CanAssignToUser             *bool   `json:"can_assign_to_user?"`
	CanHideTimeline             *bool   `json:"can_hide_timeline?"`
	CanAgreeOnGoingPublic       *bool   `json:"can_agree_on_going_public?"`
	CanBePubliclyDisclosed      *bool   `json:"can_be_publicly_disclosed?"`
	CanPostInternalComments     *bool   `json:"can_post_internal_comments?"`
	CanManageCommonResponses    *bool   `json:"can_manage_common_responses?"`
	CanChangeTitle              *bool   `json:"can_change_title?"`
	CanChangeVulnerabilityTypes *bool   `json:"can_change_vulnerability_types?"`
	CanBeManuallyDisclosed      *bool   `json:"can_be_manually_disclosed?"`
	CanClone                    *bool   `json:"can_clone?"`
	CanClose                    *bool   `json:"can_close?"`
	CanBanResearcher            *bool   `json:"can_ban_researcher?"`
	AssignableTeamMembers       []User  `json:"assignable_team_members"`
	AssignableTeamMemberGroups  []Group `json:"assignable_team_member_groups"`
}

// Report represents a report
type Report struct {
	ID                               *uint64                   `json:"id"`
	URL                              *string                   `json:"url"`
	Title                            *string                   `json:"title"`
	State                            *string                   `json:"state"`
	Substate                         *string                   `json:"substate"`
	ReadableSubstate                 *string                   `json:"readable_substate"`
	CreatedAt                        *Timestamp                `json:"created_at"`
	Assignee                         *User                     `json:"assignee"` // TODO: this is probably wrong
	CreateReferenceURL               *string                   `json:"create_reference_url"`
	Reporter                         *User                     `json:"reporter"` // TODO: There are special objects like team_context here
	PromoteBounties                  *bool                     `json:"promote_bounties"`
	Team                             *Team                     `json:"team"`
	HasBounty                        *bool                     `json:"has_bounty?"`
	CanViewTeam                      *bool                     `json:"can_view_team"`
	IsExternalBug                    *bool                     `json:"is_external_bug"`
	IsParticipant                    *bool                     `json:"is_participant"`
	Stage                            *uint                     `json:"stage"` // TODO: No idea what this is and the type
	Public                           *bool                     `json:"public"`
	CVEIDs                           []string                  `json:"cve_ids"` // TODO: Is this the correct type?
	DisclosedAt                      *Timestamp                `json:"disclosed_at"`
	BugReporterAgreedOnGoingPublicAt *Timestamp                `json:"bug_reporter_agreed_on_going_public_at"`
	TeamMemberAgreedOnGoingPublicAt  *Timestamp                `json:"team_member_agreed_on_going_public_at"`
	MediationRequested               *bool                     `json:"mediation_requested"`
	Subscribed                       *bool                     `json:"subscribed"`
	SuggestedBounty                  *uint64                   `json:"suggested_bounty"`
	VulnerabilityInformation         *string                   `json:"vulnerability_information"`
	VulnerabilityInformationHTML     *string                   `json:"vulnerability_information_html"`
	Triggers                         map[string]CommonResponse `json:"triggers"`
	VulnerabilityTypes               []VulnerabilityType       `json:"vulnerability_types"`
	Attachments                      []Attachment              `json:"attachments"`
	Abilities                        ReportAbilities           `json:"abilities"`
	IsMemberOfTeam                   *bool                     `json:"is_member_of_team"`
	Activities                       []Activity                `json:"activities"`
	Summaries                        []ReportSummary           `json:"summaries"`
}

// ReportService handles communication with the bug related methods of H1.
type ReportService service

// Get retrieves a report by ID
func (s *ReportService) Get(id uint64) (*Report, *Response, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("reports/%d.json", id), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	report := new(Report)
	resp, err := s.client.Do(req, report)
	if err != nil {
		return nil, resp, err
	}

	return report, resp, err
}

// Create creates a new report and returns the ID
func (s *ReportService) Create(handle string, report *Report) (*Response, error) {
	user, resp, err := s.client.Session.GetCurrentUser()
	if err != nil {
		return resp, err
	}

	body := url.Values{
		"report[title]":                     []string{*report.Title},
		"report[vulnerability_information]": []string{*report.VulnerabilityInformation},
		"report[vulnerability_type_ids][]":  []string{"126768"}, // TODO: Don't hardcode this to None
	}

	req, err := s.client.NewRequest("POST", fmt.Sprintf("%s/reports", handle), bytes.NewBufferString(body.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("X-CSRF-Token", *user.CSRFToken)

	reportInfo := struct {
		Redirect *string `json:"redirect"`
		ReportID *uint64 `json:"report_id"`
	}{}
	resp, err = s.client.Do(req, &reportInfo)
	if err != nil {
		return resp, err
	}
	report.ID = reportInfo.ReportID

	return resp, err
}

// AddSummary assigns a group by ID to a report
/*func (s *ReportService) AddSummary(id uint64, summary *ReportSummary) (*Response, error) {
	user, resp, err := s.client.Session.GetCurrentUser()
	if err != nil {
		return resp, err
	}

	body := struct {
		*ReportSummary
		ActionType string `json:"action_type"`
	}{
		ReportSummary: summary,
		ActionType:    "publish",
	}

	buf := bytes.NewBuffer(nil)
	dec := json.NewEncoder(buf)
	if dec.Encode(&body) != nil {
		return nil, err
	}
	req, err := s.client.NewRequest("POST", fmt.Sprintf("reports/%d/summaries", id), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("X-CSRF-Token", *user.CSRFToken)

	resp, err = s.client.Do(req, summary)
	if err != nil {
		return resp, err
	}

	return resp, err
}*/

// ReportBulkResponse is used as a response for multiple report methods
type ReportBulkResponse struct {
	Flash   *string  `json:"flash"`
	Reports []Report `json:"reports"`
}

// bulk uses the /reports/bulk endpoint a report
func (s *ReportService) bulk(id uint64, replyAction string, body url.Values) (*ReportBulkResponse, *Response, error) {
	user, resp, err := s.client.Session.GetCurrentUser()
	if err != nil {
		return nil, resp, err
	}

	body.Add("reply_action", replyAction)
	body.Add("reports_count", "1")
	body.Add("report_ids[]", fmt.Sprintf("%d", id))

	req, err := s.client.NewRequest("POST", "/reports/bulk", bytes.NewBufferString(body.Encode()))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("X-CSRF-Token", *user.CSRFToken)

	var reportBulkResponse ReportBulkResponse
	resp, err = s.client.Do(req, &reportBulkResponse)
	if err != nil {
		return nil, resp, err
	}

	return &reportBulkResponse, resp, err
}

// ReportChangeStateOptions provides optional arguments to ReportService's ChangeState method
type ReportChangeStateOptions struct {
	Reference *string
}

// ChangeState changes a report state
func (s *ReportService) ChangeState(id uint64, message string, state string, opts *ReportChangeStateOptions) (*ReportBulkResponse, *Response, error) {
	body := url.Values{
		"message":  []string{message},
		"substate": []string{state},
	}

	if opts != nil {
		if opts.Reference != nil {
			body.Add("reference", *opts.Reference)
		}
	}

	return s.bulk(id, "change-state", body)
}

// ReportCloseOptions provides optional arguments to ReportService's Close method
type ReportCloseOptions struct {
	AddReporterToOriginal *bool
	OriginalReportID      *uint64
}

// Close closes a report
func (s *ReportService) Close(id uint64, message string, state string, opts *ReportCloseOptions) (*ReportBulkResponse, *Response, error) {
	body := url.Values{
		"message":  []string{message},
		"substate": []string{state},
	}

	if opts != nil {
		if opts.OriginalReportID != nil {
			body.Add("original_report_id", fmt.Sprintf("%d", *opts.OriginalReportID))
		}
	}

	return s.bulk(id, "close-report", body)
}

// Reopen reopens a report
/*func (s *ReportService) Reopen(id uint64, message string) (*ReportBulkResponse, *Response, error) {
	body := url.Values{
		"message": []string{message},
	}

	return s.bulk(id, "reopen", body)
}*/

// Comment comments on a report
func (s *ReportService) Comment(id uint64, message string, internal bool) (*ReportBulkResponse, *Response, error) {
	body := url.Values{
		"message": []string{message},
	}
	if internal {
		body.Add("is_internal", "true")
	} else {
		body.Add("is_internal", "false")
	}

	return s.bulk(id, "add-comment", body)
}

// AssignUser assigns a user by ID to a report
/*func (s *ReportService) AssignUser(id uint64, message string, assignee uint64) (*ReportBulkResponse, *Response, error) {
	body := url.Values{
		"message":     []string{message},
		"substate":    []string{"user"},
		"assignee_id": []string{fmt.Sprintf("%d", assignee)},
	}

	return s.bulk(id, "assign-to", body)
}

// AssignGroup assigns a group by ID to a report
func (s *ReportService) AssignGroup(id uint64, message string, assignee uint64) (*ReportBulkResponse, *Response, error) {
	body := url.Values{
		"message":     []string{message},
		"substate":    []string{"group"},
		"assignee_id": []string{fmt.Sprintf("%d", assignee)},
	}

	return s.bulk(id, "assign-to", body)
}*/
