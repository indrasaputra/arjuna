OUTPUT_DIR = deploy/output

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

##@ Build
.PHONY: compile
compile: ## Compile golang code to binary.
	mkdir -p $(OUTPUT_DIR)
	(cd service/$(svc) && \
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o $(OUTPUT_DIR)/$(svc) cmd/server/main.go)

.PHONY: build.user
build.user: ## Build docker user service.
	docker build --no-cache -t indrasaputra/arjuna/user:latest -f service/user/dockerfile/user.dockerfile .

.PHONY: build.elements
build.elements: ## Build docker elements.
	docker build --no-cache -t indrasaputra/arjuna-elements:latest -f dockerfile/elements.dockerfile .

##@ Test
.PHONY: test.unit
test.unit: ## Run unit test.
	tool/script/test.sh unit

.PHONY: test.cover
test.cover: ## Run unit test.
	tool/script/test.sh cover

.PHONY: test.e2e
test.e2e: ## Run e2e test using Godog.
	tool/script/godog.sh

##@ Migration
.PHONY: migration
migration: ## Create database migration.
	migrate create -ext sql -dir service/$(svc)/db/migrations $(name)

.PHONY: migrate
migrate: ## Run database migrations.
	migrate -path service/$(svc)/db/migrations -database "$(url)" -verbose up
