package checker

import "net/http"

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
			// if string(resp.StatusCode) != ck.Value {
			// 	pass = false
			// }
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
