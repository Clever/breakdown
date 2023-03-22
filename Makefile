include golang.mk
include wag.mk

SHELL := /bin/bash
export PATH := $(PWD)/bin:$(PATH)
APP_NAME ?= breakdown
EXECUTABLE = $(APP_NAME)
PKG = github.com/Clever/$(APP_NAME)
PKGS := $(shell go list ./... | grep -v /gen-go | grep -v /tools)

WAG_VERSION := latest

$(eval $(call golang-version-check,1.18))

.PHONY: all test build run $(PKGS) generate go-generate install_deps

all: test build

test: $(call golang-setup-coverage) $(PKGS) lint
$(PKGS): golang-test-all-strict-cover-deps gen-go go-generate
	$(call golang-test-all-strict-cover,$@)

build: gen-go vendor launch.go bin/kvconfig.yml
	$(call golang-build,$(PKG),$(EXECUTABLE))

bin/kvconfig.yml: kvconfig.yml
	cp kvconfig.yml bin/kvconfig.yml

# Code generation: wag, launch-gen, and mocks.
# Use `make generate` to do all code generation.
generate: gen-go vendor bin/mockgen launch.go go-generate
	go mod vendor


go-generate:
	go generate ./...

launch.go: bin/launch-gen launch/$(APP_NAME).yml
	bin/launch-gen -o ./launch.go -p main launch/$(APP_NAME).yml

gen-go: bin/wag jsdoc2md swagger.yml
	$(call wag-generate-mod,./swagger.yml)

run: bin/reflex gen-go vendor launch.go bin/kvconfig.yml
	bin/reflex -r '\.go$$' -R '^gen-go/' -R '^vendor/' -R '^gen-js/' -R '^tools/' -R '^node_modules/' -s -- bash -c "go run ."

# Dependencies setup.
# install_deps is an alias for vendor.
# vendor needs to be redone whenever go.mod or go.sum change.
# vendor must run before building bin apps
install_deps: vendor bin/reflex bin/mockgen bin/launch-gen bin/woke bin/sqlc
vendor: go.mod go.sum
	go mod vendor
bin/mockgen:
	go build -o bin/mockgen github.com/golang/mock/mockgen
bin/launch-gen:
	go build -o bin/launch-gen github.com/Clever/launch-gen
bin/reflex:
	go build -o bin/reflex github.com/cespare/reflex
bin/woke:
	go build -o bin/woke github.com/get-woke/woke
bin/sqlc:
	go build -mod=readonly -o bin/sqlc github.com/kyleconroy/sqlc/cmd/sqlc

# go.sum doesn't always exist (i.e. when repo is first created),
go.sum:
	go mod tidy

# Woke can be customized to apply to only certain globs
# See https://docs.getwoke.tech/usage/ for details
# You can also disable failure by removing --exit-1-on-failure below
lint: incl-lang-check
incl-lang-check: bin/woke
	bin/woke --exit-1-on-failure
