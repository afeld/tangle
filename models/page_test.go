package models_test

import (
	"fmt"
	"net/http"
	"net/url"
	"net/http/httptest"

	. "github.com/afeld/tangle/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Page", func() {
	Describe("GetLinks", func() {
		It("returns the links from the page", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, `<a href="foo">`)
			}))
			defer ts.Close()

			tsURL, _ := url.Parse(ts.URL)
			page := Page{AbsURL: tsURL}
			links, err := page.GetLinks()

			Expect(err).ToNot(HaveOccurred())
			Expect(links).To(HaveLen(1))
			link := links[0]
			Expect(link.Host).To(Equal(tsURL.Host))
			Expect(link.Path).To(Equal("/foo"))
		})

		It("handles no links being present", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "")
			}))
			defer ts.Close()

			u, _ := url.Parse(ts.URL)
			page := Page{AbsURL: u}
			links, err := page.GetLinks()

			Expect(err).ToNot(HaveOccurred())
			Expect(links).To(HaveLen(0))
		})
	})
})
