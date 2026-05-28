// Package links adds OSC 8 hyperlink escape sequences to rendered
// ANSI output. OSC 8 makes URLs clickable in supported terminals
// (Windows Terminal, iTerm2, Kitty, WezTerm, GNOME Terminal).
//
// OSC 8 format:  \x1b]8;;URL\x1b\\  display text  \x1b]8;;\x1b\\
//
// This post-processes glamour's ANSI output to wrap bare URLs with
// OSC 8 sequences, making them Ctrl+clickable in the terminal.
package links

import (
	"regexp"
	"strings"
)

// urlPattern matches http/https URLs in rendered text.
// It handles URLs that may be surrounded by ANSI escape sequences.
var urlPattern = regexp.MustCompile(`https?://[^\s\x1b\)]+`)

// osc8Open wraps a URL in an OSC 8 opening sequence.
func osc8Open(url string) string {
	return "\x1b]8;;" + url + "\x1b\\"
}

// osc8Close is the OSC 8 closing sequence.
const osc8Close = "\x1b]8;;\x1b\\"

// AddOSC8Links post-processes ANSI-rendered text to wrap URLs with
// OSC 8 hyperlink sequences. URLs that are already inside OSC 8
// sequences are not double-wrapped.
func AddOSC8Links(ansiText string) string {
	// Skip if the text already contains OSC 8 sequences
	if strings.Contains(ansiText, "\x1b]8;;") {
		return ansiText
	}

	return urlPattern.ReplaceAllStringFunc(ansiText, func(url string) string {
		// Clean trailing punctuation that's likely not part of the URL
		cleaned := strings.TrimRight(url, ".,;:!?)]}>\"'")
		trailing := url[len(cleaned):]
		return osc8Open(cleaned) + cleaned + osc8Close + trailing
	})
}
