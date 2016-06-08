package reporter

import (
	"fmt"

	"github.com/afeld/tangle/models"
)

func reportBrokenLink(link models.Link) {
	source := link.SourceURL.String()
	line := link.Node.LineNumber()
	dest, _ := link.DestURL()
	fmt.Printf("%s line %d has broken link to %s.\n", source, line, dest)
}

func ReportResults(resultByLink map[models.Link]bool) {
	numBrokenLinks := 0
	for link, isValid := range resultByLink {
		if !isValid {
			numBrokenLinks++
			reportBrokenLink(link)
		}
	}

	fmt.Printf("Number of broken links: %d/%d\n", numBrokenLinks, len(resultByLink))
}
