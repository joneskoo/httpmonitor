package web

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"sync"

	"github.com/joneskoo/httpmonitor/fetcher"
)

// protects statusMap
var statusMapMutex sync.RWMutex

// protected by statusMapMutex as it is read and written concurrently
var statusMap map[string]fetcher.Result

// Status list page contents
const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>HTTP monitor</title>
	</head>
	<body>
                <table>
                <tr>
                        <th>URL</th>
                        <th>Status</th>
                        <th>Response time</th>
                </tr>
                <tr>{{range .}}
                        <td>{{ .URL }}</td>
                        <td>{{ .StatusText }}</td>
                        <td>{{ .ResponseTime }}</td>
                </tr>{{else}}
                <tr>
                        <td><strong>no data</td>
                        <td></td>
                        <td></td>
                </tr>
                {{end}}
                </table>
	</body>
</html>
`

func list(w http.ResponseWriter, r *http.Request) {
	// Reader for statusMap
	statusMapMutex.RLock()
	defer statusMapMutex.RUnlock()
	var keys []string
	for k := range statusMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	type Status struct {
		URL, StatusText, ResponseTime string
	}
	var statuses []Status
	for _, k := range keys {
		s := statusMap[k]
		statuses = append(statuses, Status{s.URL, s.StatusText(), s.Dur.String()})
	}

	t := template.Must(template.New("list").Parse(tpl))
	t.ExecuteTemplate(w, "list", statuses)
	return
}

// StartListening initializes the server that presents the latest data
// from web server status monitoring
func StartListening(addr string, events chan fetcher.Result) {
	statusMap = make(map[string]fetcher.Result)
	log.Print("Listening on http://", addr)
	http.HandleFunc("/", list)
	go http.ListenAndServe(addr, nil)
	go func(events chan fetcher.Result) {
		for {
			// Writer for statusMap
			res := <-events
			statusMapMutex.Lock()
			statusMap[res.URL] = res
			statusMapMutex.Unlock()
		}
	}(events)
}

func main() {
	addr := "[::1]:8000"
	dummy := make(chan fetcher.Result)
	StartListening(addr, dummy)
}
