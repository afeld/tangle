package models

import (
	"io/ioutil"
	"net/http"
	"net/url"

	// using fork because of https://github.com/moovweb/gokogiri/pull/93#issuecomment-215582446
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/html"
	"github.com/jbowtie/gokogiri/xml"
)

type Page struct {
	AbsURL *url.URL
}

func (p *Page) getDoc() (doc *html.HtmlDocument, err error) {
	resp, err := http.Get(p.AbsURL.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	doc, err = gokogiri.ParseHtml(body)
	return
}

func (p *Page) getAnchors() (anchors []xml.Node, err error) {
	doc, err := p.getDoc()
	if err != nil {
		return
	}

	return doc.Search("//a[@href]")
}

func (p *Page) GetLinks() (links []Link, err error) {
	anchors, err := p.getAnchors()
	if err != nil {
		return
	}
	links = make([]Link, len(anchors))
	for i, anchor := range anchors {
		links[i] = Link{
			SourceURL: *p.AbsURL,
			Node: anchor,
		}
	}
	return
}
