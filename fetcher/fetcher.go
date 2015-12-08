package fetcher

import (
	"fmt"
	"httpmonitor/checker"
	"log"
	"net/http"
	"time"
)

// Request configuration for what to check
type Request struct {
	URL      string          // URL we fetch
	Timeout  time.Duration   // How long before we timeout and abort a request
	Interval time.Duration   // At what rate we send requests (for loop)
	Checks   []checker.Check // List of checks for status 'pass'
	Error    error
}

func (r Request) String() string {
	return fmt.Sprintf("<GET '%s' every %s timeout=%s (%d checks)>",
		r.URL, r.Interval, r.Timeout, len(r.Checks))
}

// Result from a HTTP status check
type Result struct {
	URL    string        // URL we fetched
	Dur    time.Duration // Duration it took to fetch it
	Status bool          // Status check pass (true)/fail (false)
	Error  error
}

// StatusText is the pass/fail/unreachable status for check
func (r Result) StatusText() (status string) {
	if r.Error != nil {
		status = "unreachable"
	} else if r.Status {
		status = "pass"
	} else {
		status = "fail"
	}
	return
}

// FetchSingleURL retrieves a single URL based on configuration structure
// Request and returns a response structure Result.
func FetchSingleURL(req Request) (res Result) {
	res = Result{URL: req.URL}
	// Time how long it takes
	requestStartTime := time.Now()

	// Configure timeout
	client := http.Client{
		Timeout: req.Timeout,
	}

	// Perform HTTP GET request
	resp, err := client.Get(req.URL)
	if err != nil {
		log.Print("Request failed: ", err)
		res.Error = err // Store to result
		return
	}
	defer resp.Body.Close() // Close body to free connection after done

	// Run checks
	pass, err := checker.DoCheck(resp, req.Checks)
	if err != nil {
		log.Print("Check failed: ", err)
		return
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
			for range time.Tick(req.Interval) {
				res := FetchSingleURL(req)
				c <- res
			}
		}(req)
	}
	// Immediately return the channel where results will arrive
	return c
}
