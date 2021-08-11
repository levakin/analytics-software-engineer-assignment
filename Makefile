ifneq (,$(wildcard .env))
	include .env
	export
endif

.PHONY: help
help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'


# ==============================================================================
# Modules support

.PHONY: deps-reset
deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

.PHONY: tidy
tidy:
	go mod tidy
	go mod vendor

.PHONY: deps-upgrade
deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

.PHONY: deps-cleancache
deps-cleancache:
	go clean -modcache

.PHONY: test
test: ## Run go tests
	go test ./... -count=1

.PHONY: fmt
fmt: ## Format code
	gofumpt -l -w $$(go list -f {{.Dir}} ./... | grep -v /vendor/)

.PHONY: lint
lint: ## Run lint go code
	golangci-lint run
