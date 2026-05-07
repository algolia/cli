#!/usr/bin/env node
'use strict';

const { execFileSync } = require('child_process');

const pkg = `@algolia/cli-${process.platform}-${process.arch}`;
let binPath;
try {
  ({ binPath } = require(pkg));
} catch {
  console.error(
    `algolia: unsupported platform ${process.platform}/${process.arch}\n` +
    `Install the appropriate platform package manually: npm install ${pkg}`
  );
  process.exit(1);
}

try {
  execFileSync(binPath, process.argv.slice(2), { stdio: 'inherit' });
} catch (e) {
  process.exit(e.status ?? 1);
}
