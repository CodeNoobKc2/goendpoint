package request

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	matchPathParam                       = regexp.MustCompile(`^\{[[:space:]]*([[:alnum:]|_|\-]*)[[:space:]]*\}$`)
	errPathTemplateLengthShouldNotBeZero = errors.New("path template length should not be 0")
)

type PathTemplate struct {
	raw        string
	testRegexp *regexp.Regexp
	chunkSize  int
	posToKey   map[int]string
}

func NewPathTemplate(path string) (*PathTemplate, error) {
	if len(path) == 0 {
		return nil, errPathTemplateLengthShouldNotBeZero
	}

	path = strings.TrimPrefix(path, "/")
	splited := strings.Split(path, "/")
	posToKey := make(map[int]string)
	uncompiledRegexp := &strings.Builder{}

	for pos, s := range splited {
		if matchPathParam.MatchString(s) {
			posToKey[pos] = strings.Trim(s, "{} ")
			uncompiledRegexp.WriteString(`/([[:alnum:]|\-|_])*`)
		} else {
			uncompiledRegexp.WriteString("/" + s)
		}
	}

	return &PathTemplate{raw: path, testRegexp: regexp.MustCompile(uncompiledRegexp.String()), chunkSize: len(splited), posToKey: posToKey}, nil
}

func (t PathTemplate) Parse(actual string) (map[string][]string, error) {
	if !t.testRegexp.MatchString(actual) {
		return nil, fmt.Errorf("uri '%v' does not match pattern '/%v'", actual, t.raw)
	}

	splited := strings.Split(strings.TrimPrefix(actual, "/"), "/")
	ret := make(map[string][]string)
	for i, s := range splited {
		if key, has := t.posToKey[i]; has {
			ret[key] = []string{s}
		}
	}

	return ret, nil
}
