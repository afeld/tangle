package scanner_test

import (
	. "github.com/afeld/tangle/scanner"

	"net/http"
	"net/url"

	. "github.com/afeld/tangle/helpers"
	. "github.com/afeld/tangle/models"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func registerResponse(dest string, status int) {
	responder := httpmock.NewStringResponder(status, "")
	httpmock.RegisterResponder("HEAD", dest, responder)
}

func registerLink(dest string, status int) Link {
	registerResponse(dest, status)
	return CreateLink(dest)
}

var _ = Describe("Scanner", func() {
	Describe("ScanLinks", func() {
		It("returns the number of broken links", func() {
			dest1 := registerLink("http://ok.com", 200)
			dest2 := registerLink("http://not-ok.com", 404)
			links := []Link{dest1, dest2}

			resultByLink := ScanLinks(links)

			Expect(len(resultByLink)).To(Equal(2))
			Expect(resultByLink[dest1]).To(BeTrue())
			Expect(resultByLink[dest2]).To(BeFalse())
		})

		It("only checks each URL once", func() {
			hits := 0
			responder := func(req *http.Request) (*http.Response, error) {
				hits++
				return httpmock.NewStringResponse(200, ""), nil
			}
			httpmock.RegisterResponder("HEAD", "http://ok.com", responder)

			// two identical links that are different tags
			dest1 := CreateLink("http://ok.com")
			dest2 := CreateLink("http://ok.com")
			Expect(dest1).ToNot(Equal(dest2))
			links := []Link{dest1, dest2}

			resultByLink := ScanLinks(links)

			Expect(len(resultByLink)).To(Equal(2))
			Expect(resultByLink[dest1]).To(BeTrue())
			Expect(resultByLink[dest2]).To(BeTrue())
			Expect(hits).To(Equal(1))
		})
	})

	Describe("ScanPage", func() {
		It("returns the result for each link", func() {
			registerResponse("http://ok.com", 200)
			registerResponse("http://not-ok.com", 404)

			sourceStr := "http://source.com"
			source, _ := url.Parse(sourceStr)
			responder := httpmock.NewStringResponder(200, `
				<a href="http://ok.com"></a>
				<a href="http://not-ok.com"></a>
				<a href="http://not-ok.com"></a>
			`)
			httpmock.RegisterResponder("GET", sourceStr, responder)

			resultByLink, err := ScanPage(source, false)

			Expect(err).ToNot(HaveOccurred())
			Expect(len(resultByLink)).To(Equal(3))
			for link, result := range resultByLink {
				dest, _ := link.DestURL()
				if dest.Host == "ok.com" {
					Expect(result).To(BeTrue())
				} else {
					Expect(result).To(BeFalse())
				}
			}
		})

		It("ignores external links, when specified", func() {
			sourceStr := "http://source.com"
			source, _ := url.Parse(sourceStr)
			responder := httpmock.NewStringResponder(200, `
				<a href="http://external.com"></a>
			`)
			httpmock.RegisterResponder("GET", sourceStr, responder)

			resultByLink, err := ScanPage(source, true)

			Expect(err).ToNot(HaveOccurred())
			Expect(len(resultByLink)).To(Equal(0))
		})
	})
})
