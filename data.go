package go_textee

import (
	"errors"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	gematria "github.com/andreimerlescu/go-gematria"
	sema "github.com/andreimerlescu/go-sema"
)

type ITextee interface {
	ParseString(input string) *Textee
	SortedSubstrings() SortedStringQuantities
	CalculateGematria() *Textee
}

var (
	ErrEmptyInput    ArgumentError = errors.New("empty input")
	ErrGematriaParse GematriaError = errors.New("unable to parse gematria for value")
	ErrRegexpMissing RegexpError   = errors.New("regexp compile result missing")
	ErrBadParsing    ParseError    = errors.New("failed to parse the string")
)

type ArgumentError error
type GematriaError error
type RegexpError error
type ParseError error
type CleanError error

type Textee struct {
	mu            sync.RWMutex
	Input         string                       `json:"i"`
	Gematria      gematria.Gematria            `json:"g"`
	Substrings    map[string]*atomic.Int32     `json:"s"` // map[Substring]*atomic.Int32
	Gematrias     map[string]gematria.Gematria `json:"mg"`
	ScoresEnglish map[uint][]string            `json:"ge"`
	ScoresJewish  map[uint][]string            `json:"gj"`
	ScoresSimple  map[uint][]string            `json:"gs"`
}

type SubstringQuantity struct {
	Substring string `json:"s"`
	Quantity  int    `json:"q"`
}

type SortedStringQuantities []SubstringQuantity

var regCleanSubstring *regexp.Regexp
var regFindSentences *regexp.Regexp

var sem = sema.New(runtime.GOMAXPROCS(0))

// stringToSentenceSlice splits text into sentences, considering abbreviations.
func stringToSentenceSlice(text string) ([]string, error) {
	if regFindSentences == nil {
		return nil, ErrRegexpMissing
	}
	matches := regFindSentences.FindAllString(text, -1)
	for i, match := range matches {
		matches[i] = strings.TrimSpace(match)
	}

	if len(matches) == 0 {
		return []string{text}, nil
	}
	return matches, nil
}

// cleanSubstring returns the string to A-Za-z0-9\s only
func cleanSubstring(word string) (string, error) {
	if regCleanSubstring == nil {
		return "", ErrRegexpMissing
	}
	return regCleanSubstring.ReplaceAllString(word, ""), nil
}
