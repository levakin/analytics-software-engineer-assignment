ifneq (,$(wildcard .env))
	include .env
	export
endif

help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'


# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

test: ## Run go tests
	go test ./... -count=1

fmt: ## Format code
	gofumpt -l -w $$(go list -f {{.Dir}} ./... | grep -v /vendor/)

lint: ## Run lint go code
	golangci-lint run
