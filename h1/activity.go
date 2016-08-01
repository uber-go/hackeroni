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
	"encoding/json"
)

// Activity represents activities that have occured in a given report.
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#activity
type Activity struct {
	report      *Report
	ID          *string         `json:"id"`
	Type        *string         `json:"type"`
	Message     *string         `json:"message"`
	Internal    *bool           `json:"internal"`
	CreatedAt   *Timestamp      `json:"created_at"`
	UpdatedAt   *Timestamp      `json:"updated_at"`
	RawActor    json.RawMessage `json:"actor"` // Used by the Actor() method
	Attachments []Attachment    `json:"attachments,omitempty"`
	rawData     []byte          // Used by the Activity() method
}

// Helper types for JSONUnmarshal
type activity Activity // Used to avoid recursion of JSONUnmarshal
type activityUnmarshalHelper struct {
	activity
	Attributes    *activity `json:"attributes"`
	Relationships struct {
		Attachments struct {
			Data []Attachment `json:"data"`
		} `json:"attachments,omitempty"`
		RawActor struct {
			Data json.RawMessage `json:"data"`
		} `json:"actor"`
	} `json:"relationships"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *Activity) UnmarshalJSON(b []byte) error {
	var helper activityUnmarshalHelper
	helper.Attributes = &helper.activity
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*a = Activity(helper.activity)
	a.Attachments = helper.Relationships.Attachments.Data
	a.RawActor = helper.Relationships.RawActor.Data
	a.rawData = b
	return nil
}

// Actor returns returns the parsed actor. For recognized actor types, a value of the corresponding struct type will be returned.
func (a *Activity) Actor() (actor interface{}) {
	var obj unknownResource
	if err := json.Unmarshal(a.RawActor, &obj); err != nil {
		panic(err.Error())
	}
	switch *obj.Type {
	case UserType:
		actor = &User{}
	case ProgramType:
		actor = &Program{}
	}
	if err := json.Unmarshal(a.RawActor, &actor); err != nil {
		panic(err.Error())
	}
	return actor
}

// Activity returns the parsed activity. For recognized activity types, a value of the corresponding struct type will be returned.
func (a *Activity) Activity() (activity interface{}) {
	switch *a.Type {
	case ActivityBountyAwardedType:
		activity = &ActivityBountyAwarded{}
	case ActivityBountySuggestedType:
		activity = &ActivityBountySuggested{}
	case ActivityBugClonedType:
		activity = &ActivityBugCloned{}
	case ActivityExternalUserInvitationCancelledType:
		activity = &ActivityExternalUserInvitationCancelled{}
	case ActivityExternalUserInvitedType:
		activity = &ActivityExternalUserInvited{}
	case ActivityExternalUserJoinedType:
		activity = &ActivityExternalUserJoined{}
	case ActivityExternalUserRemovedType:
		activity = &ActivityExternalUserRemoved{}
	case ActivityGroupAssignedToBugType:
		activity = &ActivityGroupAssignedToBug{}
	case ActivityReferenceIDAddedType:
		activity = &ActivityReferenceIDAdded{}
	case ActivityReportTitleUpdatedType:
		activity = &ActivityReportTitleUpdated{}
	case ActivityReportVulnerabilityTypesUpdatedType:
		activity = &ActivityReportVulnerabilityTypesUpdated{}
	case ActivitySwagAwardedType:
		activity = &ActivitySwagAwarded{}
	case ActivityUserAssignedToBugType:
		activity = &ActivityUserAssignedToBug{}
	case ActivityUserBannedFromProgramType:
		activity = &ActivityUserBannedFromProgram{}
	default:
		return nil
	}
	if err := json.Unmarshal(a.rawData, activity); err != nil {
		panic(err.Error())
	}
	return activity
}

// Report returns the report this activity is a child of
func (a *Activity) Report() *Report {
	return a.report
}

// ActivityBountyAwarded occurs when a bounty is awarded.
//
// HackerOne API docs:https://api.hackerone.com/docs/v1#activity-bounty-awarded
type ActivityBountyAwarded struct {
	BountyAmount *string `json:"bounty_amount"`
	BonusAmount  *string `json:"bonus_amount"`
}

// Helper types for JSONUnmarshal
type activityBountyAwarded ActivityBountyAwarded // Used to avoid recursion of JSONUnmarshal
type activityBountyAwardedUnmarshalHelper struct {
	Attributes activityBountyAwarded `json:"attributes"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *ActivityBountyAwarded) UnmarshalJSON(b []byte) error {
	var helper activityBountyAwardedUnmarshalHelper
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*a = ActivityBountyAwarded(helper.Attributes)
	return nil
}

// ActivityBountySuggested occurs when a bounty is suggested.
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#activity-bounty-suggested
type ActivityBountySuggested struct {
	BountyAmount *string `json:"bounty_amount"`
	BonusAmount  *string `json:"bonus_amount"`
}

// Helper types for JSONUnmarshal
type activityBountySuggested ActivityBountySuggested // Used to avoid recursion of JSONUnmarshal
type activityBountySuggestedUnmarshalHelper struct {
	Attributes activityBountySuggested `json:"attributes"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *ActivityBountySuggested) UnmarshalJSON(b []byte) error {
	var helper activityBountySuggestedUnmarshalHelper
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*a = ActivityBountySuggested(helper.Attributes)
	return nil
}

// ActivityBugCloned occurs when a bug is cloned.
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#activity-bug-cloned
type ActivityBugCloned struct {
	OriginalReportID *int `json:"original_report_id"`
}

// Helper types for JSONUnmarshal
type activityBugCloned ActivityBugCloned // Used to avoid recursion of JSONUnmarshal
type activityBugClonedUnmarshalHelper struct {
	Attributes activityBugCloned `json:"attributes"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *ActivityBugCloned) UnmarshalJSON(b []byte) error {
	var helper activityBugClonedUnmarshalHelper
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*a = ActivityBugCloned(helper.Attributes)
	return nil
}

// ActivityExternalUserInvitationCancelled occurs when a external user's invitiation is cancelled.
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#activity-external-user-invitation-cancelled
type ActivityExternalUserInvitationCancelled struct {
	Email *string `json:"email"`
}

// Helper types for JSONUnmarshal
type activityExternalUserInvitationCancelled ActivityExternalUserInvitationCancelled // Used to avoid recursion of JSONUnmarshal
type activityExternalUserInvitationCancelledSuggestedUnmarshalHelper struct {
	Attributes activityExternalUserInvitationCancelled `json:"attributes"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *ActivityExternalUserInvitationCancelled) UnmarshalJSON(b []byte) error {
	var helper activityExternalUserInvitationCancelledSuggestedUnmarshalHelper
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*a = ActivityExternalUserInvitationCancelled(helper.Attributes)
	return nil
}

// ActivityExternalUserInvited occurs when a external user is invited.
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#activity-external-user-invited
type ActivityExternalUserInvited struct {
	Email *string `json:"email"`
}

// Helper types for JSONUnmarshal
type activityExternalUserInvited ActivityExternalUserInvited // Used to avoid recursion of JSONUnmarshal
type activityExternalUserInvitedUnmarshalHelper struct {
	Attributes activityExternalUserInvited `json:"attributes"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *ActivityExternalUserInvited) UnmarshalJSON(b []byte) error {
	var helper activityExternalUserInvitedUnmarshalHelper
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*a = ActivityExternalUserInvited(helper.Attributes)
	return nil
}

// ActivityExternalUserJoined occurs when a external user joins.
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#activity-external-user-joined
type ActivityExternalUserJoined struct {
	DuplicateReportID *int `json:"duplicate_report_id"`
}

// Helper types for JSONUnmarshal
type activityExternalUserJoined ActivityExternalUserJoined // Used to avoid recursion of JSONUnmarshal
type activityExternalUserJoinedUnmarshalHelper struct {
	Attributes activityExternalUserJoined `json:"attributes"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *ActivityExternalUserJoined) UnmarshalJSON(b []byte) error {
	var helper activityExternalUserJoinedUnmarshalHelper
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*a = ActivityExternalUserJoined(helper.Attributes)
	return nil
}

// ActivityExternalUserRemoved occurs when a external user is removed
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#activity-external-user-removed
type ActivityExternalUserRemoved struct {
	RemovedUser *User `json:"removed_user"`
}

// Helper types for JSONUnmarshal
type activityExternalUserRemoved ActivityExternalUserRemoved // Used to avoid recursion of JSONUnmarshal
type activityExternalUserRemovedUnmarshalHelper struct {
	Relationships struct {
		RemovedUser struct {
			Data *User `json:"data"`
		} `json:"removed_user"`
	} `json:"relationships"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *ActivityExternalUserRemoved) UnmarshalJSON(b []byte) error {
	var helper activityExternalUserRemovedUnmarshalHelper
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	a.RemovedUser = helper.Relationships.RemovedUser.Data
	return nil
}

// ActivityGroupAssignedToBug occurs when a group is assigned to a report.
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#activity-group-assigned-to-bug
type ActivityGroupAssignedToBug struct {
	Group *Group `json:"group"`
}

// Helper types for JSONUnmarshal
type activityGroupAssignedToBug ActivityGroupAssignedToBug // Used to avoid recursion of JSONUnmarshal
type activityGroupAssignedToBugUnmarshalHelper struct {
	Relationships struct {
		Group struct {
			Data *Group `json:"data"`
		} `json:"group"`
	} `json:"relationships"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *ActivityGroupAssignedToBug) UnmarshalJSON(b []byte) error {
	var helper activityGroupAssignedToBugUnmarshalHelper
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	a.Group = helper.Relationships.Group.Data
	return nil
}

// ActivityReferenceIDAdded occurs when a reference id/url is added to a report.
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#activity-reference-id-added
type ActivityReferenceIDAdded struct {
	Reference    *string `json:"reference"`
	ReferenceURL *string `json:"reference_url"`
}

// Helper types for JSONUnmarshal
type activityReferenceIDAdded ActivityReferenceIDAdded // Used to avoid recursion of JSONUnmarshal
type activityReferenceIDAddedUnmarshalHelper struct {
	Attributes activityReferenceIDAdded `json:"attributes"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *ActivityReferenceIDAdded) UnmarshalJSON(b []byte) error {
	var helper activityReferenceIDAddedUnmarshalHelper
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*a = ActivityReferenceIDAdded(helper.Attributes)
	return nil
}

// ActivityReportTitleUpdated occurs when report title is updated
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#activity-report-title-updated
type ActivityReportTitleUpdated struct {
	OldTitle *string `json:"old_title"`
	NewTitle *string `json:"new_title"`
}

// Helper types for JSONUnmarshal
type activityReportTitleUpdated ActivityReportTitleUpdated // Used to avoid recursion of JSONUnmarshal
type activityReportTitleUpdatedUnmarshalHelper struct {
	Attributes activityReportTitleUpdated `json:"attributes"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *ActivityReportTitleUpdated) UnmarshalJSON(b []byte) error {
	var helper activityReportTitleUpdatedUnmarshalHelper
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*a = ActivityReportTitleUpdated(helper.Attributes)
	return nil
}

// ActivityReportVulnerabilityTypesUpdated occurs when vulnerability types for a report are updated.
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#activity-report-vulnerability-types-updated
type ActivityReportVulnerabilityTypesUpdated struct {
	OldVulnerabilityTypes []VulnerabilityType `json:"old_vulnerability_types"`
	NewVulnerabilityTypes []VulnerabilityType `json:"new_vulnerability_types"`
}

// Helper types for JSONUnmarshal
type activityReportVulnerabilityTypesUpdated ActivityReportVulnerabilityTypesUpdated // Used to avoid recursion of JSONUnmarshal
type activityReportVulnerabilityTypesUpdatedUnmarshalHelper struct {
	Relationships struct {
		OldVulnerabilityTypes struct {
			Data []VulnerabilityType `json:"data"`
		} `json:"old_vulnerability_types"`
		NewVulnerabilityTypes struct {
			Data []VulnerabilityType `json:"data"`
		} `json:"new_vulnerability_types"`
	} `json:"relationships"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *ActivityReportVulnerabilityTypesUpdated) UnmarshalJSON(b []byte) error {
	var helper activityReportVulnerabilityTypesUpdatedUnmarshalHelper
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	a.OldVulnerabilityTypes = helper.Relationships.OldVulnerabilityTypes.Data
	a.NewVulnerabilityTypes = helper.Relationships.NewVulnerabilityTypes.Data
	return nil
}

// ActivitySwagAwarded occurs when swag is awarded
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#activity-swag-awarded
type ActivitySwagAwarded struct {
	Swag *Swag `json:"swag"`
}

// Helper types for JSONUnmarshal
type activitySwagAwarded ActivitySwagAwarded // Used to avoid recursion of JSONUnmarshal
type activitySwagAwardedUnmarshalHelper struct {
	Relationships struct {
		Swag struct {
			Data *Swag `json:"data"`
		} `json:"swag"`
	} `json:"relationships"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *ActivitySwagAwarded) UnmarshalJSON(b []byte) error {
	var helper activitySwagAwardedUnmarshalHelper
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	a.Swag = helper.Relationships.Swag.Data
	return nil
}

// ActivityUserAssignedToBug occurs when a user is assigned to a report.
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#activity-user-assigned-to-bug
type ActivityUserAssignedToBug struct {
	AssignedUser *User `json:"assigned_user"`
}

// Helper types for JSONUnmarshal
type activityUserAssignedToBug ActivityUserAssignedToBug // Used to avoid recursion of JSONUnmarshal
type activityUserAssignedToBugUnmarshalHelper struct {
	Relationships struct {
		AssignedUser struct {
			Data *User `json:"data"`
		} `json:"assigned_user"`
	} `json:"relationships"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *ActivityUserAssignedToBug) UnmarshalJSON(b []byte) error {
	var helper activityUserAssignedToBugUnmarshalHelper
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	a.AssignedUser = helper.Relationships.AssignedUser.Data
	return nil
}

// ActivityUserBannedFromProgram occurs when a user is banned from a program.
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#activity-user-banned-from-program
type ActivityUserBannedFromProgram struct {
	RemovedUser *User `json:"removed_user"`
}

// Helper types for JSONUnmarshal
type activityUserBannedFromProgram ActivityUserBannedFromProgram // Used to avoid recursion of JSONUnmarshal
type activityUserBannedFromProgramUnmarshalHelper struct {
	Relationships struct {
		RemovedUser struct {
			Data *User `json:"data"`
		} `json:"removed_user"`
	} `json:"relationships"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *ActivityUserBannedFromProgram) UnmarshalJSON(b []byte) error {
	var helper activityUserBannedFromProgramUnmarshalHelper
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	a.RemovedUser = helper.Relationships.RemovedUser.Data
	return nil
}
