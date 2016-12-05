package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/jawher/mow.cli" // Such a great library, check/star it on github !
)

func main() {
	app := cli.App("cub", "Checks urls from a 'csv' FILE and write to a DEST") // Utilization of mow.cli for creating a CLI app
	app.Version("v version", "0.0.1")

	var (
		verbose = app.BoolOpt("l logs", true, "Print result of requests")
		src = app.StringArg("SRC", "urls.csv", "File with urls")
		dest = app.StringArg("DEST", "urls-results.csv", "Destination file with results")
	)

	app.Action = func() {
		application(*src, *dest, *verbose) // here goes our code concerning checking URLs
	}

	app.Run(os.Args)
}

func application(srcUrl string, destUrl string, verbose bool) {
	// ---------------------------------- initialization
	srcFile, err := os.Open(srcUrl)
	defer srcFile.Close()

	if err != nil {
		log.Fatal(err)
	}

	responses := make(chan string, 10) // declaration of buffered channel with length of 10
	var wg sync.WaitGroup

	// ---------------------------------- execution

	go writeToFile(responses, destUrl, &wg) // magic go keyword of going concurrent

	scanner := bufio.NewScanner(srcFile)
	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}

		go checkUrl(scanner.Text(), verbose, responses, &wg) // magic go keyword of going concurrent
	}

	// ---------------------------------- finalization

	wg.Wait() // Wait for all responses to be written in file before continue
	close(responses) // We must close the channel in order to end the loop on channel responses in writeFile function
}

func writeToFile(responses chan string, destUrl string, wg *sync.WaitGroup) {
	// Takes the WaitGroup and a channel as parameters
	destFile, err := os.Create(destUrl)
	defer destFile.Sync()        // We will sync and
	defer destFile.Close()       // close the file at the end

	if err != nil {
		log.Fatal(err)
	}

	for response := range responses { // Thanks to the magic range keyword we loop on the responses added in channel from other goroutines
		_, err := destFile.WriteString(response) // write to file, nothing special

		wg.Done() // We have wrote the responses in the file we can now close the channel if it is the last response

		if err != nil {
			log.Fatal(err)
		}
	}
}

func checkUrl(url string, verbose bool, responses chan string, wg *sync.WaitGroup) { // takes a channel and WaitGroup as argument
	wg.Add(1) // We wait for this goroutine to end before closing channel
	resp, err := http.Head(url)

	if verbose {
		fmt.Printf("'%s', '%s', '%s',\n", url, resp.Status, resp.Header.Get("Content-Type"))
	}

	switch {
	case err != nil:
		responses <- fmt.Sprintf("%s, %s,\n", url, "KO") // An example of adding a response in the channel which will be read by the writeToFile function
		break
	case resp.StatusCode != 200:
		responses <- fmt.Sprintf("%s, %s,\n", url, "KO")
		break
	case strings.Contains(resp.Header.Get("Content-Type"), "text/html"):
		responses <- fmt.Sprintf("%s, %s,\n", url, "KO")
		break
	default:
		responses <- fmt.Sprintf("%s, %s,\n", url, "OK")
		break
	}
}