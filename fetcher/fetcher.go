package fetcher

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joneskoo/httpmonitor/checker"
)

// Request configuration for what to check
type Request struct {
	URL      string          // URL to fetch
	Timeout  float32         // Request timeout, seconds
	Interval float32         // Poll interval, seconds
	Checks   []checker.Check // List of checks for status 'pass'
	Error    error
}

// TimeoutDuration returns the request timeout as time.Duration
func (r Request) TimeoutDuration() time.Duration {
	return time.Duration(r.Timeout * 1e9)
}

// PollIntervalDuration returns the poll interval as time.Duration
func (r Request) PollIntervalDuration() time.Duration {
	return time.Duration(r.Interval * 1e9)
}

func (r Request) String() string {
	return fmt.Sprintf("GET '%v' every %v timeout=%v (%v checks)>",
		r.URL, r.PollIntervalDuration(), r.TimeoutDuration(), len(r.Checks))
}

// Result from a HTTP status check
type Result struct {
	URL    string        // URL we fetched
	Dur    time.Duration // Duration it took to fetch it
	Status bool          // Status check pass (true)/fail (false)
	Error  error
}

func (r Result) String() string {
	return fmt.Sprintf("%v %v in %v", r.URL, r.StatusText(), r.Dur)
}

// StatusText is the pass/fail/unreachable status for check
func (r Result) StatusText() (status string) {
	if r.Error != nil {
		status = "unreachable"
	} else if r.Status {
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
	} else if r.Status {
		return "✅"
	} else {
		return "❌"
	}

}

// FetchSingleURL retrieves a single URL based on configuration structure
// Request and returns a response structure Result.
func FetchSingleURL(req Request) (res Result) {
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
		log.Print("Request failed: ", err)
		res.Dur = time.Since(requestStartTime)
		res.Error = err // Store to result
		return
	}
	defer resp.Body.Close() // Close body to free connection after done

	// Run checks
	pass, err := checker.DoCheck(resp, req.Checks)
	if err != nil {
		log.Print("Check failed: ", err)
	}
	res.Status = pass
	res.Dur = time.Since(requestStartTime)
	return
}

// FetchUrls fetches a list of URLs all concurrently in goroutines
// and immediately return a channel streaming Result objects.
func FetchUrls(requests []Request) chan Result {
	c := make(chan Result)
	for _, req := range requests {
		// Start fetch in background and put the result
		// into channel when it is done.
		go func(req Request) {
			for range time.Tick(req.PollIntervalDuration()) {
				res := FetchSingleURL(req)
				c <- res
			}
		}(req)
	}
	// Immediately return the channel where results will arrive
	return c
}
