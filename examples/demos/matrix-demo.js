#!/usr/bin/env node
import React from 'react';
import { render } from 'ink';
import { Matrix } from 'sysc-react';

// Matrix demo for README animation capture
const MatrixDemo = () => {
  return React.createElement(Matrix, {
    width: 60,
    height: 20,
    theme: 'dracula',
    frameRate: 50
  });
};

render(React.createElement(MatrixDemo));
