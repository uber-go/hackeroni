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

// Attachment represents an attachment (typically to a report or comment).
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#attachment
type Attachment struct {
	ID          *string    `json:"id"`
	Type        *string    `json:"type"`
	FileName    *string    `json:"file_name"`
	ContentType *string    `json:"content_type"`
	FileSize    *int       `json:"file_size"`
	ExpiringURL *string    `json:"expiring_url"`
	CreatedAt   *Timestamp `json:"created_at"`
}

// Helper types for JSONUnmarshal
type attachment Attachment // Used to avoid recursion of JSONUnmarshal
type attachmentUnmarshalHelper struct {
	*attachment
	Attributes *attachment `json:"attributes"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *Attachment) UnmarshalJSON(b []byte) error {
	result := attachment{}
	var helper attachmentUnmarshalHelper
	helper.attachment = &result
	helper.Attributes = &result
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*a = Attachment(result)
	return nil
}
