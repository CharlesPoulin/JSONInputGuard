SHELL := /bin/bash

.PHONY: build run test bench bench-guard benchmark-all lambda

# Build the application
build:
	go build ./...

# Run the application
run:
	go run ./cmd/server

# Run all tests
test:
	go test ./...

# Run all benchmarks
bench:
	go test -bench=. -benchmem ./bench -run=^$

# Run only the GuardOnly64KB benchmark
bench-guard:
	go test -bench=GuardOnly64KB -benchmem ./bench -run=^$

# Run all benchmarks and tests
benchmark-all:
	go test -bench=. -benchmem ./... -run=.

# Build the lambda function
lambda:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-s -w' -o bootstrap ./cmd/lambda && zip -qj lambda.zip bootstrap && rm -f bootstrap
