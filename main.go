package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	// using fork because of https://github.com/moovweb/gokogiri/pull/93#issuecomment-215582446
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/html"
)

func getDoc(url string) (doc *html.HtmlDocument, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	fmt.Println("Request complete: ", resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	doc, err = gokogiri.ParseHtml(body)
	return
}

func getLinks(url string) (links []string, err error) {
	doc, err := getDoc(url)
	if err != nil {
		return
	}

	anchors, err := doc.Search("//a")
	if err != nil {
		return
	}
	links = make([]string, len(anchors), len(anchors))
	for i, anchor := range anchors {
		links[i] = anchor.Attr("href")
	}
	return
}

func main() {
	url := os.Args[1]
	if len(os.Args) != 2 {
		log.Fatal("Usage:\n\n\tgo run main.go <url>\n\n")
	}

	links, err := getLinks(url)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(links)
}
