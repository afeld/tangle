package scanner

import (
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/afeld/tangle/models"
)

func checkLink(link models.Link, wg *sync.WaitGroup, numBrokenLinks *uint32) {
	defer wg.Done()

	if !link.IsValid() {
		// https://gobyexample.com/atomic-counters
		atomic.AddUint32(numBrokenLinks, 1)

		dest, _ := link.DestURL()
		fmt.Printf("%s line %d has broken link to %s.\n", link.SourceURL.String(), link.Node.LineNumber(), dest)
	}
}

// scans the links in parallel
func ScanLinks(links []models.Link) (numBroken uint32) {
	var wg sync.WaitGroup
	for _, link := range links {
		wg.Add(1)
		go checkLink(link, &wg, &numBroken)
	}
	wg.Wait()

	return
}

func ScanPage(source *url.URL) (numBrokenLinks uint32, err error) {
	page := models.Page{AbsURL: source}
	links, err := page.GetLinks()
	if err != nil {
		return
	}

	fmt.Printf("Number of links found: %d\n", len(links))
	numBrokenLinks = ScanLinks(links)

	return
}
