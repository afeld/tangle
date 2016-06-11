package scanner

import (
	"net/url"
	"sync"

	"github.com/afeld/tangle/models"
)

// Scans the URLs in parallel. Takes a Map as input so that the URLs are unique (the values on input are ignored). The Map is modified with the actual values.
func scanURLs(resultByURL *map[url.URL]bool) {
	mx := &sync.Mutex{}

	var wg sync.WaitGroup
	for ur, _ := range *resultByURL {
		wg.Add(1)
		go func(u url.URL) {
			defer wg.Done()

			isValid := models.IsValidURL(u.String())
			mx.Lock()
			(*resultByURL)[u] = isValid
			mx.Unlock()
		}(ur)
	}
	wg.Wait()

	return
}

// scans the links in parallel
func ScanLinks(links []models.Link) (resultByLink map[models.Link]bool) {
	resultByURL := make(map[url.URL]bool)
	for _, link := range links {
		dest, _ := link.AbsDestURL()
		// the value is arbitrary
		resultByURL[*dest] = false
	}

	scanURLs(&resultByURL)

	resultByLink = make(map[models.Link]bool)
	for _, link := range links {
		dest, _ := link.AbsDestURL()
		resultByLink[link] = resultByURL[*dest]
	}

	return
}

func internalLinks(links []models.Link) []models.Link {
	filteredLinks := make([]models.Link, 0, len(links))
	for _, link := range links {
		isExternal, _ := link.IsExternal()
		if !isExternal {
			filteredLinks = append(filteredLinks, link)
		}
	}
	return filteredLinks
}

func ScanPage(source *url.URL, disableExternal bool) (resultByLink map[models.Link]bool, err error) {
	page := models.Page{AbsURL: source}
	links, err := page.GetLinks()
	if err != nil {
		return
	}

	var filteredLinks []models.Link
	if disableExternal {
		filteredLinks = internalLinks(links)
	} else {
		filteredLinks = links
	}

	resultByLink = ScanLinks(filteredLinks)
	return
}
