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

.PHONY: clean
clean: ## Clean build artifacts
	rm -r $(CURDIR)/reg

.PHONY: help
help:
	echo $(NAME)
	echo "make bootstrap | dev | build | clean"
