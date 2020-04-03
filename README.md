# os-go-netloc-rule

[![Build Status](https://www.travis-ci.org/cfhamlet/os-go-netloc-rule.svg?branch=master)](https://www.travis-ci.org/cfhamlet/os-go-netloc-rule)
[![codecov](https://codecov.io/gh/cfhamlet/os-go-netloc-rule/branch/master/graph/badge.svg)](https://codecov.io/gh/cfhamlet/os-go-netloc-rule)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/cfhamlet/os-go-netloc-rule?tab=overview)

A common library for netloc rule use case.

[Python Version](https://github.com/cfhamlet/os-netloc-rule)

## Install

You can get the library with ``go get``

```
go get -u github.com/cfhamlet/os-go-netloc-rule
```

## Usage

```
package main

import (
	"fmt"

	"github.com/cfhamlet/os-go-netloc-rule/matcher"
	"github.com/cfhamlet/os-go-netloc-rule/netloc"
)

func main() {
	matcher := matcher.New()
	matcher.Load(netloc.New(".google.com", "80", "http"), 1)
	fmt.Println(matcher.MatchURL("http://www.google.com:80/test/"))
}
```

## License
  MIT licensed.
