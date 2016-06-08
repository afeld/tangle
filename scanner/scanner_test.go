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


func createLink(dest string) Link {
	source, _ := url.Parse("http://source.com")
	node := CreateAnchor(dest)
	return Link{
		SourceURL: *source,
		Node: node,
	}
}

func registerResponse(dest string, status int) Link {
	responder := httpmock.NewStringResponder(status, "")
	httpmock.RegisterResponder("HEAD", dest, responder)

	return createLink(dest)
}

var _ = Describe("Scanner", func() {
	Describe("ScanLinks", func() {
		It("returns the number of broken links", func() {
			dest1 := registerResponse("http://ok.com", 200)
			dest2 := registerResponse("http://not-ok.com", 404)

			links := []Link{dest1, dest2}
			numBrokenLinks := ScanLinks(links)
			Expect(numBrokenLinks).To(Equal(uint32(1)))
		})
	})

	Describe("ScanPage", func() {
		It("returns the number of broken links", func() {
			registerResponse("http://ok.com", 200)
			registerResponse("http://not-ok.com", 404)

			source, _ := url.Parse("http://source.com")
			responder := httpmock.NewStringResponder(200, `
				<a href="http://ok.com"></a>
				<a href="http://not-ok.com"></a>
			`)
			httpmock.RegisterResponder("GET", "http://source.com", responder)

			numBrokenLinks, err := ScanPage(source)
			Expect(err).ToNot(HaveOccurred())
			Expect(numBrokenLinks).To(Equal(uint32(1)))
		})
	})
})
