package regex

import (
	"regexp"
	"sync"
)

var (
	regexMu = sync.RWMutex{}

	regexMap = make(map[string]*regexp.Regexp)
)

func getRegexp(pattern string) (regex *regexp.Regexp, err error) {
	regexMu.RLock()
	regex = regexMap[pattern]
	regexMu.RUnlock()
	if regex != nil {
		return
	}

	regex, err = regexp.Compile(pattern)
	if err != nil {
		return
	}
	regexMu.Lock()
	regexMap[pattern] = regex
	regexMu.Unlock()
	return
}
