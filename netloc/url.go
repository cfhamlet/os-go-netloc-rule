package netloc

import (
	"errors"
	"fmt"
	"net/url"
)

// ParsedURL TODO
type ParsedURL struct {
	*url.URL
	Host string
	Port string
}

// ParseURL TODO
func ParseURL(rawURL string) (parsedURL *ParsedURL, err error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return parsedURL, err
	}

	if parsed.Host == "" {
		err = errors.New("empty host")
	} else {
		host := parsed.Hostname()
		port := parsed.Port()
		if host == "" {
			err = fmt.Errorf("empty host %s", parsed.Host)
		} else {
			parsedURL = &ParsedURL{parsed, host, port}
		}
	}
	if err != nil {
		err = &url.Error{Op: "parse", URL: rawURL, Err: err}
	}
	return parsedURL, err
}
