include Makefile.help.mk

##@ Format
.PHONY: format
format: ## Format golang and proto files.
	tool/format.sh

.PHONY: tidy
tidy: ## Format golang and proto files.
	tool/tidy.sh

.PHONY: lint.cleancache
lint.cleancache: ## Clean golangci-lint cache.
	golangci-lint cache clean

.PHONY: lint
lint: ## Lint proto files using buf and golang files using golangci-lint.
lint: lint.cleancache
	tool/lint.sh ${svc}

.PHONY: check.import
check.import: ## Check if import blocks are separated accordingly.
	tool/check-import.sh

##@ Generator
.PHONY: gen.proto
gen.proto: ## Generate golang files from proto.
	tool/generate-proto.sh

##@ Test
.PHONY: test.unit
test.unit: ## Run unit test.
	tool/test.sh unit

.PHONY: test.cover
test.cover: ## Run unit test.
	tool/test.sh cover
