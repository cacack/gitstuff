# GitStuff

A comprehensive Go CLI application for managing GitLab repositories. This tool allows you to list all repositories in your GitLab instance, clone them individually or all at once, and check their local status including current working branch.

## Features

- **List Repositories**: View all repositories with hierarchical group structure
- **Clone Management**: Download single repositories or all at once  
- **Status Tracking**: See which repos are cloned and their current branch
- **Group Structure**: Maintains exact GitLab group/subgroup organization
- **Flexible Authentication**: Supports both HTTPS and SSH cloning
- **Update Support**: Pull latest changes for already cloned repositories

## Installation

### Quick Build
```bash
go build -o gitlab-cli .
```

### Using Make (Recommended)
```bash
make build
```

### Install System-Wide
```bash
make install
```

This will install `gitlab-cli` to `/usr/local/bin` so you can run it from anywhere.

## Configuration

Before using the CLI, you need to configure your GitLab connection:

```bash
./gitstuff config
```

This will prompt you for:

- **GitLab URL**: Your GitLab instance URL (e.g., `https://gitlab.com` or just `gitlab.com`)
- **Access Token**: Your GitLab personal access token
- **Base Directory**: Local directory for cloned repositories (default: `~/gitlab-repos`)
- **SSL Certificate Verification**: Whether to skip SSL verification for self-signed certificates

> **Note**: The CLI automatically adds `https://` to URLs that don't specify a protocol.

### Self-Signed Certificates

If your GitLab instance uses self-signed certificates, you'll need to use the `--insecure` flag to skip SSL certificate verification:

```bash
./gitstuff config --insecure
```

This is common in corporate environments with internal GitLab instances.

### Creating a GitLab Access Token

1. Go to your GitLab instance
2. Navigate to User Settings > Access Tokens
3. Create a token with at least `read_repository` scope
4. Copy the token for use with the CLI

### Alternative Configuration

You can also configure using command flags:

```bash
./gitstuff config --url https://gitlab.example.com --token your-token --base-dir /path/to/repos

# For GitLab instances with self-signed certificates
./gitstuff config --url https://gitlab.example.com --token your-token --insecure
```

## Usage

### List All Repositories

```bash
# Simple list view (shows folder structure and status, no URLs)
./gitstuff list

# Tree view with group structure
./gitstuff list --tree

# Show additional details like URLs
./gitstuff list --verbose

# List without status information
./gitstuff list --status=false
```

### Clone Repositories

```bash
# Clone a specific repository
./gitstuff clone group/project-name

# Clone all repositories
./gitstuff clone --all

# Clone using SSH instead of HTTPS
./gitstuff clone --ssh group/project-name

# Update already cloned repositories
./gitstuff clone --all --update
```

## Repository Structure

The CLI maintains the exact GitLab group structure on your filesystem:

```text
~/gitlab-repos/
├── group1/
│   ├── project1/
│   ├── project2/
│   └── subgroup1/
│       └── nested-project/
├── group2/
│   └── another-project/
└── standalone-project/
```

## Repository Status Information

The CLI shows comprehensive status for each repository:

- ✅ **Cloned**: Repository exists locally and is a valid git repository
- ❌ **Not cloned**: Repository doesn't exist locally
- ⚠️ **Directory exists but not git repo**: Directory exists but isn't initialized as git
- 🔄 **Has uncommitted changes**: Repository has local modifications
- **Branch information**: Current working branch name

## Configuration File

Configuration is stored in `~/.gitlab-cli.yaml`:

```yaml
gitlab:
  url: "https://gitlab.example.com"
  token: "your-access-token"
local:
  base_dir: "/path/to/gitlab-repos"
```

## Commands Reference

### `gitstuff config`

Configure GitLab connection settings.

**Flags:**

- `-u, --url`: GitLab instance URL
- `-t, --token`: GitLab access token  
- `-d, --base-dir`: Base directory for repositories
- `-k, --insecure`: Skip SSL certificate verification (for self-signed certificates)

### `gitstuff list`

List all GitLab repositories with status information.

**Flags:**

- `-t, --tree`: Display in tree structure with groups
- `-s, --status`: Show local repository status (default: true)
- `-v, --verbose`: Show additional details like URLs

### `gitstuff clone`

Clone GitLab repositories.

**Usage:**

- `gitstuff clone [repository-path]`: Clone specific repository
- `gitstuff clone --all`: Clone all repositories

**Flags:**

- `-a, --all`: Clone all repositories
- `-s, --ssh`: Use SSH for cloning (default: HTTPS)
- `-u, --update`: Pull latest changes for existing repositories

## Examples

### Basic Workflow

```bash
# 1. Configure the CLI
./gitstuff config

# 2. List all repositories to see what's available
./gitstuff list --tree

# 3. Clone all repositories
./gitstuff clone --all

# 4. Later, update all repositories
./gitstuff clone --all --update
```

### Working with Specific Repositories

```bash
# Clone a specific project
./gitstuff clone mygroup/myproject

# Update a specific project
./gitstuff clone mygroup/myproject --update

# Use SSH for cloning
./gitstuff clone mygroup/myproject --ssh
```

## Requirements

- Go 1.19 or later
- Git installed and available in PATH
- GitLab access token with appropriate permissions
- Network access to your GitLab instance

## Testing

We have several ways to run the test suite:

### Easy Way (Recommended)

```bash
# Run all tests with clear output
make test

# Run tests with detailed output
make test-verbose
```

### Manual Way
```bash
# Test individual packages
go test ./internal/config
go test ./internal/git
go test ./internal/gitlab

# Test all packages at once
go test ./internal/config ./internal/git ./internal/gitlab
```

### Why not `go test ./...`?

The `./...` pattern in Go means "all packages in current directory and subdirectories". While it works, it can be confusing because:
- It tries to test packages that don't have tests (like `cmd/`)
- The output shows "no test files" warnings
- It's not immediately clear what's being tested

Our explicit approach tests only the packages that actually have tests, giving cleaner output.

## Error Handling

The CLI provides clear error messages for common issues:

- **Missing configuration**: Prompts to run `gitlab-cli config`
- **Invalid GitLab credentials**: Clear authentication error messages
- **Network issues**: Helpful network connectivity error messages
- **Git errors**: Detailed git operation error messages

## License

This project is open source and available under the MIT License.
