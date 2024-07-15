OUTPUT_DIR			= deploy/output
PROTOGEN_IMAGE		= indrasaputra/protogen:2024-07-13
SERVICES			= blueprint gateway auth user

include Makefile.help.mk

##@ Format
.PHONY: format
format: ## Format golang and proto files.
	tool/script/format.sh

.PHONY: tidy
tidy: ## Format golang and proto files.
	tool/script/tidy.sh

.PHONY: lint.struct
lint.struct: ## Run fieldalignment.
	tool/script/lint-struct.sh

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

.PHONY: gen.proto.docker
gen.proto.docker: ## Generate proto and prettify files using docker.
	docker run -it --rm \
    --mount "type=bind,source=$(PWD),destination=/work" \
    --mount "type=volume,source=arjuna-go-mod-cache,destination=/go,consistency=cached" \
    --mount "type=volume,source=arjuna-buf-cache,destination=/home/.cache,consistency=cached" \
    -w /work $(PROTOGEN_IMAGE) make -e -f Makefile gen.proto pretty

.PHONY: gen.mock
gen.mock: ## Generate mock from all golang interfaces.
	tool/script/generate-mock.sh

.PHONY: gen.req
gen.req: ## Generate requirement document.
	tool/script/requirement.sh $(name)

##@ Build
.PHONY: compile
compile: ## Compile service.
	tool/script/compile.sh $(svc)

.PHONY: compile.all
compile.all: ## Compile all services.
	for svc in $(SERVICES); do \
		tool/script/compile.sh $$svc; \
	done

.PHONY: build
build: ## Build docker for service.
	tool/script/docker-build.sh $(svc)

.PHONY: build.all
build.all: ## Build docker for all services.
	for svc in $(SERVICES); do \
		tool/script/docker-build.sh $$svc; \
	done

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
