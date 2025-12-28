# Contributing to ScreenCraft Node.js SDK

Thank you for your interest in contributing to the ScreenCraft Node.js SDK! This document provides guidelines and instructions for contributing.

## Development Setup

1. **Fork and Clone**
   ```bash
   git clone https://github.com/your-username/screencraft-nodejs-sdk.git
   cd screencraft-nodejs-sdk
   ```

2. **Install Dependencies**
   ```bash
   npm install
   ```

3. **Build the Project**
   ```bash
   npm run build
   ```

4. **Run Tests**
   ```bash
   npm test
   ```

5. **Lint Code**
   ```bash
   npm run lint
   ```

## Project Structure

```
screencraft-nodejs-sdk/
├── src/
│   ├── index.ts          # Main entry point
│   ├── client.ts         # ScreenCraftClient class
│   ├── types.ts          # TypeScript type definitions
│   ├── errors.ts         # Error classes
│   └── *.test.ts         # Unit tests
├── examples/             # Usage examples
├── dist/                 # Compiled output (git-ignored)
├── package.json          # Package configuration
├── tsconfig.json         # TypeScript configuration
└── README.md             # Documentation
```

## Development Workflow

### 1. Create a Branch

Create a feature branch from `main`:
```bash
git checkout -b feature/your-feature-name
```

Use descriptive branch names:
- `feature/add-retry-logic`
- `fix/timeout-handling`
- `docs/update-readme`

### 2. Make Changes

- Write clean, readable code
- Follow existing code style
- Add TypeScript types for all new code
- Update tests for your changes
- Update documentation as needed

### 3. Test Your Changes

Run the full test suite:
```bash
npm test
```

Run linting:
```bash
npm run lint
```

Build the project:
```bash
npm run build
```

### 4. Commit Your Changes

Write clear, descriptive commit messages:
```bash
git commit -m "Add retry logic with exponential backoff"
```

Good commit messages:
- Start with a verb in imperative mood (Add, Fix, Update, Remove)
- Keep the first line under 72 characters
- Add detailed description if needed

### 5. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a pull request on GitHub with:
- Clear title describing the change
- Description of what was changed and why
- Reference any related issues

## Code Style Guidelines

### TypeScript

- Use strict TypeScript mode
- Always define explicit types for function parameters and return values
- Avoid `any` - use `unknown` if type is truly unknown
- Use interfaces for object shapes
- Use type aliases for unions and intersections
- Document complex types with JSDoc comments

### Naming Conventions

- Classes: PascalCase (e.g., `ScreenCraftClient`)
- Interfaces/Types: PascalCase (e.g., `ScreenshotOptions`)
- Functions/Methods: camelCase (e.g., `screenshot()`)
- Constants: UPPER_SNAKE_CASE (e.g., `DEFAULT_TIMEOUT`)
- Private methods: prefix with underscore (e.g., `_validateUrl()`)

### Code Organization

- One class per file
- Group related types in `types.ts`
- Keep functions focused and single-purpose
- Extract complex logic into private methods
- Add JSDoc comments for public APIs

### Error Handling

- Use custom error classes from `errors.ts`
- Provide meaningful error messages
- Include context in error details
- Never swallow errors silently

## Testing Guidelines

### Writing Tests

- Write tests for all new features
- Test both success and failure cases
- Test edge cases and boundary conditions
- Use descriptive test names
- Keep tests isolated and independent

### Test Structure

```typescript
describe('FeatureName', () => {
  describe('methodName', () => {
    it('should do something specific', () => {
      // Arrange
      const input = 'test';

      // Act
      const result = someFunction(input);

      // Assert
      expect(result).toBe('expected');
    });
  });
});
```

## Documentation Guidelines

### Code Documentation

- Add JSDoc comments for all public APIs
- Include parameter descriptions
- Document return types
- Add usage examples for complex features
- Explain "why" not just "what"

### README Updates

When adding features:
- Add usage example to README
- Update API reference section
- Add to feature list if applicable
- Update table of contents if needed

## Pull Request Process

1. **Before Submitting**
   - Ensure all tests pass
   - Run linting and fix issues
   - Update documentation
   - Add changelog entry if applicable

2. **PR Description**
   - Describe what changed and why
   - Include screenshots for UI changes
   - Link related issues
   - Note any breaking changes

3. **Review Process**
   - Respond to feedback promptly
   - Make requested changes
   - Keep discussions focused and professional

4. **After Approval**
   - Squash commits if requested
   - Ensure CI passes
   - Wait for maintainer to merge

## Release Process

Maintainers will handle releases:

1. Update version in `package.json`
2. Update `CHANGELOG.md`
3. Create git tag
4. Publish to npm
5. Create GitHub release

## Getting Help

- Check existing issues and PRs
- Read the documentation thoroughly
- Ask questions in issue comments
- Join community discussions

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Help others learn and grow
- Focus on what's best for the project

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
