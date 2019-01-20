package tokenizer

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []Token
	}{
		{
			"basic test",
			"text",
			[]Token{{"", 0, "text"}},
		},
		{
			"split two tokens",
			"token1 token2",
			[]Token{{"", 0, "token1"}, {" ", 7, "token2"}},
		},
		{
			"two spaces between tokens",
			"token1  token2",
			[]Token{{"", 0, "token1"}, {"  ", 8, "token2"}},
		},
		{
			"capture tail spaces, too",
			" text ",
			[]Token{{" ", 1, "text"}, {" ", len(" text "), ""}},
		},
		{
			"is newline aware",
			"token1\r\ntoken2",
			[]Token{{"", 0, "token1"}, {"\r\n", len("token1\r\n"), "token2"}},
		},
		{
			"is NBS aware",
			"token1\u00A0token2",
			[]Token{{"", 0, "token1"}, {"\u00A0", len("token1\u00A0"), "token2"}},
		},
		{
			"is Unicode aware",
			"token1\u2028token2",
			[]Token{{"", 0, "token1"}, {"\u2028", len("token1\u2028"), "token2"}},
		},
		{
			"split and preserve punctuation",
			"This, that, and him!",
			[]Token{
				{"", 0, "This"},
				{"", len("This"), ","},
				{" ", len("This, "), "that"},
				{"", len("This, that"), ","},
				{" ", len("This, that, "), "and"},
				{" ", len("This, that, and "), "him"},
				{"", len("This, that, and him"), "!"},
			},
		},
		{
			"preserve tokens with inner punctuation",
			"31.12.2000 23:59:59",
			[]Token{
				{"", 0, "31.12.2000"},
				{" ", len("31.12.2000 "), "23:59:59"},
			},
		},
		{
			"split tokens with CamelCase",
			"camelCase",
			[]Token{
				{"", 0, "camel"},
				{"", 5, "Case"},
			},
		},
		{
			"split hyphens",
			"hyphen-ate",
			[]Token{
				{"", 0, "hyphen"},
				{"-", 7, "ate"},
			},
		},
		{
			"split inner sentence marker .",
			"last.First",
			[]Token{
				{"", 0, "last"},
				{"", 4, "."},
				{"", 5, "First"},
			},
		},
		{
			"split inner sentence marker !",
			"last!First",
			[]Token{
				{"", 0, "last"},
				{"", 4, "!"},
				{"", 5, "First"},
			},
		},
		{
			"split inner sentence marker ?",
			"last?First",
			[]Token{
				{"", 0, "last"},
				{"", 4, "?"},
				{"", 5, "First"},
			},
		},
		{
			"split around open/close punctuation",
			"this(that)there",
			[]Token{
				{"", 0, "this"},
				{"", 4, "("},
				{"", 5, "that"},
				{"", 9, ")"},
				{"", 10, "there"},
			},
		},
		{
			"do not split hyphens inside digits",
			"10-11-2012 alpha-1",
			[]Token{
				{"", 0, "10-11-2012"},
				{" ", 11, "alpha"},
				{"-", 17, "1"},
			},
		},
		{
			"split around underscores like hyphens",
			"this_that 123_456",
			[]Token{
				{"", 0, "this"},
				{"_", 5, "that"},
				{" ", 10, "123_456"},
			},
		},
		{
			"split apostrophes in text",
			"He's 'tis 1'234'567 10's",
			[]Token{
				{"", 0, "He"},
				{"", 2, "'s"},
				{" ", 5, "'"},
				{"", 6, "tis"},
				{" ", 10, "1'234'567"},
				{" ", 20, "10"},
				{"", 22, "'s"},
			},
		},
		{
			"correctly split single quotes",
			"Here's it: 'They said so!'",
			[]Token{
				{"", 0, "Here"},
				{"", 4, "'s"},
				{" ", 7, "it"},
				{"", 9, ":"},
				{" ", 11, "'"},
				{"", 12, "They"},
				{" ", 17, "said"},
				{" ", 22, "so"},
				{"", 24, "!"},
				{"", 25, "'"},
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Split(tt.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefault().Split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOneSplit(t *testing.T) {
	t.Run("one-off test", func(t *testing.T) {
		text := "This, that, and him!"
		want := []Token{
			{"", 0, "This"},
			{"", len("This"), ","},
			{" ", len("This, "), "that"},
			{"", len("This, that"), ","},
			{" ", len("This, that, "), "and"},
			{" ", len("This, that, and "), "him"},
			{"", len("This, that, and him"), "!"},
		}

		if got := Split(text); !reflect.DeepEqual(got, want) {
			t.Errorf("Split() = %v, want %v", got, want)
		}
	})
}
