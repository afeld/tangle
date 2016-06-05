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
	AbsURL *url.URL
}

func (p *Page) getDoc() (doc *html.HtmlDocument, err error) {
	resp, err := http.Get(p.AbsURL.String())
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

func (p *Page) GetLinks() (links []*url.URL, err error) {
	doc, err := p.getDoc()
	if err != nil {
		return
	}

	anchors, err := doc.Search("//a")
	if err != nil {
		return
	}
	links = make([]*url.URL, 0, len(anchors))
	for _, anchor := range anchors {
		link := anchor.Attr("href")
		otherRelativeURL, _ := url.Parse(link)
		if err != nil {
			fmt.Println("Bad URL:", link)
			continue
		}
		otherAbsURL := p.AbsURL.ResolveReference(otherRelativeURL)
		links = append(links, otherAbsURL)
	}
	return
}
