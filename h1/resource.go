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

// Used to unmarshal unknown types
type unknownResource struct {
	Type *string `json:"type"`
}

// Type represent the possible values for the "Type" attribute
const (
	ActivityAgreedOnGoingPublicType             string = "activity-agreed-on-going-public"
	ActivityBountyAwardedType                   string = "activity-bounty-awarded"
	ActivityBountySuggestedType                 string = "activity-bounty-suggested"
	ActivityBugClonedType                       string = "activity-bug-cloned"
	ActivityBugDuplicateType                    string = "activity-bug-duplicate"
	ActivityBugInformativeType                  string = "activity-bug-informative"
	ActivityBugNeedsMoreInfoType                string = "activity-bug-needs-more-info"
	ActivityBugNewType                          string = "activity-bug-new"
	ActivityBugNotApplicableType                string = "activity-bug-not-applicable"
	ActivityBugReopenedType                     string = "activity-bug-reopened"
	ActivityBugResolvedType                     string = "activity-bug-resolved"
	ActivityBugSpamType                         string = "activity-bug-spam"
	ActivityBugTriagedType                      string = "activity-bug-triaged"
	ActivityCommentType                         string = "activity-comment"
	ActivityExternalUserInvitationCancelledType string = "activity-external-user-invitation-cancelled"
	ActivityExternalUserInvitedType             string = "activity-external-user-invited"
	ActivityExternalUserJoinedType              string = "activity-external-user-joined"
	ActivityExternalUserRemovedType             string = "activity-external-user-removed"
	ActivityGroupAssignedToBugType              string = "activity-group-assigned-to-bug"
	ActivityHackerRequestedMediationType        string = "activity-hacker-requested-mediation"
	ActivityManuallyDisclosedType               string = "activity-manually-disclosed"
	ActivityMediationRequestedType              string = "activity-mediation-requested"
	ActivityNotEligibleForBountyType            string = "activity-not-eligible-for-bounty"
	ActivityReferenceIDAddedType                string = "activity-reference-id-added"
	ActivityReportBecamePublicType              string = "activity-report-became-public"
	ActivityReportTitleUpdatedType              string = "activity-report-title-updated"
	ActivityReportVulnerabilityTypesUpdatedType string = "activity-report-vulnerability-types-updated"
	ActivityReportSeverityUpdatedType           string = "activity-report-severity-updated"
	ActivitySwagAwardedType                     string = "activity-swag-awarded"
	ActivityUserAssignedToBugType               string = "activity-user-assigned-to-bug"
	ActivityUserBannedFromProgramType           string = "activity-user-banned-from-program"
	AddressType                                 string = "address"
	AttachmentType                              string = "attachment"
	BountyType                                  string = "bounty"
	GroupType                                   string = "group"
	ProgramType                                 string = "program"
	ReportSummaryType                           string = "report-summary"
	MemberType                                  string = "member"
	ReportType                                  string = "report"
	SwagType                                    string = "swag"
	SeverityType                                string = "severity"
	UserType                                    string = "user"
	VulnerabilityTypeType                       string = "vulnerability-type"
)

// Bool allocates a new bool value to store v at and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// String allocates a new bool value to store v at and returns a pointer to it.
func String(v string) *string { return &v }

// Int allocates a new bool value to store v at and returns a pointer to it.
func Int(v int) *int { return &v }

// Uint64 allocates a new uint64 value to store v at and returns a pointer to it.
func Uint64(v uint64) *uint64 { return &v }

// Float64 allocates a new float64 value to store v at and returns a pointer to it.
func Float64(v float64) *float64 { return &v }
