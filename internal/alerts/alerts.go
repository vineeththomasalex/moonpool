// Package alerts transforms GitHub-style alert blockquotes into
// styled markdown that glamour can render with visual distinction.
//
// Input:
//
//	> [!NOTE]
//	> This is a note.
//
// Output:
//
//	> **ℹ️ Note**
//	>
//	> This is a note.
package alerts

import (
	"regexp"
	"strings"
)

// AlertType defines the supported GitHub alert types.
type AlertType struct {
	Tag    string // e.g., "NOTE"
	Icon   string // e.g., "ℹ️"
	Label  string // e.g., "Note"
}

var alertTypes = []AlertType{
	{"NOTE", "ℹ️", "Note"},
	{"TIP", "💡", "Tip"},
	{"IMPORTANT", "❗", "Important"},
	{"WARNING", "⚠️", "Warning"},
	{"CAUTION", "🔴", "Caution"},
}

// alertPattern matches `> [!TYPE]` at the start of a blockquote line.
// Captures the alert type keyword (case-insensitive).
var alertPattern = regexp.MustCompile(`(?mi)^>\s*\[!(NOTE|TIP|IMPORTANT|WARNING|CAUTION)\]\s*$`)

// Process transforms GitHub-style alert blockquotes in the markdown
// into styled blockquotes with icon + bold label headers.
// Non-alert content is passed through unchanged.
func Process(markdown string) string {
	if !alertPattern.MatchString(markdown) {
		return markdown
	}

	lines := strings.Split(markdown, "\n")
	var result []string
	i := 0

	for i < len(lines) {
		line := lines[i]

		// Check if this line is an alert header
		match := alertPattern.FindStringSubmatch(line)
		if match == nil {
			result = append(result, line)
			i++
			continue
		}

		// Found an alert — transform the header
		tag := strings.ToUpper(match[1])
		alert := findAlertType(tag)
		if alert == nil {
			result = append(result, line)
			i++
			continue
		}

		// Replace the [!TYPE] line with icon + bold label
		result = append(result, "> **"+alert.Icon+" "+alert.Label+"**")
		result = append(result, ">")
		i++

		// Pass through remaining blockquote lines unchanged
		for i < len(lines) {
			nextLine := lines[i]
			if strings.HasPrefix(nextLine, ">") || strings.HasPrefix(nextLine, "> ") {
				result = append(result, nextLine)
				i++
			} else if strings.TrimSpace(nextLine) == "" {
				// Blank line ends the blockquote
				result = append(result, nextLine)
				i++
				break
			} else {
				break
			}
		}
	}

	return strings.Join(result, "\n")
}

func findAlertType(tag string) *AlertType {
	for i := range alertTypes {
		if alertTypes[i].Tag == tag {
			return &alertTypes[i]
		}
	}
	return nil
}
