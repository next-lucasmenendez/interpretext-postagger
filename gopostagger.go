// HMM applied to Part-Of-Speech Tagging in Go.
package gopostagger

// Split sentence into tokens and create and process tokens.
func TagSentence(raw_sentence string, model string) (Sentence, error) {
	var sentence Sentence = Tokenize(raw_sentence)
	if err := sentence.Tag(model); err != nil {
		return Sentence{}, err
	}

	var result Sentence
	for _, t := range sentence {
		result = append(result, t)
	}
	return result, nil
}

