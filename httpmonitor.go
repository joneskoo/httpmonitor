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

// Config file format
type Config struct {
	Version int
	Monitor []fetcher.Request
	Log     string
	HTTP    string
}

func usage() {
	log.Fatal("Usage: httpmonitor <CONFIG>")
}

// Get list of URIs from command line and time how long
// it takes to retrieve them all concurrently
func main() {
	var conf Config
	if len(os.Args) != 2 {
		usage()
	}
	file, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("Failed to read config: ", err)
	}

	err = json.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal("Failed to parse JSON in config: ", err)
	}
	log.Print("Version: ", conf.Version)
	log.Print("Log file: ", conf.Log)
	log.Print("HTTP listen address: ", conf.HTTP)
	log.Print("Monitor targets:")
	for _, target := range conf.Monitor {
		log.Print(" ", target)
	}

	// Start result fetching and get channel
	resultChannel := fetcher.FetchUrls(conf.Monitor)

	// Start HTTP server for checking current status
	webChannel := make(chan fetcher.Result)
	if conf.HTTP != "" {
		web.StartListening(conf.HTTP, webChannel)
	}

	for { // Process stream of results
		res := <-resultChannel
		// Write plain text console log
		msg := fmt.Sprintf("%v  %v in %v\n",
			res.StatusEmoji(), res.URL, res.Dur)
		log.Print(msg)

		// Update status to status web server
		if conf.HTTP != "" {
			webChannel <- res
		}
	}
}
