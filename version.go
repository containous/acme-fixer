package main

import (
	"fmt"
	"runtime"
)

var (
	// Version holds the current version.
	Version = "dev"
	// BuildDate holds the build date.
	BuildDate = "I don't remember exactly"
	// ShortCommit shot commit
	ShortCommit = ""
)

// displayVersion Display acme-fixer version.
func displayVersion() {
	fmt.Printf(`acme-fixer:
 version     : %s
 commit      : %s
 build date  : %s
 go version  : %s
 go compiler : %s
 platform    : %s/%s
`, Version, ShortCommit, BuildDate, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
}
