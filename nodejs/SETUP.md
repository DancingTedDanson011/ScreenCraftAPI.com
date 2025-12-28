# ScreenCraft SDK Development Setup

This guide will help you set up the development environment for the ScreenCraft Node.js SDK.

## Prerequisites

- Node.js 14.0.0 or higher
- npm or yarn package manager
- Git

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/screencraft/screencraft-nodejs-sdk.git
cd screencraft-nodejs-sdk
```

### 2. Install Dependencies

```bash
npm install
```

This will install all required dependencies including:
- `axios` - HTTP client for API requests
- `typescript` - TypeScript compiler
- `jest` - Testing framework
- `eslint` - Code linting
- `prettier` - Code formatting

### 3. Build the Project

```bash
npm run build
```

This compiles the TypeScript source code from `src/` to JavaScript in `dist/`.

## Development Commands

### Build

Compile TypeScript to JavaScript:
```bash
npm run build
```

### Test

Run the test suite:
```bash
npm test
```

Run tests with coverage:
```bash
npm test -- --coverage
```

### Lint

Check code for issues:
```bash
npm run lint
```

Auto-fix linting issues:
```bash
npm run lint -- --fix
```

### Format

Format code with Prettier:
```bash
npm run format
```

## Project Structure

```
screencraft-nodejs-sdk/
├── src/                    # Source code
│   ├── index.ts           # Main entry point
│   ├── client.ts          # ScreenCraftClient implementation
│   ├── types.ts           # TypeScript type definitions
│   ├── errors.ts          # Custom error classes
│   └── *.test.ts          # Unit tests
├── examples/              # Usage examples
│   ├── basic-screenshot.ts
│   ├── fullpage-screenshot.ts
│   ├── generate-pdf.ts
│   ├── async-webhook.ts
│   ├── error-handling.ts
│   └── mobile-screenshot.ts
├── dist/                  # Compiled output (generated)
├── node_modules/          # Dependencies (generated)
├── package.json           # Package configuration
├── tsconfig.json          # TypeScript configuration
├── tsconfig.build.json    # Build-specific TypeScript config
├── jest.config.js         # Jest test configuration
├── .eslintrc.js          # ESLint configuration
├── .prettierrc           # Prettier configuration
├── .editorconfig         # Editor configuration
├── .gitignore            # Git ignore rules
├── .npmignore            # NPM ignore rules
├── README.md             # Main documentation
├── CHANGELOG.md          # Version history
├── CONTRIBUTING.md       # Contribution guidelines
├── LICENSE               # MIT license
└── SETUP.md              # This file
```

## Development Workflow

1. **Make Changes**
   - Edit source files in `src/`
   - Follow TypeScript and ESLint rules
   - Add tests for new features

2. **Test Changes**
   ```bash
   npm test
   ```

3. **Lint and Format**
   ```bash
   npm run lint
   npm run format
   ```

4. **Build**
   ```bash
   npm run build
   ```

5. **Test in Real Project**
   ```bash
   npm link
   cd ../your-test-project
   npm link screencraft
   ```

## Testing the SDK

### Unit Tests

Run all tests:
```bash
npm test
```

Run specific test file:
```bash
npm test -- client.test.ts
```

Watch mode:
```bash
npm test -- --watch
```

### Integration Testing

1. Set your API key:
   ```bash
   export SCREENCRAFT_API_KEY="your-api-key"
   ```

2. Run examples:
   ```bash
   ts-node examples/basic-screenshot.ts
   ```

## Configuration Files

### tsconfig.json

TypeScript compiler configuration. Key settings:
- `strict: true` - Enables all strict type checking
- `target: "ES2020"` - Compiles to ES2020
- `module: "commonjs"` - Uses CommonJS modules for Node.js

### package.json

NPM package configuration. Key scripts:
- `build` - Compiles TypeScript
- `test` - Runs Jest tests
- `lint` - Runs ESLint
- `format` - Runs Prettier

### jest.config.js

Jest test framework configuration:
- Uses `ts-jest` preset for TypeScript
- Requires 80% code coverage
- Tests files matching `*.test.ts` or `*.spec.ts`

## Publishing

### Prepare for Release

1. Update version in `package.json`
2. Update `CHANGELOG.md`
3. Build the project:
   ```bash
   npm run build
   ```
4. Run tests:
   ```bash
   npm test
   ```

### Publish to NPM

```bash
npm login
npm publish
```

### Verify Publication

```bash
npm info screencraft
```

## Troubleshooting

### Build Errors

If you encounter build errors:
1. Delete `dist/` and `node_modules/`
2. Reinstall dependencies: `npm install`
3. Rebuild: `npm run build`

### Test Failures

If tests fail:
1. Ensure all dependencies are installed
2. Check Node.js version (must be 14+)
3. Clear Jest cache: `npm test -- --clearCache`

### Type Errors

If you encounter TypeScript errors:
1. Check `tsconfig.json` settings
2. Ensure all type definitions are installed
3. Run `npm install` to update types

## IDE Setup

### Visual Studio Code

Recommended extensions:
- ESLint
- Prettier
- TypeScript and JavaScript Language Features

Settings (`.vscode/settings.json`):
```json
{
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  }
}
```

### WebStorm

WebStorm has built-in support for TypeScript, ESLint, and Prettier.

Enable auto-formatting:
1. Preferences → Languages & Frameworks → JavaScript → Prettier
2. Check "On save" and "On code reformat"

## Getting Help

- Read the [README.md](README.md) for usage documentation
- Check [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines
- Review examples in the `examples/` directory
- Open an issue on GitHub for bugs or questions

## License

MIT - See [LICENSE](LICENSE) file for details.
