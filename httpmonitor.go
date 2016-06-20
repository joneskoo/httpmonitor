package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joneskoo/httpmonitor/fetcher"
	"github.com/joneskoo/httpmonitor/web"
)

func usage() {
	log.Fatal("Usage: httpmonitor <CONFIG>")
}

var listenAddress = "127.0.0.1:3131"

// Get list of URIs from command line and time how long
// it takes to retrieve them all concurrently
func main() {
	var targets []fetcher.Request
	if len(os.Args) != 2 {
		usage()
	}

	// Read targets file
	file, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("Failed to read config: ", err)
	}

	// Parse targets file as JSON
	err = json.Unmarshal(file, &targets)
	if err != nil {
		log.Fatal("Failed to parse JSON in config: ", err)
	}

	// Print target configuration
	log.Print("Monitor targets:")
	for _, target := range targets {
		log.Print(" ", target)
	}

	// Start result fetching and get channel
	resultChannel := fetcher.FetchUrls(targets)

	// Start HTTP server for checking current status
	webChannel := make(chan fetcher.Result)
	web.StartListening(listenAddress, webChannel)

	// Process stream of results
	for {
		res := <-resultChannel
		// Write plain text console log
		msg := fmt.Sprintf("%v  %v %v in %v\n",
			res.StatusEmoji(), res.HTTPStatus, res.URL, res.Dur)
		log.Print(msg)

		// Update status to status web server
		webChannel <- res
	}
}
