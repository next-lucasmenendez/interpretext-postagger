// Classify all words into a sentence with each tag
package gopostagger

import (
	"sort"
	"regexp"
	"strings"
)

// Data struct that contains data from each word: order into their sentence, raw content and proposed tag.
type Token struct {
	Order int
	Raw, Tag string
}

// Data struct that contains orderable list of words (tokens)
type Sentence []*Token

func (s Sentence) Len() int {
	return len(s)
}

func (s Sentence) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Sentence) Less(i, j int) bool {
	return s[i].Order < s[j].Order
}

// Load model pased by arguments and call tagger to get each word tag based on that model.
// If model doesnt exists, get corpus associated and train the model.
func (sentence Sentence) Tag(modelName string) error {
	var err error
	var m *Model
	if m, err = Load(modelName); err != nil {
		if err = m.Train(modelName); err != nil {
			return err
		}
	}

	sort.Sort(sentence)
	var currentTag string = Start
	for _, token := range sentence {
		var maxScore float64
		var lowerToken string = strings.ToLower(token.Raw)
		var probs map[string]float64 = m.GetProbs(lowerToken, currentTag)
		for tag, score := range probs {
			if score > maxScore {
				currentTag = tag
				maxScore = score
			}
		}
		token.Tag = currentTag
	}
	return nil
}

func Tokenize(raw_sentence string) Sentence {
	var tokenRgx *regexp.Regexp = regexp.MustCompile(`\s`)
	var puntStart *regexp.Regexp = regexp.MustCompile(`(¡|¿|\(\[\{)(.+)`)
	var puntEnd *regexp.Regexp = regexp.MustCompile(`(.+)(\.|:|;|,|!|\?|\)|]|})`)

	var offset int
	var sentence Sentence
	var words []string = tokenRgx.Split(raw_sentence, -1)
	for i, raw := range words {
		if rawStart := puntStart.FindStringSubmatch(raw); len(rawStart) > 2 {
			sentence = append(sentence, &Token{Order: i + offset, Raw: rawStart[1]})
			offset++
			sentence = append(sentence, &Token{Order: i + offset, Raw: rawStart[2]})
		} else if rawEnd := puntEnd.FindStringSubmatch(raw); len(rawEnd) > 2 {
			sentence = append(sentence, &Token{Order: i + offset, Raw: rawEnd[1]})
			offset++
			sentence = append(sentence, &Token{Order: i + offset, Raw: rawEnd[2]})
		} else {
			sentence = append(sentence, &Token{Order: i + offset, Raw: raw})
		}
	}

	return sentence
}
