# Export System Architecture

## Overview
The BIT editor can export ASCII art to two targets:
1. **syscgo** - Local assets for animations
2. **sysc-walls** - Wallpaper/screensaver system

## Export Targets

### 1. syscgo Export (Simple)

**Path**: `./assets/<filename>.txt`

**Process**:
1. User creates ASCII art in BIT editor
2. Ctrl+S → Select "syscgo"
3. Enter filename
4. File saved to `assets/<filename>.txt`
5. Immediately available in TUI file selector

**No config modification needed** - syscgo auto-discovers files in assets/

### 2. sysc-walls Export (Advanced)

**Paths**:
- ASCII file: `~/.local/share/syscgo/walls/<filename>.txt`
- Config: `~/.config/sysc-walls/daemon.conf`

**Process**:
1. User creates ASCII art in BIT editor
2. Ctrl+S → Select "sysc-walls"
3. Enter wallpaper name (e.g., "cyberpunk-logo")
4. System performs:
   - Create `~/.local/share/syscgo/walls/` if needed
   - Save ASCII art to `~/.local/share/syscgo/walls/<name>.txt`
   - Read/parse `~/.config/sysc-walls/daemon.conf`
   - Add or update wallpaper entry
   - Write config back

**Config Format** (INI-style):

```ini
[idle]
timeout = 300

[daemon]
monitors = all

[animation]
type = print
theme = dracula
file = /home/user/.local/share/syscgo/walls/cyberpunk-logo.txt
duration = infinite

[terminal]
fullscreen = true
opacity = 0.95
```

**Config Modification Logic**:

```go
// If [animation] section exists with "file" key:
//   - Update existing file path
//   - Preserve other settings (theme, type, duration)

// If [animation] section missing or no file key:
//   - Create [animation] section
//   - Set sensible defaults:
//     - type = print (or beam-text, ring-text)
//     - theme = dracula
//     - duration = infinite
//     - file = <new file path>
```

## Implementation Files

### `tui/export.go`

```go
package tui

type ExportTarget int

const (
    ExportSyscGo ExportTarget = iota
    ExportSyscWalls
)

type ExportOptions struct {
    Target   ExportTarget
    Filename string
    Content  string
}

// ExportAsciiArt handles export to different targets
func ExportAsciiArt(opts ExportOptions) error

// exportToSyscGo saves to local assets/
func exportToSyscGo(filename, content string) error

// exportToSyscWalls saves to walls dir + updates config
func exportToSyscWalls(filename, content string) error
```

### `tui/syscwalls_config.go`

```go
package tui

import "gopkg.in/ini.v1"

type SyscWallsConfig struct {
    configPath string
    cfg        *ini.File
}

// LoadSyscWallsConfig loads daemon.conf
func LoadSyscWallsConfig() (*SyscWallsConfig, error)

// SetWallpaperFile updates the file path in [animation] section
func (c *SyscWallsConfig) SetWallpaperFile(filePath string) error

// Save writes config back to disk
func (c *SyscWallsConfig) Save() error

// EnsureDefaults sets default values if sections missing
func (c *SyscWallsConfig) EnsureDefaults() error
```

## File Paths & Directory Structure

```
# syscgo export
./assets/
  ├── myart.txt
  ├── logo.txt
  └── banner.txt

# sysc-walls export
~/.local/share/syscgo/walls/
  ├── cyberpunk.txt
  ├── matrix-style.txt
  └── company-logo.txt

~/.config/sysc-walls/
  └── daemon.conf
```

## User Experience Flow

### syscgo Export:
```
1. Edit ASCII art
2. Ctrl+S
3. Select "syscgo - Save to assets/ folder"
4. Enter: "myart"
5. ✓ Saved to assets/myart.txt
6. Back to main menu (file now in selector)
```

### sysc-walls Export:
```
1. Edit ASCII art (with BIT editor styling)
2. Ctrl+S
3. Select "sysc-walls - Save as wallpaper"
4. Enter: "cyberpunk-logo"
5. ✓ Saved to ~/.local/share/syscgo/walls/cyberpunk-logo.txt
6. ✓ Updated ~/.config/sysc-walls/daemon.conf
7. Message: "Wallpaper configured! Restart sysc-walls daemon to activate."
8. Back to editor or main menu
```

## Config Update Strategy

**Safe config editing**:
1. Load existing daemon.conf (or create if missing)
2. Parse with INI library
3. Update only [animation].file path
4. Preserve all other settings
5. Add sensible defaults only if section completely missing
6. Write atomically (temp file + rename)

**Default [animation] values** (if creating new):
```ini
[animation]
type = beam-text
theme = dracula
file = /home/user/.local/share/syscgo/walls/<filename>.txt
duration = infinite
```

## Error Handling

### syscgo export errors:
- Directory not writable → Show error, suggest checking permissions
- File exists → Ask to overwrite (y/n prompt)
- Invalid filename → Show error, ask to re-enter

### sysc-walls export errors:
- ~/.local/share/syscgo/walls/ not writable → Create or show error
- ~/.config/sysc-walls/ not writable → Show error, suggest manual config
- daemon.conf parse error → Show warning, create new config with defaults
- daemon.conf missing → Create new config with defaults

## Dependencies

```go
// Add to go.mod
require (
    gopkg.in/ini.v1 v1.67.0  // INI config parsing
)
```

## Testing

1. **Unit tests**:
   - INI config parsing/writing
   - Path resolution (~ expansion)
   - Filename validation

2. **Integration tests**:
   - Export to syscgo assets/
   - Export to sysc-walls with config update
   - Config creation from scratch
   - Config update preserving existing values

3. **Manual tests**:
   - Create art, export to syscgo, verify in file selector
   - Create art, export to sysc-walls, verify file + config
   - Check sysc-walls daemon picks up new wallpaper
