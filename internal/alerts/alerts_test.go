package alerts

import (
	"strings"
	"testing"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string    // strings that must be present
		absent   []string    // strings that must NOT be present
	}{
		{
			name:     "NOTE alert",
			input:    "> [!NOTE]\n> This is a note.",
			contains: []string{"ℹ️", "Note", "This is a note."},
			absent:   []string{"[!NOTE]"},
		},
		{
			name:     "WARNING alert",
			input:    "> [!WARNING]\n> Be careful here.",
			contains: []string{"⚠️", "Warning", "Be careful here."},
			absent:   []string{"[!WARNING]"},
		},
		{
			name:     "TIP alert",
			input:    "> [!TIP]\n> Use this trick.",
			contains: []string{"💡", "Tip", "Use this trick."},
			absent:   []string{"[!TIP]"},
		},
		{
			name:     "IMPORTANT alert",
			input:    "> [!IMPORTANT]\n> Read this first.",
			contains: []string{"❗", "Important", "Read this first."},
			absent:   []string{"[!IMPORTANT]"},
		},
		{
			name:     "CAUTION alert",
			input:    "> [!CAUTION]\n> This is dangerous.",
			contains: []string{"🔴", "Caution", "This is dangerous."},
			absent:   []string{"[!CAUTION]"},
		},
		{
			name:     "multi-line alert",
			input:    "> [!NOTE]\n> First line.\n> Second line.\n> Third line.",
			contains: []string{"ℹ️", "Note", "First line.", "Second line.", "Third line."},
			absent:   []string{"[!NOTE]"},
		},
		{
			name:     "alert with blank line after",
			input:    "> [!WARNING]\n> Content here.\n\nRegular paragraph.",
			contains: []string{"⚠️", "Warning", "Content here.", "Regular paragraph."},
			absent:   []string{"[!WARNING]"},
		},
		{
			name:     "no alert — regular blockquote",
			input:    "> This is a regular blockquote.\n> No alert here.",
			contains: []string{"This is a regular blockquote.", "No alert here."},
			absent:   []string{"ℹ️", "⚠️", "💡", "❗", "🔴"},
		},
		{
			name:     "no alert — plain text",
			input:    "Just a paragraph.\n\nAnother paragraph.",
			contains: []string{"Just a paragraph.", "Another paragraph."},
		},
		{
			name:     "mixed content — alert and regular blockquote",
			input:    "> [!NOTE]\n> Alert content.\n\n> Regular quote.",
			contains: []string{"ℹ️", "Note", "Alert content.", "Regular quote."},
			absent:   []string{"[!NOTE]"},
		},
		{
			name:     "case insensitive tag",
			input:    "> [!note]\n> Lowercase note.",
			contains: []string{"ℹ️", "Note", "Lowercase note."},
			absent:   []string{"[!note]"},
		},
		{
			name:     "unknown alert type passes through",
			input:    "> [!UNKNOWN]\n> Some content.",
			contains: []string{"[!UNKNOWN]", "Some content."},
		},
		{
			name:     "empty content",
			input:    "",
			contains: []string{""},
		},
		{
			name:     "alert header only, no body",
			input:    "> [!NOTE]",
			contains: []string{"ℹ️", "Note"},
			absent:   []string{"[!NOTE]"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Process(tc.input)

			for _, want := range tc.contains {
				if !strings.Contains(result, want) {
					t.Errorf("expected output to contain %q\ngot: %s", want, result)
				}
			}
			for _, unwanted := range tc.absent {
				if strings.Contains(result, unwanted) {
					t.Errorf("expected output NOT to contain %q\ngot: %s", unwanted, result)
				}
			}
		})
	}
}

func TestProcessPreservesNonAlertContent(t *testing.T) {
	input := "# Heading\n\nParagraph text.\n\n```go\nfmt.Println(\"hello\")\n```\n\n- List item\n"
	result := Process(input)
	if result != input {
		t.Errorf("non-alert content should pass through unchanged\ngot: %s\nwant: %s", result, input)
	}
}
