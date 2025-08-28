# GitStuff

[![CI](https://github.com/neilfarmer/gitstuff/actions/workflows/ci.yml/badge.svg)](https://github.com/neilfarmer/gitstuff/actions/workflows/ci.yml)
[![Release](https://github.com/neilfarmer/gitstuff/actions/workflows/release.yml/badge.svg)](https://github.com/neilfarmer/gitstuff/actions/workflows/release.yml)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/neilfarmer/gitstuff)](https://github.com/neilfarmer/gitstuff/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/neilfarmer/gitstuff)](https://golang.org/)

A comprehensive Go CLI application for managing GitLab repositories. This tool allows you to list all repositories in your GitLab instance, clone them individually or all at once, and check their local status including current working branch.

## Quick Start

**Linux x86_64:**
```bash
curl -L https://github.com/neilfarmer/gitstuff/releases/latest/download/gitstuff-linux-amd64 -o gitstuff
chmod +x gitstuff
sudo mv gitstuff /usr/local/bin/
```

**Linux ARM64:**
```bash
curl -L https://github.com/neilfarmer/gitstuff/releases/latest/download/gitstuff-linux-arm64 -o gitstuff
chmod +x gitstuff
sudo mv gitstuff /usr/local/bin/
```

**macOS x86_64:**
```bash
curl -L https://github.com/neilfarmer/gitstuff/releases/latest/download/gitstuff-darwin-amd64 -o gitstuff
chmod +x gitstuff
sudo mv gitstuff /usr/local/bin/
```

**macOS ARM64 (M1/M2):**
```bash
curl -L https://github.com/neilfarmer/gitstuff/releases/latest/download/gitstuff-darwin-arm64 -o gitstuff
chmod +x gitstuff
sudo mv gitstuff /usr/local/bin/
```

**Then configure and use:**
```bash
# Configure your GitLab connection
gitstuff config

# List your repositories
gitstuff list

# Clone all repositories
gitstuff clone --all
```

## Features

- **List Repositories**: View all repositories with hierarchical group structure
- **Group Filtering**: Filter repositories by GitLab group to focus on specific teams or projects
- **Clone Management**: Download single repositories or all at once
- **Status Tracking**: See which repos are cloned and their current branch
- **Group Structure**: Maintains exact GitLab group/subgroup organization
- **Flexible Authentication**: Supports both HTTPS and SSH cloning
- **Update Support**: Pull latest changes for already cloned repositories

## Installation

Choose the appropriate command for your platform:

**Linux x86_64:**

```bash
curl -L https://github.com/neilfarmer/gitstuff/releases/latest/download/gitstuff-linux-amd64 -o gitstuff
chmod +x gitstuff
sudo mv gitstuff /usr/local/bin/
```

**Linux ARM64:**

```bash
curl -L https://github.com/neilfarmer/gitstuff/releases/latest/download/gitstuff-linux-arm64 -o gitstuff
chmod +x gitstuff
sudo mv gitstuff /usr/local/bin/
```

**macOS x86_64:**

```bash
curl -L https://github.com/neilfarmer/gitstuff/releases/latest/download/gitstuff-darwin-amd64 -o gitstuff
chmod +x gitstuff
sudo mv gitstuff /usr/local/bin/
```

**macOS ARM64 (M1/M2):**

```bash
curl -L https://github.com/neilfarmer/gitstuff/releases/latest/download/gitstuff-darwin-arm64 -o gitstuff
chmod +x gitstuff
sudo mv gitstuff /usr/local/bin/
```

Or download directly from the [releases page](https://github.com/neilfarmer/gitstuff/releases/latest).

### Build from Source

**Prerequisites:** Go 1.21 or later

```bash
# Clone the repository
git clone https://github.com/neilfarmer/gitstuff.git
cd gitstuff

# Build using make
make build

# Or build directly
go build -o gitstuff .

# Install system-wide
make install
```

## Configuration

Before using the CLI, you need to configure your GitLab connection:

```bash
gitstuff config
```

This will prompt you for:

- **GitLab URL**: Your GitLab instance URL (e.g., `https://gitlab.com` or just `gitlab.com`)
- **Access Token**: Your GitLab personal access token
- **Base Directory**: Local directory for cloned repositories (default: `~/gitlab-repos`)
- **SSL Certificate Verification**: Whether to skip SSL verification for self-signed certificates
- **Default Group Filter**: Optional group to filter repositories by default (e.g., `team-a` or `team-a/backend`)

> **Note**: The CLI automatically adds `https://` to URLs that don't specify a protocol.

### Creating a GitLab Access Token

1. Go to your GitLab instance
2. Navigate to User Settings > Access Tokens
3. Create a token with at least `read_repository` scope
4. Copy the token for use with the CLI

### Alternative Configuration

You can also configure using command flags:

```bash
gitstuff config --url https://gitlab.example.com --token your-token --base-dir /path/to/repos

# For GitLab instances with self-signed certificates
gitstuff config --url https://gitlab.example.com --token your-token --insecure
```

## Usage

### List All Repositories

```bash
# Simple list view (shows folder structure and status, no URLs)
gitstuff list

# Tree view with group structure
gitstuff list --tree

# Show additional details like URLs
gitstuff list --verbose

# List without status information
gitstuff list --status=false
```

### Clone Repositories

```bash
# Clone a specific repository
gitstuff clone group/project-name

# Clone all repositories
gitstuff clone --all

# Clone using SSH instead of HTTPS
gitstuff clone --ssh group/project-name

# Update already cloned repositories
gitstuff clone --all --update
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

Configuration is stored in `~/.gitstuff.yaml`:

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
- `-g, --group`: Default group to filter repositories (optional)

### `gitstuff list`

List all GitLab repositories with status information.

**Flags:**

- `-t, --tree`: Display in tree structure with groups
- `-s, --status`: Show local repository status (default: true)
- `-v, --verbose`: Show additional details like URLs
- `-g, --group`: Filter repositories to only those in the specified group

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
gitstuff config

# 2. List all repositories to see what's available
gitstuff list --tree

# 3. Clone all repositories
gitstuff clone --all

# 4. Later, update all repositories
gitstuff clone --all --update
```

### Group Filtering

Filter repositories by GitLab group to focus on specific teams or projects:

```bash
# List repositories only in the "backend" group
gitstuff list --group backend

# List repositories in a nested group
gitstuff list --group team-a/backend

# Use tree view with group filtering
gitstuff list --tree --group team-a

# Set a default group in config (applies to all commands)
gitstuff config --group team-a

# Override config default with command flag
gitstuff list --group team-b
```

### Working with Specific Repositories

```bash
# Clone a specific project
gitstuff clone mygroup/myproject

# Update a specific project
gitstuff clone mygroup/myproject --update

# Use SSH for cloning
gitstuff clone mygroup/myproject --ssh
```

## Requirements

- Go 1.19 or later
- Git installed and available in PATH
- GitLab access token with appropriate permissions
- Network access to your GitLab instance

## Testing

```bash
# Run all tests with clear output
make test

# Run tests with detailed output
make test-verbose
```

## Error Handling

The CLI provides clear error messages for common issues:

- **Missing configuration**: Prompts to run `gitstuff config`
- **Invalid GitLab credentials**: Clear authentication error messages
- **Network issues**: Helpful network connectivity error messages
- **Git errors**: Detailed git operation error messages

## Releases

### Creating a Release

To create a new release:

1. **Tag the release:**

   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```

2. **GitHub Actions automatically:**

   - Runs all tests
   - Builds binaries for all platforms
   - Creates GitHub release with download links
   - Generates install scripts

3. **Release artifacts include:**
   - Cross-platform binaries (Linux, macOS, Windows)
   - Architecture support (x64, ARM64)
   - Automated install scripts
   - SHA256 checksums for verification

### Versioning

This project follows [Semantic Versioning](https://semver.org/):

- `v1.2.3` - Major.Minor.Patch
- `v1.2.0-beta.1` - Pre-release versions

## License

This project is open source and available under the MIT License.
