// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package create

import "fmt"

func buildAPI(dir, modPath, name string) error {
	writeFile(dir, "main.go", fmt.Sprintf(`package main

import (
	"log"
	"net/http"
	"os"

	"%s/handler"
	"%s/middleware"
)

func main() {
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /api/v1/items", handler.ListItems)
	mux.HandleFunc("POST /api/v1/items", handler.CreateItem)
	mux.HandleFunc("GET /api/v1/items/{id}", handler.GetItem)

	// Wrap with middleware
	wrapped := middleware.Logger(middleware.CORS(mux))

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	log.Printf("Starting API server on %%s", addr)
	if err := http.ListenAndServe(addr, wrapped); err != nil {
		log.Fatal(err)
	}
}
`, modPath, modPath))

	writeFile(dir, "handler/health.go", `package handler

import (
	"encoding/json"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
`)

	writeFile(dir, "handler/items.go", `package handler

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Item struct {
	ID   string `+"`"+`json:"id"`+"`"+`
	Name string `+"`"+`json:"name"`+"`"+`
}

var (
	mu    sync.RWMutex
	items = []Item{
		{ID: "1", Name: "First item"},
		{ID: "2", Name: "Second item"},
	}
	nextID = 3
)

func ListItems(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	mu.RLock()
	defer mu.RUnlock()
	for _, item := range items {
		if item.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, `+"`"+`{"error":"not found"}`+"`"+`, http.StatusNotFound)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, `+"`"+`{"error":"invalid json"}`+"`"+`, http.StatusBadRequest)
		return
	}
	mu.Lock()
	item.ID = itoa(nextID)
	nextID++
	items = append(items, item)
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
`)

	writeFile(dir, "middleware/logger.go", `package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
`)

	writeFile(dir, "middleware/cors.go", `package middleware

import "net/http"

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
`)

	writeFile(dir, "README.md", fmt.Sprintf(`# %s

A REST API server built with Go's standard library.

## Run

`+"```"+`bash
go run .
`+"```"+`

## Endpoints

| Method | Path               | Description     |
|--------|--------------------|-----------------|
| GET    | /health            | Health check    |
| GET    | /api/v1/items      | List all items  |
| POST   | /api/v1/items      | Create an item  |
| GET    | /api/v1/items/{id} | Get item by ID  |

## Test

`+"```"+`bash
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/items
curl -X POST http://localhost:8080/api/v1/items -d '{"name":"new item"}'
`+"```"+`
`, name))

	return nil
}
