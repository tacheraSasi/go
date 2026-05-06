// Copyright 2026 The Tachera Sasi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go create — project scaffolding

package create

import (
	"cmd/go/internal/base"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var CmdCreate = &base.Command{
	UsageLine: "go create <template> <module-path>",
	Short:     "scaffold a new Go project from a template",
	Long: `
Create scaffolds a new Go project using a built-in template.

The first argument is the template name. The second argument is the module path
(e.g., github.com/user/myapp). A directory named after the last element of the
module path is created in the current directory and populated with the project
files.

Available templates:

  cli          Command-line application with flags and subcommand structure
  api          HTTP REST API server with router and middleware
  microservice Lightweight microservice with health check and graceful shutdown
  web          Web server with static files and template rendering
  lib          Reusable library with exported API and examples
  grpc         gRPC service with protobuf placeholder

Run 'go create list' to see all available templates with descriptions.

See https://go.dev/ref/mod#go-mod-init for more about module paths.
`,
	Run: runCreate,
}

func init() {
	base.AddChdirFlag(&CmdCreate.Flag)
}

var templates = map[string]struct {
	desc  string
	build func(dir, modPath, name string) error
}{
	"cli":          {"Command-line application with flags and subcommands", buildCLI},
	"api":          {"HTTP REST API server with routing and middleware", buildAPI},
	"microservice": {"Lightweight microservice with health check and graceful shutdown", buildMicroservice},
	"web":          {"Web server with templates and static files", buildWeb},
	"lib":          {"Reusable library package with examples", buildLib},
	"grpc":         {"gRPC service scaffold with protobuf placeholder", buildGRPC},
}

func runCreate(ctx context.Context, cmd *base.Command, args []string) {
	if len(args) == 0 {
		base.Fatalf("go create: missing template name\nRun 'go help create' for usage.")
	}

	// Handle "go create list"
	if args[0] == "list" {
		fmt.Println("Available templates:")
		order := []string{"cli", "api", "microservice", "web", "lib", "grpc"}
		for _, name := range order {
			t := templates[name]
			fmt.Printf("  %-14s %s\n", name, t.desc)
		}
		fmt.Println("\nUsage: go create <template> <module-path>")
		return
	}

	tmplName := args[0]
	tmpl, ok := templates[tmplName]
	if !ok {
		base.Fatalf("go create: unknown template %q\nRun 'go create list' to see available templates.", tmplName)
	}

	if len(args) < 2 {
		base.Fatalf("go create: missing module path\nUsage: go create %s <module-path>", tmplName)
	}
	if len(args) > 2 {
		base.Fatalf("go create: too many arguments\nUsage: go create %s <module-path>", tmplName)
	}

	modPath := args[1]
	name := filepath.Base(modPath)
	dir := filepath.Join(".", name)

	if _, err := os.Stat(dir); err == nil {
		base.Fatalf("go create: directory %q already exists", dir)
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		base.Fatalf("go create: %v", err)
	}

	// Write go.mod
	goMod := fmt.Sprintf("module %s\n\ngo 1.27\n", modPath)
	writeFile(dir, "go.mod", goMod)

	// Build template files
	if err := tmpl.build(dir, modPath, name); err != nil {
		base.Fatalf("go create: %v", err)
	}

	fmt.Printf("Created %s project in ./%s\n", tmplName, name)
	fmt.Printf("  cd %s\n  go run .\n", name)
}

func writeFile(dir, relPath, content string) {
	full := filepath.Join(dir, relPath)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		base.Fatalf("go create: %v", err)
	}
	if err := os.WriteFile(full, []byte(strings.TrimLeft(content, "\n")), 0o644); err != nil {
		base.Fatalf("go create: %v", err)
	}
}
