# GitLab CLI - Claude Code Instructions

This file contains instructions for Claude Code to help maintain and work with this GitLab CLI project effectively.

## Project Overview

This is a comprehensive Go CLI application for managing GitLab repositories. It allows users to:
- List all repositories from their GitLab instance with hierarchical group structure
- Clone repositories individually or in bulk (HTTPS/SSH support)
- Check local repository status including current branch and uncommitted changes
- Maintain exact GitLab filesystem organization with groups/subgroups

## Project Structure

```
gitlab-cli/
├── main.go                    # Application entry point
├── cmd/                       # CLI commands using Cobra framework
│   ├── root.go               # Root command configuration
│   ├── config.go             # Configuration management command
│   ├── list.go               # Repository listing command
│   ├── clone.go              # Repository cloning command
│   └── version.go            # Version information command
├── internal/                 # Internal packages (not for external use)
│   ├── config/               # Configuration logic and tests
│   │   ├── config.go         # Config loading, validation, creation
│   │   └── config_test.go    # Comprehensive config tests
│   ├── gitlab/               # GitLab API client and tests
│   │   ├── client.go         # GitLab API integration with URL normalization
│   │   └── client_test.go    # GitLab client tests including URL validation
│   └── git/                  # Git operations and tests
│       ├── operations.go     # Git status, clone, pull operations
│       └── operations_test.go # Git operations tests
├── Makefile                  # Build and test automation
├── README.md                 # Comprehensive documentation
├── go.mod                    # Go module configuration
└── .gitignore               # Git ignore rules
```

## Key Technologies Used

- **CLI Framework**: Cobra (github.com/spf13/cobra)
- **Configuration**: Viper (github.com/spf13/viper) with YAML
- **GitLab API**: go-gitlab (github.com/xanzy/go-gitlab)
- **Testing**: Go standard testing package
- **Build**: Makefile for clear build/test commands

## Development Guidelines

### Building the Application

**Preferred method:**
```bash
make build
```

**Alternative:**
```bash
go build -o gitlab-cli .
```

### Running Tests

**IMPORTANT**: Do NOT use `go test ./...` as the user finds this confusing because it shows "no test files" warnings for packages without tests.

**Preferred methods:**
```bash
make test           # Clean output, tests only packages with tests
make test-verbose   # Detailed output when debugging
```

**Manual method (if needed):**
```bash
go test ./cmd ./internal/config ./internal/git ./internal/gitlab
```

**Why this approach:** 
- Only tests packages that actually have tests
- Provides clean, clear output
- No confusing warnings about missing test files


### Code Organization

- **Commands**: All CLI commands go in `cmd/` directory
- **Internal Logic**: Business logic goes in `internal/` packages
- **Tests**: Every package with logic must have comprehensive tests
- **URL Handling**: The GitLab client automatically normalizes URLs (adds https:// if missing)

### Testing Requirements - CRITICAL

**EVERY FEATURE AND FLAG MUST HAVE TESTS** - This project maintains excellent test coverage and you must continue this standard.

- All `internal/` packages must have comprehensive test coverage
- **All `cmd/` packages must have tests for display logic and flag behavior**
- **Every new feature requires corresponding tests before implementation**
- **Every new flag or configuration option needs test coverage**
- Tests should cover both success and error cases
- Use table-driven tests for multiple test scenarios
- Test edge cases like missing configuration, invalid URLs, git errors
- Test all configuration combinations (with/without insecure flag, different URLs, etc.)
- Test CLI flag parsing and validation
- Test interactive prompts and user input handling
- **Test display output for different flag combinations (verbose, tree, status)**

### Configuration Management

- Configuration is stored in `~/.gitstuff.yaml`
- Uses Viper for config loading with environment variable support
- Sensitive data (tokens) should be handled securely with proper file permissions (0600)
- Default base directory is `~/gitlab-repos` if not specified
- Supports `insecure` flag for skipping SSL certificate verification (self-signed certificates)

### Error Handling

- Provide clear, actionable error messages
- Use fmt.Errorf with error wrapping for context
- Guide users to solutions (e.g., "run 'gitlab-cli config' first")
- Handle network issues, authentication failures, and git errors gracefully

### Code Quality - MANDATORY CHECKS

**CRITICAL**: All code changes must pass the following quality checks before submission. This is non-negotiable.

**ALWAYS run ALL quality checks after making ANY change:**
```bash
make test      # Run all tests - MUST pass
make format    # Format all code - MANDATORY  
make lint      # Run linting - MUST pass with zero issues

# OR use the convenient combined command:
make quality   # Runs test + format + lint automatically
```

**Quality Check Requirements:**

1. **Testing (MANDATORY)**:
   - ALL tests must pass: `make test`
   - No test failures or errors allowed
   - Tests cover all new functionality and edge cases
   - Test coverage must be maintained or improved

2. **Code Formatting (MANDATORY)**:
   - **gofmt**: All Go code must pass `gofmt -l .` (no output = properly formatted)
   - **goimports**: All import statements must be properly organized and formatted
   - Use `make format` to auto-fix formatting issues
   - Check formatting: `make check-format`

3. **Linting (MANDATORY)**:
   - **golangci-lint**: All code must pass `make lint` with ZERO issues
   - No linting errors, warnings, or suggestions allowed
   - Fix all issues before proceeding - no exceptions
   - Common issues: unused variables, ineffectual assignments, error handling

**Workflow for ANY Code Change:**
```bash
# 1. Make your changes
# 2. ALWAYS run the quality checks:
make quality   # ← Runs test + format + lint (RECOMMENDED)
# OR run individually:
make test      # ← Must pass
make format    # ← Must run  
make lint      # ← Must pass with zero issues
# 3. Only then commit/submit
```

**Failure to follow this process will result in rejected changes.**

**General Code Style:**
- Follow standard Go conventions
- DO NOT ADD COMMENTS unless specifically requested by the user
- Use descriptive variable and function names
- Keep functions focused and single-purpose
- Proper error handling with meaningful messages

### Adding New Features

1. Add the command in `cmd/` directory following existing patterns
2. Implement business logic in appropriate `internal/` package
3. Add comprehensive tests for new functionality
4. Update README.md with new command documentation
5. Test using `make test` to ensure nothing breaks

### Common Tasks

**Adding a new command:**
1. Create new file in `cmd/` directory
2. Follow existing patterns (see `cmd/list.go` or `cmd/clone.go`)
3. Add command to root command in `init()` function
4. Add tests if the command has complex logic

**Modifying GitLab API calls:**
- All GitLab API logic goes in `internal/gitlab/client.go`
- URLs are automatically normalized (https:// prefix added if missing)
- Handle pagination for large result sets
- Add appropriate error handling and retry logic
- Support for insecure connections (self-signed certificates) via insecure flag

**Modifying Git operations:**
- All Git operations go in `internal/git/operations.go`
- Use `exec.Command` with proper error handling
- Always check if Git is available before operations
- Handle common Git errors gracefully

## Dependencies

- **Required**: Go 1.19 or later
- **External**: git command must be available in PATH
- **Network**: Access to GitLab instance for API calls
- **Permissions**: GitLab access token with read_repository scope minimum

## User Experience Priorities

- Clear, intuitive command structure
- Helpful error messages with actionable guidance
- Consistent output formatting
- Respect for GitLab's group/subgroup structure in filesystem layout
- Support for both beginners and power users (flags, options)

## Maintenance Notes

- The `go-gitlab` dependency is deprecated but functional - consider migration in future
- URL normalization handles common user input mistakes
- File permissions are important for security (config file is 0600)
- The application is designed to be stateless - no persistent data beyond config

## When Working on This Project

1. **MANDATORY: Run quality checks after EVERY change**: `make quality` (or `make test && make format && make lint`)
2. **All tests must pass**: Use `make test` (never `go test ./...`)
3. **Zero linting issues**: `make lint` must pass with no errors, warnings, or suggestions
4. **Code must be formatted**: `make format` before any commits
5. **Build using `make build`** for consistency
6. **Update README.md** if adding new features or changing behavior
7. **Consider security implications** especially for token handling
8. **Maintain the clean architecture** with separate concerns in `internal/` packages
9. **Test edge cases** like network failures, invalid configs, missing git repos

**Remember: `make quality` = MANDATORY for every code change**
