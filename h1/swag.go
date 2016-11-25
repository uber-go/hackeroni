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

// Swag represents swag that has/hasn't been sent to an address.
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#swag
type Swag struct {
	ID        *string    `json:"id"`
	Type      *string    `json:"type"`
	Sent      *bool      `json:"sent"`
	CreatedAt *Timestamp `json:"created_at"`
	Address   *Address   `json:"address,omitempty"`
}

// Helper types for JSONUnmarshal
type swag Swag // Used to avoid recursion of JSONUnmarshal
type swagUnmarshalHelper struct {
	swag
	Attributes    *swag `json:"attributes"`
	Relationships struct {
		Address struct {
			Data *Address `json:"data"`
		} `json:"address,omitempty"`
	} `json:"relationships"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (s *Swag) UnmarshalJSON(b []byte) error {
	var helper swagUnmarshalHelper
	helper.Attributes = &helper.swag
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*s = Swag(helper.swag)
	s.Address = helper.Relationships.Address.Data
	return nil
}
