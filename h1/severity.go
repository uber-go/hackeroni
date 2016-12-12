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

// SeverityRating represent possible severity ratings
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#severity
const (
	SeverityRatingNone              string = "none"
	SeverityRatingLow               string = "low"
	SeverityRatingMedium            string = "medium"
	SeverityRatingHigh              string = "high"
	SeverityAuthorTypeUser          string = "User"
	SeverityAuthorTypeTeam          string = "Team"
	SeverityAttackVectorNetwork     string = "network"
	SeverityAttackVectorAdjacent    string = "adjacent"
	SeverityAttackVectorLocal       string = "local"
	SeverityAttackVectorPhysical    string = "physical"
	SeverityAttackComplexityLow     string = "low"
	SeverityAttackComplexityHigh    string = "high"
	SeverityPrivilegesRequiredLow   string = "low"
	SeverityPrivilegesRequiredHigh  string = "high"
	SeverityUserInteractionNone     string = "none"
	SeverityUserInteractionRequired string = "required"
	SeverityScopeUnchanged          string = "unchanged"
	SeverityScopeChanged            string = "changed"
	SeverityConfidentialityLow      string = "low"
	SeverityConfidentialityHigh     string = "high"
	SeverityIntegrityLow            string = "low"
	SeverityIntegrityHigh           string = "high"
	SeverityAvailabilityLow         string = "low"
	SeverityAvailabilityHigh        string = "high"
)

// Severity represents a severity object
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#severity
type Severity struct {
	ID                 *string    `json:"id"`
	Type               *string    `json:"type"`
	Rating             *string    `json:"rating"`
	AuthorType         *string    `json:"author_type"`
	UserID             *int       `json:"user_id"` // TODO: This is inconsistant with the rest of the API, maybe auto-cast to string
	Score              *float64   `json:"score,omitempty"`
	AttackVector       *string    `json:"attack_vector,omitempty"`
	AttackComplexity   *string    `json:"attack_complexity,omitempty"`
	PrivilegesRequired *string    `json:"privileges_required,omitempty"`
	UserInteraction    *string    `json:"user_interaction,omitempty"`
	Scope              *string    `json:"scope,omitempty"`
	Confidentiality    *string    `json:"confidentiality,omitempty"`
	Integrity          *string    `json:"integrity,omitempty"`
	Availability       *string    `json:"availability,omitempty"`
	CreatedAt          *Timestamp `json:"created_at"`
}

// Helper types for JSONUnmarshal
type severity Severity // Used to avoid recursion of JSONUnmarshal
type severityUnmarshalHelper struct {
	severity
	Attributes *severity `json:"attributes"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (s *Severity) UnmarshalJSON(b []byte) error {
	var helper severityUnmarshalHelper
	helper.Attributes = &helper.severity
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*s = Severity(helper.severity)
	return nil
}
