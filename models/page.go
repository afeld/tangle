package models

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	// using fork because of https://github.com/moovweb/gokogiri/pull/93#issuecomment-215582446
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/html"
)

type Page struct {
	URL *url.URL
}

func (p *Page) getDoc() (doc *html.HtmlDocument, err error) {
	resp, err := http.Get(p.URL.String())
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

func (p *Page) GetLinks() (links []string, err error) {
	doc, err := p.getDoc()
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
