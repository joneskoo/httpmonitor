# HTTP Monitor #

HTTP monitor (naming is hard) is a tool to monitor for HTTP service
availability. It supports concurrent monitoring of a number of targets
and basic checks of the response for expected data.

[![Travis-CI build status](https://travis-ci.org/joneskoo/httpmonitor.svg?branch=master)](https://travis-ci.org/joneskoo/)
[![codecov](https://codecov.io/gh/joneskoo/httpmonitor/branch/master/graph/badge.svg)](https://codecov.io/gh/joneskoo/httpmonitor)

Example target configuration:

```json
[
  {
    "URL": "http://localhost:8000/",
    "Timeout": 1,
    "Interval": 1,
    "Checks": [
      {
        "BodyContains": "Directory listing for",
        "StatusCode": 200,
        "BodyRegEx": "Hello!?"
      }
    ]
  }
]
```

Timeout and Interval are specified in seconds.

The HTTP monitor supports the following status checks:

Type         | Value (for check to pass)
-------------|-------------------------------------
BodyContains | String that must be in the HTTP body
StatusCode   | Acceptable HTTP status code
BodyRegEx    | [Regular expression][regex] match to body content

There is a built in HTTP server for checking the current status.
To enable it, set "HTTP" in configuration to a string "IP:port" to set
which address to bind the server to.

## Installing and quick start ##

Assuming [Go][go] has been installed, and [GOPATH][gopath] is set up:

    $ go get github.com/joneskoo/httpmonitor
    $ $GOPATH/bin/httpmonitor
    Usage: httpmonitor [OPTIONS] config.json
      -bind string
        	bind address for HTTP server (e.g. '127.0.0.1:8000', default disabled)

Create a taget configuration file like the example above and set up
what to monitor, then launch. Optionally you may enable the HTTP interface.

## Design goals ##

Since the number of monitored servers depends on the application,
and to make this utility reusable in various contexts
and environments, this application was implemented in
[Go][go]. The Go language is exceptionally
well suited for concurrent tasks.

[go]: https://golang.org/ "Go programming language"
[regex]: https://golang.org/pkg/regexp/ "Go package regexp documentation"
[gopath]: https://golang.org/doc/code.html#GOPATH "Setting up GOPATH"
