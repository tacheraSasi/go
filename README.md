# My Go — A Custom Go Toolchain

> **By [tacherasasi](https://github.com/tacherasasi)**

A custom build of the [Go programming language](https://go.dev) with additional developer experience features built on top of the official toolchain. Syncs with upstream [golang/go](https://github.com/golang/go) so you always get the latest Go releases plus extra tools.

![Gopher image](https://golang.org/doc/gopher/fiveyears.jpg)
_Gopher image by [Renee French][rf], licensed under [Creative Commons 4.0 Attribution license][cc4-by]._

---

## What's Different?

This fork adds new subcommands to the `go` CLI that don't exist in standard Go.

### `go create` — Project Scaffolding

Scaffold a new Go project in one command. No more copy-pasting boilerplate.

```bash
go create <template> <module-path>
```

**Example:**

```bash
go create api github.com/tacherasasi/myapi
# Creates ./myapi with a full REST API project ready to go
```

**Available templates:**

| Template       | What you get                                                             |
| -------------- | ------------------------------------------------------------------------ |
| `cli`          | CLI app with subcommands, flags, help text, `.gitignore`, README         |
| `api`          | REST API with routing, CRUD handlers, logging & CORS middleware          |
| `microservice` | HTTP service with graceful shutdown, health/readiness probes, Dockerfile |
| `web`          | Web server with `embed.FS`, HTML templates, static CSS                   |
| `lib`          | Library package with tests, examples, and doc comments                   |
| `grpc`         | gRPC scaffold with `.proto` file and placeholder server                  |

Run `go create list` to see all templates from the command line.

---

## Install From Source

You need Go 1.22.6+ installed as a bootstrap compiler.

```bash
git clone https://github.com/tacherasasi/my-go.git
cd my-go/src
./make.bash
```

Then add the binary to your PATH:

```bash
export PATH=$HOME/path/to/my-go/bin:$PATH
```

Verify:

```bash
go version
go create list
```

## Staying Up to Date with Upstream Go

This repo tracks upstream Go. To pull the latest changes:

```bash
git fetch upstream
git merge upstream/master
```

Then rebuild:

```bash
cd src && ./make.bash
```

---

## Upstream Go

This is built on top of the official Go toolchain. Everything from standard Go works as expected.

- Official Go: https://go.dev
- Source: https://github.com/golang/go
- Docs: https://go.dev/doc

### Contributing

Ideas for new templates or features? Open an issue or PR.

If you want to contribute to Go itself, see the [Go contribution guide](https://go.dev/doc/contribute).

[rf]: https://reneefrench.blogspot.com/
[cc4-by]: https://creativecommons.org/licenses/by/4.0/
