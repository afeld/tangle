package models

import (
	"net/url"

	// using fork because of https://github.com/moovweb/gokogiri/pull/93#issuecomment-215582446
	"github.com/jbowtie/gokogiri/xml"
)

type Link struct {
	SourceURL url.URL
	Node      xml.Node
}

// may be relative
func (l *Link) DestURL() (*url.URL, error) {
	link := l.Node.Attr("href")
	return url.Parse(link)
}

func (l *Link) AbsDestURL() (URL *url.URL, err error) {
	relativeURL, err := l.DestURL()
	if err != nil {
		return
	}
	URL = l.SourceURL.ResolveReference(relativeURL)
	return
}
