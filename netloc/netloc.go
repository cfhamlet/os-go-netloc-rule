package netloc

import (
	"strings"
)

// Netloc TODO
type Netloc struct {
	host   string
	port   string
	scheme string
}

// Nil TODO
var Nil = New("", "", "")

// New TODO
func New(host, port, scheme string) *Netloc {
	return &Netloc{host, port, scheme}
}

func (netloc *Netloc) String() string {
	return strings.Join([]string{netloc.host, netloc.port, netloc.scheme}, "|")
}

// Equal TODO
func (netloc *Netloc) Equal(o *Netloc) bool {
	if o == netloc {
		return true
	}
	return netloc.host == o.host &&
		netloc.port == o.port &&
		netloc.scheme == o.scheme
}

// Host TODO
func (netloc *Netloc) Host() string {
	return netloc.host
}

// Port TODO
func (netloc *Netloc) Port() string {
	return netloc.port
}

// Scheme TODO
func (netloc *Netloc) Scheme() string {
	return netloc.scheme
}
