package main

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// reportEntry holds one test's rendered output for the HTML report.
type reportEntry struct {
	Name       string
	Pass       bool
	RawInput   string
	HTMLOutput string
}

var (
	reportMu      sync.Mutex
	reportEntries []reportEntry
)

// addReportEntry records a test's rendered output for inclusion in the HTML report.
func addReportEntry(name string, pass bool, rawInput, htmlOutput string) {
	reportMu.Lock()
	reportEntries = append(reportEntries, reportEntry{
		Name:       name,
		Pass:       pass,
		RawInput:   rawInput,
		HTMLOutput: htmlOutput,
	})
	reportMu.Unlock()
}

func TestMain(m *testing.M) {
	code := m.Run()
	writeHTMLReport()
	os.Exit(code)
}

func writeHTMLReport() {
	if len(reportEntries) == 0 {
		return
	}

	reportDir := filepath.Join("testdata")
	os.MkdirAll(reportDir, 0o755)
	reportPath := filepath.Join(reportDir, "report.html")

	f, err := os.Create(reportPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create report: %v\n", err)
		return
	}
	defer f.Close()

	passCount := 0
	failCount := 0
	for _, e := range reportEntries {
		if e.Pass {
			passCount++
		} else {
			failCount++
		}
	}

	fmt.Fprintf(f, `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Moonpool Render Test Report</title>
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: 'Segoe UI', system-ui, -apple-system, sans-serif;
    background: #0d1117;
    color: #e6edf3;
    padding: 24px;
    line-height: 1.6;
  }
  h1 {
    font-size: 1.8em;
    margin-bottom: 8px;
    color: #f0f6fc;
  }
  .summary {
    margin-bottom: 24px;
    padding: 16px;
    background: #161b22;
    border-radius: 8px;
    border: 1px solid #30363d;
  }
  .summary .pass { color: #3fb950; font-weight: bold; }
  .summary .fail { color: #f85149; font-weight: bold; }
  .test-card {
    margin: 16px 0;
    border: 1px solid #30363d;
    border-radius: 8px;
    overflow: hidden;
  }
  .test-header {
    padding: 12px 16px;
    font-weight: 600;
    font-size: 0.95em;
    cursor: pointer;
    user-select: none;
  }
  .test-header:hover { filter: brightness(1.1); }
  .test-pass .test-header { background: #1a2e1a; color: #3fb950; border-left: 4px solid #3fb950; }
  .test-fail .test-header { background: #2e1a1a; color: #f85149; border-left: 4px solid #f85149; }
  .test-body {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0;
    border-top: 1px solid #30363d;
  }
  .pane {
    padding: 16px;
    overflow-x: auto;
  }
  .pane-label {
    font-size: 0.75em;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: #8b949e;
    margin-bottom: 8px;
    font-weight: 600;
  }
  .source-pane {
    background: #161b22;
    border-right: 1px solid #30363d;
  }
  .source-pane pre {
    white-space: pre-wrap;
    word-wrap: break-word;
    font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace;
    font-size: 13px;
    color: #c9d1d9;
    line-height: 1.5;
  }
  .render-pane {
    background: #1e1e2e;
  }
  .render-pane .term-container {
    font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace;
    font-size: 13px;
    line-height: 1.5;
    white-space: pre-wrap;
    word-wrap: break-word;
  }
  .toggle-icon { float: right; transition: transform 0.2s; }
  .collapsed .test-body { display: none; }
  .collapsed .toggle-icon { transform: rotate(-90deg); }
</style>
</head>
<body>
<h1>📎 Moonpool Render Test Report</h1>
<div class="summary">
  <span class="pass">✅ %d passed</span> &nbsp; <span class="fail">❌ %d failed</span>
  &nbsp; | &nbsp; %d total test(s)
</div>
`, passCount, failCount, passCount+failCount)

	for i, e := range reportEntries {
		cls := "test-pass"
		icon := "✅"
		if !e.Pass {
			cls = "test-fail"
			icon = "❌"
		}
		fmt.Fprintf(f, `<div class="test-card %s" id="test-%d">
  <div class="test-header" onclick="this.parentElement.classList.toggle('collapsed')">
    %s %s <span class="toggle-icon">▼</span>
  </div>
  <div class="test-body">
    <div class="pane source-pane">
      <div class="pane-label">Markdown Source</div>
      <pre>%s</pre>
    </div>
    <div class="pane render-pane">
      <div class="pane-label">Rendered Output</div>
      <div class="term-container">%s</div>
    </div>
  </div>
</div>
`, cls, i, icon, html.EscapeString(e.Name), html.EscapeString(e.RawInput), e.HTMLOutput)
	}

	fmt.Fprintf(f, `</body></html>`)
	fmt.Fprintf(os.Stderr, "\n📎 Report written to %s\n", reportPath)
}
