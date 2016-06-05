package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/afeld/tangle/models"
)

func main() {
	rawURL := os.Args[1]
	if len(os.Args) != 2 {
		log.Fatal("Usage:\n\n\tgo run main.go <url>\n\n")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalln(err)
	}
	page := models.Page{URL: u}

	links, err := page.GetLinks()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(links)
}
