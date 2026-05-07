#!/usr/bin/env node
'use strict';

const { execFileSync } = require('child_process');

const platforms = {
  'darwin-x64':   () => require('@algolia/cli-darwin-x64'),
  'darwin-arm64': () => require('@algolia/cli-darwin-arm64'),
  'linux-x64':    () => require('@algolia/cli-linux-x64'),
  'linux-arm64':  () => require('@algolia/cli-linux-arm64'),
  'win32-x64':    () => require('@algolia/cli-win32-x64'),
  'win32-arm64':  () => require('@algolia/cli-win32-arm64'),
};

const key = `${process.platform}-${process.arch}`;
const loader = platforms[key];

if (!loader) {
  console.error(
    `algolia: unsupported platform ${process.platform}/${process.arch}\n` +
    `Install the appropriate platform package manually: npm install @algolia/cli-${key}`
  );
  process.exit(1);
}

const { binPath } = loader();

try {
  execFileSync(binPath, process.argv.slice(2), { stdio: 'inherit' });
} catch (e) {
  process.exit(e.status ?? 1);
}
