package tokenizer

import (
	"strings"
	"unicode"
)

// Subsplit the text and return the slice of all begin:end offset pairs
func Subsplit(text string) [][2]int {
	var result = make([][2]int, 0, 1)
	for pair := range Subtokenize(text) {
		result = append(result, pair)
	}
	return result
}

// Subtokenize the text by generating all begin:end offset pairs
func Subtokenize(text string) <-chan [2]int {
	gen := make(chan [2]int)
	go func() {
		defer close(gen)
		s := state{
			text: text, generator: gen,
			textLen: len(text), tokenStart: len(text)}
		s.run()
	}()
	return gen
}

type state struct {
	text           string
	generator      chan<- [2]int
	next           characterGroup
	nextOffset     int
	current        characterGroup
	currentOffset  int
	previous       characterGroup
	previousOffset int
	before         characterGroup
	beforeOffset   int
	textLen        int
	tokenStart     int
}

// update the next offset and characterGroup (dropping the before's values)
func (s *state) update(offset int, codepoint characterGroup) {
	s.before = s.previous
	s.beforeOffset = s.previousOffset
	s.previous = s.current
	s.previousOffset = s.currentOffset
	s.current = s.next
	s.currentOffset = s.nextOffset
	s.next = codepoint
	s.nextOffset = offset
}


// characterGroup is a concept much alike a Unicode category
type characterGroup int

const (
	start       characterGroup = iota
	lower        // 1
	upper        // 2
	number       // 3
	terminal     // 4
	hyphen       // 5
	punctuation  // 6
	apostrophe   // 7
	symbol       // 8
	space        // 9
	end          // 10
)

// alnumGroup returns true if that characterGroup is alphanumeric
func alnumGroup(group characterGroup) bool {
	switch group {
	case number:
		return true
	case lower:
		return true
	case upper:
		return true
	default:
		return false
	}
}

// hyphens including the underscore (which is treated alike)
const hyphens = "\u00AD\u058A\u05BE\u0F0C\u1400\u1806\u2010\u2011\u2012\u2e17\u30A0_-"

// apostrophes including the single quote (which is treated alike)
const apostrophes = "\u00B4\u02B9\u02BC\u2019\u2032'"

// run the sub-tokenizer and produce any token offsets
func (s *state) run() {
	for offset, codepoint := range s.text {
		s.produce(s.process(offset, codepoint))
	}
	s.update(s.textLen, end)
	s.produce(s.emit())
	s.update(s.textLen, end)
	s.produce(s.emit())
}

// produce the next sub-text offsets into the state generator
func (s *state) produce(tokenEnd int) {
	if tokenEnd > 0 && (tokenEnd > s.tokenStart || tokenEnd == s.textLen) {
		s.generator <- [2]int{s.tokenStart, tokenEnd}

		if s.current != symbol && (s.next == space || s.current == hyphen) {
			s.tokenStart = s.textLen
		} else {
			s.tokenStart = tokenEnd
		}
	}
}

// process determines whether to produce the next sub-tokens
func (s *state) process(offset int, codepoint rune) int {
	switch {
	case unicode.IsLower(codepoint):
		return s.lower(offset)
	case unicode.IsSpace(codepoint):
		return s.space(offset)
	case unicode.IsUpper(codepoint):
		return s.upper(offset)
	case unicode.IsNumber(codepoint):
		return s.number(offset)
	case codepoint == '.' || codepoint == '?' || codepoint == '!':
		return s.terminal(offset)
	case unicode.Is(unicode.Ps, codepoint) || unicode.Is(unicode.Pe, codepoint):
		return s.punctuation(offset)
	case strings.ContainsRune(hyphens, codepoint):
		return s.hyphen(offset)
	case strings.ContainsRune(apostrophes, codepoint):
		return s.apostrophe(offset)
	default:
		return s.symbol(offset)
	}
}

/* state-specific handling: update the current state and decide what to emit */

func (s *state) number(offset int) int {
	s.update(offset, number)
	return s.emit()
}

func (s *state) upper(offset int) int {
	s.update(offset, upper)
	return s.emit()
}

func (s *state) lower(offset int) int {
	s.update(offset, lower)
	return s.emit()
}

func (s *state) terminal(offset int) int {
	s.update(offset, terminal)
	return s.emit()
}

func (s *state) hyphen(offset int) int {
	s.update(offset, hyphen)
	return s.emit()
}

func (s *state) punctuation(offset int) int {
	s.update(offset, punctuation)
	return s.emit()
}

func (s *state) apostrophe(offset int) int {
	s.update(offset, apostrophe)
	return s.emit()
}

func (s *state) symbol(offset int) int {
	s.update(offset, symbol)
	return s.emit()
}

func (s *state) space(offset int) int {
	s.update(offset, space)
	return s.emit()
}

func (s *state) emit() int {
	if (s.previous == space || s.previous == start) && s.current != space {
		s.tokenStart = s.currentOffset
	} else if alnumGroup(s.current) && s.tokenStart == s.textLen {
		s.tokenStart = s.currentOffset
	}

	switch s.current {
	case lower:
		return s.emitAtLower()
	case upper:
		return s.emitAtUpper()
	case number:
		return s.emitAtNumber()
	case terminal:
		return s.emitAtTerminal()
	case hyphen:
		return s.emitAtHyphen()
	case punctuation:
		return s.emitAtPunctuation()
	case apostrophe:
		return s.emitAtApostrophe()
	case symbol:
		return s.emitAtSymbol()
	case space:
		return s.emitAtSpace()
	case end:
		return s.emitAtEnd()
	default:
		return 0
	}
}

/* state-specific spitting: return 0 if nothing to split, the offset to split at otherwise */

func (s *state) emitAtLower() int {
	if !alnumGroup(s.previous) &&
		s.previous != space &&
		s.previous != hyphen &&
		(s.previous != apostrophe || s.before == space) {
		return s.currentOffset
	} else {
		return 0
	}
}

func (s *state) emitAtUpper() int {
	if s.previous == lower ||
		(!alnumGroup(s.previous) &&
			s.previous != space &&
			(s.previous != apostrophe || s.before == space)) {
		return s.currentOffset
	} else {
		return 0
	}
}

func (s *state) emitAtNumber() int {
	if !alnumGroup(s.previous) && s.previous != space {
		if s.before != number {
			return s.currentOffset
		} else {
			return 0
		}
	} else {
		return 0
	}
}

func (s *state) emitAtTerminal() int {
	if s.previous != space && !(s.previous == number && s.next == number) {
		return s.currentOffset
	} else {
		return 0
	}
}

func (s *state) emitAtHyphen() int {
	if s.previous != space && !(s.previous == number && s.next == number) {
		return s.currentOffset
	} else {
		return 0
	}
}

func (s *state) emitAtPunctuation() int {
	if s.previous != space && !(s.previous == number && s.next == number) {
		return s.currentOffset
	} else {
		return 0
	}
}

func (s *state) emitAtApostrophe() int {
	if s.previous == lower && s.next == lower &&
		s.text[s.previousOffset] == 'n' && s.text[s.nextOffset] == 't' {
		return s.previousOffset
	} else if s.previous != space && !(s.previous == number && s.next == number) {
		return s.currentOffset
	} else {
		return 0
	}
}

func (s *state) emitAtSymbol() int {
	if s.previous != space && !(s.previous == number && s.next == number) {
		return s.currentOffset
	} else {
		return 0
	}
}

func (s *state) emitAtSpace() int {
	if s.previous != space && s.previous != start {
		return s.currentOffset
	} else {
		return 0
	}
}

func (s *state) emitAtEnd() int {
	return s.currentOffset
}
