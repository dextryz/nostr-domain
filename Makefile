SHELL := /usr/bin/env bash

fmt:
	go mod tidy -compat=1.17
	gofmt -l -s -w .

test:
	go test ./...
