SHELL := /bin/bash

.PHONY: build run bench bench-guard lambda

build:
	go build ./...

run:
	go run ./cmd/server

bench:
	go test -bench=. -benchmem ./bench -run=^$

bench-guard:
	go test -bench=GuardOnly64KB -benchmem ./bench -run=^$

lambda:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-s -w' -o bootstrap ./cmd/lambda && zip -qj lambda.zip bootstrap && rm -f bootstrap
