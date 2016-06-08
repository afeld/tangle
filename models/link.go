package models

import (
	"net/http"
	"net/url"

	// using fork because of https://github.com/moovweb/gokogiri/pull/93#issuecomment-215582446
	"github.com/jbowtie/gokogiri/xml"

	"github.com/PuerkitoBio/purell"
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
	URL, err = l.DestURL()
	if err != nil {
		return
	}
	// make it absolute
	URL = l.SourceURL.ResolveReference(URL)
	// normalize
	urlStr := purell.NormalizeURL(URL, purell.FlagsSafe)
	return url.Parse(urlStr)
}

func (l *Link) IsValid() bool {
	dest, err := l.AbsDestURL()
	if err != nil {
		return false
	}
	return IsValidURL(dest.String())
}

func IsValidURL(absUrl string) bool {
	resp, err := http.Head(absUrl)
	if err != nil {
		return false
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return true
	}
	return false
}
