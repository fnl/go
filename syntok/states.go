package tokenizer

type segmentationState int

const (
	beginSegmentation segmentationState = iota
	firstToken
	innerToken
	terminalToken
	endSegmentation
)

type segmenter struct {
	state      segmentationState
	buffer     []Token
	production Sentence
	generator  <-chan Sentence
}

func evaluate(s *segmenter) {
	s.production = s.buffer
}
