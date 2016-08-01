hackeroni [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov]
======
A Go interface around [api.hackerone.com](https://api.hackerone.com/).

# Usage
```go
import "github.com/uber-go/hackeroni/h1"
```

To list all reports matching a filter:
```go
reports, _, err := client.Report.List(h1.ReportListOptions{
	Program: []string{"uber"},
})
if err != nil {
	panic(err)
}
for _, report := range reports {
	fmt.Println("Report Title:", *report.Title)
}
```

To retrieve a specific report:
```go
report, _, err := client.Report.Get("123456")
if err != nil {
	panic(err)
}
fmt.Println("Report Title:", *report.Title)
```

## Authentication
The `h1` library does not directly handle authentication. Instead, when creating a new client, you can pass a http.Client that handles authentication for you. It does provide a `APIAuthTransport` structure when using API Token authentication. It is used like this:
```go
tp := h1.APIAuthTransport{
	APIIdentifier: "your-h1-api-token-identifier",
	APIToken: "big-long-api-token-from-h1",
}

client := h1.NewClient(tp.Client())
```

## Pagination
All requests for listing resources such as `Report` support pagination. Pagination options are described in the h1.ListOptions struct and passed to the list methods as an optional parameter. Pages information is available via the h1.ResponseLinks struct embedded in the h1.Response struct.
```go
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
```

[doc-img]: https://godoc.org/github.com/uber-go/hackeroni/h1?status.svg
[doc]: https://godoc.org/github.com/uber-go/hackeroni/h1
[ci-img]: https://travis-ci.org/uber-go/hackeroni.svg?branch=master
[ci]: https://travis-ci.org/uber-go/hackeroni
[cov-img]: https://coveralls.io/repos/github/uber-go/hackeroni/badge.svg?branch=master
[cov]: https://coveralls.io/github/uber-go/hackeroni?branch=master
