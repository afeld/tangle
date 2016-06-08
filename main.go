package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

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

	fmt.Printf("Number of links found: %d\n", len(resultByLink))

	numBrokenLinks := 0
	for link, isValid := range resultByLink {
		if !isValid {
			numBrokenLinks++

			source := link.SourceURL.String()
			line := link.Node.LineNumber()
			dest, _ := link.DestURL()
			fmt.Printf("%s line %d has broken link to %s.\n", source, line, dest)
		}
	}

	fmt.Printf("Number of broken links: %d\n", numBrokenLinks)
}
