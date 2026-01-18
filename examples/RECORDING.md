# Recording Animations for README

This directory contains demos and recording scripts for capturing animations as GIFs to include in the README.

## Quick Demo

To run a demo locally:

```bash
npm install
npm run demo:fire    # Run fire animation demo
npm run demo:matrix  # Run matrix animation demo
```

## Recording GIFs with VHS

We use [VHS](https://github.com/charmbracelet/vhs) to record terminal animations as GIFs.

### Install VHS

**macOS:**
```bash
brew install vhs
```

**Linux:**
```bash
# Install via Go
go install github.com/charmbracelet/vhs@latest

# Or download binary from releases
# https://github.com/charmbracelet/vhs/releases
```

### Record Animations

```bash
# Record individual animations
npm run record:fire
npm run record:matrix

# Record all animations
npm run record:all
```

The GIFs will be saved in `examples/demos/` directory.

## Alternative: Using asciinema

You can also use [asciinema](https://asciinema.org/) and [agg](https://github.com/asciinema/agg) to create GIFs:

```bash
# Install asciinema
brew install asciinema  # macOS
# or
apt install asciinema   # Linux

# Install agg for converting to GIF
cargo install --git https://github.com/asciinema/agg

# Record
asciinema rec fire.cast -c "npm run demo:fire"

# Convert to GIF
agg fire.cast fire.gif
```

## Alternative: Using ttygif

Another option is [ttygif](https://github.com/icholy/ttygif):

```bash
# Install ttygif
npm install -g ttygif

# Record (this will create a ttyrec file)
ttyrec myrecording

# In the recording session, run:
npm run demo:fire

# Exit the recording with Ctrl+D

# Convert to GIF
ttygif myrecording
```

## Alternative: Using terminalizer

[Terminalizer](https://github.com/faressoft/terminalizer) is another great option:

```bash
# Install
npm install -g terminalizer

# Record
terminalizer record fire --skip-sharing

# In the recording, run:
npm run demo:fire

# Render to GIF
terminalizer render fire
```

## Manual Recording Tips

If you want to record manually:

1. **Use high quality settings**: Set terminal to at least 60x20 characters
2. **Use appropriate themes**: Dracula, Nord, and Catppuccin look great in recordings
3. **Record for 3-5 seconds**: Enough to show the animation looping
4. **Optimize GIF size**: Use tools like `gifsicle` to reduce file size:
   ```bash
   gifsicle -O3 --colors 256 input.gif -o output.gif
   ```

## Recommended Settings

For best results in README:

- **Width**: 60-80 characters
- **Height**: 20-24 lines
- **Duration**: 3-5 seconds
- **Frame rate**: 10-15 FPS for GIFs (lower file size)
- **Colors**: 256 colors max for GIFs
- **File size**: Keep under 5MB for GitHub

## VHS Tape File Format

The `.tape` files in this directory use VHS format. You can customize them:

```tape
# Set output file
Output my-animation.gif

# Set theme
Set Theme "Dracula"

# Set dimensions
Set Width 1200
Set Height 600

# Set font size
Set FontSize 16

# Commands to run
Type "npm run demo:fire"
Enter
Sleep 5s
Ctrl+C
```

See [VHS documentation](https://github.com/charmbracelet/vhs) for more options.
