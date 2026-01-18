# Contributing to sysc-React

This document provides guidance for contributors working on the sysc-React terminal animation library.

## Development Setup

### Prerequisites

- **Bun**: This project uses Bun as its primary JavaScript runtime and package manager
  - Install from: https://bun.sh
  - Version: 1.0.0 or higher

### Installing Dependencies

```bash
# Root dependencies
bun install

# Example dependencies
cd examples
bun install
```

## AI Coding Agent Skills

This project uses Vercel's Agent Skills system to maintain consistent code quality and follow React best practices.

### Installing Skills (Bun)

When available, use `bunx` instead of `npx` for better performance:

```bash
# List available skills
bunx add-skill vercel-labs/agent-skills --list

# Install React best practices
bunx add-skill vercel-labs/agent-skills --skill vercel-react-best-practices

# Install web design guidelines
bunx add-skill vercel-labs/agent-skills --skill web-design-guidelines
```

### Currently Installed Skills

- **vercel-react-best-practices**: React and Next.js performance optimization guidelines
  - Located in `.github/skills/vercel-react-best-practices/`
  - 45 rules across 8 categories
  - Used by GitHub Copilot and other AI coding agents

## Recording Animations

The project uses `asciinema` and `agg` for recording terminal animations.

### Setup Recording Tools

```bash
# Install with Bun
bun add -d asciinema @abstr/agg
```

### Recording Process

1. **Create a simple demo script** in `examples/demos/`:
```javascript
#!/usr/bin/env node
import { FireEffect, getFirePalette } from '../../dist/index.js';

const fire = new FireEffect(80, 24, getFirePalette('dracula'));

setInterval(() => {
  fire.update();
  console.clear();
  console.log(fire.render());
}, 50);
```

2. **Record with asciinema**:
```bash
cd examples
asciinema rec --overwrite -c "timeout 5 node demos/fire-simple.js" recordings/fire.cast
```

3. **Convert to GIF with agg**:
```bash
agg recordings/fire.cast recordings/fire.gif
```

### Recording Best Practices

- Keep recordings short (3-5 seconds) to minimize file size
- Use `timeout` command to auto-stop the recording
- Record at standard terminal dimensions (80x24)
- Store .cast files for editability, GIFs for documentation

## Code Quality

### Linting and Type Checking

```bash
# Run TypeScript compiler
bun run build

# The project uses TypeScript for type safety
# Make sure your code compiles without errors
```

### Testing

```bash
# Run examples to validate changes
cd examples
bun run fire
bun run matrix
```

## Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test your changes with the examples
5. Ensure code compiles without errors
6. Submit a pull request

## Project Structure

```
sysc-React/
├── src/                    # TypeScript source code
│   ├── animations/         # Animation classes
│   ├── components/         # React components
│   ├── effects/            # Effect metadata
│   └── palettes/           # Color palettes
├── examples/               # Example usage
│   ├── demos/              # Simple demo scripts for recording
│   ├── recordings/         # Animation recordings (.cast and .gif)
│   └── *.js                # Full-featured examples
├── .github/skills/         # AI agent skills
└── dist/                   # Compiled output
```

## License

This project is MIT licensed. See LICENSE for details.
