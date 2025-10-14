# Visual Improvements for Mozzy 🎨

Inspired by VHS terminal recordings and modern CLI tools, here are ideas to make mozzy's output even more beautiful and user-friendly.

## 1. Enhanced Help Screen with ASCII Art Banner

### Current:
```
Usage:
  mozzy [command]
```

### Improved:
```
╭─────────────────────────────────────────────────────────╮
│                                                         │
│  🦣  MOZZY - Postman for your Terminal                 │
│      Beautiful HTTP client with superpowers            │
│                                                         │
╰─────────────────────────────────────────────────────────╯

Usage:
  mozzy [command]

Quick Start:
  mozzy GET https://api.example.com --color
  mozzy POST /users --json '{"name":"Alice"}'
  mozzy save my-request GET /api
```

## 2. Better Fonts & Icons

### Use Unicode Box Drawing & Emojis
```
┌─ Request ─────────────────────────────────────────┐
│ GET https://api.example.com/users/1               │
│ ⏱️  Response Time: 123ms                          │
│ 📊 Status: 200 OK                                 │
└───────────────────────────────────────────────────┘
```

### Response Headers Section
```
╭─ Response Headers ────────────────────────────────╮
│ Content-Type: application/json                    │
│ Cache-Control: max-age=3600                       │
│ X-RateLimit-Remaining: 999                        │
╰───────────────────────────────────────────────────╯
```

## 3. Color Schemes & Themes

Add theme support with preset color schemes:

```bash
mozzy GET /api --theme dracula
mozzy GET /api --theme monokai
mozzy GET /api --theme nord
mozzy GET /api --theme solarized
```

**Implementation:**
- Store themes in `~/.mozzy/themes/`
- JSON format for custom themes
- Color codes for keys, values, numbers, booleans, nulls

## 4. Progress Bars & Spinners

### Current Upload:
```
Uploading... done
```

### Improved:
```
┌─ Uploading files ─────────────────────────────────┐
│ ⣾ product-image.jpg                              │
│ ████████████████░░░░  80% (4.2MB / 5.2MB)        │
│ ⏱️  2.3s elapsed • 400KB/s • ~1s remaining         │
└───────────────────────────────────────────────────┘
```

## 5. Table Output for Collections

### Current `mozzy list`:
```
my-api GET https://api.example.com/users
  Get all users

test GET https://example.com/test
  Test endpoint
```

### Improved:
```
╭──────────────┬────────┬──────────────────────────────┬─────────────────╮
│ Name         │ Method │ URL                          │ Description     │
├──────────────┼────────┼──────────────────────────────┼─────────────────┤
│ my-api       │ GET    │ https://api.example.com/...  │ Get all users   │
│ test         │ GET    │ https://example.com/test     │ Test endpoint   │
│ create-user  │ POST   │ https://api.example.com/...  │ Create new user │
╰──────────────┴────────┴──────────────────────────────┴─────────────────╯

💡 Tip: Run 'mozzy exec <name>' to execute a saved request
```

## 6. Syntax Highlighting for JSON

### Enhanced JSON with better visual hierarchy:
```json
{
  "user": {                    ← Dimmed brackets
    "id": 1,                   ← Cyan key, Yellow number
    "name": "Alice",           ← Green string
    "active": true,            ← Magenta boolean
    "role": null,              ← Red null
    "email": "alice@ex.com"    ← Green + underline for URLs/emails
  }
}
```

## 7. Success/Error Banners

### Success:
```
╭───────────────────────────────────────────────────╮
│  ✅ SUCCESS                                       │
│  Status: 200 OK • Time: 234ms                    │
╰───────────────────────────────────────────────────╯
```

### Error:
```
╭───────────────────────────────────────────────────╮
│  ❌ ERROR                                         │
│  Status: 404 Not Found                           │
│  URL: https://api.example.com/missing            │
╰───────────────────────────────────────────────────╯
```

## 8. Workflow Visualization

### Running `mozzy run workflow.yaml`:
```
🔄 Running Workflow: User Onboarding

┌─ Step 1/4: Create user ──────────────────────────┐
│ POST https://api.example.com/register            │
│ ✓ Complete (201) • 456ms                         │
│ Captured: userId = 12345                         │
└───────────────────────────────────────────────────┘

┌─ Step 2/4: Verify email ─────────────────────────┐
│ POST https://api.example.com/verify              │
│ ✓ Complete (200) • 234ms                         │
└───────────────────────────────────────────────────┘

✨ Workflow completed successfully in 1.2s
```

## 9. Help Screen Categories

### Improved command grouping:
```
╭─ HTTP Commands ──────────────────────────────────╮
│  GET, POST, PUT, PATCH, DELETE                   │
│  Send HTTP requests with various methods         │
╰──────────────────────────────────────────────────╯

╭─ Collection Management ──────────────────────────╮
│  save     Save request to collection             │
│  list     List saved requests                    │
│  exec     Execute saved request                  │
╰──────────────────────────────────────────────────╯

╭─ Advanced Features ──────────────────────────────╮
│  run      Execute YAML workflows                 │
│  test     Run workflow as test suite             │
│  jwt      JWT decode/verify/sign                 │
│  diff     Compare JSON responses                 │
╰──────────────────────────────────────────────────╯
```

## 10. Interactive Mode

Add an interactive TUI (Terminal User Interface):

```bash
mozzy interactive
# or
mozzy -i
```

Features:
- Navigate collections with arrow keys
- Edit requests in-place
- View response in split pane
- Syntax highlighting
- History navigation
- Search functionality

**Libraries to use:**
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/charmbracelet/bubbles` - Components (table, list, etc.)

## 11. Request/Response Diff View

### Side-by-side comparison:
```
╭─ Request ────────────────╮  ╭─ Response ───────────────╮
│ POST /api/users          │  │ 201 Created              │
│                          │  │                          │
│ {                        │  │ {                        │
│   "name": "Alice"        │  │   "id": 123,             │
│   "email": "a@ex.com"    │  │   "name": "Alice",       │
│ }                        │  │   "email": "a@ex.com",   │
│                          │  │   "created": "2025-..."  │
│                          │  │ }                        │
╰──────────────────────────╯  ╰──────────────────────────╯
```

## 12. Smart Defaults Based on TTY

Auto-detect terminal capabilities:
- Check if terminal supports colors (TERM variable)
- Detect terminal width for responsive layouts
- Adjust output based on pipe vs interactive
- Unicode support detection

## Implementation Priority

### High Priority (Quick Wins):
1. ✅ Box drawing for help screen
2. ✅ Better icons and emojis
3. ✅ Enhanced error messages
4. ✅ Table output for collections

### Medium Priority:
5. Theme support
6. Progress bar improvements
7. Success/error banners
8. Workflow visualization

### Low Priority (Nice to Have):
9. Interactive TUI mode
10. Diff view
11. Custom fonts
12. Animation effects

## Libraries to Consider

```go
// Styling & Colors
"github.com/charmbracelet/lipgloss"

// TUI Framework
"github.com/charmbracelet/bubbletea"

// Components (tables, lists, progress bars)
"github.com/charmbracelet/bubbles"

// Syntax highlighting
"github.com/alecthomas/chroma"

// Box drawing
"github.com/jedib0t/go-pretty/v6/table"
"github.com/jedib0t/go-pretty/v6/progress"

// Spinners
"github.com/briandowns/spinner"
```

## Examples from Popular CLIs

### gh (GitHub CLI):
- Clean table output
- Color-coded status
- Interactive prompts

### httpie:
- Beautiful syntax highlighting
- Clear request/response separation
- Smart defaults

### lazygit:
- Full TUI interface
- Keyboard shortcuts
- Split panes

### k9s:
- Live updating
- Color themes
- Resource management

## Next Steps

1. Create `pkg/ui/` package for UI components
2. Add lipgloss for styling
3. Implement box-drawing help screen
4. Add table output for collections
5. Create theme system
6. Add progress bars for uploads/downloads
7. Consider TUI mode for v2.0
