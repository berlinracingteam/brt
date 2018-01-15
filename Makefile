#!/usr/bin/env bash

GOCMD	:= go
GOBUILD	:= $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST	:= $(GOCMD) test
GORUN   := $(GOCMD) run

BINARY_NAME := server

.PHONY: all build clean run test todo

all: test build
build:	; @$(GOBUILD) -o $(BINARY_NAME) -v cmd/server/server.go
clean:	; @$(GOCLEAN) && rm -rf $(BINARY_NAME) $(HANDLER).so $(HANDLER).zip
run:	; @$(GORUN) -v cmd/server/server.go
test:	; @$(GOTEST) -cover ./...

latest:
	curl -o latest.dump `heroku pg:backups public-url --app brt-backend`
	pg_restore --verbose --clean --no-acl --no-owner -d brt latest.dump

todo:
	@grep \
		--exclude-dir=./vendor \
		--exclude=./Makefile \
		--text \
		--color \
		-nRo ' TODO:.*' .
