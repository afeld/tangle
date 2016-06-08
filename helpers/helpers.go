package helpers

import (
	"fmt"

	// using fork because of https://github.com/moovweb/gokogiri/pull/93#issuecomment-215582446
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"

	. "github.com/onsi/gomega"
)

func CreateAnchor(URL string) xml.Node {
	markup := fmt.Sprintf(`<a href="%s"></a>`, URL)
	doc, err := gokogiri.ParseHtml([]byte(markup))
	Expect(err).NotTo(HaveOccurred())

	nodes, err := doc.Search("/html/body/*")
	Expect(err).NotTo(HaveOccurred())
	Expect(nodes).To(HaveLen(1))

	return nodes[0]
}
