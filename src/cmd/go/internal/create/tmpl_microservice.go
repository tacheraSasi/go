// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package create

import "fmt"

func buildMicroservice(dir, modPath, name string) error {
	writeFile(dir, "main.go", fmt.Sprintf(`package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"%s/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /ready", handler.Ready)
	mux.HandleFunc("GET /api/v1/ping", handler.Ping)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting %s on %%s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-done
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %%v", err)
	}
	log.Println("Server stopped")
}
`, modPath, name))

	writeFile(dir, "handler/handler.go", `package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

var startTime = time.Now()

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ready",
		"uptime": time.Since(startTime).String(),
	})
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}
`)

	writeFile(dir, "Dockerfile", fmt.Sprintf(`FROM golang:1.27 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /%s .

FROM gcr.io/distroless/static
COPY --from=builder /%s /
EXPOSE 8080
ENTRYPOINT ["/%s"]
`, name, name, name))

	writeFile(dir, "README.md", fmt.Sprintf(`# %s

A lightweight microservice with health checks and graceful shutdown.

## Run

`+"```"+`bash
go run .
`+"```"+`

## Docker

`+"```"+`bash
docker build -t %s .
docker run -p 8080:8080 %s
`+"```"+`

## Endpoints

| Method | Path          | Description         |
|--------|---------------|---------------------|
| GET    | /health       | Liveness probe      |
| GET    | /ready        | Readiness probe     |
| GET    | /api/v1/ping  | Ping/pong           |
`, name, name, name))

	return nil
}
