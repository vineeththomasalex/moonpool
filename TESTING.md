# Testing Guide

## Quick Reference

```powershell
# Run all tests (fast, offline, ~8s)
go test ./...

# Run with verbose output
go test -v ./...

# Run with real Windows Terminal screenshots (~100s)
go build -o moonpool.exe .
go test ./... -screenshots

# Update golden files after intentional rendering changes
go test ./... -update

# View visual test report (generated after every test run)
start testdata\report.html
```

## Test Architecture

```
┌──────────────────────────────────────────────────────────────┐
│  Layer 1: Unit Tests (no TUI, no terminal, no network)       │
│  • Utility functions: frontmatter, file detection, code wrap │
│  • Glamour rendering: markdown → ANSI string assertions      │
│  • Golden file regression detection                          │
│  Run: go test ./...  (~8s, fully offline)                    │
├──────────────────────────────────────────────────────────────┤
│  Layer 2: Visual Report (ANSI → HTML or WT screenshots)      │
│  • buildkite/terminal-to-html converts ANSI to HTML          │
│  • internal/screencap captures real WT windows (-screenshots) │
│  • TestMain writes testdata/report.html after every run      │
│  • Developer opens report to spot-check visual quality       │
│  Run: go test ./... -screenshots  (~100s)                    │
└──────────────────────────────────────────────────────────────┘
```

## Test Coverage

### Utility Functions (`utils/utils_test.go`)

| Test | What it covers |
|------|---------------|
| `TestRemoveFrontmatter` | YAML frontmatter stripping: valid, missing, malformed, empty, CRLF |
| `TestIsMarkdownFile` | Extension detection: .md, .markdown, .mdown, .mkd, case sensitivity, no extension, non-markdown |
| `TestWrapCodeBlock` | Code fence wrapping: with language, without, empty content |
| `TestRemoveFrontmatterWithFixture` | End-to-end frontmatter stripping on a real fixture file |

### Markdown Rendering (`render_test.go`)

Each test renders a fixture through glamour (dark style, 80-char width) and asserts:
1. **Content assertions** — ANSI-stripped text contains expected strings
2. **Golden file regression** — exact ANSI output matches `testdata/*.golden`
3. **Visual report entry** — rendered output added to `testdata/report.html`

| Test | Fixture | What it verifies |
|------|---------|-----------------|
| `TestRenderHeadings` | `headings.md` | h1–h6, headings with bold/code/links |
| `TestRenderInlineFormatting` | `inline_formatting.md` | **bold**, *italic*, ~~strike~~, `code`, nesting |
| `TestRenderCodeBlocks` | `code_blocks.md` | Go/Python/JSON/HTML fenced blocks, plain fenced, indented |
| `TestRenderLists` | `lists.md` | Unordered (- * +), ordered, nested, mixed, multi-line items |
| `TestRenderLinksImages` | `links_images.md` | Inline links, reference links, autolinks, images (alt text) |
| `TestRenderTables` | `tables.md` | Simple, aligned, inline formatting in cells, wide tables |
| `TestRenderBlockquotes` | `blockquotes.md` | Simple, nested (3 levels), blockquotes with lists/code/headings |
| `TestRenderHorizontalRules` | `horizontal_rules.md` | `---`, `***`, `___` separators |
| `TestRenderTaskLists` | `task_lists.md` | `- [x]` checked, `- [ ]` unchecked, mixed with regular items |
| `TestRenderEmoji` | `emoji.md` | Shortcodes (:rocket:), Unicode passthrough (🚀 📎) |
| `TestRenderFrontmatter` | `frontmatter.md` | Frontmatter present in raw render, body content intact |
| `TestRenderCombinedReadme` | `combined_readme.md` | Realistic README with multiple features combined |

### Negative Tests — Unsupported/Edge Cases

| Test | What it verifies |
|------|-----------------|
| `TestNegativeMermaid` | ` ```mermaid ` blocks render as code, no crash |
| `TestNegativeMath` | `$E=mc^2$` and `$$...$$` appear as text, no crash |
| `TestNegativeHTML` | `<div>`, `<details>`, `<table>` HTML tags handled gracefully |
| `TestNegativeFootnotes` | `[^1]` footnote syntax doesn't crash |
| `TestNegativeEmptyContent` | Empty string, whitespace, single char, frontmatter-only |
| `TestNegativeLongContent` | 200-section document renders without error |
| `TestNegativeLongLine` | 1000-char line wraps correctly |

### Style & Width Variations

| Test | What it verifies |
|------|-----------------|
| `TestRenderDarkVsLight` | Dark and light styles produce different output |
| `TestRenderWidthVariations` | Width 40/80/120 produce different wrapping |

### Screencap (`internal/screencap/`)

| Test | What it verifies |
|------|-----------------|
| `TestCaptureCommand` | Win32 window capture: launches moonpool in WT, captures, saves PNG |

## Golden Files

Golden files in `testdata/*.golden` store the exact ANSI output of each rendering test. They act as regression detectors — if the rendering changes, the tests fail with a diff.

**Update golden files** after intentional changes:
```powershell
go test ./... -update
```

**Review the diff** before committing updated golden files:
```powershell
git diff testdata/*.golden
```

## Visual Test Report

Every test run generates `testdata/report.html` — an HTML file showing:
- **Left pane**: Markdown source
- **Right pane**: Rendered output (HTML or real WT screenshot)
- Pass/fail status for each test
- Collapsible test cards

### Without `-screenshots` (default)
Uses `buildkite/terminal-to-html` to convert ANSI escape codes to colored HTML. Fast (~8s) but the rendering approximates what WT shows.

### With `-screenshots`
Uses `internal/screencap` to open a real Windows Terminal window for each test, capture the screen via Win32 APIs, and embed the PNG. Slower (~100s) but shows exactly what a user would see.

**Prerequisites for screenshots:**
```powershell
go build -o moonpool.exe .    # screenshots run the built binary
```

## Adding New Tests

1. **Create a fixture** in `testdata/fixtures/your_feature.md`
2. **Add a test function** in `render_test.go`:
   ```go
   func TestRenderYourFeature(t *testing.T) {
       raw, ansiOut, htmlOut := renderFixture(t, "your_feature.md")
       addReportEntryWithFixture(t.Name(), true, raw, htmlOut, "your_feature.md")

       assertContains(t, ansiOut, "expected text")
       golden.RequireEqual(t, []byte(ansiOut))
   }
   ```
3. **Generate golden file**: `go test ./... -update`
4. **Review**: open `testdata/report.html` and verify visual output
5. **Commit** the fixture, test, and golden file together

## Dependencies (test only)

| Package | Purpose |
|---------|---------|
| `github.com/charmbracelet/x/exp/golden` | Golden file assertions |
| `github.com/buildkite/terminal-to-html/v3` | ANSI → HTML for report |
| `internal/screencap` | Real WT window capture (Win32) |
