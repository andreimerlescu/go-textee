package go_textee

import (
	`regexp`
	`runtime`
	`strings`
	`sync`
	`sync/atomic`

	gematria "github.com/andreimerlescu/go-gematria"
	sema "github.com/andreimerlescu/go-sema"
)

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

var reg_clean_substring = regexp.MustCompile(`[^a-zA-Z0-9\s]`)
var reg_find_sentences = regexp.MustCompile(`(?m)([^.!?]*[.!?])(?:\s|$)`)

var sem = sema.New(runtime.GOMAXPROCS(0))

// string_to_sentence_slice splits text into sentences, considering abbreviations.
func string_to_sentence_slice(text string) []string {
	matches := reg_find_sentences.FindAllString(text, -1)
	for i, match := range matches {
		matches[i] = strings.TrimSpace(match)
	}

	if len(matches) == 0 {
		return []string{text}
	}
	return matches
}

// clean_substring returns the string to A-Za-z0-9\s only
func clean_substring(word string) string {
	return reg_clean_substring.ReplaceAllString(word, "")
}
