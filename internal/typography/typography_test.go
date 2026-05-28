package typography

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
			name:     "em dash",
			input:    "Something --- very important.",
			contains: []string{"Something — very important."},
			absent:   []string{"---"},
		},
		{
			name:     "en dash",
			input:    "Pages 10--20 of the book.",
			contains: []string{"Pages 10–20 of the book."},
			absent:   []string{"--"},
		},
		{
			name:     "ellipsis",
			input:    "Wait for it...",
			contains: []string{"Wait for it…"},
			absent:   []string{"..."},
		},
		{
			name:     "all three together",
			input:    "Well... he said --- or maybe -- who knows.",
			contains: []string{"…", "—", "–"},
			absent:   []string{"...", "---"},
		},
		{
			name:     "code span preserved",
			input:    "Use `--flag` for options.",
			contains: []string{"`--flag`"},
		},
		{
			name:     "code block preserved",
			input:    "Text with --dash--.\n```\ncode with -- dashes\n```\nMore --text--.",
			contains: []string{"Text with –dash–.", "code with -- dashes", "More –text–."},
		},
		{
			name:     "no transformations needed",
			input:    "Plain text without any special characters.",
			contains: []string{"Plain text without any special characters."},
		},
		{
			name:     "horizontal rule preserved",
			input:    "---",
			contains: []string{"---"},
		},
		{
			name:     "front matter delimiters preserved",
			input:    "---\ntitle: Test\n---\nContent with -- dash.",
			contains: []string{"---\ntitle: Test\n---\nContent with – dash."},
		},
		{
			name:     "empty string",
			input:    "",
			contains: []string{""},
		},
		{
			name:     "multiple code spans on same line",
			input:    "Use `--verbose` or `--quiet` for output -- choose one.",
			contains: []string{"`--verbose`", "`--quiet`", "– choose one."},
		},
		{
			name:     "table separator preserved",
			input:    "| Col1 | Col2 |\n|------|------|\n| a    | b    |",
			contains: []string{"|------|------|"},
		},
		{
			name:     "table separator with alignment preserved",
			input:    "|:------|:------:|------:|",
			contains: []string{"|:------|:------:|------:|"},
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

func TestProcessPreservesCodeBlocks(t *testing.T) {
	input := "```python\na = b -- c\nprint(\"...\")\n```\n"
	result := Process(input)
	if !strings.Contains(result, "a = b -- c") {
		t.Errorf("code block content should not be transformed\ngot: %q", result)
	}
	if !strings.Contains(result, "print(\"...\")") {
		t.Errorf("code block content should not be transformed\ngot: %q", result)
	}
}
