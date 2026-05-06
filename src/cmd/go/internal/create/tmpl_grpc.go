// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package create

import "fmt"

func buildGRPC(dir, modPath, name string) error {
	writeFile(dir, "main.go", fmt.Sprintf(`package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"%s/server"
)

func main() {
	addr := ":50051"
	if port := os.Getenv("GRPC_PORT"); port != "" {
		addr = ":" + port
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %%v", err)
	}

	srv := server.New()
	log.Printf("Starting gRPC server on %%s", addr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %%v", err)
		}
	}()

	<-done
	log.Println("Shutting down...")
	srv.GracefulStop()
}
`, modPath))

	writeFile(dir, "server/server.go", `package server

// TODO: This is a placeholder. To use gRPC:
//
// 1. Install protoc and the Go gRPC plugins:
//    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
//    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
//
// 2. Define your service in proto/service.proto
//
// 3. Generate Go code:
//    protoc --go_out=. --go-grpc_out=. proto/service.proto
//
// 4. Implement the generated interface in this package.

import "net"

// Server wraps the gRPC server (placeholder).
type Server struct{}

// New creates a new Server.
func New() *Server {
	return &Server{}
}

// Serve starts serving on the listener.
func (s *Server) Serve(lis net.Listener) error {
	// Replace with: grpcServer.Serve(lis)
	// For now, block until the listener is closed.
	for {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}
		conn.Close()
	}
}

// GracefulStop stops the server gracefully.
func (s *Server) GracefulStop() {}
`)

	writeFile(dir, "proto/service.proto", fmt.Sprintf(`syntax = "proto3";

package %s;

option go_package = "%s/proto";

service %sService {
  rpc Ping (PingRequest) returns (PingResponse);
}

message PingRequest {
  string message = 1;
}

message PingResponse {
  string message = 1;
}
`, name, modPath, capitalize(name)))

	writeFile(dir, "README.md", fmt.Sprintf(`# %s

A gRPC service scaffold.

## Setup

Install the protoc compiler and Go plugins:

`+"```"+`bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
`+"```"+`

## Generate

`+"```"+`bash
protoc --go_out=. --go-grpc_out=. proto/service.proto
`+"```"+`

## Run

`+"```"+`bash
go run .
`+"```"+`
`, name))

	return nil
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	b := []byte(s)
	if b[0] >= 'a' && b[0] <= 'z' {
		b[0] -= 'a' - 'A'
	}
	return string(b)
}
