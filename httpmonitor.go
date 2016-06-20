package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joneskoo/httpmonitor/fetcher"
	"github.com/joneskoo/httpmonitor/web"
)

var binaryName = "httpmonitor"

func myUsage() {
	fmt.Printf("Usage: %s [OPTIONS] config.json\n", binaryName)
	flag.PrintDefaults()
}

// HTTP server listen address
var listenAddress string

func init() {
	const (
		listenUsage   = "bind address for HTTP server (e.g. '127.0.0.1:8000', default disabled)"
		listenDefault = ""
	)
	flag.StringVar(&listenAddress, "bind", listenDefault, listenUsage)
}

// Get list of URIs from command line and time how long
// it takes to retrieve them all concurrently
func main() {
	flag.Usage = myUsage
	flag.Parse()
	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Read targets file
	file, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal("Failed to read config: ", err)
	}

	// Parse targets file as JSON
	var targets []fetcher.Request
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
	if listenAddress != "" {
		web.StartListening(listenAddress, webChannel)
	}

	// Process stream of results
	for {
		res := <-resultChannel
		// Write plain text console log
		msg := fmt.Sprintf("%v  %v %v in %v\n",
			res.StatusEmoji(), res.HTTPStatus, res.URL, res.Dur)
		log.Print(msg)

		// Update status to status web server
		if listenAddress != "" {
			webChannel <- res
		}
	}
}
