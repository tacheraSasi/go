// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package create

import "fmt"

func buildWeb(dir, modPath, name string) error {
	writeFile(dir, "main.go", fmt.Sprintf(`package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

var tmpl *template.Template

func main() {
	var err error
	tmpl, err = template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(staticFS))
	mux.HandleFunc("GET /", handleHome)
	mux.HandleFunc("GET /about", handleAbout)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	log.Printf("Starting %s on %%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

type pageData struct {
	Title string
	Name  string
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl.ExecuteTemplate(w, "home.html", pageData{Title: "Home", Name: "%s"})
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "about.html", pageData{Title: "About", Name: "%s"})
}
`, modPath, name, name, name))

	writeFile(dir, "templates/base.html", fmt.Sprintf(`{{define "base"}}<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - %s</title>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <nav>
        <a href="/">Home</a>
        <a href="/about">About</a>
    </nav>
    <main>
        {{template "content" .}}
    </main>
</body>
</html>{{end}}
`, name))

	writeFile(dir, "templates/home.html", `{{template "base" .}}
{{define "content"}}
<h1>Welcome to {{.Name}}</h1>
<p>Your Go web application is running.</p>
{{end}}
`)

	writeFile(dir, "templates/about.html", `{{template "base" .}}
{{define "content"}}
<h1>About {{.Name}}</h1>
<p>Built with Go and the standard library.</p>
{{end}}
`)

	writeFile(dir, "static/style.css", `* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: system-ui, -apple-system, sans-serif;
    line-height: 1.6;
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
    color: #333;
}

nav {
    margin-bottom: 2rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid #ddd;
}

nav a {
    margin-right: 1rem;
    text-decoration: none;
    color: #0066cc;
}

nav a:hover {
    text-decoration: underline;
}

h1 {
    margin-bottom: 1rem;
}
`)

	writeFile(dir, "README.md", fmt.Sprintf(`# %s

A web server with HTML templates and static files, built with Go's standard library.

## Run

`+"```"+`bash
go run .
`+"```"+`

Then open http://localhost:8080 in your browser.
`, name))

	return nil
}
