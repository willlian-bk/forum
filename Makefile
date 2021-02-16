.SILENT:

.PHONY: build
build:
	go build -o forum -v ./cmd/app

.PHONY: git
git:
	git add .
	git commit -m "$(comment)"
	git push

.PHONY: git-conf
git-conf:
	git config --global user.email "$(email)"
	git config --global user.name "$(name)"

.DEFAULT_GOAL := build