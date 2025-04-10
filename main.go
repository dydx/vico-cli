// Package main implements the entry point for the Vicohome CLI application.
//
// This application provides a command-line interface for interacting with the Vicohome API,
// allowing users to manage and query Vicohome devices and events.
package main

import "github.com/dydx/vico-cli/cmd"

func main() {
	cmd.Execute()
}
