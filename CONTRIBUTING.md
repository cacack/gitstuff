# Contributing to GitStuff

Thank you for your interest in contributing to GitStuff! This document provides guidelines for contributing to the project.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/neilfarmer/gitstuff.git
   cd gitstuff
   ```
3. **Install dependencies**:
   ```bash
   go mod download
   ```

## Development Workflow

### Setting up Development Environment

1. **Go version**: Ensure you have Go 1.21 or later installed
2. **Build the project**:
   ```bash
   make build
   ```
3. **Run tests**:
   ```bash
   make test
   ```

### Making Changes

1. **Create a feature branch**:

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Write your code** following the project conventions:

   - Follow Go best practices and idioms
   - Keep functions focused and single-purpose
   - Use descriptive variable and function names
   - **Do NOT add comments** unless absolutely necessary

3. **Write tests** for your changes:

   - Every new feature must have tests
   - Every new flag or configuration option needs test coverage
   - Tests should cover both success and error cases
   - Use table-driven tests for multiple test scenarios

4. **Test your changes**:

   ```bash
   # Run all tests
   make test

   # Run tests with verbose output
   make test-verbose

   # Build and test the binary
   make build
   ./gitstuff --help
   ```

5. **Lint your code**:
   ```bash
   golangci-lint run
   ```

### Project Structure

```
gitstuff/
├── cmd/                    # CLI commands (Cobra)
├── internal/
│   ├── config/            # Configuration management
│   ├── gitlab/            # GitLab API client
│   └── git/               # Git operations
├── .github/workflows/     # GitHub Actions
└── README.md
```

### Testing Standards

**CRITICAL**: Every feature and flag must have tests. This project maintains excellent test coverage.

- All `internal/` packages must have comprehensive test coverage
- All `cmd/` packages must have tests for display logic and flag behavior
- Test edge cases like missing configuration, invalid URLs, git errors
- Test all configuration combinations
- Test CLI flag parsing and validation
- Test display output for different flag combinations

### Code Style

- Follow standard Go conventions
- Use `gofmt` and `goimports`
- No comments unless specifically requested
- Use descriptive names over comments
- Keep functions focused and single-purpose
- Handle errors gracefully with meaningful messages

### Commit Messages

Write clear, descriptive commit messages:

```
Add verbose flag to list command

- Hides URLs by default for cleaner output
- Shows URLs with --verbose/-v flag
- Maintains backward compatibility
- Includes comprehensive tests
```

## Submitting Changes

1. **Push your changes**:

   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a Pull Request** on GitHub with:

   - Clear description of the changes
   - Reference to any related issues
   - Screenshots if UI changes are involved

3. **Ensure CI passes**:
   - All tests must pass
   - Code must pass linting
   - Build must succeed

## Release Process

Releases are automated via GitHub Actions:

1. **Tag the release**:

   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```

2. **GitHub Actions automatically**:
   - Runs all tests
   - Builds cross-platform binaries
   - Creates GitHub release
   - Generates install scripts

## Getting Help

- **Issues**: Create a [GitHub issue](https://github.com/neilfarmer/gitstuff/issues)
- **Discussions**: Use [GitHub Discussions](https://github.com/neilfarmer/gitstuff/discussions)
- **Questions**: Tag your issue with "question"

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help maintain a welcoming environment
- Follow GitHub's community guidelines

Thank you for contributing to GitStuff!
