# BIT-like Editor Design Document

## Overview
Implement a FIGlet-style banner text editor in the TUI that generates ASCII art from text using .bit font files (JSON format).

## Architecture

### 1. Font System (`tui/bitfont.go`)

```go
type BitFont struct {
    Name       string
    Author     string
    License    string
    Height     int
    Characters map[rune][]string // char -> lines of ASCII art
}

func LoadBitFont(path string) (*BitFont, error)
func ListAvailableFonts() []string
func (f *BitFont) RenderText(text string) []string
```

### 2. Editor State Extensions (`tui/model.go`)

Add to Model struct:
```go
// BIT Editor Mode
bitEditorMode     bool
bitInputText      textinput.Model  // Single line text input
bitFonts          []string          // Available font names
bitSelectedFont   int
bitAlignment      int               // 0=left, 1=center, 2=right
bitColor          string            // Hex color
bitScale          float64           // 0.5, 1.0, 2.0, 3.0, 4.0
bitShadow         bool
bitCharSpacing    int               // Extra spaces between chars
bitLineSpacing    int               // Extra lines between rows
bitFocusedControl int               // Which control has focus
bitCurrentFont    *BitFont          // Loaded font
bitPreviewLines   []string          // Rendered output
```

### 3. Text Renderer (`tui/bitrender.go`)

```go
type RenderOptions struct {
    Font         *BitFont
    Text         string
    Alignment    int    // 0=left, 1=center, 2=right
    Color        string // Hex color
    Scale        float64
    Shadow       bool
    CharSpacing  int
    LineSpacing  int
}

func RenderBitText(opts RenderOptions) []string
func ApplyColor(lines []string, color string) string
func ApplyScale(lines []string, scale float64) []string
func ApplyShadow(lines []string) []string
func ApplyAlignment(lines []string, width int, align int) []string
```

### 4. UI Components (`tui/biteditor_view.go`)

```go
func (m Model) renderBitEditorView() string {
    // Layout:
    // 1. Preview canvas (top, ~60% height)
    // 2. Text input field
    // 3. Control row 1: Font selector, Alignment buttons, Color picker
    // 4. Control row 2: Scale, Shadow, Spacing
    // 5. Help text
}

func (m Model) renderFontSelector() string
func (m Model) renderAlignmentButtons() string
func (m Model) renderColorPicker() string
func (m Model) renderScaleSelector() string
func (m Model) renderShadowToggle() string
func (m Model) renderSpacingControls() string
```

### 5. Input Handling (`tui/biteditor_update.go`)

```go
func (m Model) handleBitEditorKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd)
func (m Model) updateBitPreview() Model // Regenerate preview when settings change
func (m Model) saveBitArt() (Model, tea.Cmd)
```

## Implementation Plan

### Phase 1: Font System (Core)
- [ ] Create `tui/bitfont.go` with JSON parser
- [ ] Add sample .bit fonts to `assets/fonts/` directory
- [ ] Implement basic text rendering (no scaling/shadow)
- [ ] Test font loading and rendering

### Phase 2: Editor UI Structure
- [ ] Add bit editor state to Model
- [ ] Create `tui/biteditor_view.go` with layout
- [ ] Implement text input field
- [ ] Create preview canvas area

### Phase 3: Font Browser
- [ ] Font dropdown selector with navigation
- [ ] Load font on selection
- [ ] Update preview when font changes

### Phase 4: Alignment Controls
- [ ] Left/Center/Right buttons
- [ ] Visual indication of active alignment
- [ ] Apply alignment to preview

### Phase 5: Color Selection
- [ ] Color picker with theme colors + hex input
- [ ] Apply color to preview using lipgloss
- [ ] Show current color visually

### Phase 6: Advanced Options
- [ ] Scale selector (0.5x - 4.0x)
- [ ] Shadow toggle and implementation
- [ ] Character/word/line spacing controls

### Phase 7: Integration
- [ ] Wire up Ctrl+S save flow
- [ ] Export to syscgo assets
- [ ] Add to main TUI navigation
- [ ] Testing and polish

## File Structure

```
tui/
├── bitfont.go         # Font loading/parsing
├── bitrender.go       # Text rendering engine
├── biteditor_view.go  # UI rendering
├── biteditor_update.go # Input handling
└── model.go           # State (extended)

assets/
└── fonts/
    ├── 3d-ascii.bit
    ├── banner.bit
    ├── block.bit
    └── ...
```

## Navigation Flow

```
Main TUI
  ↓ [Select "BIT Text Editor" from files]
  ↓
BIT Editor Mode
  ↓ [Edit text + adjust options]
  ↓ [Ctrl+S]
  ↓
Export Prompt
  ↓ [Choose syscgo]
  ↓
Save Prompt
  ↓ [Enter filename]
  ↓
Saved to assets/ + Return to Main TUI
```

## Sample .bit Font Format

```json
{
  "name": "3D ASCII",
  "author": "Unknown",
  "license": "Public Domain",
  "height": 6,
  "characters": {
    "A": [
      "   ▄████████ ",
      "  ███    ███ ",
      "  ███    ███ ",
      "  ███    ███ ",
      "▀███████████ ",
      "  ███    ███ "
    ],
    "B": [
      "▀█████████▄  ",
      "  ███    ███ ",
      "  ███    ███ ",
      "  ███    ███ ",
      "▄█████████▀  ",
      "▀█████████▀  "
    ]
  }
}
```

## Key Features

1. **Font Selection** - Browse .bit fonts from assets/fonts/
2. **Live Preview** - See styled output as you type
3. **Alignment** - Left/Center/Right text positioning
4. **Color** - Pick from themes or custom hex
5. **Scale** - Resize text (0.5x to 4.0x)
6. **Shadow** - Add drop shadow effect
7. **Spacing** - Adjust char/word/line spacing
8. **Export** - Save to assets for use in animations

## Dependencies

All existing:
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles` - Input components
- `github.com/charmbracelet/lipgloss` - Styling/colors
- `encoding/json` - Font file parsing

## Testing Strategy

1. Unit test font loading from JSON
2. Test text rendering with different fonts
3. Test each option (alignment, color, scale, shadow)
4. Integration test: full editor workflow
5. Manual TUI testing for UX polish
