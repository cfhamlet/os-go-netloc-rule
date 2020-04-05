package matcher_test

import (
	"strings"
	"testing"

	"github.com/cfhamlet/os-go-netloc-rule/matcher"
	"github.com/cfhamlet/os-go-netloc-rule/netloc"
	"github.com/stretchr/testify/assert"
)

func fromString(k string) *netloc.Netloc {
	s := strings.Split(k, "|")
	return netloc.New(s[0], s[1], s[2])
}

func createMatcher(data map[string]interface{}) *matcher.Matcher {
	matcher := matcher.New()
	for k, v := range data {
		matcher.Load(fromString(k), v)
	}
	return matcher
}

func Test001(t *testing.T) {
	matcher := createMatcher(
		map[string]interface{}{
			"||":                 0,
			"google.com||":       1,
			"www.google.com||":   2,
			"xxx.google.com||":   3,
			".google.com|77|":    4,
			".google.com||ftp":   5,
			"||ftp":              6,
			"|99|":               7,
			"b.google.com|88|":   8,
			"a.b.google.com|88|": 9,
		},
	)

	tests := []struct {
		url      string
		expected int
	}{
		{"http://001.google.com.cn/", 0},
		{"http://001.google.com/", 1},
		{"http://www.google.com/", 2},
		{"http://001.xxx.google.com/", 3},
		{"http://001.xxx.google.com:77/", 4},
		{"ftp://001.xxx.google.com/", 5},
		{"ftp://001.xxx.google.com.hk/", 6},
		{"ftp://001.xxx.google.com.hk:99/", 7},
		{"ftp://1.b.google.com:88/", 8},
		{"ftp://a.b.google.com:88/", 9},
	}
	for _, test := range tests {
		_, v := matcher.MatchURL(test.url)
		assert.Equal(t, test.expected, v)
	}
}

func Test002(t *testing.T) {
	data := map[string]interface{}{
		"google.com||":     1,
		"www.google.com||": 2,
		"xxx.google.com||": 3,
		".google.com|77|":  4,
		".google.com||ftp": 5,
		"||ftp":            6,
		"|99|":             7,
	}
	matcher := createMatcher(data)
	out := map[string]interface{}{}
	matcher.Iter(
		func(nl *netloc.Netloc, rule interface{}) bool {
			out[nl.String()] = rule
			return true
		},
	)
	assert.Equal(t, data, out)
}

func Test003(t *testing.T) {
	data := map[string]interface{}{
		"google.com||":        1,
		"www.google.com||":    2,
		"xxx.google.com||":    3,
		".google.com|77|":     4,
		".google.com|77|http": 4,
		".google.com|77|ftp":  4,
		".google.com||ftp":    5,
		"||ftp":               6,
		"|99|":                7,
	}
	matcher := createMatcher(data)
	assert.Equal(t, len(data), matcher.Size())
	for k := range data {
		n := fromString(k)
		matcher.Delete(n)
	}
	assert.Equal(t, 0, matcher.Size())
}

func Test004(t *testing.T) {
	data := map[string]interface{}{
		"google.com||":     1,
		"www.google.com||": 2,
		"xxx.google.com||": 3,
	}
	matcher := createMatcher(data)
	assert.Equal(t, 3, matcher.Size())
	data = map[string]interface{}{
		"google.com||":     4,
		"www.google.com||": 5,
		"xxx.google.com||": 6,
	}
	for k, v := range data {
		n := fromString(k)
		matcher.Load(n, v)
	}
	assert.Equal(t, 3, matcher.Size())
	_, v := matcher.MatchURL("http://xxx.google.com/")
	assert.Equal(t, 6, v)
	matcher.LoadWithCmp(fromString("xxx.google.com||"), 1,
		func(old, new interface{}) bool {
			a := old.(int)
			b := new.(int)
			return a < b
		},
	)
	_, v = matcher.MatchURL("http://xxx.google.com/")
	assert.Equal(t, 1, v)
}

func Test006(t *testing.T) {
	matcher := createMatcher(
		map[string]interface{}{
			"www.google.com|80|ftp":  1,
			"www.google.com|80|http": 2,
			"google.com|80|http":     3,
		},
	)

	tests := []struct {
		url      string
		expected int
	}{
		{"http://www.google.com:80/", 2},
		{"ftp://www.google.com:80/", 1},
		{"ftp://abc.www.google.com:80/", 1},
	}
	for _, test := range tests {
		_, v := matcher.MatchURL(test.url)
		assert.Equal(t, test.expected, v)
	}
}
func Test007(t *testing.T) {
	data := map[string]interface{}{
		"google.com||":     1,
		"www.google.com||": 2,
		"xxx.google.com||": 3,
		".google.com|77|":  4,
		".google.com||ftp": 5,
		"||ftp":            6,
		"|99|":             7,
	}
	matcher := createMatcher(data)
	for _, s := range []string{
		"abc.google.com||",
		".google.com|88|http",
		".google.com|77|http",
		".google.com|77|ftp",
	} {
		assert.Equal(t, 7, matcher.Size())
		sp := strings.Split(s, "|")
		r1, r2 := matcher.Delete(netloc.New(sp[0], sp[1], sp[2]))
		var r *netloc.Netloc
		assert.Equal(t, r, r1)
		assert.Equal(t, nil, r2)
	}
}

func BenchmarkMatch(b *testing.B) {
	data := map[string]interface{}{
		"google.com||":     1,
		"www.google.com||": 2,
		"xxx.google.com||": 3,
		".google.com|77|":  4,
		".google.com||ftp": 5,
		"||ftp":            6,
		"|99|":             7,
	}
	matcher := createMatcher(data)
	for i := 0; i < b.N; i++ {
		matcher.Match("a.b.c.d.google.com", "", "")
	}
}

func BenchmarkParseURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = netloc.ParseURL("http://www.google.com:80/a/b/c/?k=1")
	}
}
