# ScreenCraft Node.js SDK - Complete Index

Welcome to the ScreenCraft Node.js/TypeScript SDK! This document serves as a comprehensive guide to all available resources.

## Quick Navigation

### Getting Started
- [QUICKSTART.md](QUICKSTART.md) - Get up and running in 5 minutes
- [README.md](README.md) - Complete usage documentation
- [SETUP.md](SETUP.md) - Development environment setup

### Development
- [CONTRIBUTING.md](CONTRIBUTING.md) - How to contribute
- [CHANGELOG.md](CHANGELOG.md) - Version history
- [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) - Complete project overview

### Examples
All examples are located in the [examples/](examples/) directory:
- [basic-screenshot.ts](examples/basic-screenshot.ts) - Simple screenshot
- [fullpage-screenshot.ts](examples/fullpage-screenshot.ts) - Full-page capture
- [generate-pdf.ts](examples/generate-pdf.ts) - PDF generation
- [async-webhook.ts](examples/async-webhook.ts) - Async operations
- [error-handling.ts](examples/error-handling.ts) - Error handling
- [mobile-screenshot.ts](examples/mobile-screenshot.ts) - Mobile emulation
- [README.md](examples/README.md) - Examples documentation

## Documentation Structure

### Main Documentation (Read First)

1. **QUICKSTART.md**
   - Installation
   - Basic usage
   - Complete example
   - Environment setup
   - Common use cases

2. **README.md**
   - Features overview
   - Installation instructions
   - Comprehensive examples
   - API reference
   - Error handling
   - Rate limiting
   - TypeScript support

3. **SETUP.md**
   - Development prerequisites
   - Installation steps
   - Build commands
   - Project structure
   - Testing guide
   - IDE configuration

### Developer Documentation

4. **CONTRIBUTING.md**
   - Development setup
   - Workflow guidelines
   - Code style
   - Testing guidelines
   - Pull request process
   - Release process

5. **PROJECT_SUMMARY.md**
   - Complete project overview
   - File structure
   - Architecture details
   - Core features
   - Design patterns
   - Dependencies

6. **CHANGELOG.md**
   - Version history
   - Release notes
   - Breaking changes
   - New features

## Source Code

### Core Files (`src/`)

| File | Description | Lines |
|------|-------------|-------|
| [index.ts](src/index.ts) | Main entry point, exports all public APIs | ~40 |
| [client.ts](src/client.ts) | ScreenCraftClient class, HTTP communication | ~350 |
| [types.ts](src/types.ts) | TypeScript type definitions and interfaces | ~200 |
| [errors.ts](src/errors.ts) | Custom error classes | ~150 |
| [client.test.ts](src/client.test.ts) | Unit tests | ~200 |

### Configuration Files

| File | Purpose |
|------|---------|
| [package.json](package.json) | NPM package configuration |
| [tsconfig.json](tsconfig.json) | TypeScript compiler settings |
| [tsconfig.build.json](tsconfig.build.json) | Build-specific TypeScript config |
| [jest.config.js](jest.config.js) | Jest test framework configuration |
| [.eslintrc.js](.eslintrc.js) | ESLint code quality rules |
| [.prettierrc](.prettierrc) | Prettier code formatting |
| [.editorconfig](.editorconfig) | Editor configuration |
| [.gitignore](.gitignore) | Git ignore patterns |
| [.npmignore](.npmignore) | NPM publish ignore patterns |

### CI/CD

| File | Purpose |
|------|---------|
| [.github/workflows/ci.yml](.github/workflows/ci.yml) | GitHub Actions CI/CD pipeline |

### Utilities

| File | Purpose |
|------|---------|
| [verify-setup.js](verify-setup.js) | Setup verification script |

## API Reference Quick Links

### Main Class

**ScreenCraftClient** - Main API client class

Constructor:
```typescript
new ScreenCraftClient(config: ScreenCraftConfig)
```

Methods:
- `screenshot(options: ScreenshotOptions): Promise<SyncResponse>`
- `screenshotAsync(options: AsyncScreenshotOptions): Promise<AsyncResponse>`
- `pdf(options: PDFOptions): Promise<SyncResponse>`
- `pdfAsync(options: AsyncPDFOptions): Promise<AsyncResponse>`

### Type Definitions

**Core Types**:
- `ScreenCraftConfig` - Client configuration
- `ScreenshotOptions` - Screenshot parameters
- `PDFOptions` - PDF generation parameters
- `SyncResponse` - Synchronous response
- `AsyncResponse` - Asynchronous job response
- `Viewport` - Viewport configuration
- `WebhookConfig` - Webhook settings

**Image & Format Types**:
- `ImageFormat` - 'png' | 'jpeg' | 'webp'
- `PDFFormat` - 'A4' | 'Letter' | 'Legal' | etc.

### Error Classes

- `ScreenCraftError` - Base error class
- `AuthenticationError` - API key issues (401/403)
- `ValidationError` - Invalid parameters (400)
- `RateLimitError` - Rate limiting (429)
- `CaptureError` - Capture failures
- `NetworkError` - Network issues
- `TimeoutError` - Request timeouts
- `ServerError` - Server errors (5xx)
- `UnknownError` - Unexpected errors

## Usage Examples

### Installation
```bash
npm install screencraft
```

### Basic Screenshot
```typescript
import { ScreenCraftClient } from 'screencraft';

const client = new ScreenCraftClient({ apiKey: 'your-api-key' });
const screenshot = await client.screenshot({
  url: 'https://example.com',
  format: 'png'
});
```

### PDF Generation
```typescript
const pdf = await client.pdf({
  url: 'https://example.com',
  format: 'A4',
  printBackground: true
});
```

### Error Handling
```typescript
import { RateLimitError, ValidationError } from 'screencraft';

try {
  await client.screenshot({ url: 'https://example.com' });
} catch (error) {
  if (error instanceof RateLimitError) {
    console.log('Rate limit exceeded');
  }
}
```

## Development Workflow

1. **Install Dependencies**
   ```bash
   npm install
   ```

2. **Build Project**
   ```bash
   npm run build
   ```

3. **Run Tests**
   ```bash
   npm test
   ```

4. **Lint Code**
   ```bash
   npm run lint
   ```

5. **Format Code**
   ```bash
   npm run format
   ```

6. **Verify Setup**
   ```bash
   node verify-setup.js
   ```

## Testing

### Run All Tests
```bash
npm test
```

### Run with Coverage
```bash
npm test -- --coverage
```

### Run Specific Test
```bash
npm test -- client.test.ts
```

### Watch Mode
```bash
npm test -- --watch
```

## Common Tasks

### Add New Feature
1. Create branch: `git checkout -b feature/my-feature`
2. Add implementation in `src/`
3. Add tests in `src/*.test.ts`
4. Update documentation
5. Run tests: `npm test`
6. Build: `npm run build`
7. Create pull request

### Fix Bug
1. Create branch: `git checkout -b fix/bug-description`
2. Write failing test first
3. Implement fix
4. Verify all tests pass
5. Update CHANGELOG.md
6. Create pull request

### Update Documentation
1. Edit relevant `.md` file
2. Ensure examples work
3. Update version if needed
4. Submit changes

## Publishing

### Prepare Release
1. Update `package.json` version
2. Update `CHANGELOG.md`
3. Build: `npm run build`
4. Test: `npm test`

### Publish to NPM
```bash
npm login
npm publish
```

### Create GitHub Release
1. Tag version: `git tag v1.0.0`
2. Push tag: `git push origin v1.0.0`
3. GitHub Actions will auto-publish

## Support & Resources

### Official Resources
- **API Docs**: https://screencraftapi.com/docs
- **GitHub**: https://github.com/screencraft/screencraft-nodejs-sdk
- **NPM**: https://www.npmjs.com/package/screencraft

### Getting Help
- Check documentation first
- Search existing GitHub issues
- Read examples in `examples/` directory
- Contact: support@screencraftapi.com

### Community
- Report bugs via GitHub Issues
- Request features via GitHub Issues
- Contribute via Pull Requests
- Follow contribution guidelines in CONTRIBUTING.md

## File Count Summary

- **Source files**: 5 TypeScript files
- **Examples**: 6 example files + README
- **Documentation**: 8 markdown files
- **Configuration**: 10 config files
- **Tests**: 1 test file (expandable)
- **Total files**: 30+ files

## Lines of Code

- **Source code**: ~950 lines
- **Tests**: ~200 lines
- **Examples**: ~450 lines
- **Documentation**: ~2,500 lines
- **Total**: ~4,100 lines

## Technology Stack

- **Language**: TypeScript 5.0
- **Runtime**: Node.js 14+
- **HTTP Client**: Axios 1.6+
- **Testing**: Jest 29
- **Linting**: ESLint 8
- **Formatting**: Prettier 3
- **CI/CD**: GitHub Actions

## Next Steps

1. **For Users**: Read [QUICKSTART.md](QUICKSTART.md)
2. **For Developers**: Read [SETUP.md](SETUP.md)
3. **For Contributors**: Read [CONTRIBUTING.md](CONTRIBUTING.md)
4. **For Examples**: Check [examples/README.md](examples/README.md)

## License

MIT License - See [LICENSE](LICENSE) file

## Version

Current version: **1.0.0**

See [CHANGELOG.md](CHANGELOG.md) for version history.

---

**Last Updated**: December 28, 2025
**Maintained By**: ScreenCraft Team
