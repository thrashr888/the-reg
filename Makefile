NAME ?= $(shell basename "$(CURDIR)")

default: help

.PHONY: bootstrap
bootstrap:
	psql -f STRUCTURE.sql thereg
	go get github.com/cespare/reflex
	go get -u github.com/ddollar/forego

.PHONY: dev
dev:
	reflex -r '\.(go|html)' -s -- sh -c 'go build -o reg && forego start'

.PHONY: build
build:
	go build -o reg
	go build -o build/current/reg
	env GOOS=darwin GOARCH=amd64 go build -o build/darwin-amd64/reg
	env GOOS=linux GOARCH=amd64 go build -o build/linux-amd64/reg
	env GOOS=linux GOARCH=arm go build -o build/linux-arm/reg
	env GOOS=linux GOARCH=arm64 go build -o build/linux-arm64/reg
	env GOOS=windows GOARCH=amd64 go build -o build/windows-amd64/reg

.PHONY: release
release:
	# TODO

.PHONY: deploy
deploy:
	# TODO

.PHONY: clean
clean: ## Clean build artifacts
	rm -r $(CURDIR)/reg

.PHONY: help
help:
	echo $(NAME)
	echo "make bootstrap | dev | build | clean"
