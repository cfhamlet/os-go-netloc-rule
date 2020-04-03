package main

import (
	"fmt"

	"github.com/cfhamlet/os-go-netloc-rule/matcher"
	"github.com/cfhamlet/os-go-netloc-rule/netloc"
)

func main() {
	matcher := matcher.New()
	matcher.Load(netloc.New("www.google.com", "80", "http"), 1)
	fmt.Println(matcher.MatchURL("http://www.google.com:80/test/"))
}
