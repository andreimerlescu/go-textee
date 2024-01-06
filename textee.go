package go_textee

import (
	`fmt`
	`log`
	`sort`
	`strings`
	`sync`
	`sync/atomic`

	gematria "github.com/andreimerlescu/go-gematria"
)

type ITextee interface {
	ParseString(input string) *Textee
	SortedSubstrings() SortedStringQuantities
	CalculateGematria() *Textee
}

func NewTextee(opts ...string) *Textee {
	input := strings.Join(opts, " ")
	gem, err := gematria.NewGematria(input)
	if err != nil {
		log.Printf("failed to get NewGematria due to err %v", err)
		return nil
	}
	tt := &Textee{
		Input:         input,
		Gematria:      gem,
		Substrings:    make(map[string]*atomic.Int32),
		Gematrias:     make(map[string]gematria.Gematria),
		ScoresEnglish: make(map[uint][]string),
		ScoresJewish:  make(map[uint][]string),
		ScoresSimple:  make(map[uint][]string),
	}
	tt.ParseString(strings.Join(opts, " "))
	return tt
}

func (tt *Textee) ParseString(input string) *Textee {
	sentences := string_to_sentence_slice(input)

	var wg sync.WaitGroup
	for _, sentence := range sentences {
		wg.Add(1)
		go func(sentence string) {
			defer wg.Done()
			words := strings.Fields(sentence)

			for i := 0; i < len(words); i++ {
				for j := i + 1; j <= i+3 && j <= len(words); j++ {
					substring := strings.Join(words[i:j], " ")
					cleaned_substring := clean_substring(substring)
					cleaned_substring = strings.ToLower(cleaned_substring)

					if cleaned_substring != "" {
						tt.mu.Lock()
						if _, ok := tt.Substrings[cleaned_substring]; !ok {
							tt.Substrings[cleaned_substring] = new(atomic.Int32)
						}
						tt.Substrings[cleaned_substring].Add(1)
						tt.mu.Unlock()
					}
				}
			}
		}(sentence)
	}
	wg.Wait()
	log.Printf("done waiting")
	return tt
}

func (tt *Textee) String() string {
	if len(tt.Substrings) == 0 {
		return fmt.Sprintf("Error: No .ParseString() invocation found for tt (type=%T) (address=%p).", tt, &tt)
	}
	hasGematria := len(tt.ScoresEnglish) > 0 || len(tt.ScoresJewish) > 0 || len(tt.ScoresSimple) > 0
	var output string
	for _, data := range tt.SortedSubstrings() {
		hasGematria := hasGematria
		if hasGematria {
			addition := fmt.Sprintf("\"%v\": %d [English %d] [Jewish %d] [Simple %d]\n",
				data.Substring, data.Quantity,
				tt.Gematrias[data.Substring].English,
				tt.Gematrias[data.Substring].Jewish,
				tt.Gematrias[data.Substring].Simple)
			output = output + addition
		} else {
			addition := fmt.Sprintf("\"%v\": %d\n", data.Substring, data.Quantity)
			output = output + addition
		}
	}
	return output
}

func (tt *Textee) SortedSubstrings() SortedStringQuantities {
	tt.mu.RLock()
	defer tt.mu.RUnlock()
	var sortedQuantities SortedStringQuantities

	for k, v := range tt.Substrings {
		quantity := int(v.Load())
		sortedQuantities = append(sortedQuantities, SubstringQuantity{Substring: k, Quantity: quantity})
	}
	sort.Sort(sortedQuantities)

	return sortedQuantities
}

func (tt *Textee) CalculateGematria() *Textee {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	substrings := tt.Substrings
	english_results := make(map[uint][]string)
	jewish_results := make(map[uint][]string)
	simple_results := make(map[uint][]string)
	errorCounter := atomic.Int32{}
	for substring, _ := range substrings {
		gemscore, err := gematria.NewGematria(substring)
		if err != nil {
			errorCounter.Add(1)
			continue
		}
		english_results[gemscore.English] = append(english_results[gemscore.English], substring)
		jewish_results[gemscore.Jewish] = append(jewish_results[gemscore.Jewish], substring)
		simple_results[gemscore.Simple] = append(simple_results[gemscore.Simple], substring)
		if tt.Gematrias == nil {
			tt.Gematrias = make(map[string]gematria.Gematria)
		}
		tt.Gematrias[substring] = gemscore
	}
	substrings = nil
	tt.ScoresEnglish = english_results
	tt.ScoresJewish = jewish_results
	tt.ScoresSimple = simple_results
	english_results = nil
	jewish_results = nil
	simple_results = nil
	return tt
}
