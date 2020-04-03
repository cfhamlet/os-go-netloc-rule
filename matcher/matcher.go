package matcher

import (
	"strings"
	"sync"

	"github.com/cfhamlet/os-go-netloc-rule/netloc"
)

// GT TODO
type GT func(interface{}, interface{}) bool

// Nil TODO
var Nil = newNetlocRule(nil, nil)

type netlocRule struct {
	*netloc.Netloc
	rule interface{}
}

func newNetlocRule(nl *netloc.Netloc, rule interface{}) *netlocRule {
	return &netlocRule{nl, rule}
}

// matchUnit
type matchUnit struct {
	nlcRules map[string][]*netlocRule
}

func (unit *matchUnit) Add(new *netloc.Netloc, rule interface{}, gt GT) (*netloc.Netloc, interface{}) {
	port := new.Port()
	l, ok := unit.nlcRules[port]
	if !ok {
		unit.nlcRules[port] = []*netlocRule{newNetlocRule(new, rule)}
		return nil, nil
	}
	for i := 0; i < len(l); i++ {
		old := l[i]
		if old.Scheme() != new.Scheme() {
			continue
		}
		if gt == nil || !gt(old.rule, rule) {
			l[i] = newNetlocRule(new, rule)
			return old.Netloc, old.rule
		}
		return new, rule
	}
	l = append(l, newNetlocRule(new, rule))
	unit.nlcRules[port] = l
	return nil, nil
}

func newMatchUnit() *matchUnit {
	return &matchUnit{make(map[string][]*netlocRule)}
}

// Matcher TODO
type Matcher struct {
	units map[string]*matchUnit
	size  int
	*sync.RWMutex
}

// New TODO
func New() *Matcher {
	return &Matcher{
		make(map[string]*matchUnit),
		0,
		&sync.RWMutex{},
	}
}

// Size TODO
func (matcher *Matcher) Size() int {
	matcher.RLock()
	defer matcher.RUnlock()
	return matcher.size
}

// MatchHost TODO
func (matcher *Matcher) MatchHost(host string) (*netloc.Netloc, interface{}) {
	return matcher.Match(Empty, host, Empty)
}

// MatchHostPort TODO
func (matcher *Matcher) MatchHostPort(host, port string) (*netloc.Netloc, interface{}) {
	return matcher.Match(Empty, host, port)
}

func betterMatch(n1, n2 *netlocRule, port, scheme string) *netlocRule {
	if n1.Netloc == nil {
		return n2
	} else if n2.Netloc == nil {
		return n1
	}
	if port != Empty {
		if n1.Port() == n2.Port() && n1.Port() == port {
			if len(n1.Host()) > len(n2.Host()) {
				return n1
			}
			return n2
		}
		if port == n2.Port() {
			return n2
		} else if port == n1.Port() {
			return n1
		}
	}
	if scheme != Empty {
		if n1.Scheme() == n2.Scheme() && n1.Scheme() == scheme {
			if len(n1.Host()) > len(n2.Host()) {
				return n1
			}
			return n2
		}
		if scheme == n2.Scheme() {
			return n2
		} else if scheme == n1.Scheme() {
			return n1
		}
	}
	if len(n2.Host()) > len(n1.Host()) {
		return n2
	}

	return n1
}

func (matcher *Matcher) matchPicec(piece, port, scheme string) (*netlocRule, bool) {
	unit, ok := matcher.units[piece]
	var bestMatch *netlocRule = Nil
	if !ok {
		return bestMatch, false
	}
	if port != Empty {
		nlcs, ok := unit.nlcRules[port]
		if ok {
			for _, nlc := range nlcs {
				if scheme != Empty && nlc.Scheme() == scheme {
					return nlc, true
				} else if nlc.Scheme() == scheme || nlc.Scheme() == Empty {
					bestMatch = betterMatch(bestMatch, nlc, port, scheme)
				}
			}
		}
	}
	nlcs, ok := unit.nlcRules[Empty]
	if ok {
		for _, nlc := range nlcs {
			if nlc.Scheme() == scheme || nlc.Scheme() == Empty {
				bestMatch = betterMatch(bestMatch, nlc, port, scheme)
			}
		}
	}
	return bestMatch, false
}

// MatchURL TODO
func (matcher *Matcher) MatchURL(URL string) (*netloc.Netloc, interface{}) {
	parsed, err := netloc.ParseURL(URL)
	if err != nil {
		return nil, nil
	}
	return matcher.Match(parsed.Host, parsed.Port, parsed.Parsed.Scheme)
}

// Match TODO
func (matcher *Matcher) Match(host, port, scheme string) (*netloc.Netloc, interface{}) {
	matcher.RLock()
	defer matcher.RUnlock()
	piece := host
	bestMatch := Nil
	for {
		nlr, exact := matcher.matchPicec(piece, port, scheme)
		if exact {
			bestMatch = nlr
			break
		} else {
			bestMatch = betterMatch(bestMatch, nlr, port, scheme)
		}
		if piece == Empty {
			break
		}
		piece = nextPiece(piece)
	}
	return bestMatch.Netloc, bestMatch.rule
}

// Delete TODO
func (matcher *Matcher) Delete(netloc *netloc.Netloc) (*netloc.Netloc, interface{}) {
	matcher.Lock()
	defer matcher.Unlock()
	host := netloc.Host()
	unit, ok := matcher.units[host]
	if !ok {
		return nil, nil
	}
	port := netloc.Port()
	nlrs, ok := unit.nlcRules[port]
	if !ok {
		return nil, nil
	}
	nlrsLen := len(nlrs)
	var i = 0
	for ; i < nlrsLen; i++ {
		nlr := nlrs[i]
		if nlr.Equal(netloc) {
			break
		}
	}
	if i >= len(nlrs) {
		return nil, nil
	}
	nlr := nlrs[i]

	if i == nlrsLen-1 {
	} else {
		nlrs[i] = nlrs[nlrsLen]
		nlrs[nlrsLen] = nil
	}
	nlrs = nlrs[:nlrsLen-1]
	if len(nlrs) <= 0 {
		delete(unit.nlcRules, port)
		if len(unit.nlcRules) <= 0 {
			delete(matcher.units, host)
		}
	} else {
		unit.nlcRules[port] = nlrs
	}
	matcher.size--
	return nlr.Netloc, nlr.rule
}

// Load TODO
func (matcher *Matcher) Load(netloc *netloc.Netloc, rule interface{}) (*netloc.Netloc, interface{}) {
	return matcher.LoadWithCmp(netloc, rule, nil)
}

// LoadWithCmp TODO
func (matcher *Matcher) LoadWithCmp(netloc *netloc.Netloc, rule interface{}, cmp GT) (*netloc.Netloc, interface{}) {
	matcher.Lock()
	defer matcher.Unlock()
	host := netloc.Host()
	_, ok := matcher.units[host]
	if !ok {
		matcher.units[host] = newMatchUnit()
	}
	n, v := matcher.units[host].Add(netloc, rule, cmp)
	if n == nil {
		matcher.size++
	}
	return n, v
}

// IterFunc TODO
type IterFunc func(*netloc.Netloc, interface{}) bool

// Iter TODO
func (matcher *Matcher) Iter(f func(*netloc.Netloc, interface{}) bool) {
	matcher.RLock()
	defer matcher.RUnlock()

	for _, unit := range matcher.units {
		for _, nlrs := range unit.nlcRules {
			for _, nlr := range nlrs {
				if !f(nlr.Netloc, nlr.rule) {
					goto BREAK
				}
			}
		}
	}
BREAK:
}

func nextPiece(piece string) string {
	l := len(piece)
	i := 0
	for ; i < l && piece[i] == ByteDot; i++ {
	}
	if i != 0 {
		piece = piece[i:]
	} else {
		i := strings.Index(piece, Dot)
		if i < 0 {
			return Empty
		}
		for ; i < l && piece[i] == ByteDot; i++ {
		}
		piece = piece[i-1:]
	}
	return piece
}
