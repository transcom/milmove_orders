DB_NAME_DEV = dev_db
DB_NAME_DEPLOYED_MIGRATIONS = deployed_migrations
DB_NAME_TEST = test_db
DB_DOCKER_CONTAINER_DEV = orders-db-dev
DB_DOCKER_CONTAINER_DEPLOYED_MIGRATIONS = orders-db-deployed-migrations
DB_DOCKER_CONTAINER_TEST = orders-db-test
# The version of the postgres container should match production as closely
# as possible.
# https://github.com/transcom/ppp-infra/blob/7ba2e1086ab1b2a0d4f917b407890817327ffb3d/modules/aws-app-environment/database/variables.tf#L48
DB_DOCKER_CONTAINER_IMAGE = postgres:10.10
export PGPASSWORD=mysecretpassword


# Convenience for LDFLAGS
WEBSERVER_LDFLAGS=-X main.gitBranch=$(shell git branch | grep \* | cut -d ' ' -f2) -X main.gitCommit=$(shell git rev-list -1 HEAD)
GC_FLAGS=-trimpath=$(GOPATH)
DB_PORT_DEV=7432
DB_PORT_TEST=7433
DB_PORT_DEPLOYED_MIGRATIONS=7434
DB_PORT_DOCKER=5432
ifdef CIRCLECI
	DB_PORT_DEV=5432
	DB_PORT_TEST=5432
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		LDFLAGS=-linkmode external -extldflags -static
	endif
endif

ifdef GOLAND
	GOLAND_GC_FLAGS=all=-N -l
endif


.PHONY: help
help:  ## Print the help documentation
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


#
# ----- END PREAMBLE -----
#

#
# ----- START CHECK TARGETS -----
#

# This target ensures that the pre-commit hook is installed and kept up to date
# if pre-commit updates.
.PHONY: ensure_pre_commit
ensure_pre_commit: .git/hooks/pre-commit ## Ensure pre-commit is installed
.git/hooks/pre-commit: /usr/local/bin/pre-commit
	pre-commit install
	pre-commit install-hooks


.PHONY: check_hosts
check_hosts: .check_hosts.stamp ## Check that hosts are in the /etc/hosts file
.check_hosts.stamp: scripts/check-hosts-file
ifndef CIRCLECI
	scripts/check-hosts-file
else
	@echo "Not checking hosts on CircleCI."
endif
	touch .check_hosts.stamp

.PHONY: check_go_version
check_go_version: .check_go_version.stamp ## Check that the correct Golang version is installed
.check_go_version.stamp: scripts/check-go-version
	scripts/check-go-version
	touch .check_go_version.stamp

.PHONY: check_gopath
check_gopath: .check_gopath.stamp ## Check that $GOPATH exists in $PATH
.check_gopath.stamp:
	scripts/check-gopath
	touch .check_gopath.stamp


.PHONY: test
test: server_test ## Run all tests

.PHONY: check_log_dir
check_log_dir: ## Make sure we have a log directory
	mkdir -p log

#
# ----- END CHECK TARGETS -----
#

#
# ----- START BIN TARGETS -----
#

### Go Tool Targets

bin/gin: .check_go_version.stamp .check_gopath.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/gin github.com/codegangsta/gin
bin/swagger: .check_go_version.stamp .check_gopath.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/swagger github.com/go-swagger/go-swagger/cmd/swagger

# No static linking / $(LDFLAGS) because go-junit-report is only used for building the CirlceCi test report
bin/go-junit-report: .check_go_version.stamp .check_gopath.stamp
	go build -o bin/go-junit-report github.com/jstemmer/go-junit-report

### Cert Targets

bin/rds-ca-2019-root.pem:
	mkdir -p bin/
	curl -sSo bin/rds-ca-2019-root.pem https://s3.amazonaws.com/rds-downloads/rds-ca-2019-root.pem

### Orders Targets

bin/orders:
	go build -gcflags="$(GOLAND_GC_FLAGS) $(GC_FLAGS)" -asmflags=-trimpath=$(GOPATH) -ldflags "$(LDFLAGS) $(WEBSERVER_LDFLAGS)" -o bin/orders ./cmd/orders

#
# ----- END BIN TARGETS -----
#

#
# ----- START SERVER TARGETS -----
#

.PHONY: server_generate
server_generate: .check_go_version.stamp .check_gopath.stamp pkg/gen/ ## Generate golang server code from Swagger files
pkg/gen/: bin/swagger $(shell find swagger -type f -name *.yaml)
	scripts/gen-server

.PHONY: server_build
server_build: bin/orders ## Build the server

# This command is for running the server by itself, it will serve the compiled frontend on its own
# Note: Don't double wrap with aws-vault because the pkg/cli/vault.go will handle it
server_run_standalone: check_log_dir server_build db_dev_run
	DEBUG_LOGGING=true ./bin/orders serve 2>&1 | tee -a log/dev.log

# This command will rebuild the swagger go code and rerun server on any changes
server_run:
	find ./swagger -type f -name "*.yaml" | entr -c -r make server_run_default
# This command runs the server behind gin, a hot-reload server
# Note: Gin is not being used as a proxy so assigning odd port and laddr to keep in IPv4 space.
# Note: The INTERFACE envar is set to configure the gin build, orders_gin, local IP4 space with default port 8080.
server_run_default: .check_hosts.stamp .check_go_version.stamp .check_gopath.stamp check_log_dir bin/gin server_generate db_dev_run
	INTERFACE=localhost DEBUG_LOGGING=true \
	$(AWS_VAULT) ./bin/gin \
		--build ./cmd/orders \
		--bin /bin/orders_gin \
		--laddr 127.0.0.1 --port 9001 \
		--excludeDir node_modules \
		--immediate \
		--buildArgs "-i -ldflags=\"$(WEBSERVER_LDFLAGS)\"" \
		serve \
		2>&1 | tee -a log/dev.log

.PHONY: build_tools
build_tools: bin/gin \
	bin/swagger \
	bin/rds-ca-2019-root.pem  ## Build all tools

.PHONY: server_test
server_test: db_test_reset db_test_migrate server_test_standalone ## Run server unit tests

.PHONY: server_test_standalone
server_test_standalone: ## Run server unit tests with no deps
	NO_DB=1 scripts/run-server-test

.PHONY: server_test_build
server_test_build:
	NO_DB=1 DRY_RUN=1 scripts/run-server-test

.PHONY: server_test_all
server_test_all: db_dev_reset db_dev_migrate ## Run all server unit tests
	# Like server_test but runs extended tests that may hit external services.
	LONG_TEST=1 scripts/run-server-test

.PHONY: server_test_coverage_generate
server_test_coverage_generate: db_test_reset db_test_migrate server_test_coverage_generate_standalone ## Run server unit test coverage

.PHONY: server_test_coverage_generate_standalone
server_test_coverage_generate_standalone: ## Run server unit tests with coverage and no deps
	# Add coverage tracker via go cover
	NO_DB=1 COVERAGE=1 scripts/run-server-test

.PHONY: server_test_coverage
server_test_coverage: db_test_reset db_test_migrate server_test_coverage_generate ## Run server unit test coverage with html output
	DB_PORT=$(DB_PORT_TEST) go tool cover -html=coverage.out

.PHONY: server_test_docker
server_test_docker:
	docker-compose -f docker-compose.circle.yml --compatibility up --remove-orphans --abort-on-container-exit

.PHONY: server_test_docker_down
server_test_docker_down:
	docker-compose -f docker-compose.circle.yml --compatibility down

#
# ----- END SERVER TARGETS -----
#

#
# ----- START DB_DEV TARGETS -----
#

.PHONY: db_dev_destroy
db_dev_destroy: ## Destroy Dev DB
ifndef CIRCLECI
	@echo "Destroying the ${DB_DOCKER_CONTAINER_DEV} docker database container..."
	docker rm -f $(DB_DOCKER_CONTAINER_DEV) || echo "No database container"
	rm -fr mnt/db_dev # delete mount directory if exists
else
	@echo "Relying on CircleCI's database setup to destroy the DB."
endif

.PHONY: db_dev_start
db_dev_start: ## Start Dev DB
ifndef CIRCLECI
	brew services stop postgresql 2> /dev/null || true
endif
	@echo "Starting the ${DB_DOCKER_CONTAINER_DEV} docker database container..."
	# If running do nothing, if not running try to start, if can't start then run
	docker start $(DB_DOCKER_CONTAINER_DEV) || \
		docker run -d --name $(DB_DOCKER_CONTAINER_DEV) \
			-e POSTGRES_PASSWORD=$(PGPASSWORD) \
			-p $(DB_PORT_DEV):$(DB_PORT_DOCKER)\
			$(DB_DOCKER_CONTAINER_IMAGE)

.PHONY: db_dev_create
db_dev_create: ## Create Dev DB
	@echo "Create the ${DB_NAME_DEV} database..."
	DB_NAME=postgres scripts/wait-for-db && DB_NAME=postgres psql-wrapper "CREATE DATABASE $(DB_NAME_DEV);" || true

.PHONY: db_dev_run
db_dev_run: db_dev_start db_dev_create ## Run Dev DB (start and create)

.PHONY: db_dev_reset
db_dev_reset: db_dev_destroy db_dev_run ## Reset Dev DB (destroy and run)

.PHONY: db_dev_migrate_standalone ## Migrate Dev DB directly
db_dev_migrate_standalone: bin/orders
	@echo "Migrating the ${DB_NAME_DEV} database..."
	DB_DEBUG=0 bin/orders migrate -p "file://migrations/${APPLICATION}/secure;file://migrations/${APPLICATION}/schema" -m "migrations/${APPLICATION}/migrations_manifest.txt"

.PHONY: db_dev_migrate
db_dev_migrate: db_dev_migrate_standalone ## Migrate Dev DB

.PHONY: db_dev_psql
db_dev_psql: ## Open PostgreSQL shell for Dev DB
	scripts/psql-dev

#
# ----- END DB_DEV TARGETS -----
#

#

#
# ----- START DB_TEST TARGETS -----
#

.PHONY: db_test_destroy
db_test_destroy: ## Destroy Test DB
ifndef CIRCLECI
	@echo "Destroying the ${DB_DOCKER_CONTAINER_TEST} docker database container..."
	docker rm -f $(DB_DOCKER_CONTAINER_TEST) || \
		echo "No database container"
else
	@echo "Relying on CircleCI's database setup to destroy the DB."
	psql postgres://postgres:$(PGPASSWORD)@localhost:$(DB_PORT_TEST)?sslmode=disable -c 'DROP DATABASE IF EXISTS $(DB_NAME_TEST);'
endif

.PHONY: db_test_start
db_test_start: ## Start Test DB
ifndef CIRCLECI
	brew services stop postgresql 2> /dev/null || true
	@echo "Starting the ${DB_DOCKER_CONTAINER_TEST} docker database container..."
	docker start $(DB_DOCKER_CONTAINER_TEST) || \
		docker run --name $(DB_DOCKER_CONTAINER_TEST) \
			-e \
			POSTGRES_PASSWORD=$(PGPASSWORD) \
			-d \
			-p $(DB_PORT_TEST):$(DB_PORT_DOCKER)\
			$(DB_DOCKER_CONTAINER_IMAGE)\
			-c fsync=off\
			-c full_page_writes=off
else
	@echo "Relying on CircleCI's database setup to start the DB."
endif

.PHONY: db_test_create
db_test_create: ## Create Test DB
ifndef CIRCLECI
	@echo "Create the ${DB_NAME_TEST} database..."
	DB_NAME=postgres DB_PORT=$(DB_PORT_TEST) scripts/wait-for-db && \
		createdb -p $(DB_PORT_TEST) -h localhost -U postgres $(DB_NAME_TEST) || true
else
	@echo "Relying on CircleCI's database setup to create the DB."
	psql postgres://postgres:$(PGPASSWORD)@localhost:$(DB_PORT_TEST)?sslmode=disable -c 'CREATE DATABASE $(DB_NAME_TEST);'
endif

.PHONY: db_test_run
db_test_run: db_test_start db_test_create ## Run Test DB

.PHONY: db_test_reset
db_test_reset: db_test_destroy db_test_run ## Reset Test DB (destroy and run)

.PHONY: db_test_migrate_standalone
db_test_migrate_standalone: bin/orders ## Migrate Test DB directly
ifndef CIRCLECI
	@echo "Migrating the ${DB_NAME_TEST} database..."
	DB_DEBUG=0 DB_NAME=$(DB_NAME_TEST) DB_PORT=$(DB_PORT_TEST) bin/orders migrate -p "file://migrations/${APPLICATION}/secure;file://migrations/${APPLICATION}/schema" -m "migrations/${APPLICATION}/migrations_manifest.txt"
else
	@echo "Migrating the ${DB_NAME_TEST} database..."
	DB_DEBUG=0 DB_NAME=$(DB_NAME_TEST) DB_PORT=$(DB_PORT_DEV) bin/orders migrate -p "file://migrations/${APPLICATION}/secure;file://migrations/${APPLICATION}/schema" -m "migrations/${APPLICATION}/migrations_manifest.txt"
endif

.PHONY: db_test_migrate
db_test_migrate: db_test_migrate_standalone ## Migrate Test DB

.PHONY: db_test_migrations_build
db_test_migrations_build: .db_test_migrations_build.stamp ## Build Test DB Migrations Docker Image
.db_test_migrations_build.stamp:
	@echo "Build the docker migration container..."
	docker build -f Dockerfile.migrations_local --tag e2e_migrations:latest .

.PHONY: db_test_psql
db_test_psql: ## Open PostgreSQL shell for Test DB
	scripts/psql-test

#
# ----- END DB_TEST TARGETS -----
#

# ----- START RANDOM TARGETS -----
#

.PHONY: gofmt
gofmt:  ## Run go fmt over all Go files
	go fmt $$(go list ./...) >> /dev/null

.PHONY: pre_commit_tests
pre_commit_tests: .client_deps.stamp bin/swagger ## Run pre-commit tests
	pre-commit run --all-files

.PHONY: pretty
pretty: gofmt ## Run code through Golang formatters

.PHONY: docker_circleci
docker_circleci:
	docker run -it --rm=true -v $(PWD):$(PWD) -w $(PWD) trussworks/circleci-docker-primary:latest bash

.PHONY: prune_images
prune_images:  ## Prune docker images
	@echo '****************'
	docker image prune -a

.PHONY: prune_containers
prune_containers:  ## Prune docker containers
	@echo '****************'
	docker container prune

.PHONY: prune_volumes
prune_volumes:  ## Prune docker volumes
	@echo '****************'
	docker volume prune

.PHONY: prune
prune: prune_images prune_containers prune_volumes ## Prune docker containers, images, and volumes

.PHONY: clean
clean: ## Clean all generated files
	rm -f .*.stamp
	rm -rf ./bin
	rm -rf ./tmp/secure_migrations
	rm -rf ./log

.PHONY: spellcheck
spellcheck: ## Run interactive spellchecker
	@which mdspell -s || (echo "Install mdspell with yarn global add markdown-spellcheck" && exit 1)
	/usr/local/bin/mdspell --ignore-numbers --ignore-acronyms --en-us --no-suggestions \
		`find . -type f -name "*.md" \
			-not -path "./node_modules/*" \
			-not -path "./vendor/*" \
			-not -path "./docs/adr/index.md" | sort`

#
# ----- END RANDOM TARGETS -----
#
