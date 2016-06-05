package models_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/afeld/tangle/models"
	// using fork because of https://github.com/moovweb/gokogiri/pull/93#issuecomment-215582446
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func createAnchor(URL string) xml.Node {
	markup := fmt.Sprintf(`<a href="%s"></a>`, URL)
	doc, err := gokogiri.ParseHtml([]byte(markup))
	Expect(err).NotTo(HaveOccurred())

	nodes, err := doc.Search("/html/body/*")
	Expect(err).NotTo(HaveOccurred())

	return nodes[0]
}

var _ = Describe("Link", func() {
	Describe("AbsDestURL", func() {
		It("returns the full URL if already absolute", func() {
			sourceURL, _ := url.Parse("http://example.com/")
			link := Link{
				SourceURL: *sourceURL,
				Node:      createAnchor("http://example2.com/"),
			}

			dest, err := link.AbsDestURL()
			Expect(err).NotTo(HaveOccurred())
			Expect(dest.String()).To(Equal("http://example2.com/"))
		})

		It("resolves relative URLs", func() {
			sourceURL, _ := url.Parse("http://example.com/")
			link := Link{
				SourceURL: *sourceURL,
				Node:      createAnchor("foo"),
			}

			dest, err := link.AbsDestURL()
			Expect(err).NotTo(HaveOccurred())
			Expect(dest.String()).To(Equal("http://example.com/foo"))
		})
	})

	Describe("IsValid", func() {
		It("returns `true` when the URL exists", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()

			sourceURL, _ := url.Parse("http://example.com/")

			link := Link{
				SourceURL: *sourceURL,
				Node:      createAnchor(ts.URL),
			}
			Expect(link.IsValid()).To(BeTrue())
		})

		It("returns `false` when the URL 404s", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}))
			defer ts.Close()

			sourceURL, _ := url.Parse("http://example.com/")

			link := Link{
				SourceURL: *sourceURL,
				Node:      createAnchor(ts.URL),
			}
			Expect(link.IsValid()).To(BeFalse())
		})
	})
})
