package netloc

import (
	"strings"
)

// Netloc TODO
type Netloc struct {
	Host   string `json:"host"`
	Port   string `json:"port"`
	Scheme string `json:"scheme"`
}

// New TODO
func New(host, port, scheme string) *Netloc {
	return &Netloc{host, port, scheme}
}

func (netloc *Netloc) String() string {
	return strings.Join([]string{netloc.Host, netloc.Port, netloc.Scheme}, "|")
}

// Equal TODO
func (netloc *Netloc) Equal(o *Netloc) bool {
	if o == netloc {
		return true
	}
	return netloc.Host == o.Host &&
		netloc.Port == o.Port &&
		netloc.Scheme == o.Scheme
}
