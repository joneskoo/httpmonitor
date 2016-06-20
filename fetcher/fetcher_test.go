package fetcher

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joneskoo/httpmonitor/checker"
)

const successBodyContent = "Hello, client"

// HTTP 200 ok with simple text body
func successHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, successBodyContent)
}

// HTTP 200 ok with simple text body
func forbiddenHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Forbidden", http.StatusForbidden)
}

func TestFetchChecks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(successHandler))
	defer ts.Close()

	// Test cases for HTTP 200 OK with simple text response
	cases := []struct {
		in   []checker.Check
		want bool
	}{
		// No checks (default) should pass
		{[]checker.Check{}, true},
		// Body check should find strings
		{[]checker.Check{{BodyContains: "Hello"}}, true},
		{[]checker.Check{{BodyRegEx: ".ello"}}, true},
		{[]checker.Check{{BodyRegEx: "H.{3}o"}}, true},
		{[]checker.Check{{BodyRegEx: ".allo"}}, false},
		{[]checker.Check{{BodyContains: "client"}}, true},
		{[]checker.Check{{BodyContains: "Client"}}, false},
		// Check status code check
		{[]checker.Check{{StatusCode: 200}}, true},
		{[]checker.Check{{StatusCode: 201}}, false},
		// Many checks together
		{[]checker.Check{{BodyContains: "Hello", StatusCode: 200}}, true},
		{[]checker.Check{{BodyContains: "hello", StatusCode: 200}}, false},
		{[]checker.Check{{BodyContains: "Hello", StatusCode: 201}}, false},
		{[]checker.Check{{BodyContains: "hello", StatusCode: 201}}, false},
		{[]checker.Check{{BodyContains: "Hello", StatusCode: 200, BodyRegEx: "H.{3}o"}}, true},
		{[]checker.Check{{BodyContains: "hello", StatusCode: 200, BodyRegEx: "H.{3}o"}}, false},
		{[]checker.Check{{BodyContains: "Hello", StatusCode: 200, BodyRegEx: "H.{4}o"}}, false},
	}
	for _, c := range cases {
		res := FetchSingleURL(Request{
			URL:      ts.URL,
			Timeout:  0.1,
			Interval: 0.001,
			Checks:   c.in, // checks from test case
		})
		got := res.Passed
		if got != c.want {
			t.Errorf("Check: %#v Expected pass=%v, got pass=%v. Server returned body: %#v", c.in, c.want, got, successBodyContent)
		}
	}

	// Test cases for error status
	ts.Config.Handler = http.HandlerFunc(forbiddenHandler)
	cases = []struct {
		in   []checker.Check
		want bool
	}{
		{[]checker.Check{}, false},
		{[]checker.Check{{BodyContains: "Hello"}}, false},
		{[]checker.Check{{BodyContains: "Forbidden"}}, false},
		{[]checker.Check{{StatusCode: 200}}, false},
		{[]checker.Check{{StatusCode: 403}}, true},
		{[]checker.Check{{StatusCode: 403, BodyContains: "Forbidden"}}, true},
	}
	for _, c := range cases {
		res := FetchSingleURL(Request{
			URL:      ts.URL,
			Timeout:  0.1,
			Interval: 0.001,
			Checks:   c.in, // checks from test case
		})
		got := res.Passed
		if got != c.want {
			t.Errorf("Check: %#v Expected pass=%v, got pass=%v. Server returned error Forbidden", c.in, c.want, got)
		}
	}

}
