package scanner_test

import (
	. "github.com/afeld/tangle/scanner"

	"net/url"

	. "github.com/afeld/tangle/models"
	. "github.com/afeld/tangle/helpers"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func registerResponse(dest string, status int) {
	responder := httpmock.NewStringResponder(status, "")
	httpmock.RegisterResponder("HEAD", dest, responder)
}

func createLink(dest string) Link {
	source, _ := url.Parse("http://source.com")
	node := CreateAnchor(dest)
	return Link{
		SourceURL: *source,
		Node: node,
	}
}

func registerLink(dest string, status int) Link {
	registerResponse(dest, status)
	return createLink(dest)
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
	})

	Describe("ScanPage", func() {
		It("returns the result for each link", func() {
			registerResponse("http://ok.com", 200)
			registerResponse("http://not-ok.com", 404)

			source, _ := url.Parse("http://source.com")
			responder := httpmock.NewStringResponder(200, `
				<a href="http://ok.com"></a>
				<a href="http://not-ok.com"></a>
				<a href="http://not-ok.com"></a>
			`)
			httpmock.RegisterResponder("GET", "http://source.com", responder)

			resultByLink, err := ScanPage(source)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(resultByLink)).To(Equal(3))
		})
	})
})
