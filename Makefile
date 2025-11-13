.PHONY: build
build:
	go build -v ./cmd/app/main.go

.DEFAULT_GOAL := build