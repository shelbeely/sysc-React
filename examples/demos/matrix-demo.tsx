#!/usr/bin/env node
import React from 'react';
import { render } from 'ink';
import { Matrix } from '../../dist/index.js';

// Matrix demo for README animation capture
const MatrixDemo = () => {
  return <Matrix width={60} height={20} theme="dracula" frameRate={50} />;
};

render(<MatrixDemo />);
