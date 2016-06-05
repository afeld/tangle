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
	return
}

func startPage() (page models.Page) {
	u := startURL()
	page = models.Page{URL: u}
	return
}

func main() {
	page := startPage()
	links, err := page.GetLinks()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(links)
}
