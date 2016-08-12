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

	"testing"
)

func Test_Report(t *testing.T) {
	var actual Report
	loadResource(t, &actual, "tests/resources/report.json")
	expected := Report{
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
	assert.Equal(t, expected, actual)

	var activitiesReport Report
	loadResource(t, &activitiesReport, "tests/resources/report_activities.json")
	assert.Equal(t, activitiesReport.Activities[0].Report(), &activitiesReport)

	participantsActual := activitiesReport.Participants(false)
	participantsExpected := []User{
		User{
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
		User{
			ID:       String("1338"),
			Type:     String(UserType),
			Disabled: Bool(false),
			Username: String("hackeroni-example"),
			Name:     String("Hackeroni Example"),
			ProfilePicture: UserProfilePicture{
				Size62x62:   String("/assets/avatars/default.png"),
				Size82x82:   String("/assets/avatars/default.png"),
				Size110x110: String("/assets/avatars/default.png"),
				Size260x260: String("/assets/avatars/default.png"),
			},
			CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		},
	}
	assert.Equal(t, participantsActual, participantsExpected)

	internalParticipantsActual := activitiesReport.Participants(true)
	internalParticipantsExpected := []User{
		User{
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
		User{
			ID:       String("1339"),
			Type:     String(UserType),
			Disabled: Bool(false),
			Username: String("hackeroni-example2"),
			Name:     String("Hackeroni Example 2"),
			ProfilePicture: UserProfilePicture{
				Size62x62:   String("/assets/avatars/default.png"),
				Size82x82:   String("/assets/avatars/default.png"),
				Size110x110: String("/assets/avatars/default.png"),
				Size260x260: String("/assets/avatars/default.png"),
			},
			CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		},
		User{
			ID:       String("1338"),
			Type:     String(UserType),
			Disabled: Bool(false),
			Username: String("hackeroni-example"),
			Name:     String("Hackeroni Example"),
			ProfilePicture: UserProfilePicture{
				Size62x62:   String("/assets/avatars/default.png"),
				Size82x82:   String("/assets/avatars/default.png"),
				Size110x110: String("/assets/avatars/default.png"),
				Size260x260: String("/assets/avatars/default.png"),
			},
			CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		},
	}
	assert.Equal(t, internalParticipantsActual, internalParticipantsExpected)

	var assigneeUserReport Report
	loadResource(t, &assigneeUserReport, "tests/resources/report_assignee-user.json")
	assigneeUserActual := assigneeUserReport.Assignee().(*User)
	assigneeUserExpected := &User{
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
	}
	assert.Equal(t, assigneeUserExpected, assigneeUserActual)

	var assigneeGroupReport Report
	loadResource(t, &assigneeGroupReport, "tests/resources/report_assignee-group.json")
	assigneeGroupActual := assigneeGroupReport.Assignee().(*Group)
	assigneeGroupExpected := &Group{
		ID:        String("1337"),
		Type:      String(GroupType),
		Name:      String("Admin"),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
	}
	assert.Equal(t, assigneeGroupExpected, assigneeGroupActual)

	var assigneeNilReport Report
	assigneeNilReport.RawAssignee = []byte("{}")
	assigneeNil := assigneeNilReport.Assignee()
	assert.Nil(t, assigneeNil)

	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Report.Assignee() with invalid JSON should panic")
			}
		}()
		var assigneeInvalidReport Report
		assigneeInvalidReport.RawAssignee = []byte(`Invalid JSON`)
		assigneeInvalidReport.Assignee()
	}()

	func() {
		defer func() {
			if recover() == nil {
				assert.Fail(t, "Report.Assignee() with incorrect JSON should panic")
			}
		}()
		var assigneeInvalidReport Report
		assigneeInvalidReport.RawAssignee = []byte(`{"type": "group", "id": 123}`)
		assigneeInvalidReport.Assignee()
	}()
}
