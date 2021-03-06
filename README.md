# Tangle [![Build Status](https://travis-ci.org/afeld/tangle.svg?branch=master)](https://travis-ci.org/afeld/tangle)

A broken link checker for web sites.

## Usage

Requires [Go 1.6+](https://golang.org).

```bash
# install
go get github.com/afeld/tangle
cd $GOPATH/src/github.com/afeld/tangle

# run
go run main.go <url>

# optionally, disable external link checking
go run main.go -disable-external <url>
```

## See also

* [404finder](https://fourohfourtracker.herokuapp.com/)
* [HTML Proofer](https://github.com/gjtorikian/html-proofer)
* [LinkChecker](https://wummel.github.io/linkchecker/)
* [W3C Link Checker](https://validator.w3.org/checklink)
