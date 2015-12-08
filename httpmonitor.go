package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/joneskoo/httpmonitor/fetcher"
	"github.com/joneskoo/httpmonitor/web"
)

func usage() {
	log.Fatal("Usage: httpmonitor <CONFIG>")
}

// Get list of URIs from command line and time how long
// it takes to retrieve them all concurrently
func main() {
	// Configure
	type Config struct {
		Version int
		Monitor []fetcher.Request
		Log     string
		HTTP    string
	}
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
		log.Print("  ", target)
	}

	// Open CSV and write CSV header
	var w *csv.Writer
	if conf.Log != "" {
		file, err := os.Create(conf.Log)
		if err != nil {
			log.Fatal("Can't create log file: ", err)
		}
		defer file.Close()

		w = csv.NewWriter(file)
		w.Write([]string{
			"timestamp",
			"target URL",
			"response time",
			"status check"})
		w.Flush()
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

		// Write CSV log output
		if conf.Log != "" {
			err := w.Write([]string{
				time.Now().Format(time.RFC3339), // timestamp
				res.URL, // target URL
				fmt.Sprintf("%0.3f", res.Dur.Seconds()), // duration
				res.StatusText()})                       // status check
			w.Flush()
			if err != nil {
				log.Fatal(err)
			}
		}

		// Write plain text console log
		msg := fmt.Sprintf("%s status=%s in %s\n",
			res.URL, res.StatusText(), res.Dur)
		log.Print(msg)

		// Update status to status web server
		if conf.HTTP != "" {
			webChannel <- res
		}
	}
}
