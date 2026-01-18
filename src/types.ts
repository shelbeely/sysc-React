/**
 * Animation interface that all effects implement
 */
export interface Animation {
  /**
   * Update advances the animation by one frame
   */
  update(): void;

  /**
   * Render returns the current frame as a string
   */
  render(): string;

  /**
   * Reset restarts the animation from the beginning
   */
  reset(): void;
}

/**
 * Configuration for animation effects
 */
export interface Config {
  /** Terminal width in characters */
  width: number;
  /** Terminal height in characters */
  height: number;
  /** Color theme name */
  theme: string;
}

/**
 * Effect metadata describing an animation effect
 */
export interface EffectMetadata {
  /** Effect name (e.g., "fire", "matrix") */
  name: string;
  /** Whether effect requires text input */
  requiresText: boolean;
  /** Brief description */
  description: string;
  /** Version when effect was added */
  versionAdded: string;
  /** Effect category (e.g., "particle", "text", "abstract") */
  category: string;
}

/**
 * Theme metadata describing a color theme
 */
export interface ThemeMetadata {
  /** Theme name (e.g., "nord", "dracula") */
  name: string;
  /** Alternative names for the theme */
  aliases: string[];
  /** Brief description */
  description: string;
  /** Version when theme was added */
  versionAdded: string;
}
