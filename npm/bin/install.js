#!/usr/bin/env node
const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const os = require('os');

const platform = os.platform();
const arch = os.arch();

const platformMap = {
  'darwin-arm64': 'oculus-darwin-arm64',
  'darwin-x64': 'oculus-darwin-amd64',
  'linux-x64': 'oculus-linux-amd64',
  'linux-arm64': 'oculus-linux-arm64',
  'win32-x64': 'oculus-windows-amd64.exe',
};

const key = `${platform}-${arch}`;
const binary = platformMap[key];

if (!binary) {
  console.error(`Unsupported platform: ${key}`);
  process.exit(1);
}

const src = path.join(__dirname, binary);
const dest = path.join(__dirname, platform === 'win32' ? 'oculus.exe' : 'oculus');

if (fs.existsSync(src)) {
  fs.copyFileSync(src, dest);
  if (platform !== 'win32') {
    fs.chmodSync(dest, 0o755);
  }
  console.log(`Oculus installed for ${key}`);
} else {
  console.error(`Binary not found: ${src}`);
  console.error('Try installing from source: go install github.com/howlerops/oculus/cmd/oculus@latest');
  process.exit(1);
}
