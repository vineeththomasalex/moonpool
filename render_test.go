package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	terminal "github.com/buildkite/terminal-to-html/v3"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/styles"
	"github.com/charmbracelet/x/exp/golden"
)

// renderFixture loads a fixture file and renders it through glamour.
// Returns the raw markdown, rendered ANSI, and rendered HTML.
func renderFixture(t *testing.T, name string) (raw, ansiOut, htmlOut string) {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("testdata", "fixtures", name))
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", name, err)
	}
	raw = string(data)

	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(styles.DarkStyleConfig),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	ansiOut, err = r.Render(raw)
	if err != nil {
		t.Fatalf("failed to render %s: %v", name, err)
	}

	htmlOut = string(terminal.Render([]byte(ansiOut)))
	return raw, ansiOut, htmlOut
}

// ansiEscape matches ANSI escape sequences including CSI, OSC, and single-byte escapes.
var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]|\x1b\][^\x07]*\x07|\x1b[^[\]a-zA-Z]?`)

// stripANSI removes ANSI escape sequences for plain text assertions.
func stripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

// assertContains checks that ANSI-stripped output contains the expected text.
func assertContains(t *testing.T, output, expected string) {
	t.Helper()
	plain := stripANSI(output)
	if !strings.Contains(plain, expected) {
		t.Errorf("expected output to contain %q, but it did not.\nPlain output (first 500 chars): %.500s", expected, plain)
	}
}

// assertNotContains checks that output does NOT contain the text.
func assertNotContains(t *testing.T, output, unexpected string) {
	t.Helper()
	if strings.Contains(output, unexpected) {
		t.Errorf("expected output NOT to contain %q, but it did", unexpected)
	}
}

// --- Rendering Tests: Markdown Features ---

func TestRenderHeadings(t *testing.T) {
	raw, ansiOut, htmlOut := renderFixture(t, "headings.md")
addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "headings.md")

	assertContains(t, ansiOut, "Heading Level 1")
	assertContains(t, ansiOut, "Heading Level 2")
	assertContains(t, ansiOut, "Heading Level 3")
	assertContains(t, ansiOut, "Heading Level 4")
	assertContains(t, ansiOut, "Heading Level 5")
	assertContains(t, ansiOut, "Heading Level 6")
	assertContains(t, ansiOut, "bold")
	assertContains(t, ansiOut, "inline code")

	golden.RequireEqual(t, []byte(ansiOut))
}

func TestRenderInlineFormatting(t *testing.T) {
	raw, ansiOut, htmlOut := renderFixture(t, "inline_formatting.md")
addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "inline_formatting.md")

	assertContains(t, ansiOut, "bold text")
	assertContains(t, ansiOut, "italic text")
	assertContains(t, ansiOut, "strikethrough text")
	assertContains(t, ansiOut, "inline code")
	assertContains(t, ansiOut, "bold and italic")

	golden.RequireEqual(t, []byte(ansiOut))
}

func TestRenderCodeBlocks(t *testing.T) {
	raw, ansiOut, htmlOut := renderFixture(t, "code_blocks.md")
addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "code_blocks.md")

	assertContains(t, ansiOut, "fmt.Println")
	assertContains(t, ansiOut, "Hello, World!")
	assertContains(t, ansiOut, "def greet")
	assertContains(t, ansiOut, "moonpool")
	assertContains(t, ansiOut, "plain code block")

	golden.RequireEqual(t, []byte(ansiOut))
}

func TestRenderLists(t *testing.T) {
	raw, ansiOut, htmlOut := renderFixture(t, "lists.md")
addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "lists.md")

	assertContains(t, ansiOut, "Item one")
	assertContains(t, ansiOut, "First item")
	assertContains(t, ansiOut, "Parent item")
	assertContains(t, ansiOut, "Child item")

	golden.RequireEqual(t, []byte(ansiOut))
}

func TestRenderLinksImages(t *testing.T) {
	raw, ansiOut, htmlOut := renderFixture(t, "links_images.md")
addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "links_images.md")

	assertContains(t, ansiOut, "a link to example")
	assertContains(t, ansiOut, "reference link")
	assertContains(t, ansiOut, "Alt text for image")

	golden.RequireEqual(t, []byte(ansiOut))
}

func TestRenderTables(t *testing.T) {
	raw, ansiOut, htmlOut := renderFixture(t, "tables.md")
addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "tables.md")

	assertContains(t, ansiOut, "Alice")
	assertContains(t, ansiOut, "Bob")
	assertContains(t, ansiOut, "Charlie")
	assertContains(t, ansiOut, "New York")

	golden.RequireEqual(t, []byte(ansiOut))
}

func TestRenderBlockquotes(t *testing.T) {
	raw, ansiOut, htmlOut := renderFixture(t, "blockquotes.md")
addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "blockquotes.md")

	assertContains(t, ansiOut, "This is a blockquote")
	assertContains(t, ansiOut, "Level one")

	golden.RequireEqual(t, []byte(ansiOut))
}

func TestRenderHorizontalRules(t *testing.T) {
	raw, ansiOut, htmlOut := renderFixture(t, "horizontal_rules.md")
addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "horizontal_rules.md")

	assertContains(t, ansiOut, "Content above the rule")
	assertContains(t, ansiOut, "Content between rules")
	assertContains(t, ansiOut, "Final content")

	golden.RequireEqual(t, []byte(ansiOut))
}

func TestRenderTaskLists(t *testing.T) {
	raw, ansiOut, htmlOut := renderFixture(t, "task_lists.md")
addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "task_lists.md")

	assertContains(t, ansiOut, "Completed task")
	assertContains(t, ansiOut, "Incomplete task")

	golden.RequireEqual(t, []byte(ansiOut))
}

func TestRenderEmoji(t *testing.T) {
	raw, ansiOut, htmlOut := renderFixture(t, "emoji.md")
addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "emoji.md")

	assertContains(t, ansiOut, "🚀")
	assertContains(t, ansiOut, "📎")

	golden.RequireEqual(t, []byte(ansiOut))
}

func TestRenderFrontmatter(t *testing.T) {
	raw, ansiOut, htmlOut := renderFixture(t, "frontmatter.md")
addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "frontmatter.md")

	// Frontmatter is NOT stripped by glamour — it's stripped by utils.RemoveFrontmatter
	// before rendering. Here we test glamour's raw handling.
	assertContains(t, ansiOut, "Document After Frontmatter")
	assertContains(t, ansiOut, "Section Two")

	golden.RequireEqual(t, []byte(ansiOut))
}

func TestRenderCombinedReadme(t *testing.T) {
	raw, ansiOut, htmlOut := renderFixture(t, "combined_readme.md")
addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "combined_readme.md")

	assertContains(t, ansiOut, "Moonpool")
	assertContains(t, ansiOut, "Syntax Highlighting")
	assertContains(t, ansiOut, "go install")
	assertContains(t, ansiOut, "style")

	golden.RequireEqual(t, []byte(ansiOut))
}

// --- Negative Tests: Unsupported Features ---

func TestNegativeMermaid(t *testing.T) {
	raw, ansiOut, htmlOut := renderFixture(t, "negative_mermaid.md")
addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "negative_mermaid.md")

	// Mermaid should render as a code block, not crash
	assertContains(t, ansiOut, "graph LR")
	assertContains(t, ansiOut, "Decision")
	assertContains(t, ansiOut, "should render as a plain code block")
}

func TestNegativeMath(t *testing.T) {
	raw, ansiOut, htmlOut := renderFixture(t, "negative_math.md")
addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "negative_math.md")

	// Math should appear as text, not crash
	assertContains(t, ansiOut, "E = mc")
	assertContains(t, ansiOut, "should not crash")
}

func TestNegativeHTML(t *testing.T) {
	raw, ansiOut, htmlOut := renderFixture(t, "negative_html.md")
addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "negative_html.md")

	// HTML should be handled gracefully
	assertContains(t, ansiOut, "Raw HTML Elements")
}

func TestNegativeFootnotes(t *testing.T) {
	raw, ansiOut, htmlOut := renderFixture(t, "negative_footnotes.md")
addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "negative_footnotes.md")

	// Footnotes should not crash
	assertContains(t, ansiOut, "footnote")
}

func TestNegativeEmptyContent(t *testing.T) {
	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(styles.DarkStyleConfig),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"whitespace only", "   \n\n   "},
		{"single character", "x"},
		{"only frontmatter", "---\ntitle: test\n---\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := r.Render(tc.input)
			if err != nil {
				t.Fatalf("Render(%q) failed: %v", tc.input, err)
			}
			htmlOut := string(terminal.Render([]byte(out)))
			addReportEntry(t.Name(), true, tc.input, htmlOut)
		})
	}
}

func TestNegativeLongContent(t *testing.T) {
	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(styles.DarkStyleConfig),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	// Generate a long document
	var sb strings.Builder
	sb.WriteString("# Long Document\n\n")
	for i := range 200 {
		sb.WriteString(fmt.Sprintf("## Section %d\n\nParagraph %d with some content.\n\n", i+1, i+1))
	}
	longInput := sb.String()

	out, err := r.Render(longInput)
	if err != nil {
		t.Fatalf("failed to render long content: %v", err)
	}

	assertContains(t, out, "Section 1")
	assertContains(t, out, "Section 200")

	htmlOut := string(terminal.Render([]byte(out)))
	addReportEntry(t.Name(), true, "# Long Document\n\n## Section 1 ... ## Section 200\n(200 sections, truncated for report)", htmlOut)
}

func TestNegativeLongLine(t *testing.T) {
	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(styles.DarkStyleConfig),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	longLine := "# Heading\n\n" + strings.Repeat("word ", 200)
	out, err := r.Render(longLine)
	if err != nil {
		t.Fatalf("failed to render long line: %v", err)
	}

	assertContains(t, out, "word")
	htmlOut := string(terminal.Render([]byte(out)))
	addReportEntry(t.Name(), true, longLine[:100]+"... (1000 chars)", htmlOut)
}

// --- Style Variation Tests ---

func TestRenderDarkVsLight(t *testing.T) {
	input := "# Hello World\n\nThis is **bold** and *italic* text.\n\n```go\nfmt.Println(\"hi\")\n```\n"

	darkR, _ := glamour.NewTermRenderer(glamour.WithStyles(styles.DarkStyleConfig), glamour.WithWordWrap(80))
	lightR, _ := glamour.NewTermRenderer(glamour.WithStyles(styles.LightStyleConfig), glamour.WithWordWrap(80))

	darkOut, err := darkR.Render(input)
	if err != nil {
		t.Fatalf("dark render failed: %v", err)
	}
	lightOut, err := lightR.Render(input)
	if err != nil {
		t.Fatalf("light render failed: %v", err)
	}

	// Both should contain the text
	assertContains(t, darkOut, "Hello World")
	assertContains(t, lightOut, "Hello World")

	// They should be different (different ANSI codes)
	if darkOut == lightOut {
		t.Error("dark and light styles produced identical output")
	}

	darkHTML := string(terminal.Render([]byte(darkOut)))
	lightHTML := string(terminal.Render([]byte(lightOut)))
	addReportEntryWithCommand(t.Name()+" (dark)", true, input, darkHTML, fmt.Sprintf("%s -s dark -w 80 testdata/fixtures/combined_readme.md", moonpoolPath))
	addReportEntryWithCommand(t.Name()+" (light)", true, input, lightHTML, fmt.Sprintf("%s -s light -w 80 testdata/fixtures/combined_readme.md", moonpoolPath))
}

// --- Width Tests ---

func TestRenderWidthVariations(t *testing.T) {
	input := "# Width Test\n\nThis is a paragraph with enough text that it should wrap differently at different widths. The quick brown fox jumps over the lazy dog.\n"

	widths := []int{40, 80, 120}
	for _, w := range widths {
		t.Run(fmt.Sprintf("width_%d", w), func(t *testing.T) {
			r, _ := glamour.NewTermRenderer(glamour.WithStyles(styles.DarkStyleConfig), glamour.WithWordWrap(w))
			out, err := r.Render(input)
			if err != nil {
				t.Fatalf("render at width %d failed: %v", w, err)
			}
			assertContains(t, out, "Width Test")
			htmlOut := string(terminal.Render([]byte(out)))
			addReportEntry(fmt.Sprintf("%s (width=%d)", t.Name(), w), true, input, htmlOut)
		})
	}
}
