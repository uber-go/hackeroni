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

/*
Package h1 provides a client for the HackerOne API.
Usage:
	import "github.com/uber-go/hackeroni/h1"
Construct a new HackerOne client, then use the various services on the client to
access different parts of the HackerOne API. For example:
	authTransport := h1.APIAuthTransport{
		APIIdentifier: "your-h1-api-token-identifier",
		APIToken: "big-long-api-token-from-h1",
	}
	client := authTransport.NewClient(tp.Client())

	report, _, err := client.Report.Get("123456")
	if err != nil {
		panic(err)
	}

	fmt.Println("Report Title:", *report.Title)

Authentication

The h1 library does not directly handle authentication. Instead, when creating a new client, you can pass a http.Client that handles authentication for you. It does provide a APIAuthTransport structure when using API Token authentication.

Pagination

All requests for listing resources such as `Report` support pagination. Pagination options are described in the h1.ListOptions struct and passed to the list methods as an optional parameter. Pages information is available via the h1.ResponseLinks struct embedded in the h1.Response struct.
	filter := h1.ReportListFilter{
		Program: []string{"uber"},
	}
	var listOpts h1.ListOptions

	var allReports []h1.Report
	for {
		reports, resp, err := client.Report.List(filter, &listOpts)
		if err != nil {
			panic(err)
		}
		allReports = append(allReports, reports...)
		if resp.Links.Next == "" {
			break
		}
		listOpts.Page = resp.Links.NextPageNumber()
	}
*/
package h1
