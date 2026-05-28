# AGENTS.md — Moonpool

Instructions for AI agents working on this codebase.

## Project Overview

Moonpool is a fork of [charmbracelet/glow](https://github.com/charmbracelet/glow) — a terminal-based markdown renderer. The goal is to add Mermaid diagram support and optimize for Windows Terminal.

- **Language**: Go
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Markdown Renderer**: [Glamour](https://github.com/charmbracelet/glamour) (wraps Goldmark)
- **Target Platform**: Windows Terminal (but code is cross-platform)

## Build & Test

```powershell
# Build
go build -o moonpool.exe .

# Run all tests (fast, offline, ~8s)
go test ./...

# Run with real WT screenshots (~100s)
go build -o moonpool.exe .
go test ./... -screenshots

# Update golden files after rendering changes
go test ./... -update

# Lint
golangci-lint run
```

Go is expected at `$env:TEMP\go-sdk` with GOPATH at `$env:TEMP\gopath`. Set these env vars per session:
```powershell
$env:GOROOT = "$env:TEMP\go-sdk"
$env:GOPATH = "$env:TEMP\gopath"
$env:PATH = "$env:GOROOT\bin;$env:GOPATH\bin;$env:PATH"
```

## Repository Structure

```
moonpool/
├── main.go                  # CLI entry point, flag parsing, render dispatch
├── config_cmd.go            # Config file management
├── console_windows.go       # Windows ANSI VT processing enablement
├── github.go                # GitHub repo → README URL resolution
├── gitlab.go                # GitLab repo → README URL resolution
├── style.go                 # Style definitions
├── url.go                   # URL parsing for remote markdown
├── log.go                   # Logging setup
├── man_cmd.go               # Man page generation
│
├── ui/                      # TUI (Bubble Tea) interface
│   ├── ui.go                # Main TUI model and program setup
│   ├── pager.go             # Pager view (scrollable markdown display)
│   ├── stash.go             # File browser / stash view
│   ├── stashhelp.go         # Help overlay for stash
│   ├── stashitem.go         # List item rendering
│   ├── markdown.go          # Document metadata type
│   ├── styles.go            # TUI style constants
│   ├── keys.go              # Keybinding definitions
│   ├── config.go            # TUI config
│   ├── editor.go            # External editor integration
│   ├── sort.go              # Sort options
│   └── ignore_*.go          # Platform-specific ignore paths
│
├── utils/
│   └── utils.go             # RemoveFrontmatter, IsMarkdownFile, WrapCodeBlock, GlamourStyle
│
├── internal/
│   └── screencap/           # Win32 window screenshot capture
│       └── screencap_windows.go
│
├── testdata/
│   ├── fixtures/            # Markdown test fixtures (one per feature)
│   ├── *.golden             # Golden files for regression detection
│   └── report.html          # Generated visual test report (gitignored)
│
├── render_test.go           # Rendering tests (content + golden + visual report)
├── render_report_test.go    # Test report infrastructure (TestMain, HTML generation)
├── glow_test.go             # Original flag tests
├── url_test.go              # URL parser tests (network, self-skipping)
└── utils/utils_test.go      # Utility function tests
```

## Rendering Pipeline

```
Source file/URL → read bytes
    → utils.RemoveFrontmatter()
    → utils.IsMarkdownFile() ? content : utils.WrapCodeBlock()
    → glamour.NewTermRenderer(style, width, ...)
    → renderer.Render(content) → ANSI string
    → stdout (pipe mode) or pager TUI (interactive mode)
```

Key function: `glamourRender()` in `ui/pager.go` (TUI path) and `executeCLI()` in `main.go` (pipe path).

## Testing Approach

See [TESTING.md](TESTING.md) for the full testing guide.

**Key points for agents:**
- All tests run offline — no network required
- Tests use `glamour.WithStyles(styles.DarkStyleConfig)` and `glamour.WithWordWrap(80)` for deterministic output
- Use `assertContains(t, ansiOut, "text")` which strips ANSI before comparing
- Golden files detect rendering regressions — run `go test ./... -update` after intentional changes
- Every rendering test adds an entry to `testdata/report.html` for visual verification

## Key Decisions

- **Windows Terminal only** (for now): We don't need to support iTerm, Kitty, or tmux
- **No external tools for testing**: Screenshots use Win32 APIs directly, no npm/pip dependencies
- **Glamour for rendering**: Don't bypass glamour — extend it or post-process its output
- **Golden files are source of truth**: If a golden file changes, that's a rendering change — review the diff

## What Not to Do

- Don't modify files outside the `moonpool/` directory
- Don't install packages globally — keep everything in `$env:TEMP`
- Don't push to GitHub without explicit user confirmation
- Don't add dependencies for testing that aren't Go packages
- Don't modify the upstream glow code unnecessarily — prefer additive changes

## Planned Work

- [ ] Mermaid diagram rendering (the primary goal)
- [ ] Strip unnecessary platform code (macOS, Linux-specific paths)
- [ ] Improve test coverage for TUI interactions (teatest)
