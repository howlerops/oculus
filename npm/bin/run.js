#!/usr/bin/env node
const { execFileSync } = require('child_process');
const path = require('path');
const os = require('os');

const ext = os.platform() === 'win32' ? '.exe' : '';
const binary = path.join(__dirname, `oculus${ext}`);

try {
  execFileSync(binary, process.argv.slice(2), { stdio: 'inherit' });
} catch (e) {
  if (e.status !== undefined) process.exit(e.status);
  console.error('Failed to run oculus:', e.message);
  process.exit(1);
}
