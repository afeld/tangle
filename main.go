package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/afeld/tangle/reporter"
	"github.com/afeld/tangle/scanner"
)

func startURL() (u *url.URL) {
	rawURL := os.Args[1]
	if len(os.Args) != 2 {
		log.Fatal("Usage:\n\n\tgo run main.go <url>\n\n")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalln(err)
	}
	if !u.IsAbs() {
		log.Fatalln("Must be an absolute URL.")
	}
	return
}

func main() {
	source := startURL()

	fmt.Println("Checking for broken links...")
	resultByLink, err := scanner.ScanPage(source)
	if err != nil {
		log.Fatalln(err)
	}
	reporter.ReportResults(resultByLink)
}
