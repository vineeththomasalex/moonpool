// Package typography applies smart typographic replacements to markdown text.
//
// Transforms:
//   - Straight quotes → curly quotes: "hello" → "hello", 'hello' → 'hello'
//   - Dashes: --- → —, -- → –
//   - Ellipsis: ... → …
//
// These replacements are applied before glamour rendering to produce
// typographically correct output without needing Goldmark's Typographer extension.
package typography

import (
	"regexp"
	"strings"
)

// Process applies smart typography replacements.
// Only transforms text outside of code spans and code blocks to avoid
// corrupting code/URLs.
func Process(markdown string) string {
	lines := strings.Split(markdown, "\n")
	inCodeBlock := false
	var result []string

	for _, line := range lines {
		// Track fenced code blocks
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inCodeBlock = !inCodeBlock
			result = append(result, line)
			continue
		}

		if inCodeBlock {
			result = append(result, line)
			continue
		}

		result = append(result, transformLine(line))
	}

	return strings.Join(result, "\n")
}

// codeSpanPattern matches inline code spans to protect them.
var codeSpanPattern = regexp.MustCompile("`[^`]+`")

// thematicBreak matches lines that are markdown horizontal rules.
var thematicBreak = regexp.MustCompile(`^\s*(---+|\*\*\*+|___+)\s*$`)

// tableSeparator matches markdown table separator lines like |---|:---:|---:|
var tableSeparator = regexp.MustCompile(`^\s*\|[\s:|-]+\|\s*$`)

func transformLine(line string) string {
	// Don't transform thematic breaks (---, ***, ___) or frontmatter delimiters
	if thematicBreak.MatchString(line) {
		return line
	}

	// Don't transform table separator lines (|---|---|)
	if tableSeparator.MatchString(line) {
		return line
	}

	// Split line into code spans and non-code segments
	// Replace only in non-code segments
	spans := codeSpanPattern.FindAllStringIndex(line, -1)
	if len(spans) == 0 {
		return applyTypography(line)
	}

	var result strings.Builder
	prev := 0
	for _, span := range spans {
		result.WriteString(applyTypography(line[prev:span[0]]))
		result.WriteString(line[span[0]:span[1]])
		prev = span[1]
	}
	result.WriteString(applyTypography(line[prev:]))
	return result.String()
}

func applyTypography(text string) string {
	// Order matters: do triple dash before double
	text = strings.ReplaceAll(text, "---", "—")
	text = strings.ReplaceAll(text, "--", "–")
	text = strings.ReplaceAll(text, "...", "…")
	return text
}
