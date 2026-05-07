#!/usr/bin/env node
'use strict';

const { execFileSync } = require('child_process');

let binPath;
switch (`${process.platform}-${process.arch}`) {
  case 'darwin-x64':   ({ binPath } = require('@algolia/cli-darwin-x64'));   break;
  case 'darwin-arm64': ({ binPath } = require('@algolia/cli-darwin-arm64')); break;
  case 'linux-x64':    ({ binPath } = require('@algolia/cli-linux-x64'));    break;
  case 'linux-arm64':  ({ binPath } = require('@algolia/cli-linux-arm64'));  break;
  case 'win32-x64':    ({ binPath } = require('@algolia/cli-win32-x64'));    break;
  case 'win32-arm64':  ({ binPath } = require('@algolia/cli-win32-arm64'));  break;
  default:
    console.error(
      `algolia: unsupported platform ${process.platform}/${process.arch}\n` +
      `Install the appropriate platform package manually: npm install @algolia/cli-${process.platform}-${process.arch}`
    );
    process.exit(1);
}

try {
  execFileSync(binPath, process.argv.slice(2), { stdio: 'inherit' });
} catch (e) {
  process.exit(e.status ?? 1);
}
