package highlight

import (
	"strings"
	"testing"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
		absent   []string
	}{
		{
			name:     "basic highlight",
			input:    "This is ==highlighted text== in a sentence.",
			contains: []string{"🟡", "highlighted text", "**"},
			absent:   []string{"==highlighted text=="},
		},
		{
			name:     "multiple highlights",
			input:    "Both ==first== and ==second== are highlighted.",
			contains: []string{"🟡 first", "🟡 second"},
			absent:   []string{"==first==", "==second=="},
		},
		{
			name:     "highlight with special chars",
			input:    "Check ==important: read this!== now.",
			contains: []string{"🟡", "important: read this!"},
		},
		{
			name:     "no highlight — plain text",
			input:    "No highlights here.",
			contains: []string{"No highlights here."},
			absent:   []string{"🟡"},
		},
		{
			name:     "single equals — not highlight",
			input:    "a = b and c = d",
			contains: []string{"a = b and c = d"},
			absent:   []string{"🟡"},
		},
		{
			name:     "empty highlight — ignored",
			input:    "This is ==== not a highlight.",
			contains: []string{"===="},
		},
		{
			name:     "highlight at start of line",
			input:    "==Start== of line.",
			contains: []string{"🟡 Start"},
			absent:   []string{"==Start=="},
		},
		{
			name:     "highlight at end of line",
			input:    "End of ==line==",
			contains: []string{"🟡 line"},
			absent:   []string{"==line=="},
		},
		{
			name:     "multiline — no match across lines",
			input:    "==start\nend==",
			contains: []string{"==start\nend=="},
			absent:   []string{"🟡"},
		},
		{
			name:     "empty string",
			input:    "",
			contains: []string{""},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Process(tc.input)
			for _, want := range tc.contains {
				if !strings.Contains(result, want) {
					t.Errorf("expected %q in output\ngot: %q", want, result)
				}
			}
			for _, unwanted := range tc.absent {
				if strings.Contains(result, unwanted) {
					t.Errorf("expected %q NOT in output\ngot: %q", unwanted, result)
				}
			}
		})
	}
}

func TestProcessPreservesNonHighlight(t *testing.T) {
	input := "# Heading\n\nPlain paragraph with no highlights.\n"
	result := Process(input)
	if result != input {
		t.Errorf("non-highlight content should pass through unchanged\ngot: %q\nwant: %q", result, input)
	}
}
