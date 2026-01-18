import React, { useEffect, useState } from 'react';
import { Box, Text } from 'ink';
import { MatrixEffect } from '../animations/matrix';
import { getFirePalette } from '../palettes';

export interface MatrixProps {
  /** Terminal width in characters (default: 80) */
  width?: number;
  /** Terminal height in characters (default: 24) */
  height?: number;
  /** Theme name (default: 'dracula') */
  theme?: string;
  /** Frame rate in milliseconds (default: 50) */
  frameRate?: number;
}

/**
 * Matrix component - Renders Matrix digital rain effect
 */
export const Matrix: React.FC<MatrixProps> = ({
  width = 80,
  height = 24,
  theme = 'dracula',
  frameRate = 50,
}) => {
  const [frame, setFrame] = useState('');
  const [effect] = useState(() => {
    const palette = getFirePalette(theme);
    return new MatrixEffect(width, height, palette);
  });

  useEffect(() => {
    const interval = setInterval(() => {
      effect.update();
      const rendered = effect.render();
      setFrame(rendered);
    }, frameRate);

    return () => clearInterval(interval);
  }, [effect, frameRate]);

  useEffect(() => {
    // Update palette when theme changes
    const palette = getFirePalette(theme);
    effect.updatePalette(palette);
  }, [theme, effect]);

  useEffect(() => {
    // Resize when dimensions change
    effect.resize(width, height);
  }, [width, height, effect]);

  return (
    <Box flexDirection="column">
      <Text>{frame}</Text>
    </Box>
  );
};
