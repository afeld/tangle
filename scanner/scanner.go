package scanner

import (
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/afeld/tangle/models"
)

func checkLink(source *url.URL, link models.Link, wg *sync.WaitGroup, numBrokenLinks *uint32) {
	defer wg.Done()

	if !link.IsValid() {
		// https://gobyexample.com/atomic-counters
		atomic.AddUint32(numBrokenLinks, 1)

		dest, _ := link.DestURL()
		fmt.Printf("%s line %d has broken link to %s.\n", source.String(), link.Node.LineNumber(), dest)
	}
}

func scanLinks(source *url.URL, links []models.Link) (numBroken uint32) {
	var wg sync.WaitGroup

	for _, link := range links {
		wg.Add(1)
		go checkLink(source, link, &wg, &numBroken)
	}
	wg.Wait()

	return
}

func ScanPage(source *url.URL) (err error) {
	fmt.Println("Checking for broken links...")

	page := models.Page{AbsURL: source}
	links, err := page.GetLinks()
	if err != nil {
		return
	}

	fmt.Printf("Number of links found: %d\n", len(links))
	numBrokenLinks := scanLinks(source, links)
	fmt.Printf("Number of broken links: %d\n", numBrokenLinks)

	return
}
