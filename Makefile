include Makefile.help.mk

##@ Format
.PHONY: format
format: ## Format golang and proto files.
	tool/script/format.sh

.PHONY: tidy
tidy: ## Format golang and proto files.
	tool/script/tidy.sh

.PHONY: lint.cleancache
lint.cleancache: ## Clean golangci-lint cache.
	golangci-lint cache clean

.PHONY: lint
lint: ## Lint proto files using buf and golang files using golangci-lint.
lint: lint.cleancache
	tool/script/lint.sh ${svc}

.PHONY: pretty
pretty: ## Prettify golang and proto files. Basically, it runs tidy, format, and lint command.
pretty: tidy gen.mock format lint

.PHONY: check.import
check.import: ## Check if import blocks are separated accordingly.
	tool/script/check-import.sh

##@ Generator
.PHONY: gen.proto
gen.proto: ## Generate golang files from proto.
	tool/script/generate-proto.sh

.PHONY: gen.mock
gen.mock: ## Generate mock from all golang interfaces.
	tool/script/generate-mock.sh

.PHONY: gen.req
gen.req: ## Generate requirement document.
	tool/script/requirement.sh $(name)

##@ Test
.PHONY: test.unit
test.unit: ## Run unit test.
	tool/script/test.sh unit

.PHONY: test.cover
test.cover: ## Run unit test.
	tool/script/test.sh cover

##@ Migration
.PHONY: migration
migration: ## Create database migration.
	migrate create -ext sql -dir service/$(svc)/db/migrations $(name)

.PHONY: migrate
migrate: ## Run database migrations.
	migrate -path service/$(svc)/db/migrations -database "$(url)" -verbose up
