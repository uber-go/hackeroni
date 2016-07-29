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

package polling

import (
	"github.com/uber-go/hackeroni/h1"

	"github.com/robmccoll/mitlru"

	"time"
)

// pollingCache is used to track a given
type pollingCache struct {
	Client             *h1.Client           // The h1.Client to use when making requests
	Filter               h1.ReportListFilter // The h1.ReportListOptions to use when making requests
	Window             time.Duration        // How long to look back, recommended 2*Interval
	ReportLastActivity map[string]time.Time // The last activity we know about on that report
	ActivitySeen       *mitlru.TTLRUCache   // If we've seen this particular activity id or not
	ErrorChan          chan error
	ReportChan         chan *h1.Report
	ActivityChan       chan h1.Activity
}

// Start begins polling for events. It returns an error, report and activity channel which emit their respective objects when they occur
func Start(client *h1.Client, filter h1.ReportListFilter, interval time.Duration, window time.Duration) (chan error, chan *h1.Report, chan h1.Activity) {
	// Create a cache
	cache := pollingCache{
		Client:             client,
		Filter:               filter,
		Window:             window,
		ReportLastActivity: make(map[string]time.Time),
		ActivitySeen:       mitlru.NewTTLRUCache(100000, interval + window), // Expire known activities after the interval+window
		ErrorChan:          make(chan error),
		ReportChan:         make(chan *h1.Report),
		ActivityChan:       make(chan h1.Activity),
	}

	// Start polling at the interval in the background
	go func(pollInterval time.Duration) {
		cache.update()
		// Loop the provided interval
		for range time.Tick(pollInterval) {
			cache.update()
		}
	}(interval)

	// Return the channels
	return cache.ErrorChan, cache.ReportChan, cache.ActivityChan
}

// Perform a poll
func (c *pollingCache) update() {
	// We want all reports updated since now minus the window
	updatedAt := time.Now().UTC().Add(-c.Window)

	// Loop all pages to get the reports
	var allReports []h1.Report
	var listOptions h1.ListOptions
	filter := c.Filter
	filter.LastActivityAtGreaterThan = updatedAt
	for {
		reports, resp, err := c.Client.Report.List(filter, &listOptions)
		if err != nil {
			c.ErrorChan <- err
			return
		}
		allReports = append(allReports, reports...)
		if resp.Links.Next == "" {
			break
		}
		listOptions.Page = resp.Links.NextPageNumber()
	}

	// Loop each updated report
	for _, report := range allReports {
		// Get the time we last saw that report
		lastActivityAt, seen := c.ReportLastActivity[*report.ID]
		// If we've seen it and the last activity updated time is equal, skip it
		if seen && lastActivityAt.Equal(report.LastActivityAt.Time) {
			continue
		}
		c.ReportLastActivity[*report.ID] = report.LastActivityAt.Time

		// In order to check the activities we have to pull the full report
		report, _, err := c.Client.Report.Get(*report.ID)
		if err != nil {
			c.ErrorChan <- err
			continue
		}

		// If we hadn't seen the report before, emit the event
		if !seen && report.CreatedAt.After(updatedAt) {
			c.ReportChan <- report
		}

		// Loop all activity in the report
		for _, activity := range report.Activities {
			// If the activity was last updated before the time we updated at, ignore it
			if activity.UpdatedAt.Time.Before(updatedAt) {
				continue
			}

			// If we have seen the activity before, ignore it
			_, seen := c.ActivitySeen.Get(*activity.ID)
			if seen {
				continue
			}

			// Emit the activity
			c.ActivityChan <- activity
		}
	}
}
