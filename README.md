# sysc-React

Terminal animation library for React with Ink support. Pure React animations ready to use in your TUI applications.

This is a React/TypeScript port of [sysc-Go](https://github.com/Nomadcxx/sysc-Go), designed to work seamlessly with [Ink](https://github.com/vadimdemedes/ink) - the React renderer for interactive command-line applications.

## Features

- 🔥 **Fire Effect** - DOOM PSX-style fire animation
- 🌧️ **Matrix Rain** - Classic Matrix digital rain effect
- 🎨 **12+ Color Themes** - Dracula, Nord, Catppuccin, Tokyo Night, Gruvbox, and more
- ⚛️ **React Components** - Drop-in components for Ink applications
- 🎯 **TypeScript** - Full type safety and IntelliSense support
- 🚀 **Bun Compatible** - Optimized for Bun.js runtime

## Installation

```bash
npm install sysc-react react ink
```

Or with Bun:

```bash
bun add sysc-react react ink
```

## Quick Start

### Using Ink Components

```tsx
import React from 'react';
import { render } from 'ink';
import { Fire, Matrix } from 'sysc-react';

const App = () => {
  return (
    <>
      <Fire width={80} height={24} theme="dracula" />
      {/* or */}
      <Matrix width={80} height={24} theme="nord" />
    </>
  );
};

render(<App />);
```

### Using Animation Classes Directly

```typescript
import { FireEffect, getFirePalette } from 'sysc-react';

const palette = getFirePalette('dracula');
const fire = new FireEffect(80, 24, palette);

// In your render loop
setInterval(() => {
  fire.update();
  const frame = fire.render();
  console.log(frame);
}, 50);
```

## Available Effects

### Fire

Classic DOOM PSX-style fire effect with rising flames.

```tsx
import { Fire } from 'sysc-react';

<Fire 
  width={80}
  height={24}
  theme="dracula"
  frameRate={50}
/>
```

**Props:**
- `width` (number, default: 80) - Terminal width in characters
- `height` (number, default: 24) - Terminal height in characters
- `theme` (string, default: 'dracula') - Color theme name
- `frameRate` (number, default: 50) - Frame rate in milliseconds

### Matrix

Matrix-style digital rain with falling character streaks.

```tsx
import { Matrix } from 'sysc-react';

<Matrix 
  width={80}
  height={24}
  theme="nord"
  frameRate={50}
/>
```

**Props:**
- `width` (number, default: 80) - Terminal width in characters
- `height` (number, default: 24) - Terminal height in characters
- `theme` (string, default: 'dracula') - Color theme name
- `frameRate` (number, default: 50) - Frame rate in milliseconds

## Available Themes

- `dracula` - Dracula dark theme with purple and pink accents
- `catppuccin` - Catppuccin Mocha soothing pastel theme
- `nord` - Nord arctic, north-bluish color palette
- `tokyo-night` - Tokyo Night dark theme inspired by Tokyo
- `gruvbox` - Gruvbox retro groove color scheme
- `material` - Material Design color palette
- `solarized` - Solarized precision colors
- `monochrome` - Grayscale monochrome theme
- `transishardjob` - Trans pride colors
- `rama` - Rama custom color scheme
- `eldritch` - Eldritch dark theme with purple and cyan
- `dark` - Simple dark theme with grayscale

## Examples

See the [`examples/`](./examples) directory for complete working examples:

```bash
cd examples
npm install
npm run demo:fire    # Run fire animation demo
npm run demo:matrix  # Run matrix animation demo
```

## Recording Animations

Want to capture animations as GIFs for your README or documentation? See [examples/RECORDING.md](./examples/RECORDING.md) for a comprehensive guide.

### Quick Recording with VHS

[VHS](https://github.com/charmbracelet/vhs) is the recommended tool for recording terminal animations as GIFs.

**Install:**
```bash
# macOS
brew install vhs

# Linux (with Go)
go install github.com/charmbracelet/vhs@latest
```

**Record:**
```bash
cd examples
npm run record:fire    # Creates demos/fire.gif
npm run record:matrix  # Creates demos/matrix.gif
```

### Other Recording Tools

The repository includes instructions for multiple recording tools:

- **VHS** (recommended) - High-quality terminal recorder with tape files
- **asciinema + agg** - Record as asciicast, convert to GIF
- **ttygif** - Simple npm-based recorder
- **terminalizer** - Feature-rich terminal recorder

Each tool has different strengths. See [examples/RECORDING.md](./examples/RECORDING.md) for detailed instructions.

### Manual Recording Tips

For best results in your README:
- Use 60-80 character width, 20-24 lines height
- Record 3-5 seconds (enough to show animation looping)
- Use 10-15 FPS for GIFs to keep file size manageable
- Optimize with `gifsicle -O3 --colors 256 input.gif -o output.gif`
- Keep file size under 5MB for GitHub

## API Reference

### Animation Interface

All effects implement the `Animation` interface:

```typescript
interface Animation {
  update(): void;           // Advance animation by one frame
  render(): string;         // Get current frame as colored string
  reset(): void;            // Restart animation from beginning
}
```

### Palette Functions

```typescript
import { getFirePalette, hexToAnsi, hexToRGB } from 'sysc-react';

// Get color palette for a theme
const palette = getFirePalette('dracula'); // string[]

// Convert hex color to ANSI escape code
const ansiCode = hexToAnsi('#ff5555'); // '\x1b[38;2;255;85;85m'

// Convert hex to RGB values
const [r, g, b] = hexToRGB('#ff5555'); // [255, 85, 85]
```

### Registry Functions

```typescript
import { 
  getEffectNames, 
  getEffectMetadata, 
  getThemeNames,
  getThemeMetadata 
} from 'sysc-react';

// Get all available effect names
const effects = getEffectNames(); // ['fire', 'matrix']

// Get effect metadata
const fireInfo = getEffectMetadata('fire');

// Get all theme names
const themes = getThemeNames();

// Get theme metadata
const draculaInfo = getThemeMetadata('dracula');
```

## Building from Source

```bash
# Clone the repository
git clone https://github.com/shelbeely/sysc-React.git
cd sysc-React

# Install dependencies
npm install  # or: bun install

# Build
npm run build  # or: bun run build

# Watch mode for development
npm run watch
```

## Development

This project uses:
- **TypeScript** for type safety
- **React** for component model
- **Ink** for terminal rendering
- **Bun** as the recommended runtime (also works with Node.js)

## Acknowledgements

- [sysc-Go](https://github.com/Nomadcxx/sysc-Go) - Original Go implementation
- [Ink](https://github.com/vadimdemedes/ink) - React for CLIs
- [terminaltexteffects](https://github.com/ChrisBuilds/terminaltexteffects) - Inspiration for terminal visual effects

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
