package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"sync/atomic"

	"github.com/afeld/tangle/models"
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

func checkLink(source *url.URL, link models.Link, wg *sync.WaitGroup, numBrokenLinks *uint32) {
	defer wg.Done()

	if !link.IsValid() {
		// https://gobyexample.com/atomic-counters
		atomic.AddUint32(numBrokenLinks, 1)

		dest, _ := link.DestURL()
		fmt.Printf("%s line %d has broken link to %s.\n", source.String(), link.Node.LineNumber(), dest)
	}
}

func main() {
	source := startURL()
	page := models.Page{AbsURL: source}

	fmt.Println("Checking for broken links...")

	links, err := page.GetLinks()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Number of links found: %d\n", len(links))

	var numBrokenLinks uint32 = 0
	var wg sync.WaitGroup

	for _, link := range links {
		wg.Add(1)
		go checkLink(source, link, &wg, &numBrokenLinks)
	}
	wg.Wait()

	fmt.Printf("Number of broken links: %d\n", numBrokenLinks)
}
