import { Animation } from '../types';
import { hexToAnsi, ANSI_RESET } from '../palettes';

/**
 * MatrixStreak represents a single vertical streak falling down the screen
 */
interface MatrixStreak {
  x: number;       // X position (column)
  y: number;       // Y position of head
  length: number;  // Length of streak
  speed: number;   // Movement speed (frames per pixel)
  counter: number; // Frame counter for movement
  active: boolean; // Whether streak is active
}

/**
 * MatrixEffect implements Matrix digital rain animation using particle-based streaks
 */
export class MatrixEffect implements Animation {
  private width: number;
  private height: number;
  private palette: string[];
  private chars: string[];
  private streaks: MatrixStreak[];
  private frame: number;

  constructor(width: number, height: number, palette: string[]) {
    this.width = width;
    this.height = height;
    this.palette = palette;
    // Use a mix of Latin, Greek, and Japanese characters like the original Matrix effect
    this.chars = [
      '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
      'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
      'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
      'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
      'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
      'α', 'β', 'γ', 'δ', 'ε', 'ζ', 'η', 'θ', 'ι', 'κ', 'λ', 'μ',
      'ν', 'ξ', 'ο', 'π', 'ρ', 'σ', 'τ', 'υ', 'φ', 'χ', 'ψ', 'ω',
      'А', 'Б', 'В', 'Г', 'Д', 'Е', 'Ж', 'З', 'И', 'Й', 'К', 'Л', 'М',
      'Н', 'О', 'П', 'Р', 'С', 'Т', 'У', 'Ф', 'Х', 'Ц', 'Ч', 'Ш', 'Щ',
      '░', '▒', '▓', '█', '▀', '▄', '▌', '▐', '■', '□', '▪', '▫',
    ];
    this.streaks = [];
    this.frame = 0;
    this.init();
  }

  /**
   * Initialize Matrix effect with some initial streaks
   */
  private init(): void {
    this.streaks = [];
    // Create initial streaks across width
    for (let i = 0; i < this.width; i++) {
      if (Math.random() < 0.1) { // 10% chance of initial streak
        const streak: MatrixStreak = {
          x: i,
          y: -Math.floor(Math.random() * this.height), // Start above screen
          length: Math.floor(Math.random() * 15) + 5,   // Length 5-20
          speed: Math.floor(Math.random() * 3) + 1,     // Speed 1-3
          counter: 0,
          active: true,
        };
        this.streaks.push(streak);
      }
    }
  }

  /**
   * Update the Matrix color palette (for theme switching)
   */
  updatePalette(palette: string[]): void {
    this.palette = palette;
  }

  /**
   * Resize reinitializes the Matrix effect with new dimensions
   */
  resize(width: number, height: number): void {
    this.width = width;
    this.height = height;
    this.init();
  }

  /**
   * Get the bright color for the head of the streak
   */
  private getHeadColor(): string {
    if (this.palette.length === 0) {
      return '#ffffff'; // Default white if no palette
    }
    // Use the brightest color from the palette for heads
    return this.palette[this.palette.length - 1];
  }

  /**
   * Get a dimmer color for trail positions
   */
  private getTrailColor(position: number, length: number): string {
    if (this.palette.length === 0) {
      return '#00aa00'; // Default dimmer green
    }

    // Calculate fade factor (0.0 = head, 1.0 = tail)
    const fadeFactor = position / length;

    // Use different colors based on position in trail
    if (fadeFactor < 0.2) {
      // Bright trail near head
      return this.palette[this.palette.length - 1];
    } else if (fadeFactor < 0.5) {
      // Medium trail
      if (this.palette.length > 2) {
        return this.palette[this.palette.length - 2];
      }
      return this.palette[0];
    } else {
      // Dim trail
      return this.palette[0];
    }
  }

  /**
   * Update advances the Matrix simulation by one frame
   */
  update(): void {
    this.frame++;

    // Update existing streaks
    const activeStreaks: MatrixStreak[] = [];
    for (const streak of this.streaks) {
      if (!streak.active) {
        continue;
      }

      // Update streak movement counter
      streak.counter++;

      // Move streak when counter reaches speed threshold
      if (streak.counter >= streak.speed) {
        streak.y++;
        streak.counter = 0;

        // Deactivate streak when it moves completely off screen
        if (streak.y - streak.length > this.height) {
          streak.active = false;
          continue;
        }
      }

      // Add updated streak to active list
      activeStreaks.push(streak);
    }

    // Replace streak list with active streaks
    this.streaks = activeStreaks;

    // Add new streaks randomly
    for (let i = 0; i < this.width; i++) {
      // Low probability to create new streaks
      if (Math.random() < 0.02 && this.streaks.length < 150) { // Limit total streaks
        const streak: MatrixStreak = {
          x: i,
          y: -Math.floor(Math.random() * 5),     // Start just above screen
          length: Math.floor(Math.random() * 15) + 5, // Length 5-20
          speed: Math.floor(Math.random() * 3) + 1,   // Speed 1-3
          counter: 0,
          active: true,
        };
        this.streaks.push(streak);
      }
    }
  }

  /**
   * Render converts the Matrix streaks to colored text output
   */
  render(): string {
    // Create empty canvas
    const canvas: string[][] = Array(this.height)
      .fill(null)
      .map(() => Array(this.width).fill(' '));
    const colors: string[][] = Array(this.height)
      .fill(null)
      .map(() => Array(this.width).fill(''));

    // Render each active streak
    for (const streak of this.streaks) {
      if (!streak.active) {
        continue;
      }

      // Render the streak - from head downward
      for (let i = 0; i < streak.length; i++) {
        const yPos = streak.y + i; // Head at streak.y, trail going down
        if (yPos >= 0 && yPos < this.height && streak.x >= 0 && streak.x < this.width) {
          // Get random character
          const char = this.chars[Math.floor(Math.random() * this.chars.length)];

          // Get color based on position in streak
          let color: string;
          if (i === 0) {
            // Head is brightest
            color = this.getHeadColor();
          } else {
            // Trail fades
            color = this.getTrailColor(i, streak.length);
          }

          // Place character on canvas
          canvas[yPos][streak.x] = char;
          colors[yPos][streak.x] = color;
        }
      }
    }

    // Convert to colored string
    let output = '';
    for (let y = 0; y < this.height; y++) {
      for (let x = 0; x < this.width; x++) {
        const char = canvas[y][x];
        if (char !== ' ' && colors[y][x] !== '') {
          // Render colored character
          output += hexToAnsi(colors[y][x]) + char + ANSI_RESET;
        } else {
          output += char;
        }
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
