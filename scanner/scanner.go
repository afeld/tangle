package scanner

import (
	"net/url"
	"sync"

	"github.com/afeld/tangle/models"
)

type Options struct {
	DisableExternal bool
}

type resultSet struct {
	ResultByURL map[url.URL]bool
	Mx sync.Mutex
	Wg sync.WaitGroup
}

// not thread-safe
func (r *resultSet) resultFor(link models.Link) bool {
	dest, _ := link.AbsDestURL()
	return r.ResultByURL[*dest]
}

func newResultSet() resultSet {
	return resultSet{
		ResultByURL: make(map[url.URL]bool),
	}
}

func scanURL(u url.URL, results *resultSet) {
	defer results.Wg.Done()

	isValid := models.IsValidURL(u.String())
	results.Mx.Lock()
	results.ResultByURL[u] = isValid
	results.Mx.Unlock()
}

// Scans the URLs in parallel. Takes a Map as input so that the URLs are unique (the values on input are ignored). The Map is modified with the actual values.
func scanURLs(results *resultSet) {
	results.Wg.Add(len(results.ResultByURL))
	for ur, _ := range results.ResultByURL {
		go scanURL(ur, results)
	}
	results.Wg.Wait()

	return
}

// scans the links in parallel
func ScanLinks(links []models.Link) (resultByLink map[models.Link]bool) {
	results := newResultSet()
	for _, link := range links {
		dest, _ := link.AbsDestURL()
		// the value is arbitrary
		results.ResultByURL[*dest] = false
	}

	scanURLs(&results)

	resultByLink = make(map[models.Link]bool)
	for _, link := range links {
		resultByLink[link] = results.resultFor(link)
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

func ScanPage(source *url.URL, options Options) (resultByLink map[models.Link]bool, err error) {
	page := models.Page{AbsURL: source}
	links, err := page.GetLinks()
	if err != nil {
		return
	}

	var filteredLinks []models.Link
	if options.DisableExternal {
		filteredLinks = internalLinks(links)
	} else {
		filteredLinks = links
	}

	resultByLink = ScanLinks(filteredLinks)
	return
}
