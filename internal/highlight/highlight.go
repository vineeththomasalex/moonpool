// Package highlight transforms ==highlighted text== syntax into
// styled markdown that glamour renders with visual distinction.
//
// Since glamour's Goldmark instance is unexported and can't be extended
// with a highlight plugin, we pre-process the markdown to replace
// ==text== with bold + emoji marker for a visible highlight effect.
package highlight

import (
	"regexp"
)

// highlightPattern matches ==text== where text doesn't span newlines.
var highlightPattern = regexp.MustCompile(`==((?:[^=\n]|=[^=])+)==`)

// Process replaces ==highlighted text== with a bolded + marked representation.
func Process(markdown string) string {
	return highlightPattern.ReplaceAllString(markdown, "**🟡 $1**")
}
