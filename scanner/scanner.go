package scanner

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/afeld/tangle/models"
)

type Options struct {
	DisableExternal bool
	Recursive       bool
}

type linkResult struct {
	Link models.Link
	Found bool
}

// slow method
func scanLink(link models.Link, results chan<- linkResult) {
	// slow method
	valid := link.IsValid()
	result := linkResult{
		Link: link,
		Found: valid,
	}
	results <- result
}

// slow method
func findLinks(page models.Page, newLinks chan<- models.Link, errors chan<- error) {
	// slow method
	links, err := page.GetLinks()
	if err != nil {
		errors <- err
		return
	}

	for _, link := range links {
		newLinks <- link
	}
}

func newPagesLoop(wg *sync.WaitGroup, newPages <-chan models.Page, newLinks chan<- models.Link, errors chan<- error) {
	for page := range newPages {
		wg.Add(1)
		go func(p models.Page) {
			findLinks(p, newLinks, errors)
			wg.Done()
		}(page)
	}
}

func recurse(link models.Link, newPages chan<- models.Page) (err error) {
	dest, err := link.AbsDestURL()
	if err != nil {
		return err
	}
	// TODO make sure this page hasn't been seen before
	page := models.Page{AbsURL: dest}
	newPages <- page

	return nil
}

func handleNewLink(wg *sync.WaitGroup, link models.Link, options Options, results chan<- linkResult, newPages chan<- models.Page) error {
	isExternal, err := link.IsExternal()
	if err != nil {
		return err
	}

	if !options.DisableExternal || !isExternal {
		wg.Add(1)
		go func() {
			scanLink(link, results)
			wg.Done()
		}()

		if options.Recursive {
			return recurse(link, newPages)
		}
	}

	return nil
}

func newLinksLoop(wg *sync.WaitGroup, options Options, newLinks <-chan models.Link, newPages chan<- models.Page, results chan<- linkResult, errors chan<- error) {
	for link := range newLinks {
		err := handleNewLink(wg, link, options, results, newPages)
		if err != nil {
			errors <- err
		}
	}
}

func recursiveScan(startPage models.Page, options Options) (resultByLink map[models.Link]bool) {
	resultByLink = make(map[models.Link]bool)
	var wg sync.WaitGroup

	newPages := make(chan models.Page)
	newLinks := make(chan models.Link)
	newResults := make(chan linkResult)
	errors := make(chan error)

	defer close(newPages)
	defer close(newLinks)
	defer close(newResults)
	defer close(errors)

	go newPagesLoop(&wg, newPages, newLinks, errors)
	go newLinksLoop(&wg, options, newLinks, newPages, newResults, errors)
	go func() {
		for result := range newResults {
			resultByLink[result.Link] = result.Found
		}
	}()
	go func() {
		for err := range errors {
			// TODO do something with these
			fmt.Println(err)
		}
	}()

	// kick off the scan
	newPages <- startPage
	wg.Wait()

	return
}

func ScanPage(source *url.URL, options Options) (resultByLink map[models.Link]bool, err error) {
	page := models.Page{AbsURL: source}
	// TODO return error(s)
	return recursiveScan(page, options), nil
}
