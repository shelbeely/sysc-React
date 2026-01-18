/**
 * Color palettes for terminal animations
 */

/**
 * Get theme-specific fire colors
 */
export function getFirePalette(themeName: string): string[] {
  switch (themeName.toLowerCase()) {
    case 'dracula':
      return [
        '#282a36', // Background
        '#44475a', // Current line
        '#6272a4', // Comment
        '#8be9fd', // Cyan
        '#50fa7b', // Green
        '#f1fa8c', // Yellow
        '#ffb86c', // Orange
        '#ff79c6', // Pink
        '#ff5555', // Red (hottest)
      ];
    case 'catppuccin':
    case 'catppuccin-mocha':
      return [
        '#1e1e2e', // Base
        '#181825', // Mantle
        '#313244', // Surface0
        '#45475a', // Surface1
        '#f38ba8', // Red
        '#fab387', // Peach
        '#f9e2af', // Yellow
        '#a6e3a1', // Green (hot tip)
      ];
    case 'nord':
      return [
        '#2e3440', // Polar Night
        '#3b4252',
        '#434c5e',
        '#4c566a',
        '#bf616a', // Aurora Red
        '#d08770', // Aurora Orange
        '#ebcb8b', // Aurora Yellow
        '#a3be8c', // Aurora Green
      ];
    case 'tokyo-night':
    case 'tokyonight':
      return [
        '#1a1b26', // Background
        '#24283b', // Background Dark
        '#414868', // Foreground Gutter
        '#f7768e', // Red
        '#ff9e64', // Orange
        '#e0af68', // Yellow
        '#9ece6a', // Green
      ];
    case 'gruvbox':
      return [
        '#282828', // Background
        '#3c3836', // BG1
        '#504945', // BG2
        '#cc241d', // Red
        '#d65d0e', // Orange
        '#d79921', // Yellow
        '#fabd2f', // Bright Yellow
        '#b8bb26', // Green (hot)
      ];
    case 'material':
      return [
        '#263238', // Background
        '#37474f', // Lighter bg
        '#546e7a', // Selection
        '#f07178', // Red
        '#f78c6c', // Orange
        '#ffcb6b', // Yellow
        '#c3e88d', // Green
      ];
    case 'solarized':
      return [
        '#002b36', // Base03 - darkest
        '#073642', // Base02
        '#586e75', // Base01
        '#dc322f', // Red
        '#cb4b16', // Orange
        '#b58900', // Yellow
        '#859900', // Green
      ];
    case 'monochrome':
      return [
        '#1a1a1a', // Dark gray
        '#2a2a2a',
        '#3a3a3a',
        '#4a4a4a',
        '#5a5a5a',
        '#7a7a7a',
        '#9a9a9a',
        '#bababa',
        '#dadada', // Light gray (hottest)
      ];
    case 'transishardjob':
      return [
        '#55cdfc', // Trans blue
        '#f7a8b8', // Trans pink
        '#ffffff', // White
        '#f7a8b8', // Pink again
        '#55cdfc', // Blue again
        '#ffffff', // White (hottest)
      ];
    case 'rama':
      return [
        '#2b2d42', // Space cadet (background)
        '#8d99ae', // Cool gray
        '#d90429', // Fire engine red
        '#ef233c', // Red Pantone
        '#edf2f4', // Anti-flash white (hottest)
      ];
    case 'eldritch':
      return [
        '#212337', // Background
        '#292e42', // Current line
        '#7081d0', // Comment
        '#04d1f9', // Cyan
        '#37f499', // Green
        '#f1fc79', // Yellow
        '#f7c67f', // Orange
        '#f265b5', // Pink
        '#f16c75', // Red (hottest)
      ];
    case 'dark':
      return [
        '#000000', // True black
        '#333333', // Dark gray
        '#666666', // Mid gray
        '#999999', // Light gray
        '#cccccc', // Lighter gray
        '#ffffff', // True white (hottest)
      ];
    default:
      return getDefaultFirePalette();
  }
}

/**
 * Get classic DOOM-style fire palette
 */
export function getDefaultFirePalette(): string[] {
  return [
    '#000000', '#1a0000', '#330000', '#4d0000',
    '#660000', '#7f0000', '#990000', '#b30000',
    '#cc0000', '#e60000', '#ff0000', '#ff1a1a',
    '#ff3333', '#ff4d4d', '#ff6600', '#ff7f00',
    '#ff9900', '#ffb300', '#ffcc00', '#ffe600',
    '#ffff00', '#ffff33', '#ffff66', '#ffff99',
    '#ffffcc', '#ffffff',
  ];
}

/**
 * Convert hex color to RGB values
 */
export function hexToRGB(hex: string): [number, number, number] {
  // Remove # if present
  if (hex.startsWith('#')) {
    hex = hex.slice(1);
  }

  // Parse RGB
  if (hex.length === 6) {
    const r = parseInt(hex.slice(0, 2), 16);
    const g = parseInt(hex.slice(2, 4), 16);
    const b = parseInt(hex.slice(4, 6), 16);
    return [r, g, b];
  }

  return [0, 0, 0];
}

/**
 * Get ANSI color code from RGB values
 */
export function rgbToAnsi(r: number, g: number, b: number): string {
  return `\x1b[38;2;${r};${g};${b}m`;
}

/**
 * Get ANSI color code from hex color
 */
export function hexToAnsi(hex: string): string {
  const [r, g, b] = hexToRGB(hex);
  return rgbToAnsi(r, g, b);
}

/**
 * Reset ANSI color
 */
export const ANSI_RESET = '\x1b[0m';
