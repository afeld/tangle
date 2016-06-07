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

func main() {
	u := startURL()
	page := models.Page{AbsURL: u}

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
		go func(l models.Link) {
			defer wg.Done()

			if !l.IsValid() {
				// https://gobyexample.com/atomic-counters
				atomic.AddUint32(&numBrokenLinks, 1)

				dest, _ := l.DestURL()
				fmt.Printf("%s line %d has broken link to %s.\n", u.String(), l.Node.LineNumber(), dest)
			}
		}(link)
	}
	wg.Wait()

	fmt.Printf("Number of broken links: %d\n", numBrokenLinks)
}
