.SILENT:

.PHONY: build
build:
	go build -o forum -v ./cmd/app

.PHONY: git
git:
	git add .
	git commit -m "$(comment)"
	git push

.DEFAULT_GOAL := build