// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package create

import "fmt"

func buildCLI(dir, modPath, name string) error {
	writeFile(dir, "main.go", fmt.Sprintf(`package main

import (
	"fmt"
	"os"

	"%s/cmd"
)

func main() {
	if err := cmd.Execute(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %%v\n", err)
		os.Exit(1)
	}
}
`, modPath))

	writeFile(dir, "cmd/root.go", fmt.Sprintf(`package cmd

import (
	"flag"
	"fmt"
)

var version = "0.1.0"

func Execute(args []string) error {
	if len(args) == 0 {
		return runHelp()
	}

	switch args[0] {
	case "version":
		fmt.Printf("%s v%%s\n", version)
		return nil
	case "greet":
		return runGreet(args[1:])
	case "help", "-h", "--help":
		return runHelp()
	default:
		return fmt.Errorf("unknown command: %%s\nRun '%s help' for usage", args[0])
	}
}

func runHelp() error {
	fmt.Println(`+"`"+`%s - A command-line tool

Usage:
  %s <command> [arguments]

Commands:
  greet    Print a greeting
  version  Print version
  help     Show this help`+"`"+`)
	return nil
}

func runGreet(args []string) error {
	fs := flag.NewFlagSet("greet", flag.ExitOnError)
	nameFlag := fs.String("name", "World", "name to greet")
	if err := fs.Parse(args); err != nil {
		return err
	}
	fmt.Printf("Hello, %%s!\n", *nameFlag)
	return nil
}
`, name, name, name, name))

	writeFile(dir, ".gitignore", fmt.Sprintf(`# Binaries
%s
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test
*.test
*.out

# Go workspace
go.work
go.work.sum
`, name))

	writeFile(dir, "README.md", fmt.Sprintf(`# %s

A command-line tool built with Go.

## Build

`+"```"+`bash
go build -o %s .
`+"```"+`

## Usage

`+"```"+`bash
./%s greet --name Gopher
./%s version
`+"```"+`
`, name, name, name, name))

	return nil
}
