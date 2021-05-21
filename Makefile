.DEFAULT_GOAL := help

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z-]+:.*?## .*$$/ {printf "\033[32m%-15s\033[0m %s\n", $$1, $$2}' Makefile | sort

.PHONY: build
build: ## Build app
	CGO_ENABLED=0 go build -o bluetooth-auto-connect cmd/main.go