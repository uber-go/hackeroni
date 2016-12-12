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
	"github.com/stretchr/testify/require"

	"io/ioutil"
	"testing"
)

func Test_ActivityActor_Panic(t *testing.T) {
	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Activity.Actor() with invalid JSON should panic")
			}
		}()
		actual := Activity{
			RawActor: []byte("Invalid JSON"),
		}
		actual.Actor()
	}()
	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Activity.Actor() with incorrect JSON should panic")
			}
		}()
		actual := Activity{
			RawActor: []byte(`{"type":"user","id":123}`),
		}
		actual.Actor()
	}()
}

func Test_ActivityActor_User(t *testing.T) {
	actor, err := ioutil.ReadFile("tests/resources/user.json")
	require.Nil(t, err)
	actual := Activity{
		RawActor: actor,
	}
	actualActor := actual.Actor().(*User)
	expectedActor := &User{
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
		Reputation: Uint64(7),
		Signal:     Float64(7.0),
		Impact:     Float64(30.0),
		CreatedAt:  NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	assert.Equal(t, expectedActor, actualActor)
}

func Test_ActivityActor_Program(t *testing.T) {
	actor, err := ioutil.ReadFile("tests/resources/program.json")
	require.Nil(t, err)
	actual := Activity{
		RawActor: actor,
	}
	actualActor := actual.Actor().(*Program)
	expectedActor := &Program{
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
	assert.Equal(t, expectedActor, actualActor)
}

func Test_ActivityActivity_Nil(t *testing.T) {
	actual := Activity{
		Type: String("unknown"),
	}
	actualActor := actual.Activity()
	assert.Nil(t, actualActor)
}

func Test_ActivityActivity_Panic(t *testing.T) {
	defer func() {
		if recover() == nil {
			assert.Fail(t, "Activity.Activity() with invalid JSON should panic")
		}
	}()
	actual := Activity{
		Type:    String(ActivityBountyAwardedType),
		rawData: []byte("Invalid JSON"),
	}
	actual.Activity()
}

func Test_ActivityAgreedOnGoingPublic(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-agreed-on-going-public.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityAgreedOnGoingPublicType),
		Message:   String("Agreed On Going Public!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityBountyAwarded(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-bounty-awarded.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityBountyAwardedType),
		Message:   String("Bounty Awarded!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}

	actualActivity := actual.Activity().(*ActivityBountyAwarded)
	expectedActivity := &ActivityBountyAwarded{
		BountyAmount: String("500"),
		BonusAmount:  String("50"),
	}
	assert.Equal(t, expectedActivity, actualActivity)

	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Activity.Activity() with incorrect JSON should panic")
			}
		}()
		actual.rawData = []byte(`{"attributes":123}`)
		actual.Activity()
	}()

	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)

}

func Test_ActivityBountySuggested(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-bounty-suggested.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityBountySuggestedType),
		Message:   String("Bounty Suggested!"),
		Internal:  Bool(true),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}

	actualActivity := actual.Activity().(*ActivityBountySuggested)
	expectedActivity := &ActivityBountySuggested{
		BountyAmount: String("500"),
		BonusAmount:  String("50"),
	}
	assert.Equal(t, expectedActivity, actualActivity)

	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Activity.Activity() with incorrect JSON should panic")
			}
		}()
		actual.rawData = []byte(`{"attributes":123}`)
		actual.Activity()
	}()

	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityBugCloned(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-bug-cloned.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityBugClonedType),
		Message:   String("Bug Cloned!"),
		Internal:  Bool(true),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}

	actualActivity := actual.Activity().(*ActivityBugCloned)
	expectedActivity := &ActivityBugCloned{
		OriginalReportID: Int(1336),
	}
	assert.Equal(t, expectedActivity, actualActivity)

	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Activity.Activity() with incorrect JSON should panic")
			}
		}()
		actual.rawData = []byte(`{"attributes":123}`)
		actual.Activity()
	}()

	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityBugDuplicate(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-bug-duplicate.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityBugDuplicateType),
		Message:   String("Bug Duplicate!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityBugInformative(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-bug-informative.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityBugInformativeType),
		Message:   String("Bug Informative!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityBugNeedsMoreInfo(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-bug-needs-more-info.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityBugNeedsMoreInfoType),
		Message:   String("Bug Needs More Info!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityBugNew(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-bug-new.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityBugNewType),
		Message:   String("Bug New!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityBugNotApplicable(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-bug-not-applicable.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityBugNotApplicableType),
		Message:   String("Bug Not Applicable!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityBugReopened(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-bug-reopened.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityBugReopenedType),
		Message:   String("Bug Reopened!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityBugResolved(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-bug-resolved.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityBugResolvedType),
		Message:   String("Bug Resolved!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityBugSpam(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-bug-spam.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityBugSpamType),
		Message:   String("Bug Spam!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityBugTriaged(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-bug-triaged.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityBugTriagedType),
		Message:   String("Bug Triaged!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityComment(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-comment.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityCommentType),
		Message:   String("Comment!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		Attachments: []Attachment{
			Attachment{
				ID:          String("1337"),
				Type:        String(AttachmentType),
				FileName:    String("root.rb"),
				ContentType: String("text/x-c++"),
				FileSize:    Int(2873),
				ExpiringURL: String("/system/attachments/files/000/001/337/original/root.rb?1454385906"),
				CreatedAt:   NewTimestamp("2016-02-02T04:05:06.000Z"),
			},
		},
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityExternalUserInvitationCancelled(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-external-user-invitation-cancelled.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityExternalUserInvitationCancelledType),
		Message:   String("External User Invitation Cancelled!"),
		Internal:  Bool(true),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}

	actualActivity := actual.Activity().(*ActivityExternalUserInvitationCancelled)
	expectedActivity := &ActivityExternalUserInvitationCancelled{
		Email: String("hacker@example.com"),
	}
	assert.Equal(t, expectedActivity, actualActivity)

	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Activity.Activity() with incorrect JSON should panic")
			}
		}()
		actual.rawData = []byte(`{"attributes":123}`)
		actual.Activity()
	}()

	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityExternalUserInvited(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-external-user-invited.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityExternalUserInvitedType),
		Message:   String("External User Invited!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}

	actualActivity := actual.Activity().(*ActivityExternalUserInvited)
	expectedActivity := &ActivityExternalUserInvited{
		Email: String("hacker@example.com"),
	}
	assert.Equal(t, expectedActivity, actualActivity)

	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Activity.Activity() with incorrect JSON should panic")
			}
		}()
		actual.rawData = []byte(`{"attributes":123}`)
		actual.Activity()
	}()

	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityExternalUserJoined(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-external-user-joined.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityExternalUserJoinedType),
		Message:   String("External User Joined!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}

	actualActivity := actual.Activity().(*ActivityExternalUserJoined)
	expectedActivity := &ActivityExternalUserJoined{
		DuplicateReportID: Int(10),
	}
	assert.Equal(t, expectedActivity, actualActivity)

	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Activity.Activity() with incorrect JSON should panic")
			}
		}()
		actual.rawData = []byte(`{"attributes":123}`)
		actual.Activity()
	}()

	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityExternalUserRemoved(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-external-user-removed.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityExternalUserRemovedType),
		Message:   String("External User Removed!"),
		Internal:  Bool(true),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}

	actualActivity := actual.Activity().(*ActivityExternalUserRemoved)
	expectedActivity := &ActivityExternalUserRemoved{
		RemovedUser: &User{
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
	}
	assert.Equal(t, expectedActivity, actualActivity)

	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Activity.Activity() with incorrect JSON should panic")
			}
		}()
		actual.rawData = []byte(`{"relationships":123}`)
		actual.Activity()
	}()

	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityGroupAssignedToBug(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-group-assigned-to-bug.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityGroupAssignedToBugType),
		Message:   String("Group Assigned To Bug!"),
		Internal:  Bool(true),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}

	actualActivity := actual.Activity().(*ActivityGroupAssignedToBug)
	expectedActivity := &ActivityGroupAssignedToBug{
		Group: &Group{
			ID:        String("1337"),
			Type:      String(GroupType),
			Name:      String("Admin"),
			CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		},
	}
	assert.Equal(t, expectedActivity, actualActivity)

	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Activity.Activity() with incorrect JSON should panic")
			}
		}()
		actual.rawData = []byte(`{"relationships":123}`)
		actual.Activity()
	}()

	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityHackerRequestedMediation(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-hacker-requested-mediation.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityHackerRequestedMediationType),
		Message:   String("Hacker Requested Mediation!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityManuallyDisclosed(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-manually-disclosed.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityManuallyDisclosedType),
		Message:   String("Manually Disclosed!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityMediationRequested(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-mediation-requested.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityMediationRequestedType),
		Message:   String("Mediation Requested!"),
		Internal:  Bool(true),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityNotEligibleForBounty(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-not-eligible-for-bounty.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityNotEligibleForBountyType),
		Message:   String("Not Eligible For Bounty!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityReferenceIDAdded(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-reference-id-added.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityReferenceIDAddedType),
		Message:   String("Reference Id Added!"),
		Internal:  Bool(true),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}

	actualActivity := actual.Activity().(*ActivityReferenceIDAdded)
	expectedActivity := &ActivityReferenceIDAdded{
		Reference:    String("reference"),
		ReferenceURL: String("example.com/reference"),
	}
	assert.Equal(t, expectedActivity, actualActivity)

	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Activity.Activity() with incorrect JSON should panic")
			}
		}()
		actual.rawData = []byte(`{"attributes":123}`)
		actual.Activity()
	}()

	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityReportBecamePublic(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-report-became-public.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityReportBecamePublicType),
		Message:   String("Report Became Public!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityReportTitleUpdated(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-report-title-updated.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityReportTitleUpdatedType),
		Message:   String("Report Title Updated!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}

	actualActivity := actual.Activity().(*ActivityReportTitleUpdated)
	expectedActivity := &ActivityReportTitleUpdated{
		OldTitle: String("xss"),
		NewTitle: String("XSS in login form"),
	}
	assert.Equal(t, expectedActivity, actualActivity)

	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Activity.Activity() with incorrect JSON should panic")
			}
		}()
		actual.rawData = []byte(`{"attributes":123}`)
		actual.Activity()
	}()

	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

// Disabled for now because of in-consistant results from API
func Test_ActivityReportVulnerabilityTypesUpdated(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-report-vulnerability-types-updated.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityReportVulnerabilityTypesUpdatedType),
		Message:   String("Report Vulnerability Types Updated!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}

	actualActivity := actual.Activity().(*ActivityReportVulnerabilityTypesUpdated)
	expectedActivity := &ActivityReportVulnerabilityTypesUpdated{
		OldVulnerabilityTypes: []VulnerabilityType{
			VulnerabilityType{
				ID:          String("1337"),
				Type:        String(VulnerabilityTypeType),
				Name:        String("Cross-Site Scripting (XSS)"),
				Description: String("Failure of a site to validate, filter, or encode user input before returning it to another user's web client."),
				CreatedAt:   NewTimestamp("2016-02-02T04:05:06.000Z"),
			},
		},
		NewVulnerabilityTypes: []VulnerabilityType{
			VulnerabilityType{
				ID:          String("1338"),
				Type:        String(VulnerabilityTypeType),
				Name:        String("UI Redressing (Clickjacking)"),
				Description: String("Tricking the user to click on a disguised user interface element, ultimately to reveal sensitive information."),
				CreatedAt:   NewTimestamp("2016-02-02T04:05:06.000Z"),
			},
		},
	}
	assert.Equal(t, expectedActivity, actualActivity)

	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Activity.Activity() with incorrect JSON should panic")
			}
		}()
		actual.rawData = []byte(`{"relationships":123}`)
		actual.Activity()
	}()

	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityReportSeverityUpdated(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-report-severity-updated.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityReportSeverityUpdatedType),
		Message:   String("Report Severity Updated!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivitySwagAwarded(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-swag-awarded.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivitySwagAwardedType),
		Message:   String("Swag Awarded!"),
		Internal:  Bool(false),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}

	actualActivity := actual.Activity().(*ActivitySwagAwarded)
	expectedActivity := &ActivitySwagAwarded{
		Swag: &Swag{
			ID:        String("1337"),
			Type:      String(SwagType),
			Sent:      Bool(false),
			CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
			Address: &Address{
				ID:          String("1337"),
				Type:        String(AddressType),
				Name:        String("Jane Doe"),
				Street:      String("535 Mission Street"),
				City:        String("San Francisco"),
				PostalCode:  String("94105"),
				State:       String("CA"),
				Country:     String("United States of America"),
				TShirtSize:  String("Large"),
				PhoneNumber: String("+1-510-000-0000"),
				CreatedAt:   NewTimestamp("2016-02-02T04:05:06.000Z"),
			},
		},
	}
	assert.Equal(t, expectedActivity, actualActivity)

	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Activity.Activity() with incorrect JSON should panic")
			}
		}()
		actual.rawData = []byte(`{"relationships":123}`)
		actual.Activity()
	}()

	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityUserAssignedToBug(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-user-assigned-to-bug.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityUserAssignedToBugType),
		Message:   String("User Assigned To Bug!"),
		Internal:  Bool(true),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}

	actualActivity := actual.Activity().(*ActivityUserAssignedToBug)
	expectedActivity := &ActivityUserAssignedToBug{
		AssignedUser: &User{
			ID:       String("1336"),
			Type:     String(UserType),
			Disabled: Bool(false),
			Username: String("other_user"),
			Name:     String("Other User"),
			ProfilePicture: UserProfilePicture{
				Size62x62:   String("/assets/avatars/default.png"),
				Size82x82:   String("/assets/avatars/default.png"),
				Size110x110: String("/assets/avatars/default.png"),
				Size260x260: String("/assets/avatars/default.png"),
			},
			CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		},
	}
	assert.Equal(t, expectedActivity, actualActivity)

	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Activity.Activity() with incorrect JSON should panic")
			}
		}()
		actual.rawData = []byte(`{"relationships":123}`)
		actual.Activity()
	}()

	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}

func Test_ActivityUserBannedFromProgram(t *testing.T) {
	var actual Activity
	loadResource(t, &actual, "tests/resources/activity-user-banned-from-program.json")
	expected := Activity{
		ID:        String("1337"),
		Type:      String(ActivityUserBannedFromProgramType),
		Message:   String("User Banned From Program!"),
		Internal:  Bool(true),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}

	actualActivity := actual.Activity().(*ActivityUserBannedFromProgram)
	expectedActivity := &ActivityUserBannedFromProgram{
		RemovedUser: &User{
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
	}
	assert.Equal(t, expectedActivity, actualActivity)

	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Activity.Activity() with incorrect JSON should panic")
			}
		}()
		actual.rawData = []byte(`{"relationships":123}`)
		actual.Activity()
	}()

	actual.rawData = nil
	actual.RawActor = nil
	assert.Equal(t, expected, actual)
}
