// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package create

import "fmt"

func buildLib(dir, modPath, name string) error {
	writeFile(dir, fmt.Sprintf("%s.go", name), fmt.Sprintf(`// Package %s provides ...
package %s

// Hello returns a greeting for the given name.
func Hello(name string) string {
	if name == "" {
		name = "World"
	}
	return "Hello, " + name + "!"
}

// Add returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}
`, name, name))

	writeFile(dir, fmt.Sprintf("%s_test.go", name), fmt.Sprintf(`package %s

import "testing"

func TestHello(t *testing.T) {
	tests := []struct {
		name string
		input string
		want string
	}{
		{"with name", "Gopher", "Hello, Gopher!"},
		{"empty", "", "Hello, World!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Hello(tt.input)
			if got != tt.want {
				t.Errorf("Hello(%%q) = %%q, want %%q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	if got := Add(2, 3); got != 5 {
		t.Errorf("Add(2, 3) = %%d, want 5", got)
	}
}
`, name))

	writeFile(dir, fmt.Sprintf("example_test.go"), fmt.Sprintf(`package %s_test

import (
	"fmt"

	"%s"
)

func ExampleHello() {
	fmt.Println(%s.Hello("Gopher"))
	// Output: Hello, Gopher!
}

func ExampleAdd() {
	fmt.Println(%s.Add(2, 3))
	// Output: 5
}
`, name, modPath, name, name))

	writeFile(dir, "README.md", fmt.Sprintf(`# %s

A Go library.

## Install

`+"```"+`bash
go get %s
`+"```"+`

## Usage

`+"```"+`go
import "%s"

msg := %s.Hello("Gopher")
`+"```"+`

## Test

`+"```"+`bash
go test ./...
`+"```"+`
`, name, modPath, modPath, name))

	return nil
}
