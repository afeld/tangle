package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/afeld/tangle/reporter"
	"github.com/afeld/tangle/scanner"
)

func startURL() (u *url.URL) {
	if flag.NArg() != 1 {
		fmt.Print("Not enough arguments. ")
		showUsage()
		os.Exit(1)
	}

	rawURL := flag.Arg(0)
	u, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalln(err)
	}
	if !u.IsAbs() {
		log.Fatalln("Must be an absolute URL.")
	}
	return
}

func showUsage() {
	fmt.Printf("Usage:\n\n  %s [options] <url>\n\nOptions:\n\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Println("")
}

func main() {
	disableExternal := flag.Bool("disable-external", false, "Disables external link checking.")
	flag.Usage = showUsage
	flag.Parse()

	source := startURL()

	fmt.Println("Checking for broken links...")
	resultByLink, err := scanner.ScanPage(source, *disableExternal)
	if err != nil {
		log.Fatalln(err)
	}
	reporter.ReportResults(resultByLink)
}
