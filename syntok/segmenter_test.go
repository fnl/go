package tokenizer

import (
	"reflect"
	"testing"
)

func TestAnalyze(t *testing.T) {
	tests := []struct {
		name      string
		paragraph string
		want      []Sentence
	}{
		{
			"basic test",
			"text",
			[]Sentence{{{"", 0, "text"}}},
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Analyze(tt.paragraph); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Analyze(%v) = %v, want %v", tt.paragraph, got, tt.want)
			}
		})
	}
}
