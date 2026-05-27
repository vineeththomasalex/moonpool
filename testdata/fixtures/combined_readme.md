# Moonpool — A Terminal Markdown Viewer

> Render markdown on the CLI, with *pizzazz*! 🚀

## Features

- **Syntax Highlighting** — Code blocks rendered with colors
- **Tables** — Full table support
- **Task Lists** — Track your todos
- **Multiple Styles** — Dark, light, and custom themes

## Installation

```bash
go install github.com/example/moonpool@latest
```

## Quick Start

1. Open a markdown file:
   ```bash
   moonpool README.md
   ```
2. Pipe content:
   ```bash
   cat doc.md | moonpool -
   ```
3. Choose a style:
   ```bash
   moonpool -s dark README.md
   ```

## Configuration

| Option      | Default | Description             |
|-------------|---------|-------------------------|
| `style`     | `auto`  | Rendering style         |
| `width`     | `80`    | Output width            |
| `pager`     | `false` | Use external pager      |
| `mouse`     | `false` | Enable mouse support    |

## API Example

```go
package main

import (
    "fmt"
    "github.com/charmbracelet/glamour"
)

func main() {
    out, _ := glamour.Render("# Hello", "dark")
    fmt.Print(out)
}
```

## Links

- [GitHub](https://github.com/example/moonpool)
- [Documentation](https://example.com/docs)
- [Issue Tracker](https://github.com/example/moonpool/issues)

---

> **Note:** This is a test fixture combining multiple markdown features.

- [x] Headers ✅
- [x] Code blocks ✅
- [x] Tables ✅
- [x] Lists ✅
- [ ] Mermaid diagrams (coming soon)
