.PHONY: build
build:
	go build -o forum -v ./cmd/app

.DEFAULT_GOAL := build