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

This project uses Vercel's Agent Skills system to maintain consistent code quality and follow React best practices. Skills are located in `.github/skills/` for GitHub Copilot integration.

### Currently Installed Skills

#### From Vercel

- **vercel-react-best-practices**: React and Next.js performance optimization guidelines
  - Located in `.github/skills/vercel-react-best-practices/`
  - 45 rules across 8 categories (eliminating waterfalls, bundle optimization, server-side performance, etc.)
  - Used by GitHub Copilot and other AI coding agents
  - Automatically triggered when working with React components or Next.js pages

- **web-design-guidelines**: Review UI code for web interface best practices
  - Located in `.github/skills/web-design-guidelines/`
  - 100+ rules covering accessibility, performance, and UX
  - Audits for ARIA labels, semantic HTML, keyboard handlers, focus states, animations, typography, images, and more
  - Triggered by "Review my UI", "Check accessibility", or "Audit design" tasks

#### From OpenCode Community

- **frontend-ui-ux**: Designer-turned-developer who crafts stunning UI/UX
  - From [code-yeongyu/oh-my-opencode](https://github.com/code-yeongyu/oh-my-opencode)
  - Creates visually stunning, emotionally engaging interfaces
  - Obsesses over pixel-perfect details, smooth animations, and intuitive interactions

- **git-master**: Git expert for commit architecture, rebasing, and history archaeology
  - From [code-yeongyu/oh-my-opencode](https://github.com/code-yeongyu/oh-my-opencode)
  - Handles atomic commits, rebase/squash, history search (blame, bisect, log -S)
  - Triggered by git operations: 'commit', 'rebase', 'squash', 'who wrote', 'when was X added'

- **test-skill**: Test skill for skill system validation
  - From [anomalyco/opencode](https://github.com/anomalyco/opencode)

### Available Skills from Vercel

Additional skills can be installed from the [vercel-labs/agent-skills](https://github.com/vercel-labs/agent-skills) repository:

1. **web-design-guidelines**: Review UI code for compliance with web interface best practices
   - 100+ rules covering accessibility, performance, and UX
   - Audits for ARIA labels, semantic HTML, keyboard handlers, focus states, animations, typography, images, and more
   - Useful for "Review my UI", "Check accessibility", or "Audit design" tasks

2. **vercel-deploy-claimable**: Deploy applications to Vercel directly from conversations
   - Auto-detects 40+ frameworks from `package.json`
   - Returns preview URL and claim URL for ownership transfer
   - Useful for "Deploy my app" or "Push this live" tasks

### Installing Additional Skills

This project uses Bun. Install skills using `bunx`:

```bash
# Add skills from Vercel
bunx skills add vercel-labs/agent-skills

# Add skills from OpenCode community
bunx skills add anomalyco/opencode
bunx skills add code-yeongyu/oh-my-opencode

# Or add from any GitHub repository
bunx skills add owner/repo
```

Skills are automatically installed to `.github/skills/` directory for GitHub Copilot.

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
