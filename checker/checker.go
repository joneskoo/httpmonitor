package checker

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

// BodySizeLimit is the number of bytes of body to read
// when checking rules for body content
const BodySizeLimit = 64 * 1024

// Check defines what is required for test to pass
type Check struct {
	StatusCode   int    // HTTP status code must be...
	BodyContains string // Response body must contain
}

// DoCheck processes all configured checks for a HTTP response.
// It returns a boolean passed if all tests pass
func DoCheck(resp *http.Response, checks []Check) (passed bool, err error) {
	for _, ck := range checks {
		// Check body contents
		if ck.BodyContains != "" {
			var bodyBytes []byte
			lr := io.LimitReader(resp.Body, BodySizeLimit)
			bodyBytes, err = ioutil.ReadAll(lr)
			if err != nil {
				return false, err
			}
			if !bytes.Contains(bodyBytes, []byte(ck.BodyContains)) {
				return false, nil
			}
		}

		// Check status code
		if ck.StatusCode > 0 {
			if resp.StatusCode != ck.StatusCode {
				return false, nil
			}
		} else {
			// HTTP error fails status check unless accepted by check
			if resp.StatusCode >= 400 {
				return false, nil
			}
		}
	}
	return true, nil
}
