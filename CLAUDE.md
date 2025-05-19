# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Build
```bash
# Build for current platform
go build -o vico-cli main.go

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o vico-cli-linux-amd64 main.go
GOOS=darwin GOARCH=arm64 go build -o vico-cli-darwin-arm64 main.go
GOOS=windows GOARCH=amd64 go build -o vico-cli-windows-amd64.exe main.go
```

### Lint
```bash
# Check formatting - the CI/CD pipeline will fail if there are formatting issues
gofmt -l .

# Fix formatting issues
gofmt -w .

# Run golint (install first: go install golang.org/x/lint/golint@latest)
golint -set_exit_status ./...
```

### Test
```bash
# Run all tests
go test ./...
```

### Release
To create a new release:
```bash
git tag v1.0.0
git push origin v1.0.0
```
The GitHub Actions workflow will automatically build binaries for all platforms and create a release.

## Architecture Overview

### Command Structure
The CLI uses the Cobra library with a hierarchical command structure:
- Root command (`cmd/root.go`)
  - `devices` subcommand (`cmd/devices/`)
    - `list` - List all devices
    - `get` - Get details for a specific device
  - `events` subcommand (`cmd/events/`)
    - `list` - List recent events
    - `get` - Get details for a specific event
    - `search` - Search events by field

### Key Packages
- `pkg/auth/` - Handles API authentication with automatic token refresh and caching
- `pkg/cache/` - Token caching to minimize authentication requests
- `pkg/models/` - Data models for API responses (Event, Device structures)
- `pkg/output/` - Output formatting (table and JSON formats)

### Authentication Flow
1. Credentials are read from environment variables: `VICOHOME_EMAIL` and `VICOHOME_PASSWORD`
2. Token is obtained from the API and cached for future requests
3. Token is automatically refreshed when expired (based on error codes -1024 to -1027)
4. Uses authentication middleware pattern in `pkg/auth/auth.go`

### Testing Approach
The `TESTING.md` file contains manual acceptance tests that define the expected CLI interface. 
When implementing new features, tests are marked as `[FAIL]` and implementation continues until
all tests pass `[PASS]`.

### CI/CD Pipeline
GitHub Actions workflow (`.github/workflows/workflow.yml`) runs on:
- Pull requests: Linting and build tests
- Main branch pushes: Multi-arch Docker builds, documentation generation
- Tag pushes: Full release with binaries for all platforms

The workflow ensures:
- Code is properly formatted (`gofmt`)
- Code passes linting (`golint`)
- Tests pass
- Multi-platform builds succeed