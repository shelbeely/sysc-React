import { EffectMetadata, ThemeMetadata } from './types';

/**
 * Library version
 */
export const LIBRARY_VERSION = '1.0.0';

/**
 * Effect registry containing metadata for all available effects
 */
export const EFFECT_REGISTRY: EffectMetadata[] = [
  {
    name: 'matrix',
    requiresText: false,
    description: 'Classic Matrix digital rain effect',
    versionAdded: '1.0.0',
    category: 'particle',
  },
  {
    name: 'fire',
    requiresText: false,
    description: 'Doom-style fire effect',
    versionAdded: '1.0.0',
    category: 'particle',
  },
  {
    name: 'fire-text',
    requiresText: true,
    description: 'Fire effect with text as negative space',
    versionAdded: '1.0.0',
    category: 'text',
  },
];

/**
 * Theme registry containing metadata for all available themes
 */
export const THEME_REGISTRY: ThemeMetadata[] = [
  {
    name: 'dracula',
    aliases: [],
    description: 'Dracula dark theme with purple and pink accents',
    versionAdded: '1.0.0',
  },
  {
    name: 'catppuccin',
    aliases: ['catppuccin-mocha'],
    description: 'Catppuccin Mocha - soothing pastel theme',
    versionAdded: '1.0.0',
  },
  {
    name: 'nord',
    aliases: [],
    description: 'Nord arctic, north-bluish color palette',
    versionAdded: '1.0.0',
  },
  {
    name: 'tokyo-night',
    aliases: ['tokyonight'],
    description: 'Tokyo Night dark theme inspired by Tokyo',
    versionAdded: '1.0.0',
  },
  {
    name: 'gruvbox',
    aliases: [],
    description: 'Gruvbox retro groove color scheme',
    versionAdded: '1.0.0',
  },
  {
    name: 'material',
    aliases: [],
    description: 'Material Design color palette',
    versionAdded: '1.0.0',
  },
  {
    name: 'solarized',
    aliases: [],
    description: 'Solarized precision colors for machines and people',
    versionAdded: '1.0.0',
  },
  {
    name: 'monochrome',
    aliases: [],
    description: 'Grayscale monochrome theme',
    versionAdded: '1.0.0',
  },
  {
    name: 'transishardjob',
    aliases: [],
    description: 'Trans pride colors',
    versionAdded: '1.0.0',
  },
  {
    name: 'rama',
    aliases: [],
    description: 'Rama custom color scheme',
    versionAdded: '1.0.0',
  },
  {
    name: 'eldritch',
    aliases: [],
    description: 'Eldritch dark theme with purple and cyan',
    versionAdded: '1.0.0',
  },
  {
    name: 'dark',
    aliases: [],
    description: 'Simple dark theme with grayscale',
    versionAdded: '1.0.0',
  },
];

/**
 * Get all available effect names
 */
export function getEffectNames(): string[] {
  return EFFECT_REGISTRY.map(effect => effect.name);
}

/**
 * Get metadata for a specific effect
 */
export function getEffectMetadata(name: string): EffectMetadata | null {
  return EFFECT_REGISTRY.find(effect => effect.name === name) || null;
}

/**
 * Check if an effect requires text input
 */
export function isTextBasedEffect(name: string): boolean {
  const meta = getEffectMetadata(name);
  return meta !== null && meta.requiresText;
}

/**
 * Get all available theme names (including aliases)
 */
export function getThemeNames(): string[] {
  const names: string[] = [];
  for (const theme of THEME_REGISTRY) {
    names.push(theme.name);
    names.push(...theme.aliases);
  }
  return names;
}

/**
 * Get metadata for a specific theme
 */
export function getThemeMetadata(name: string): ThemeMetadata | null {
  for (const theme of THEME_REGISTRY) {
    if (theme.name === name) {
      return theme;
    }
    // Check aliases
    if (theme.aliases.includes(name)) {
      return theme;
    }
  }
  return null;
}
