#!/usr/bin/env node

/**
 * ScreenCraft SDK Setup Verification Script
 *
 * This script verifies that the SDK is correctly set up and ready for development.
 */

const fs = require('fs');
const path = require('path');

const COLORS = {
  reset: '\x1b[0m',
  green: '\x1b[32m',
  red: '\x1b[31m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
};

function log(message, color = 'reset') {
  console.log(`${COLORS[color]}${message}${COLORS.reset}`);
}

function checkFileExists(filePath, description) {
  const fullPath = path.join(__dirname, filePath);
  const exists = fs.existsSync(fullPath);

  if (exists) {
    log(`✓ ${description}`, 'green');
  } else {
    log(`✗ ${description} - MISSING: ${filePath}`, 'red');
  }

  return exists;
}

function checkNodeVersion() {
  const version = process.version;
  const major = parseInt(version.slice(1).split('.')[0]);

  if (major >= 14) {
    log(`✓ Node.js version ${version} (>= 14.0.0)`, 'green');
    return true;
  } else {
    log(`✗ Node.js version ${version} (requires >= 14.0.0)`, 'red');
    return false;
  }
}

function checkPackageJson() {
  try {
    const pkgPath = path.join(__dirname, 'package.json');
    const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf8'));

    if (pkg.name === 'screencraft' && pkg.version) {
      log(`✓ package.json valid (${pkg.name}@${pkg.version})`, 'green');
      return true;
    } else {
      log('✗ package.json invalid', 'red');
      return false;
    }
  } catch (error) {
    log(`✗ package.json error: ${error.message}`, 'red');
    return false;
  }
}

function checkTypeScriptConfig() {
  try {
    const tsconfigPath = path.join(__dirname, 'tsconfig.json');
    const tsconfig = JSON.parse(fs.readFileSync(tsconfigPath, 'utf8'));

    if (tsconfig.compilerOptions && tsconfig.compilerOptions.strict) {
      log('✓ tsconfig.json valid (strict mode enabled)', 'green');
      return true;
    } else {
      log('✗ tsconfig.json missing strict mode', 'yellow');
      return true; // Non-critical
    }
  } catch (error) {
    log(`✗ tsconfig.json error: ${error.message}`, 'red');
    return false;
  }
}

function checkSourceFiles() {
  const files = [
    'src/index.ts',
    'src/client.ts',
    'src/types.ts',
    'src/errors.ts',
  ];

  let allExist = true;
  files.forEach((file) => {
    if (!checkFileExists(file, `Source file: ${file}`)) {
      allExist = false;
    }
  });

  return allExist;
}

function checkExamples() {
  const examples = [
    'examples/basic-screenshot.ts',
    'examples/fullpage-screenshot.ts',
    'examples/generate-pdf.ts',
    'examples/async-webhook.ts',
    'examples/error-handling.ts',
    'examples/mobile-screenshot.ts',
  ];

  let allExist = true;
  examples.forEach((file) => {
    if (!checkFileExists(file, `Example: ${path.basename(file)}`)) {
      allExist = false;
    }
  });

  return allExist;
}

function checkDocumentation() {
  const docs = [
    'README.md',
    'QUICKSTART.md',
    'SETUP.md',
    'CONTRIBUTING.md',
    'CHANGELOG.md',
    'LICENSE',
  ];

  let allExist = true;
  docs.forEach((file) => {
    if (!checkFileExists(file, `Documentation: ${file}`)) {
      allExist = false;
    }
  });

  return allExist;
}

function checkConfigFiles() {
  const configs = [
    '.eslintrc.js',
    '.prettierrc',
    '.editorconfig',
    '.gitignore',
    '.npmignore',
    'jest.config.js',
  ];

  let allExist = true;
  configs.forEach((file) => {
    if (!checkFileExists(file, `Config: ${file}`)) {
      allExist = false;
    }
  });

  return allExist;
}

function printSummary(checks) {
  const total = checks.length;
  const passed = checks.filter((c) => c.passed).length;
  const failed = total - passed;

  console.log('\n' + '='.repeat(60));
  log('\nVerification Summary:', 'blue');
  console.log('='.repeat(60));

  if (failed === 0) {
    log(`\n✓ All ${total} checks passed!`, 'green');
    log('\nThe ScreenCraft SDK is correctly set up.', 'green');
    log('\nNext steps:', 'blue');
    console.log('  1. Run "npm install" to install dependencies');
    console.log('  2. Run "npm run build" to compile TypeScript');
    console.log('  3. Run "npm test" to run tests');
    console.log('  4. Check examples/ directory for usage examples');
    console.log('  5. Read QUICKSTART.md to get started\n');
  } else {
    log(`\n✗ ${failed} of ${total} checks failed`, 'red');
    log('\nPlease fix the issues above before proceeding.\n', 'yellow');
  }
}

function main() {
  log('\n' + '='.repeat(60), 'blue');
  log('  ScreenCraft SDK - Setup Verification', 'blue');
  log('='.repeat(60) + '\n', 'blue');

  const checks = [
    { name: 'Node.js version', fn: checkNodeVersion },
    { name: 'Package configuration', fn: checkPackageJson },
    { name: 'TypeScript configuration', fn: checkTypeScriptConfig },
    { name: 'Source files', fn: checkSourceFiles },
    { name: 'Example files', fn: checkExamples },
    { name: 'Documentation files', fn: checkDocumentation },
    { name: 'Configuration files', fn: checkConfigFiles },
  ];

  const results = checks.map((check) => {
    log(`\nChecking ${check.name}...`, 'blue');
    const passed = check.fn();
    return { ...check, passed };
  });

  printSummary(results);
}

main();
