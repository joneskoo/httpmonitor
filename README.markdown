# HTTP Monitor #

HTTP monitor (naming is hard) is a tool to monitor for HTTP service
availability. It supports concurrent monitoring of a number of targets
and basic checks of the response for expected data.

![Travis-CI build status](https://travis-ci.org/joneskoo/httpmonitor.svg?branch=master)

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

## Design goals ##

Since the number of monitored servers depends on the application,
and to make this utility reusable in various contexts
and environments, this application was implemented in
[Go](https://golang.org/). The Go language is exceptionally
well suited for concurrent tasks.

[regex]: https://golang.org/pkg/regexp/ "Go package regexp documentation"
