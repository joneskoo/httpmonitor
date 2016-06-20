package checker

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

const kB = 1024

// BodySizeLimit is the number of bytes of body to read
// when checking rules for body content
const BodySizeLimit = 64 * kB

// Check defines what is required for test to pass
type Check struct {
	StatusCode   int    // HTTP status code must be...
	BodyContains string // Response body must contain
	BodyRegEx    string // Regular expression to match body
}

// DoCheck processes all configured checks for a HTTP response.
// It returns passed=true if all tests pass
func DoCheck(resp *http.Response, checks []Check) (passed bool, err error) {
	statusChecked := false

	// Read body bytes to memory
	var bodyBytes []byte
	lr := io.LimitReader(resp.Body, BodySizeLimit)
	bodyBytes, err = ioutil.ReadAll(lr)
	if err != nil {
		return false, err
	}

	// Go through all checks. If any check fails, return immediately.
	for _, ck := range checks {
		// Check body contents
		if ck.BodyContains != "" {
			if !bytes.Contains(bodyBytes, []byte(ck.BodyContains)) {
				return false, nil
			}
		}

		if ck.BodyRegEx != "" {
			var matched bool
			matched, err = regexp.Match(ck.BodyRegEx, bodyBytes)
			if err != nil {
				return false, err
			}
			if !matched {
				return false, fmt.Errorf("Body did not match %v", ck.BodyRegEx)
			}
		}

		// Check status code
		if ck.StatusCode > 0 {
			got := resp.StatusCode
			if got != ck.StatusCode {
				return false, fmt.Errorf("Expected HTTP status %v, got %v", ck.StatusCode, got)
			}
			statusChecked = true
		}
	}

	// If there was no check for status yet, check for error status.
	// Otherwise respect the status check earlier.
	if !statusChecked {
		// Implicit default status check for HTTP error responses
		if resp.StatusCode >= 400 {
			return false, fmt.Errorf("HTTP Error status %v", resp.StatusCode)
		}
	}

	// No problems detected, return passed=true
	return true, err
}
