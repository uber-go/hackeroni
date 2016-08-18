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

package legacy

import (
	"time"
)

// Timestamp represents a time generated from a JSON string
type Timestamp struct {
	time.Time
}

// NewTimestamp creates a new Timestamp object from a ISO8601 date string
func NewTimestamp(date string) *Timestamp {
	ts, err := time.Parse(time.RFC3339, date)
	if err != nil {
		panic(err)
	}
	return &Timestamp{
		Time: ts,
	}
}

// String calls time.Time's String method
func (t Timestamp) String() string {
	return t.Time.String()
}

// UnmarshalJSON helps unmarshal ISO8601 dates in JSON
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var err error
	(*t).Time, err = time.Parse(`"`+time.RFC3339+`"`, string(data))
	return err
}
