#!/usr/bin/env node
import React from 'react';
import { render } from 'ink';
import { Fire } from '../../dist/index.js';

// Fire demo for README animation capture
const FireDemo = () => {
  return <Fire width={60} height={20} theme="dracula" frameRate={50} />;
};

render(<FireDemo />);
