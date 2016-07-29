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

// AddressTShirtSize represent possible T-Shirt sizes for an address
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#address
const (
	AddressTShirtSizeSmall   string = "Small"
	AddressTShirtSizeMedium  string = "Medium"
	AddressTShirtSizeLarge   string = "Large"
	AddressTShirtSizeXLarge  string = "X-Large"
	AddressTShirtSizeXXLarge string = "XX-Large"
)

// Address represents an address for a user.
//
// HackerOne API docs: https://api.hackerone.com/docs/v1#address
type Address struct {
	ID          *string    `json:"id"`
	Type        *string    `json:"type"`
	Name        *string    `json:"name"`
	Street      *string    `json:"street"`
	City        *string    `json:"city"`
	PostalCode  *string    `json:"postal_code"`
	State       *string    `json:"state"`
	Country     *string    `json:"country"`
	TShirtSize  *string    `json:"tshirt_size,omitempty"`
	PhoneNumber *string    `json:"phone_number,omitempty"`
	CreatedAt   *Timestamp `json:"created_at"`
}

// Helper types for JSONUnmarshal
type address Address // Used to avoid recursion of JSONUnmarshal
type addressUnmarshalHelper struct {
	*address
	Attributes *address `json:"attributes"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (a *Address) UnmarshalJSON(b []byte) error {
	result := address{}
	var helper addressUnmarshalHelper
	helper.address = &result
	helper.Attributes = &result
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*a = Address(result)
	return nil
}
