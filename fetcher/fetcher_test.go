package fetcher

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
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
		in   []Check
		want bool
	}{
		// No checks (default) should pass
		{[]Check{}, true},
		// Body check should find strings
		{[]Check{{BodyContains: "Hello"}}, true},
		{[]Check{{BodyRegEx: ".ello"}}, true},
		{[]Check{{BodyRegEx: "H.{3}o"}}, true},
		{[]Check{{BodyRegEx: ".allo"}}, false},
		{[]Check{{BodyContains: "client"}}, true},
		{[]Check{{BodyContains: "Client"}}, false},
		// Check status code check
		{[]Check{{StatusCode: 200}}, true},
		{[]Check{{StatusCode: 201}}, false},
		// Many checks together
		{[]Check{{BodyContains: "Hello", StatusCode: 200}}, true},
		{[]Check{{BodyContains: "hello", StatusCode: 200}}, false},
		{[]Check{{BodyContains: "Hello", StatusCode: 201}}, false},
		{[]Check{{BodyContains: "hello", StatusCode: 201}}, false},
		{[]Check{{BodyContains: "Hello", StatusCode: 200, BodyRegEx: "H.{3}o"}}, true},
		{[]Check{{BodyContains: "hello", StatusCode: 200, BodyRegEx: "H.{3}o"}}, false},
		{[]Check{{BodyContains: "Hello", StatusCode: 200, BodyRegEx: "H.{4}o"}}, false},
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
		in   []Check
		want bool
	}{
		{[]Check{}, false},
		{[]Check{{BodyContains: "Hello"}}, false},
		{[]Check{{BodyContains: "Forbidden"}}, false},
		{[]Check{{StatusCode: 200}}, false},
		{[]Check{{StatusCode: 403}}, true},
		{[]Check{{StatusCode: 403, BodyContains: "Forbidden"}}, true},
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

func BenchmarkCheckNothing(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(successHandler))
	defer ts.Close()

	for i := 0; i < b.N; i++ {
		FetchSingleURL(Request{
			URL:      ts.URL,
			Timeout:  0.1,
			Interval: 0.001,
			Checks:   []Check{{}},
		})

	}
}

func BenchmarkCheckStatus(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(successHandler))
	defer ts.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FetchSingleURL(Request{
			URL:      ts.URL,
			Timeout:  0.1,
			Interval: 0.001,
			Checks:   []Check{{StatusCode: 200}},
		})

	}
}

func BenchmarkCheckBodyContains(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(successHandler))
	defer ts.Close()

	for i := 0; i < b.N; i++ {
		FetchSingleURL(Request{
			URL:      ts.URL,
			Timeout:  0.1,
			Interval: 0.001,
			Checks:   []Check{{BodyContains: "Hello"}},
		})

	}
}

func BenchmarkCheckStatusAndBody(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(successHandler))
	defer ts.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FetchSingleURL(Request{
			URL:      ts.URL,
			Timeout:  0.1,
			Interval: 0.001,
			Checks: []Check{
				{
					StatusCode:   200,
					BodyContains: "Hello",
				}},
		})

	}
}

func BenchmarkCheckBodyRegEx(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(successHandler))
	defer ts.Close()

	for i := 0; i < b.N; i++ {
		FetchSingleURL(Request{
			URL:      ts.URL,
			Timeout:  0.1,
			Interval: 0.001,
			Checks:   []Check{{BodyRegEx: "Hell."}},
		})

	}
}
