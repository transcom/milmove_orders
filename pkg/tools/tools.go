// +build tools

// This file exists to track tool dependencies. This is one of the recommended practices
// for handling tool dependencies in a Go module as outlined here:
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package tools

import (
	// Install for hot reloading server
	_ "github.com/codegangsta/gin"
	// Install for generating swagger code
	_ "github.com/go-swagger/go-swagger/cmd/swagger"

	// Test packages
	_ "github.com/stretchr/objx"
)
