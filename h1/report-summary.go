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

// ReportSummaryCategory represent possible categorys for a report summary
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#report-summary
const (
	ReportSummaryCategoryResearcher string = "researcher"
	ReportSummaryCategoryTeam       string = "team"
)

// ReportSummary represents a summary of a report.
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#report-summary
type ReportSummary struct {
	ID        *string    `json:"id"`
	Type      *string    `json:"type"`
	Content   *string    `json:"content"`
	Category  *string    `json:"category"`
	CreatedAt *Timestamp `json:"created_at"`
	UpdatedAt *Timestamp `json:"updated_at"`
	User      *User      `json:"user"`
}

// Helper types for JSONUnmarshal
type reportSummary ReportSummary // Used to avoid recursion of JSONUnmarshal
type reportSummaryUnmarshalHelper struct {
	reportSummary
	Attributes    *reportSummary `json:"attributes"`
	Relationships struct {
		User struct {
			Data *User `json:"data"`
		} `json:"user"`
	} `json:"relationships"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (r *ReportSummary) UnmarshalJSON(b []byte) error {
	var helper reportSummaryUnmarshalHelper
	helper.Attributes = &helper.reportSummary
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*r = ReportSummary(helper.reportSummary)
	r.User = helper.Relationships.User.Data
	return nil
}
