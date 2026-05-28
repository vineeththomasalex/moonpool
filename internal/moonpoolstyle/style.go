// Package moonpoolstyle provides the custom moonpool rendering style.
// Based on glamour's dark style with improved heading rendering:
// - No raw ## prefix — clean heading text
// - Blue gradient: bright blue for H1 → light/dim blue for H6
// - Underline decoration for H1 and H2
// - Increasing indentation for deeper heading levels
package moonpoolstyle

import (
	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/glamour/styles"
)

func boolPtr(b bool) *bool       { return &b }
func stringPtr(s string) *string { return &s }
func uintPtr(u uint) *uint       { return &u }

// Config returns the moonpool style config — dark base with custom headings.
func Config() ansi.StyleConfig {
	s := styles.DarkStyleConfig

	// H1: Bright blue, bold, underlined
	s.H1 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix:    " ",
			Suffix:    " ",
			Color:     stringPtr("#5fafff"),
			Bold:      boolPtr(true),
			Underline: boolPtr(true),
		},
		Indent: uintPtr(1),
	}

	// H2: Medium blue, bold, underlined
	s.H2 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix:    " ",
			Suffix:    " ",
			Color:     stringPtr("#5f87d7"),
			Bold:      boolPtr(true),
			Underline: boolPtr(true),
		},
		Indent: uintPtr(2),
	}

	// H3: Slate blue, bold
	s.H3 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: " ",
			Color:  stringPtr("#5f87af"),
			Bold:   boolPtr(true),
		},
		Indent: uintPtr(3),
	}

	// H4: Steel blue, bold
	s.H4 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: " ",
			Color:  stringPtr("#87afd7"),
			Bold:   boolPtr(true),
		},
		Indent: uintPtr(4),
	}

	// H5: Light blue
	s.H5 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: " ",
			Color:  stringPtr("#87afdf"),
		},
		Indent: uintPtr(5),
	}

	// H6: Dim light blue
	s.H6 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: " ",
			Color:  stringPtr("#87afdf"),
			Faint:  boolPtr(true),
		},
		Indent: uintPtr(6),
	}

	return s
}

// NoTTYConfig returns a style for non-terminal (piped) output with clean headings.
func NoTTYConfig() ansi.StyleConfig {
	s := styles.NoTTYStyleConfig

	s.H1 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{Prefix: " ", Bold: boolPtr(true), Upper: boolPtr(true)},
	}
	s.H2 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{Prefix: " ", Bold: boolPtr(true)},
		Indent:         uintPtr(1),
	}
	s.H3 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{Prefix: " ", Bold: boolPtr(true)},
		Indent:         uintPtr(2),
	}
	s.H4 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{Prefix: " "},
		Indent:         uintPtr(3),
	}
	s.H5 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{Prefix: " "},
		Indent:         uintPtr(4),
	}
	s.H6 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{Prefix: " ", Faint: boolPtr(true)},
		Indent:         uintPtr(5),
	}

	return s
}
