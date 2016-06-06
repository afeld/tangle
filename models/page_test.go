package models_test

import (
	"net/url"

	. "github.com/afeld/tangle/models"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Page", func() {
	Describe("GetLinks", func() {
		It("returns the links from the page", func() {
			responder := httpmock.NewStringResponder(200, `<a href="foo"></a>`)
			httpmock.RegisterResponder("GET", "http://example.com", responder)

			dest, _ := url.Parse("http://example.com")
			page := Page{AbsURL: dest}
			links, err := page.GetLinks()

			Expect(err).ToNot(HaveOccurred())
			Expect(links).To(HaveLen(1))

			actual, err := links[0].AbsDestURL()
			Expect(err).ToNot(HaveOccurred())
			Expect(actual.Host).To(Equal(dest.Host))
			Expect(actual.Path).To(Equal("/foo"))
		})

		It("handles no links being present", func() {
			responder := httpmock.NewStringResponder(200, "")
			httpmock.RegisterResponder("GET", "http://example.com", responder)

			u, _ := url.Parse("http://example.com")
			page := Page{AbsURL: u}
			links, err := page.GetLinks()

			Expect(err).ToNot(HaveOccurred())
			Expect(links).To(HaveLen(0))
		})
	})
})
