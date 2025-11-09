# SYSC-Go TUI

An interactive Terminal User Interface for browsing and launching sysc-go animations.

## Features

- **Visual Selection**: Browse animations, themes, and files with arrow keys
- **Live Preview**: See your selected configuration before launching
- **Dark Theme**: Beautiful dark-themed interface using Nord color palette
- **Keyboard Navigation**: Efficient vim-style navigation

## Usage

```bash
# Launch the TUI
./syscgo-tui
```

## Controls

### Main Selection Screen

- `↑/↓` or `k/j`: Navigate within current selector
- `←/→` or `h/l`: Switch between selectors (Animation, Theme, File, Duration)
- `Enter`: Show preview of selected configuration
- `Esc` or `q`: Quit

### Preview Screen

- `Enter`: Launch animation (exits TUI)
- `Esc`: Go back to selection screen
- `Ctrl+C` or `q`: Quit

## Selectors

1. **Animation**: Choose from 13 available animation effects
   - fire, matrix, matrix-art, rain, rain-art, fireworks, pour, print, beams, beam-text, ring-text, blackhole-text, aquarium

2. **Theme**: Select from 13 color themes
   - dracula, gruvbox, nord, tokyo-night, catppuccin, material, solarized, monochrome, transishardjob, rama, eldritch, dark, default

3. **File**: Choose ASCII art file from assets folder
   - Automatically discovers .txt files in assets/

4. **Duration**: Set animation duration
   - 5s, 10s, 30s, 60s, infinite

## Design

The TUI follows a similar design philosophy to [bit](https://github.com/superstarryeyes/bit) with:
- Large canvas area for content display (welcome screen / preview)
- Compact control area at bottom with focused selectors
- Clear visual feedback for current selection
- Minimal, distraction-free interface

## Architecture

```
cmd/syscgo-tui/main.go    # Entry point
tui/
  ├── model.go             # Bubbletea model and state
  ├── view.go              # Rendering logic
  ├── update.go            # Update/event handling
  ├── animation.go         # Animation preview/launch logic
  └── files.go             # Asset file discovery
```

## Future Enhancements

- [ ] Live animation preview in canvas (currently shows static config)
- [ ] Animation launching via os.Exec
- [ ] Custom duration input
- [ ] Recent/favorite animations
- [ ] Animation descriptions/help
- [ ] Save/load presets
