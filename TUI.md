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

- `↑/↓` or `k/j`: Navigate within current selector (change selected value)
- `←/→` or `h/l`: Switch between selectors (Animation, Theme, File, Duration)
- `Enter`: Launch animation immediately (exits TUI and starts animation)
- `Esc` or `q`: Quit TUI
- `Ctrl+C`: Force quit

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

## Implementation Details

The TUI launches animations by:
1. Finding the `syscgo` binary (checks current dir, parent dir, PATH, system locations)
2. Building command with selected parameters
3. Executing `syscgo` with proper flags
4. Exiting TUI to return terminal control to the animation

## Future Enhancements

- [ ] Live animation preview in canvas
- [ ] Custom duration input
- [ ] Recent/favorite animations
- [ ] Animation descriptions/help text
- [ ] Save/load animation presets
- [ ] Search/filter animations
- [ ] Keybindings customization
