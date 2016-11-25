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

func Test_Swag(t *testing.T) {
	var actual Swag
	loadResource(t, &actual, "tests/resources/swag.json")
	expected := Swag{
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
	}
	assert.Equal(t, expected, actual)
}
