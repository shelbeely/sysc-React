#!/usr/bin/env node
import { MatrixEffect } from '../../dist/animations/matrix.js';
import { getFirePalette } from '../../dist/palettes.js';

// Simple matrix demo without React - just renders to console
const width = 60;
const height = 20;
const palette = getFirePalette('dracula');
const matrix = new MatrixEffect(width, height, palette);

// Clear screen and hide cursor
process.stdout.write('\x1b[2J\x1b[?25l');

let frameCount = 0;
const maxFrames = 100; // About 5 seconds at 50ms per frame

const interval = setInterval(() => {
  matrix.update();
  const frame = matrix.render();
  
  // Move cursor to top-left and render
  process.stdout.write('\x1b[H' + frame);
  
  frameCount++;
  if (frameCount >= maxFrames) {
    // Show cursor again and exit
    process.stdout.write('\x1b[?25h');
    clearInterval(interval);
    process.exit(0);
  }
}, 50);

// Handle Ctrl+C
process.on('SIGINT', () => {
  process.stdout.write('\x1b[?25h');
  process.exit(0);
});
