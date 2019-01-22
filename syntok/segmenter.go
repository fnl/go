package tokenizer

type Sentence []Token

// Segment a consecutive paragraph into Sentences.
func Segment(paragraph string) <-chan Sentence {
	var generator = make(chan Sentence)

	go func() {
		defer close(generator)
		var currentProduction = make(Sentence, 0, len(paragraph)/5)

		for token := range Tokenize(paragraph) {
			currentProduction = append(currentProduction, token)
		}

		if len(currentProduction) > 0 {
			generator <- currentProduction
		}
	}()

	return generator
}

// Analyze a paragraph and report all Sentences.
func Analyze(paragraph string) []Sentence {
	var result = make([]Sentence, 0, len(paragraph)/50)

	for segment := range Segment(paragraph) {
		result = append(result, segment)
	}

	return result
}
