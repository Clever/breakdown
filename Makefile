include golang.mk
include wag.mk

SHELL := /bin/bash
export PATH := $(PWD)/bin:$(PATH)
APP_NAME ?= breakdown
EXECUTABLE = $(APP_NAME)
CLI_EXECUTABLE = breakdowncli
PKG = github.com/Clever/$(APP_NAME)
PKGS := $(shell go list ./... | grep -v /gen-go | grep -v /tools)

CLI_RAWVERSION :=$(shell head -n 1 cmd/cli/VERSION)
CLI_VERSION := $(CLI_RAWVERSION)$(shell if [[ -z "$(CI)" ]]; then echo "-dev"; fi)

POSTGRES_USER=postgres
POSTGRES_PASSWORD=supersecret
POSTGRES_DB=breakdown

POSTGRES_TEST_PORT=5432
ifneq ($(CI), true)
	POSTGRES_TEST_PORT=5433
endif

WAG_VERSION := latest

$(eval $(call golang-version-check,1.19))

.PHONY: all test build run $(PKGS) generate go-generate install_deps start start-postgres \
	stop-postgres migrate-local gen-sql test-postgres release

all: test build

test: test-postgres $(call golang-setup-coverage) $(PKGS) lint
$(PKGS): golang-test-all-strict-cover-deps gen-go go-generate
	$(call golang-test-all-strict-cover,$@)

build: gen-go vendor launch.go bin/kvconfig.yml
	$(call golang-build,$(PKG),$(EXECUTABLE))

build-cli: gen-go vendor
	go build -o=bin/breakdown-cli ./cmd/cli

release:
	@GOOS=linux GOARCH=amd64 go build -tags netcgo -ldflags="-s -w -X main.version=$(CLI_VERSION)" \
		-o="$@/$(CLI_EXECUTABLE)-$(CLI_VERSION)-linux-amd64" ./cmd/cli
	@GOOS=darwin GOARCH=amd64 go build -tags netcgo -ldflags="-s -w -X main.version=$(CLI_VERSION)" \
		-o="$@/$(CLI_EXECUTABLE)-$(CLI_VERSION)-darwin-amd64" ./cmd/cli
	@GOOS=darwin GOARCH=arm64 go build -tags netcgo -ldflags="-s -w -X main.version=$(CLI_VERSION)" \
		-o="$@/$(CLI_EXECUTABLE)-$(CLI_VERSION)-darwin-arm64" ./cmd/cli

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

run: bin/reflex gen-go vendor launch.go bin/kvconfig.yml start-postgres
	POSTGRES_USERNAME=$(POSTGRES_USER) \
	POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	POSTGRES_HOST=localhost \
	POSTGRES_DB=$(POSTGRES_DB) \
	bin/reflex -r '\.go$$' -R '^gen-go/' -R '^cmd/cli/' -R '^vendor/' -R '^gen-js/' -R '^tools/' -R '^node_modules/' -s -- bash -c "go run ."

# Dependencies setup.
# install_deps is an alias for vendor.
# vendor needs to be redone whenever go.mod or go.sum change.
# vendor must run before building bin apps
install_deps: vendor bin/reflex bin/mockgen bin/launch-gen bin/woke
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

test-postgres:
ifneq ($(CI), true)
	@if docker ps --format '{{.Names}}' | grep -w breakdown-test-postgres &> /dev/null; then\
		docker stop breakdown-test-postgres; \
	fi;
	docker run \
		--name breakdown-test-postgres \
		--rm \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB)_test \
		-p 5433:5432 \
		-d \
		postgres:14.6;
	@while ! PGPASSWORD=$(POSTGRES_PASSWORD) docker exec breakdown-test-postgres psql -h 127.0.0.1 -p 5432 -U $(POSTGRES_USER) -d $(POSTGRES_DB)_test -c "SELECT 1" >/dev/null; do\
		sleep 1;\
	done
endif
	goose -dir ./db/migrations/ \
	postgres "host=127.0.0.1 port=$(POSTGRES_TEST_PORT) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB)_test" \
	up

start-postgres:
	@if ! docker ps --format '{{.Names}}' | grep -w breakdown-postgres &> /dev/null; then\
		docker run \
		--name breakdown-postgres \
		--rm \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-v breakdown-postgres-vol:/var/lib/postgresql/data \
		-p 5432:5432 \
		-d \
		postgres:14.6; \
	fi;
	@while ! PGPASSWORD=$(POSTGRES_PASSWORD) docker exec breakdown-postgres psql -h 127.0.0.1 -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "SELECT 1" >/dev/null; do\
		sleep 1;\
	done

migrate-local: start-postgres
	goose -dir ./db/migrations/ postgres "host=127.0.0.1 user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB)" up

gen-sql: bin/sqlc sqlc.yaml
	bin/sqlc generate

