package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

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

	numBrokenLinks := 0
	for _, link := range links {
		if !link.IsValid() {
			numBrokenLinks++
			dest, _ := link.DestURL()
			fmt.Printf("%s line %d has broken link to %s.\n", u.String(), link.Node.LineNumber(), dest)
		}
	}
	fmt.Printf("Number of broken links: %d\n", numBrokenLinks)
}
