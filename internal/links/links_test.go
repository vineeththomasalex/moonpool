package links

import (
	"strings"
	"testing"
)

func TestAddOSC8Links(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
		absent   []string
	}{
		{
			name:     "bare https URL",
			input:    "Visit https://example.com for more.",
			contains: []string{"\x1b]8;;https://example.com\x1b\\", "https://example.com", "\x1b]8;;\x1b\\"},
		},
		{
			name:     "bare http URL",
			input:    "See http://example.com/page for info.",
			contains: []string{"\x1b]8;;http://example.com/page\x1b\\"},
		},
		{
			name:     "URL with path and query",
			input:    "Go to https://example.com/path?q=1&b=2 now.",
			contains: []string{"\x1b]8;;https://example.com/path?q=1&b=2\x1b\\"},
		},
		{
			name:     "URL with trailing period",
			input:    "See https://example.com.",
			contains: []string{"\x1b]8;;https://example.com\x1b\\https://example.com\x1b]8;;\x1b\\."},
		},
		{
			name:     "multiple URLs",
			input:    "Visit https://one.com and https://two.com for info.",
			contains: []string{"\x1b]8;;https://one.com\x1b\\", "\x1b]8;;https://two.com\x1b\\"},
		},
		{
			name:     "no URLs — pass through",
			input:    "No URLs here, just plain text.",
			contains: []string{"No URLs here, just plain text."},
			absent:   []string{"\x1b]8;;"},
		},
		{
			name:     "already has OSC 8 — skip",
			input:    "\x1b]8;;https://example.com\x1b\\click\x1b]8;;\x1b\\",
			contains: []string{"\x1b]8;;https://example.com\x1b\\click\x1b]8;;\x1b\\"},
		},
		{
			name:     "URL with ANSI styling around it",
			input:    "\x1b[38;5;30;4mhttps://example.com\x1b[0m",
			contains: []string{"\x1b]8;;https://example.com\x1b\\"},
		},
		{
			name:     "empty string",
			input:    "",
			contains: []string{""},
		},
		{
			name:     "URL with fragment",
			input:    "See https://example.com/page#section for details.",
			contains: []string{"\x1b]8;;https://example.com/page#section\x1b\\"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := AddOSC8Links(tc.input)
			for _, want := range tc.contains {
				if !strings.Contains(result, want) {
					t.Errorf("expected output to contain %q\ngot: %q", want, result)
				}
			}
			for _, unwanted := range tc.absent {
				if strings.Contains(result, unwanted) {
					t.Errorf("expected output NOT to contain %q\ngot: %q", unwanted, result)
				}
			}
		})
	}
}

func TestAddOSC8LinksPreservesText(t *testing.T) {
	input := "Just text without any URLs whatsoever."
	result := AddOSC8Links(input)
	if result != input {
		t.Errorf("text without URLs should pass through unchanged\ngot: %q\nwant: %q", result, input)
	}
}
