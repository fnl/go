package tokenizer

// Token has a Prefix string found previous it, and the Offset of its actual Value.
type Token struct {
	Prefix string
	Offset int
	Value  string
}

// Tokenize text by generating Tokens.
func Tokenize(text string) <-chan Token {
	var generator = make(chan Token)

	go func() {
		defer close(generator)
		var last = 0
		for off := range Subtokenize(text) {
			generator <- Token{text[last:off[0]], off[0], text[off[0]:off[1]]}
			last = off[1]
		}
		if last < len(text) {
			generator <- Token{text[last:], last, ""}
		}
	}()

	return generator
}

// Split text into Tokens.
func Split(text string) []Token {
	var result = make([]Token, 0, len(text)/5)

	for token := range Tokenize(text) {
		result = append(result, token)
	}

	return result
}
