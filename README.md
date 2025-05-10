# Vicohome CLI

[![Code Coverage](https://codecov.io/gh/dydx/vico-cli/branch/main/graph/badge.svg)](https://codecov.io/gh/dydx/vico-cli)

A command-line interface tool for interacting with the Vicohome API to fetch and manage events.

## Features

- List devices and view device details
- Get event history and view detailed event information
- Identify bird species from events

## Installation

### Quick Install (Recommended)

Install the latest version with our installation script:

```bash
curl -fsSL https://raw.githubusercontent.com/dydx/vico-cli/main/scripts/install.sh | bash
```

The script will automatically detect your operating system and architecture, then install the appropriate binary.

To install a specific version:

```bash
curl -fsSL https://raw.githubusercontent.com/dydx/vico-cli/main/scripts/install.sh | bash -s v1.0.0
```

### Docker

The CLI is available as a multi-architecture Docker image (supports amd64 and arm64):

```bash
# Run the CLI with version command
docker run --rm ghcr.io/dydx/vicohome:latest version

# Run other commands
docker run --rm -e VICOHOME_EMAIL="your.email@example.com" -e VICOHOME_PASSWORD="your-password" ghcr.io/dydx/vicohome:latest devices list

# You can specify a version tag
docker run --rm ghcr.io/dydx/vicohome:v1.0.0 events list --format json
```

### Download Binary Manually

Download the pre-built binary for your platform from the [Releases page](https://github.com/dydx/vico-cli/releases).

### Build from Source

```bash
go build -o vico-cli main.go
```

## Usage

Before using this tool, set your Vicohome credentials as environment variables:

```bash
export VICOHOME_EMAIL="your.email@example.com"
export VICOHOME_PASSWORD="your-password"
```

### Devices

List all of your devices:

```bash
./vicohome devices list
```

Get details for a specific device:

```bash
./vicohome devices get [serialNumber]
```

### Events

List recent events:

```bash
./vicohome events list
```

List events for a specific time period:

```bash
./vicohome events list --hours 1
```

Get details for a specific event:

```bash
./vicohome events get [traceId]
```

## Output Formats

All commands support both table (default) and JSON output formats:

```bash
./vicohome devices list --format json
./vicohome events get [traceId] --format json
```

## Releasing a New Version

1. Tag the repository with a new version number:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. The GitHub Actions workflow will automatically:
   - Build binaries for multiple platforms (Windows, macOS, Linux on amd64 and arm64)
   - Create a new release with the binaries attached
   - Publish Docker images to GitHub Packages

## Documentation

You can view the API documentation in several ways:

### Online Documentation

The latest API documentation is automatically published to GitHub Pages:

[View Vicohome CLI Documentation](https://dydx.github.io/vico-cli/)

### Local Documentation

To generate and view documentation locally:

```bash
# Install godoc if you haven't already
go install golang.org/x/tools/cmd/godoc@latest

# Run godoc server
godoc -http=:6060

# View documentation in your browser at:
# http://localhost:6060/pkg/github.com/dydx/vico-cli/
```

## License

[MIT](LICENSE)