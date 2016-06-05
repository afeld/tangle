package models_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/afeld/tangle/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Page", func() {
	Describe("GetLinks", func() {
		It("returns the links from the page", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, `<a href="foo"></a>`)
			}))
			defer ts.Close()

			tsURL, _ := url.Parse(ts.URL)
			page := Page{AbsURL: tsURL}
			links, err := page.GetLinks()

			Expect(err).ToNot(HaveOccurred())
			Expect(links).To(HaveLen(1))

			dest, err := links[0].AbsDestURL()
			Expect(err).ToNot(HaveOccurred())
			Expect(dest.Host).To(Equal(tsURL.Host))
			Expect(dest.Path).To(Equal("/foo"))
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
