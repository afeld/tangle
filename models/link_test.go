package models_test

import (
	"net/http"
	"net/url"

	. "github.com/afeld/tangle/helpers"
	. "github.com/afeld/tangle/models"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Link", func() {
	Describe("AbsDestURL", func() {
		It("returns the full URL if already absolute", func() {
			sourceURL, _ := url.Parse("http://example.com/")
			link := Link{
				SourceURL: *sourceURL,
				Node:      CreateAnchor("http://example2.com/"),
			}

			dest, err := link.AbsDestURL()
			Expect(err).NotTo(HaveOccurred())
			Expect(dest.String()).To(Equal("http://example2.com/"))
		})

		It("resolves relative URLs", func() {
			sourceURL, _ := url.Parse("http://example.com/")
			link := Link{
				SourceURL: *sourceURL,
				Node:      CreateAnchor("foo"),
			}

			dest, err := link.AbsDestURL()
			Expect(err).NotTo(HaveOccurred())
			Expect(dest.String()).To(Equal("http://example.com/foo"))
		})
	})

	Describe("IsValid", func() {
		It("returns `true` when the URL exists", func() {
			responder := httpmock.NewStringResponder(200, "")
			httpmock.RegisterResponder("HEAD", "http://example2.com", responder)

			sourceURL, _ := url.Parse("http://example.com/")

			link := Link{
				SourceURL: *sourceURL,
				Node:      CreateAnchor("http://example2.com"),
			}
			Expect(link.IsValid()).To(BeTrue())
		})

		It("returns `false` when the URL 404s", func() {
			responder := httpmock.NewStringResponder(404, "")
			httpmock.RegisterResponder("HEAD", "http://example2.com", responder)

			sourceURL, _ := url.Parse("http://example.com/")

			link := Link{
				SourceURL: *sourceURL,
				Node:      CreateAnchor("http://example2.com"),
			}
			Expect(link.IsValid()).To(BeFalse())
		})

		It("returns `false` for a connection failure", func() {
			responder := func(req *http.Request) (*http.Response, error) {
				return httpmock.ConnectionFailure(req)
			}
			httpmock.RegisterResponder("HEAD", "http://example2.com", responder)

			sourceURL, _ := url.Parse("http://example.com/")

			link := Link{
				SourceURL: *sourceURL,
				Node:      CreateAnchor("http://example2.com"),
			}
			Expect(link.IsValid()).To(BeFalse())
		})
	})
})
