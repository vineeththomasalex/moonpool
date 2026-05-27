package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRemoveFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid yaml frontmatter",
			input:    "---\ntitle: Test\n---\n# Body content",
			expected: "# Body content",
		},
		{
			name:     "no frontmatter",
			input:    "# Just a heading\n\nSome text.",
			expected: "# Just a heading\n\nSome text.",
		},
		{
			name:     "frontmatter with blank line after delimiter",
			input:    "---\n\ntitle: Test\n---\n\nBody",
			expected: "Body",
		},
		{
			name:     "empty content",
			input:    "",
			expected: "",
		},
		{
			name:     "only frontmatter delimiters",
			input:    "---\n---\nContent after",
			expected: "Content after",
		},
		{
			name:     "single delimiter (no closing)",
			input:    "---\ntitle: Test\nNo closing delimiter",
			expected: "---\ntitle: Test\nNo closing delimiter",
		},
		{
			name:     "frontmatter with CRLF line endings",
			input:    "---\r\ntitle: Test\r\n---\r\nBody content",
			expected: "Body content",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := string(RemoveFrontmatter([]byte(tc.input)))
			if result != tc.expected {
				t.Errorf("RemoveFrontmatter(%q)\ngot:  %q\nwant: %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsMarkdownFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		// Markdown extensions
		{"README.md", true},
		{"doc.markdown", true},
		{"notes.mdown", true},
		{"file.mkdn", true},
		{"file.mkd", true},

		// Case insensitive
		{"FILE.MD", true},
		{"Doc.Markdown", true},

		// No extension → defaults to true
		{"README", true},
		{"Makefile", true},

		// Non-markdown
		{"main.go", false},
		{"data.json", false},
		{"style.css", false},
		{"script.py", false},
		{"file.txt", false},
		{"image.png", false},
	}

	for _, tc := range tests {
		t.Run(tc.filename, func(t *testing.T) {
			result := IsMarkdownFile(tc.filename)
			if result != tc.expected {
				t.Errorf("IsMarkdownFile(%q) = %v, want %v", tc.filename, result, tc.expected)
			}
		})
	}
}

func TestWrapCodeBlock(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		language string
		expected string
	}{
		{
			name:     "go code",
			content:  "fmt.Println(\"hello\")\n",
			language: "go",
			expected: "```go\nfmt.Println(\"hello\")\n```",
		},
		{
			name:     "no language",
			content:  "plain text\n",
			language: "",
			expected: "```\nplain text\n```",
		},
		{
			name:     "empty content",
			content:  "",
			language: "python",
			expected: "```python\n```",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := WrapCodeBlock(tc.content, tc.language)
			if result != tc.expected {
				t.Errorf("WrapCodeBlock(%q, %q)\ngot:  %q\nwant: %q", tc.content, tc.language, result, tc.expected)
			}
		})
	}
}

func TestRemoveFrontmatterWithFixture(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "testdata", "fixtures", "frontmatter.md"))
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	result := string(RemoveFrontmatter(data))

	if strings.Contains(result, "title: Test Document") {
		t.Error("frontmatter title should be stripped but was found in output")
	}
	if strings.Contains(result, "author: Moonpool Tests") {
		t.Error("frontmatter author should be stripped but was found in output")
	}
	if !strings.Contains(result, "Document After Frontmatter") {
		t.Error("body content should be present after frontmatter removal")
	}
}
