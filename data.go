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
	Input         string
	Gematria      gematria.Gematria
	Substrings    map[string]*atomic.Int32 // map[Substring]*atomic.Int32
	Gematrias     map[string]gematria.Gematria
	ScoresEnglish map[uint][]string
	ScoresJewish  map[uint][]string
	ScoresSimple  map[uint][]string
}

type SubstringQuantity struct {
	Substring string
	Quantity  int
}

type SortedStringQuantities []SubstringQuantity

var reg_clean_substring = regexp.MustCompile(`[^a-zA-Z0-9\s]`)

var sem = sema.New(runtime.GOMAXPROCS(0))

// string_to_sentence_slice splits text into sentences, considering abbreviations.
func string_to_sentence_slice(text string) []string {
	regex := regexp.MustCompile(`(?m)([^.!?]*[.!?])(?:\s|$)`)
	matches := regex.FindAllString(text, -1)

	// Clean up the sentences.
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
