package tokenizer

type Sentence []Token

// Segment a consecutive paragraph into Sentences.
func Segment(paragraph string) <-chan Sentence {
	var generator = make(chan Sentence)

	go func() {
		defer close(generator)

		var state = &segmenter{
			beginSegmentation,
			make([]Token, 0, 3),
			make(Sentence, 0, len(paragraph)/5),
			generator}

		for token := range Tokenize(paragraph) {
			state.buffer = append(state.buffer, token)
			evaluate(state)
		}

		if len(state.production) > 0 {
			generator <- state.production
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
