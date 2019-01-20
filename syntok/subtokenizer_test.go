package tokenizer

import (
	"reflect"
	"testing"
)

func TestSubtokenize(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  [][2]int
	}{
		{
			"simple example",
			"word",
			[][2]int{{0, 4}},
		},
		{
			"split trailing punctuation",
			"word,",
			[][2]int{{0, 4}, {4, 5}},
		},
		{
			"split leading punctuation",
			"*word",
			[][2]int{{0, 1}, {1, 5}},
		},
		{
			"split trailing terminal",
			"end!",
			[][2]int{{0, 3}, {3, 4}},
		},
		{
			"preserve date tokens",
			"31.12.2000",
			[][2]int{{0, len("31.12.2000")}},
		},
		{
			"preserve iso-date tokens",
			"2000-12-31",
			[][2]int{{0, len("2000-12-31")}},
		},
		{
			"preserve time tokens",
			"23:59:59",
			[][2]int{{0, len("23:59:59")}},
		},
		{
			"split around hyphenation",
			"hy-phen",
			[][2]int{{0, 2}, {3, len("hy-phen")}},
		},
		{
			"split around Unicode hyphenation",
			"hy\u2010phen",
			[][2]int{{0, 2}, {len("hy\u2010"), len("hy\u2010phen")}},
		},
		{
			"split around hyphenation with digit on left side",
			"alpha-1",
			[][2]int{{0, len("alpha")}, {len("alpha-"), len("alpha-1")}},
		},
		{
			"split around hyphenation with digit on right side",
			"1-alpha",
			[][2]int{{0, 1}, {len("1-"), len("1-alpha")}},
		},
		{
			"don't split around hyphenation with digits on both sides",
			"1-2",
			[][2]int{{0, 3}},
		},
		{
			"split around underscores",
			"under_score",
			[][2]int{{0, 5}, {len("under_"), len("under_score")}},
		},
		{
			"split hidden sentence terminals",
			"end.Sentence",
			[][2]int{{0, 3}, {3, 4}, {4, len("end.Sentence")}},
		},
		{
			"split initial open/close punctuation",
			"(that)",
			[][2]int{{0, 1}, {1, 5}, {5, 6}},
		},
		{
			"split inner open/close punctuation",
			"this(that)thus",
			[][2]int{{0, 4}, {4, 5}, {5, 9}, {9, 10}, {10, 14}},
		},
		{
			"split at apostrophes",
			"he's",
			[][2]int{{0, 2}, {2, 4}},
		},
		{
			"split at Unicode apostrophes",
			"he\u2019s",
			[][2]int{{0, 2}, {2, len("he\u2019s")}},
		},
		{
			"split apostrophe negations with the n attached",
			"don't",
			[][2]int{{0, 2}, {2, 5}},
		},
		{
			"don't split apostrophes in numbers",
			"1'000",
			[][2]int{{0, 5}},
		},
		{
			"split at camel case",
			"CamelCase",
			[][2]int{{0, 5}, {5, 9}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Subsplit(tt.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Subtokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOneSubtokenize(t *testing.T) {
	t.Run("one-off test", func(t *testing.T) {
		token := " 'not"
		want := [][2]int{{1, 2}, {2, len(" 'not")}}

		if got := Subsplit(token); !reflect.DeepEqual(got, want) {
			t.Errorf("Subtokenize() = %v, want %v", got, want)
		}
	})
}