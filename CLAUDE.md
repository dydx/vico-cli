# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands
- Build: `go build -o vico-cli main.go`
- Run: `./vico-cli [command] [flags]`
- Test: `go test ./...`
- Single test: `go test -v [package path] -run [TestName]`
- Run with auth: `VICOHOME_EMAIL=email VICOHOME_PASSWORD=password ./vico-cli [command]`

## Code Style Guidelines
- Formatting: Use `gofmt` for code formatting
- Imports: Group standard library, third-party, and internal imports
- Error handling: Use descriptive errors with fmt.Errorf and %w for wrapping
- Variable naming: camelCase for variables, PascalCase for exported symbols
- JSON struct tags: Use snake_case for API requests/responses
- Error checking: Check errors immediately after function calls
- Comments: Use complete sentences with proper punctuation for exported symbols
- HTTP requests: Close response bodies with defer resp.Body.Close()