import { Animation } from '../types';
import { hexToAnsi, ANSI_RESET } from '../palettes';

/**
 * FireEffect implements PSX DOOM-style fire algorithm with enhanced character gradient
 */
export class FireEffect implements Animation {
  private width: number;
  private height: number;
  private buffer: number[];
  private palette: string[];
  private chars: string[];

  constructor(width: number, height: number, palette: string[]) {
    this.width = width;
    this.height = height;
    this.palette = palette;
    // Enhanced 8-character gradient for smoother fire rendering
    this.chars = [' ', '░', '░', '▒', '▒', '▓', '▓', '█'];
    this.buffer = [];
    this.init();
  }

  /**
   * Initialize fire buffer with bottom row as heat source
   */
  private init(): void {
    this.buffer = new Array(this.width * this.height).fill(0);

    // Set bottom row to maximum heat (fire source)
    for (let i = 0; i < this.width; i++) {
      this.buffer[(this.height - 1) * this.width + i] = 65;
    }
  }

  /**
   * Update the fire color palette (for theme switching)
   */
  updatePalette(palette: string[]): void {
    this.palette = palette;
  }

  /**
   * Resize reinitializes the fire effect with new dimensions
   */
  resize(width: number, height: number): void {
    this.width = width;
    this.height = height;
    this.init();
  }

  /**
   * Spread fire propagates heat upward with random decay (DOOM algorithm)
   */
  private spreadFire(from: number): void {
    // Random horizontal offset (0-3) for flickering effect
    const offset = Math.floor(Math.random() * 4);
    const to = from - this.width - offset + 1;

    // Bounds check
    if (to < 0 || to >= this.buffer.length) {
      return;
    }

    // Random decay (0-3) for natural fade
    const decay = Math.floor(Math.random() * 4);

    let newHeat = this.buffer[from] - decay;
    if (newHeat < 0) {
      newHeat = 0;
    }

    this.buffer[to] = newHeat;
  }

  /**
   * Update advances the fire simulation by one frame
   */
  update(): void {
    // Process all pixels from bottom to top
    // (Fire spreads upward, must process bottom row first)
    for (let y = this.height - 1; y > 0; y--) {
      for (let x = 0; x < this.width; x++) {
        const index = y * this.width + x;
        this.spreadFire(index);
      }
    }
  }

  /**
   * Render returns the current frame as a colored string
   */
  render(): string {
    let output = '';

    for (let y = 0; y < this.height; y++) {
      for (let x = 0; x < this.width; x++) {
        const index = y * this.width + x;
        const heat = this.buffer[index];

        // Map heat (0-65) to palette index
        const paletteIndex = Math.floor((heat / 65) * (this.palette.length - 1));
        const color = this.palette[paletteIndex];

        // Map heat to character density
        const charIndex = Math.floor((heat / 65) * (this.chars.length - 1));
        const char = this.chars[charIndex];

        // Apply color and character
        output += hexToAnsi(color) + char + ANSI_RESET;
      }
      output += '\n';
    }

    return output;
  }

  /**
   * Reset restarts the animation from the beginning
   */
  reset(): void {
    this.init();
  }
}
