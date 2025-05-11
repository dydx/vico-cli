# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building the Application
```bash
# Build the CLI application
go build -o vicohome main.go

# Build with a specific version (used for releases)
go build -ldflags="-X 'github.com/dydx/vico-cli/cmd.Version=v1.0.0'" -o vicohome main.go
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./pkg/output/stdout

# Run tests with verbose output
go test -v ./...

# Run tests with coverage summary
go test -cover ./...

# Run tests with detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run tests with race condition detection
go test -race ./...

# Run all tests with coverage and race detection (as used in CI)
go test -race -coverprofile=coverage.out -covermode=atomic ./...
```

### Test Coverage Targets
- Aim for at least 80% overall code coverage
- Critical path functions should have 100% coverage
- Focus on testing edge cases and error conditions

### Docker Development
```bash
# Build the multi-architecture Docker image
docker build -f Dockerfile.multi -t vicohome:dev .

# Run the Docker image
docker run --rm vicohome:dev
```

### Releasing
```bash
# Create a new release
git tag v1.0.0
git push origin v1.0.0
```

The GitHub Actions workflow will automatically:
- Build binaries for multiple platforms (Windows, macOS, Linux on amd64 and arm64)
- Create a new release with the binaries attached
- Publish Docker images to GitHub Packages

## Architecture Overview

The Vicohome CLI is a command-line tool for interacting with the Vicohome API, primarily focused on managing and querying devices and events related to bird identification.

### Key Components

1. **Command Structure (`cmd/`)**:
   - Uses Cobra library for command-line interface
   - Root command (`cmd/root.go`) serves as the entry point and command router
   - Subcommands organized in subdirectories (`devices/`, `events/`)
   - Each command has its own implementation file (e.g., `list.go`, `get.go`)

2. **Authentication (`pkg/auth/`)**:
   - Handles API authentication with email/password credentials
   - Implements token caching to minimize authentication requests
   - Provides automatic token refresh on expiry
   - Credentials are stored in environment variables (`VICOHOME_EMAIL` and `VICOHOME_PASSWORD`)

3. **Caching (`pkg/cache/`)**:
   - Implements file-based token caching in the user's home directory (`~/.vicohome/auth.json`)
   - Handles token expiration and persistence

4. **Data Models (`pkg/models/`)**:
   - Defines structures representing API data like events and devices
   - Maps JSON responses to Go structures

5. **Output Formatting (`pkg/output/`)**:
   - Provides interfaces for different output formats (table and JSON)
   - Implements stdout handlers for displaying results

### Authentication Flow

1. User credentials are read from environment variables
2. The system first checks for a cached valid token
3. If no valid token exists, it authenticates with the API
4. Tokens are cached for future use with a 24-hour expiration
5. API requests use the token and handle auto-refresh when needed

### Command Execution Flow

1. Main entry point (`main.go`) delegates to the command executor
2. Root command routes to appropriate subcommand
3. Subcommand authenticates, makes API requests, and formats output
4. Results are displayed in the selected format (table or JSON)

## Development Guidelines

### Testing Best Practices

1. **Use Table-Driven Tests**:
   - Define a slice of test cases with inputs and expected outputs
   - Run each test case in a subtests using `t.Run()`
   - Separate test data from test logic to improve readability and maintainability
   - Group test cases by functionality and edge cases

2. **Test Helper Functions**:
   - Create helper functions for common test setup and teardown
   - Use functions to capture or redirect I/O when testing command line output
   - Isolate tests from the environment using temporary directories or mocks

3. **Cover Edge Cases**:
   - Test both successful and error paths
   - Include empty inputs, invalid inputs, and boundary conditions
   - Ensure each branch of conditional logic is tested

4. **Test Independence**:
   - Each test should be independent and not rely on state from other tests
   - Tests should clean up after themselves (e.g., remove temporary files)
   - Avoid global state that could affect other tests

### RULES

1. Always do the simplest single thing that could work
2. Use table-driven testing to improve test maintainability and coverage
3. Prefer testing small, focused units of functionality
4. Include both happy path and edge case tests
5. Properly isolate tests from external dependencies

# PRIORITIES

Your top priority right now is properly and thoroughly executing on the plans in @TESTING_PROJECT_TASKS.md

1. First, focus on making command functions mockable as described in TODO.md
2. Then, fix the skipped tests in the command packages
3. Finally, improve test coverage in low-coverage areas

Check TESTING_SUMMARY.md for an overview of current test coverage and next steps.
