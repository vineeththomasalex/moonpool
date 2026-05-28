# Moonpool Feature Gap Analysis

> What's missing compared to a "perfect markdown viewer", why, and how hard is it to fix?

## How to Read This Document

Each feature is rated on three dimensions:

- **Relevance** — How commonly is this feature used in real markdown files?
  - 🔴 Critical — appears in nearly every README/doc
  - 🟡 Common — appears in many docs, especially technical ones
  - 🟢 Niche — specialized or rarely used

- **Difficulty** — How hard is it to implement in moonpool?
  - ⬜ Trivial (< 1 day, few lines of code)
  - 🟦 Easy (1–2 days, existing libraries/patterns)
  - 🟨 Medium (3–5 days, new components needed)
  - 🟧 Hard (1–2 weeks, significant architecture work)
  - 🟥 Very Hard (weeks+, major dependencies or platform blockers)

- **Blocker** — What makes it hard?
  - ⚙️ Glamour limitation (unexported internals, needs fork or pre/post-processing)
  - 🖥️ Terminal limitation (character grid, no pixel rendering)
  - 📺 Windows Terminal limitation (no Sixel/Kitty graphics protocol)
  - 📦 External dependency (Node.js, headless browser, etc.)

---

## Currently Supported ✅

These features work today in moonpool — no gaps.

| Feature | Quality | Notes |
|---------|---------|-------|
| ATX Headings (H1–H6) | ✅ Excellent | Color + prefix per level, H1 gets background bar |
| Setext Headings | ✅ Good | Parsed by Goldmark, rendered same as ATX |
| Bold / Italic / Bold+Italic | ✅ Excellent | ANSI bold, italic, combined |
| Strikethrough `~~text~~` | ✅ Good | ANSI crossedout attribute |
| Inline code `` `code` `` | ✅ Excellent | Colored foreground on dark background |
| Fenced code blocks + syntax highlighting | ✅ Excellent | Chroma supports 200+ languages |
| Indented code blocks | ✅ Good | Basic rendering |
| Ordered / Unordered lists | ✅ Excellent | Bullet/number prefix, all markers (- * +) |
| Nested lists | ✅ Good | Proper indentation tracking |
| Block quotes | ✅ Good | `│` left border |
| Thematic break (horizontal rule) | ✅ Good | `────────` line |
| Links (inline + reference) | ✅ Good | Text shown + URL displayed |
| Autolinks | ✅ Good | Auto-detected by GFM |
| Task lists `- [x]` / `- [ ]` | ✅ Good | `[✓]` / `[ ]` checkboxes |
| Emoji shortcodes `:rocket:` | ✅ Good | Unicode substitution via goldmark-emoji |
| Front matter (YAML) | ✅ Good | Stripped before rendering |
| Tables (GFM) | ✅ Good | Unicode box borders, column alignment |
| Hard/soft line breaks | ✅ Good | Standard handling |
| Backslash escapes | ✅ Good | Standard handling |
| Entity references `&amp;` | ✅ Good | Unicode substitution |
| Pager (scroll, j/k, PgUp/PgDn) | ✅ Excellent | Full vim-style TUI pager |
| File browser / stash | ✅ Good | Fuzzy search, directory scanning |
| Multiple styles (dark/light/dracula/etc.) | ✅ Good | 6+ built-in + custom JSON |
| Remote URL / GitHub fetching | ✅ Good | `moonpool https://...` works |
| Word wrap control | ✅ Good | `-w` flag |
| Line numbers | ✅ Good | Config option |
| Mouse support | ✅ Good | Scroll wheel, configurable |
| Config file | ✅ Good | `glow.yml` |

---

## Missing Features

### Tier 1 — High Impact, Low Difficulty (Quick Wins)

These are commonly used features that can be implemented with minimal effort.

#### ~~1. GitHub Alerts / Admonitions~~ ✅ IMPLEMENTED
```markdown
> [!NOTE]
> This is a note callout.

> [!WARNING]
> This is a warning.
```

| Dimension | Rating |
|-----------|--------|
| Relevance | 🔴 Critical — used heavily in GitHub READMEs, docs, and wikis |
| Difficulty | 🟦 Easy |
| Blocker | ⚙️ Glamour renders these as plain blockquotes. Needs a Goldmark extension or pre-processing to detect `[!TYPE]` and apply color + icon |

**Current behavior:** Renders as a plain `│ [!NOTE]` blockquote with no color differentiation.
**Target behavior:** Color-coded left border + icon prefix (ℹ️ NOTE, ⚠️ WARNING, 💡 TIP, ❗ IMPORTANT, 🔴 CAUTION).
**Approach:** Pre-process markdown to detect `> [!TYPE]` pattern, or write a Goldmark AST transformer. Apply ANSI color based on alert type.
**Supported by:** GitHub ✅, VS Code ✅, Obsidian ✅, Typora ✅ — none of the terminal viewers support this yet, so moonpool would be first.

---

#### ~~2. OSC 8 Clickable Hyperlinks~~ ✅ IMPLEMENTED
```
Regular link: [click me](https://example.com)
Currently shows: click me (https://example.com)
Should emit: \e]8;;https://example.com\aclick me\e]8;;\a
```

| Dimension | Rating |
|-----------|--------|
| Relevance | 🔴 Critical — every doc has links |
| Difficulty | ⬜ Trivial |
| Blocker | None — Windows Terminal supports OSC 8 natively |

**Current behavior:** Links show as `text (URL)` in plain text.
**Target behavior:** Ctrl+clickable hyperlinks in Windows Terminal.
**Approach:** Post-process glamour output to wrap URLs in OSC 8 escape sequences. Glamour v2 may already support this — needs verification.
**Supported by:** mdcat ✅ — this is mdcat's key differentiator over glow.

---

#### ~~3. Live Reload / File Watching~~ ✅ ALREADY IMPLEMENTED
| Dimension | Rating |
|-----------|--------|
| Relevance | 🟡 Common — useful when editing markdown |
| Difficulty | ⬜ Trivial |
| Blocker | None — `fsnotify` is already in `go.mod` |

**Current behavior:** Static view. Must quit and re-open to see changes.
**Target behavior:** Auto-refresh when file changes on disk.
**Approach:** Wire `fsnotify` file watcher to pager re-render. The dependency is already present but unused.

---

#### 4. Clipboard Copy
| Dimension | Rating |
|-----------|--------|
| Relevance | 🟡 Common — useful for copying code snippets |
| Difficulty | ⬜ Trivial |
| Blocker | None — `atotto/clipboard` is already in `go.mod` |

**Current behavior:** Pager has OSC 52 copy but it's limited.
**Target behavior:** `y` keybinding copies current selection or visible content to system clipboard.
**Approach:** Wire `atotto/clipboard` in pager keybindings.

---

#### 5. Highlight / Mark `==text==`
```markdown
This is ==highlighted text== in a sentence.
```

| Dimension | Rating |
|-----------|--------|
| Relevance | 🟢 Niche — Obsidian/Typora users, not in GFM spec |
| Difficulty | 🟦 Easy |
| Blocker | ⚙️ Needs Goldmark extension + glamour style for background color |

**Current behavior:** `==text==` renders as literal `==text==`.
**Target behavior:** Yellow/colored background highlight via ANSI.

---

#### 6. Smart Typography (Typographer)
| Dimension | Rating |
|-----------|--------|
| Relevance | 🟢 Niche — nice to have |
| Difficulty | ⬜ Trivial |
| Blocker | None — Goldmark has `extension.Typographer` built-in |

**Current behavior:** `"hello"` stays as `"hello"`, `--` stays as `--`.
**Target behavior:** `"hello"` → `"hello"`, `--` → `–`, `---` → `—`.
**Approach:** Add `extension.Typographer` to Goldmark initialization.

---

### Tier 2 — High Impact, Medium Difficulty

#### 7. Footnotes
```markdown
Here is a statement[^1].

[^1]: This is the footnote content.
```

| Dimension | Rating |
|-----------|--------|
| Relevance | 🟡 Common — used in academic/technical docs, supported by GitHub |
| Difficulty | 🟨 Medium |
| Blocker | ⚙️ Goldmark has `extension.Footnote` but glamour's unexported `md` field blocks adding it without a fork or pre-processing |

**Current behavior:** `[^1]` renders as literal text.
**Target behavior:** Superscript number inline + numbered footnote list at document end.
**Approach:** Either fork glamour to add `WithFootnotes()` option (trivial 3-line change), or pre-process markdown to convert footnotes to a rendered format glamour can handle.

---

#### 8. In-Document Search
| Dimension | Rating |
|-----------|--------|
| Relevance | 🔴 Critical — expected in any pager/viewer |
| Difficulty | 🟨 Medium |
| Blocker | None — pure TUI feature |

**Current behavior:** No search. Stash has fuzzy file search, but the pager has no text search.
**Target behavior:** `/` to search, `n`/`N` to next/prev match, highlighted matches.
**Approach:** Bubble Tea pager enhancement. Search the rendered text, scroll viewport to match, apply ANSI highlight to matches.

---

#### 9. Table of Contents / Outline Navigation
| Dimension | Rating |
|-----------|--------|
| Relevance | 🟡 Common — useful for long documents |
| Difficulty | 🟨 Medium |
| Blocker | None — pure TUI feature |

**Current behavior:** No TOC. No way to jump to headings.
**Target behavior:** `o` or sidebar pane showing headings, Enter to jump. Or `[TOC]` placeholder expansion.
**Approach:** Parse headings from the markdown AST, build a jump list, integrate with pager viewport.

---

#### 10. Definition Lists
```markdown
Term
: Definition text here
```

| Dimension | Rating |
|-----------|--------|
| Relevance | 🟢 Niche — Pandoc/PHP Markdown Extra, not GFM |
| Difficulty | ⬜ Trivial |
| Blocker | ⚙️ Goldmark has `extension.DefinitionList` which glamour already enables, but moonpool may not be rendering it properly |

**Current behavior:** Goldmark extension is enabled in glamour. Should work — needs testing.
**Target behavior:** Bold term + indented definition with `🠶` prefix.

---

#### 11. Color Preview Swatches
```markdown
Use `#FF5733` for the accent color.
```

| Dimension | Rating |
|-----------|--------|
| Relevance | 🟢 Niche — useful in design/CSS docs |
| Difficulty | 🟦 Easy |
| Blocker | None — can post-process code spans |

**Current behavior:** `#FF5733` renders as plain inline code.
**Target behavior:** Inline code with an ANSI background color swatch matching the hex value.
**Approach:** Post-process glamour output — detect hex color patterns in code spans, prepend a colored `██` block.

---

### Tier 3 — High Impact, Hard (Platform Constrained)

#### 12. Mermaid Diagrams ⭐ PRIMARY GOAL
```markdown
```mermaid
graph LR
    A --> B --> C
```​
```

| Dimension | Rating |
|-----------|--------|
| Relevance | 🔴 Critical — used extensively in GitHub READMEs and technical docs |
| Difficulty | 🟥 Very Hard |
| Blocker | 📺 Windows Terminal has no graphics protocol (no Sixel, no Kitty). 📦 Mermaid requires Node.js/browser for rendering. ⚙️ Glamour has no extension hook. |

**Current behavior:** Renders as a syntax-highlighted code block showing the raw mermaid DSL.
**Target behavior (ideal):** Rendered diagram as an inline terminal image.
**Target behavior (realistic for WT):** One or more of:
- ASCII art approximation (no good Go library exists today)
- "Press Enter to open in browser" → launches mermaid.live with encoded diagram
- External render to PNG → open in default image viewer
- When running in a terminal WITH graphics support (WezTerm, Kitty), render inline via Sixel/Kitty

**Approach options:**
| Option | Pros | Cons |
|--------|------|------|
| Pre-process + `mmdc` → PNG → Sixel | High quality | Requires Node.js, won't display on WT |
| Pre-process + `mermaid-rs-renderer` (Rust) → PNG → Sixel | No Node.js | Still won't display on WT, needs CGO or subprocess |
| ASCII art via `mermaid-ascii` Go library | Works everywhere | Only supports flowcharts + sequence diagrams |
| Detect + launch browser | Zero dependencies | Not inline, disruptive UX |
| Detect + display styled placeholder | Simple | User can't see the diagram |

---

#### 13. Images (Inline Rendering)
```markdown
![Architecture diagram](./diagram.png)
```

| Dimension | Rating |
|-----------|--------|
| Relevance | 🔴 Critical — images are everywhere in docs |
| Difficulty | 🟥 Very Hard |
| Blocker | 📺 Windows Terminal does not support Sixel or Kitty Graphics Protocol |

**Current behavior:** Shows `Image: alt-text →` placeholder.
**Target behavior:** Inline pixel rendering (when terminal supports it), graceful fallback otherwise.
**Approach:** Detect terminal capabilities at runtime. If Kitty/Sixel/iTerm2 supported → render inline. For WT → show styled alt text placeholder + option to open in external viewer.
**Note:** WT has an open tracking issue for Sixel support but no timeline.

---

#### 14. Math / LaTeX
```markdown
Inline: $E = mc^2$
Display: $$\int_0^\infty e^{-x^2} dx = \frac{\sqrt{\pi}}{2}$$
```

| Dimension | Rating |
|-----------|--------|
| Relevance | 🟡 Common — critical for academic/scientific docs, supported by GitHub |
| Difficulty | 🟧 Hard (basic) / 🟥 Very Hard (full) |
| Blocker | 🖥️ No JavaScript runtime. Unicode math symbols cover only basic expressions. Complex expressions (fractions, matrices, integrals) have no character-grid representation. |

**Current behavior:** `$E = mc^2$` renders as literal text with dollar signs.
**Target behavior (basic):** Unicode substitution: `E = mc²`, Greek letters (α, β, γ), common operators (∑, ∫, ∂, ∞, ≤, ≥, ≠).
**Target behavior (full):** Would need image rendering (same blocker as images) or `pandoc --to ansi` integration.

**Phased approach:**
1. **Phase 1 (Easy):** Detect `$...$` blocks, apply basic Unicode substitution for superscripts, subscripts, Greek, and operators
2. **Phase 2 (Hard):** Shell out to `pandoc --to ansi` for full LaTeX rendering (requires pandoc installed)

---

#### 15. SVG Rendering
| Dimension | Rating |
|-----------|--------|
| Relevance | 🟡 Common — many GitHub repos use SVG badges and diagrams |
| Difficulty | 🟥 Very Hard |
| Blocker | 📺 Same as images — no graphics protocol in WT. No pure-Go SVG rasterizer. |

**Current behavior:** SVG images show alt text only.
**Approach:** Same as images — needs terminal graphics protocol support.

---

### Tier 4 — Lower Priority / Very Niche

| Feature | Relevance | Difficulty | Blocker | Notes |
|---------|-----------|------------|---------|-------|
| Subscript/Superscript `<sub>`/`<sup>` | 🟢 Niche | 🟦 Easy | ⚙️ HTML stripped | Unicode mapping for digits/common letters |
| WikiLinks `[[Page]]` | 🟢 Niche | 🟦 Easy | ⚙️ Goldmark extension | Obsidian-specific |
| Abbreviations `*[HTML]: ...` | 🟢 Niche | 🟨 Medium | ⚙️ Goldmark extension | PHP Markdown Extra feature |
| Custom containers `:::` | 🟢 Niche | 🟨 Medium | ⚙️ Goldmark extension | VuePress/MkDocs feature |
| Grid tables (colspan/rowspan) | 🟢 Niche | 🟧 Hard | 🖥️ Character grid | Pandoc-only feature |
| @Mentions highlighting | 🟢 Niche | ⬜ Trivial | Post-process | GitHub-specific |
| SHA/Issue autolinks `#123` | 🟢 Niche | 🟨 Medium | Needs repo context | GitHub-specific |
| GeoJSON/TopoJSON/STL | 🟢 Niche | 🟥 Very Hard | 🖥️ 📺 | GitHub-specific, impossible in terminal |
| Video/Audio embeds | 🟢 Niche | 🟥 Very Hard | 🖥️ | Impossible in terminal |
| Export to HTML/PDF | 🟡 Common | 🟨 Medium | None | glamour renders ANSI; would need HTML renderer |
| Print support | 🟢 Niche | 🟦 Easy | None | Strip ANSI, pipe to printer |
| Split view (source + preview) | 🟡 Common | 🟧 Hard | None | TUI layout work |
| Multiple tabs | 🟡 Common | 🟧 Hard | None | TUI architecture change |

---

## Nested Blockquote Bug

Worth noting: glamour has a known bug (issue #313) where nested blockquotes are not visually distinguished. Inner blockquotes render with the same style as outer ones. This affects existing "supported" functionality.

---

## Summary Matrix

| # | Feature | Relevance | Difficulty | Quick Win? |
|---|---------|-----------|------------|------------|
| ~~1~~ | ~~GitHub Alerts/Admonitions~~ | 🔴 Critical | 🟦 Easy | ✅ Done |
| ~~2~~ | ~~OSC 8 Clickable Links~~ | 🔴 Critical | ⬜ Trivial | ✅ Done |
| ~~3~~ | ~~Live Reload~~ | 🟡 Common | ⬜ Trivial | ✅ Already existed |
| 4 | Clipboard Copy | 🟡 Common | ⬜ Trivial | ✅ Yes |
| 5 | Highlight `==text==` | 🟢 Niche | 🟦 Easy | ✅ Yes |
| 6 | Smart Typography | 🟢 Niche | ⬜ Trivial | ✅ Yes |
| 7 | Footnotes | 🟡 Common | 🟨 Medium | |
| 8 | In-Document Search | 🔴 Critical | 🟨 Medium | |
| 9 | TOC / Outline Nav | 🟡 Common | 🟨 Medium | |
| 10 | Definition Lists | 🟢 Niche | ⬜ Trivial | ✅ Yes |
| 11 | Color Preview Swatches | 🟢 Niche | 🟦 Easy | ✅ Yes |
| 12 | **Mermaid Diagrams** | 🔴 Critical | 🟥 Very Hard | ⭐ Primary goal |
| 13 | Images (inline) | 🔴 Critical | 🟥 Very Hard | 📺 WT blocker |
| 14 | Math / LaTeX | 🟡 Common | 🟧–🟥 | Phased |
| 15 | SVG Rendering | 🟡 Common | 🟥 Very Hard | 📺 WT blocker |

### The 80/20
If you implement just **items 1, 2, 3, 8, and 12**, moonpool would cover ~80% of what users actually notice is missing compared to GitHub's markdown preview — and it would be the first terminal markdown viewer to support GitHub Alerts and Mermaid.
