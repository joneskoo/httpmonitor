# HTTP Monitor #

HTTP monitor (naming is hard) is a tool to monitor for HTTP service
availability. It supports concurrent monitoring of a number of targets
and basic checks of the response for expected data.

Example configuration:

```json
{
  "Version": 1,
  "Monitor": [
    {
      "URL": "http://localhost:8000/",
      "Timeout": 1.0,
      "Interval": 1.0,
      "Checks": [
        {"Type": "contains", "Value": "Directory listing for"},
        {"Type": "status",   "Value": "200"}
      ]
    }
  ],
  "Log": "checks.csv",
  "HTTP": "127.0.0.1:3131"
}
```

Timeout and Interval are specified in seconds.

The HTTP monitor supports the following status checks:

Type     | Value (for check to pass)
---------|-------------------------------------
contains | String that must be in the HTTP body
status   | Acceptable HTTP status code

The log is written in CSV format with header. The log
contains the following columns:

* URL checked
* status (pass/fail/unreachable)
* Response time (in milliseconds)

There is a built in HTTP server for checking the current status.
To enable it, set "HTTP" in configuration to a string "IP:port" to set
which address to bind the server to.

## Design goals ##

Since the number of monitored servers depends on the application,
and to make this utility reusable in various contexts
and environments, this application was implemented in
[Go](https://golang.org/). The Go language is exceptionally
well suited for concurrent tasks.
