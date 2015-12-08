# HTTP Monitor #

Copyright 2015 Joonas Kuorilehto. License: MIT-style.

This program that monitors web servers for availability.
It's intended for web server administrators to use to
monitor their service.

Example configuration:

```json
{
  "Version": 1,
  "Monitor": [
    {
      "URL": "http://localhost:8000/",
      "Timeout": 200,
      "Interval": 200,
      "Checks": [
        {
          "Type": "contains",
          "Value": "Directory listing for"
        },
        {
          "Type": "status",
          "Value": "200"
        }
      ]
    }
  ],
  "Log": "checks.csv",
  "HTTP": "127.0.0.1:3131"
}
```

Timeout and Interval are specified in units of nanoseconds.

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

## Design goals ##

Since the number of monitored servers depends on the application,
and to make this utility reusable in various contexts
and environments, this application was implemented in
[Go](https://golang.org/). The Go language is exceptionally
well suited for concurrent tasks.
