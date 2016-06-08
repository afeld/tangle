package scanner

import (
	"net/url"
	"sync"

	"github.com/afeld/tangle/models"
)

// scans the links in parallel
func ScanLinks(links []models.Link) (resultByLink map[models.Link]bool) {
	resultByLink = make(map[models.Link]bool)
	mx := &sync.Mutex{}

	var wg sync.WaitGroup
	for _, link := range links {
		wg.Add(1)
		go func(l models.Link) {
			defer wg.Done()

			isValid := l.IsValid()
			mx.Lock()
			resultByLink[l] = isValid
			mx.Unlock()
		}(link)
	}
	wg.Wait()

	return
}

func ScanPage(source *url.URL) (resultByLink map[models.Link]bool, err error) {
	page := models.Page{AbsURL: source}
	links, err := page.GetLinks()
	if err != nil {
		return
	}

	resultByLink = ScanLinks(links)
	return
}
