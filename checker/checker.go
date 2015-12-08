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
	Type  string
	Value interface{}
}

// DoCheck processes all configured checks for a HTTP response.
// It returns a boolean pass if all tests pass
func DoCheck(resp *http.Response, checks []Check) (pass bool, err error) {
	// Defaults
	err = nil
	pass = true

	for _, ck := range checks {
		switch ck.Type {
		case "contains":
			lr := io.LimitReader(resp.Body, BodySizeLimit)
			bodyBytes, e := ioutil.ReadAll(lr)
			if e != nil {
				pass = false
				return
			}
			if !bytes.Contains(bodyBytes, []byte(ck.Value.(string))) {
				pass = false
				return
			}
		case "status":
			if resp.StatusCode != int(ck.Value.(float64)) {
				pass = false
			}
		default:
			panic("Unknown check type")
		}

	}
	return
}
