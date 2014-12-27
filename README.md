httpmock.go [![Build Status](https://travis-ci.org/theypsilon/httpmock.svg?branch=master)](https://travis-ci.org/theypsilon/httpmock)
======

Simple utility, for mocking http connections. Work 90% based on this gist: https://gist.github.com/jarcoal/8940980

```go
    // in a gopkg.in/check.v1 test method:
    RegisterResponder("GET", "http://example.com/", func(*http.Request) (*http.Response, error) {
        return GetResponseWithBody("xin chao"), nil
    })

    response, _ := http.Get("http://example.com/")

    c.Assert(BodyToString(response), checker.Equals, "xin chao")
```
