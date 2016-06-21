package fetcher

import (
	"fmt"
	"net/http"
	"time"
)

// Target configuration for what to check
type Target struct {
	URL      string  // URL to fetch
	Timeout  float32 // Request timeout, seconds
	Interval float32 // Poll interval, seconds
	Checks   []Check // List of checks for status 'pass'
	Error    error
}

// TimeoutDuration returns the request timeout as time.Duration
func (r Target) TimeoutDuration() time.Duration {
	return time.Duration(r.Timeout * 1e9)
}

// PollIntervalDuration returns the poll interval as time.Duration
func (r Target) PollIntervalDuration() time.Duration {
	return time.Duration(r.Interval * 1e9)
}

func (r Target) String() string {
	return fmt.Sprintf("GET '%v' every %v timeout=%v (%v checks)>",
		r.URL, r.PollIntervalDuration(), r.TimeoutDuration(), len(r.Checks))
}

// Result from a HTTP status check
type Result struct {
	URL        string        // URL we fetched
	Dur        time.Duration // Duration it took to fetch it
	Passed     bool          // Status check pass (true)/fail (false)
	Error      error         // URL fetching error
	HTTPStatus int           // Response HTTP status code
}

func (r Result) String() string {
	return fmt.Sprintf("%v %v in %v", r.URL, r.StatusText(), r.Dur)
}

// StatusText is the pass/fail/unreachable status for check
func (r Result) StatusText() (status string) {
	if r.Error != nil {
		status = "unreachable"
	} else if r.Passed {
		status = "passed"
	} else {
		status = "failed"
	}
	return
}

// StatusEmoji is the pass/fail/unreachable as an emoji symbol
func (r Result) StatusEmoji() string {
	if r.Error != nil {
		return "💤"
	} else if r.Passed {
		return "✅"
	} else {
		return "❌"
	}

}

// FetchSingleURL retrieves a single URL based on configuration structure
// Target and returns a response structure Result.
func FetchSingleURL(req Target) (res Result) {
	res = Result{URL: req.URL}
	// Time how long it takes
	requestStartTime := time.Now()

	// Configure timeout
	client := http.Client{
		Timeout: req.TimeoutDuration(),
	}

	// Perform HTTP GET request
	resp, err := client.Get(req.URL)
	if err != nil {
		res.Dur = time.Since(requestStartTime)
		res.Error = err // Store to result
		return
	}
	res.HTTPStatus = resp.StatusCode
	defer resp.Body.Close() // Close body to free connection after done

	// Run checks
	pass, err := DoCheck(resp, req.Checks)
	res.Error = err
	res.Passed = pass
	res.Dur = time.Since(requestStartTime)
	return
}

// FetchUrls fetches a list of URLs all concurrently in goroutines
// and immediately return a channel streaming Result objects.
func FetchUrls(targets []Target) chan Result {
	c := make(chan Result)
	for _, t := range targets {
		// Start fetch in background and put the result
		// into channel when it is done.
		go func(t Target) {
			for range time.Tick(t.PollIntervalDuration()) {
				res := FetchSingleURL(t)
				c <- res
			}
		}(t)
	}
	// Immediately return the channel where results will arrive
	return c
}
