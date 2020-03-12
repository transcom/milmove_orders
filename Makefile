DB_NAME_DEV = dev_db
DB_NAME_DEPLOYED_MIGRATIONS = deployed_migrations
DB_NAME_TEST = test_db
# The version of the postgres container should match production as closely as possible.
DB_DOCKER_CONTAINER_IMAGE = postgres:10.10
DB_PORT_DEV=5432
DB_PORT_TEST=5432
export PGPASSWORD=mysecretpassword

# Convenience for LDFLAGS
WEBSERVER_LDFLAGS=-X main.gitBranch=$(shell git branch | grep \* | cut -d ' ' -f2) -X main.gitCommit=$(shell git rev-list -1 HEAD)
GC_FLAGS=-trimpath=$(GOPATH)
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	LDFLAGS=-linkmode external -extldflags -static
endif

.PHONY: help
help:  ## Print the help documentation
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: dev
dev:
	docker-compose -f docker-compose.dev.yml build --pull
	aws-vault exec "${AWS_PROFILE}" -- docker-compose -f docker-compose.dev.yml run --service-ports --rm --name orders_dev dev bash

.PHONY: dev_destroy
dev_destroy:
	docker-compose -f docker-compose.dev.yml down

.PHONY: clean
clean: ## Clean all generated files
	rm -rf ./bin
	rm -rf ./tmp/secure_migrations
	rm -rf ./log

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

.PHONY: pre_commit_tests
pre_commit_tests: ## Run pre-commit tests
	pre-commit run --all-files

.PHONY: check_hosts
check_hosts: scripts/check-hosts-file ## Check that hosts are in the /etc/hosts file
	scripts/check-hosts-file

#
# ----- END CHECK TARGETS -----
#

#
# ----- START BIN TARGETS -----
#

### Go Tool Targets

bin/gin:
	go build -ldflags "$(LDFLAGS)" -o bin/gin github.com/codegangsta/gin

# No static linking / $(LDFLAGS) because go-junit-report is only used for building the CirlceCi test report
bin/go-junit-report:
	go build -o bin/go-junit-report github.com/jstemmer/go-junit-report

### Cert Targets

bin/rds-ca-2019-root.pem:
	mkdir -p bin/
	curl -sSo bin/rds-ca-2019-root.pem https://s3.amazonaws.com/rds-downloads/rds-ca-2019-root.pem

### Orders Targets

bin/ecs-deploy:
	go build -ldflags "$(LDFLAGS)" -o bin/ecs-deploy github.com/transcom/mymove/cmd/ecs-deploy

bin/orders:
	go build -gcflags="$(GC_FLAGS)" -asmflags=-trimpath=$(GOPATH) -ldflags "$(LDFLAGS) $(WEBSERVER_LDFLAGS)" -o bin/orders ./cmd/orders

#
# ----- END BIN TARGETS -----
#

#
# ----- START SERVER TARGETS -----
#

.PHONY: check_log_dir
check_log_dir: ## Make sure we have a log directory
	mkdir -p log

.PHONY: server_generate
server_generate: pkg/gen/ ## Generate golang server code from Swagger files
pkg/gen/: $(shell find swagger -type f -name *.yaml)
	scripts/gen-server

# This command is for running the server by itself, it will serve the compiled frontend on its own
# Note: Don't double wrap with aws-vault because the pkg/cli/vault.go will handle it
server_run_standalone: check_log_dir bin/orders db_dev_run
	DEBUG_LOGGING=true ./bin/orders serve 2>&1 | tee -a log/dev.log

# This command will rebuild the swagger go code and rerun server on any changes
server_run:
	find ./swagger -type f -name "*.yaml" | entr -c -r make server_run_default
# This command runs the server behind gin, a hot-reload server
# Note: Gin is not being used as a proxy so assigning odd port and laddr to keep in IPv4 space.
# Note: The INTERFACE envar is set to configure the gin build, orders_gin, local IP4 space with default port GIN_PORT.
server_run_default: check_log_dir bin/gin server_generate db_dev_run
	INTERFACE=localhost DEBUG_LOGGING=true \
	$(AWS_VAULT) ./bin/gin \
		--build ./cmd/orders \
		--bin /bin/orders_gin \
		--laddr 127.0.0.1 --port "$(GIN_PORT)" \
		--excludeDir node_modules \
		--immediate \
		--buildArgs "-i -ldflags=\"$(WEBSERVER_LDFLAGS)\"" \
		serve \
		2>&1 | tee -a log/dev.log

# acceptance_test runs a few acceptance tests against a local or remote environment.
# This can help identify potential errors before deploying a container.
.PHONY: acceptance_test
acceptance_test: bin/rds-ca-2019-root.pem ## Run acceptance tests
ifndef TEST_ACC_ENV
	@echo "Running acceptance tests for webserver using local environment."
	@echo "* Use environment XYZ by setting environment variable to TEST_ACC_ENV=XYZ."
	TEST_ACC_CWD=$(PWD) \
	SERVE_ADMIN=true \
	SERVE_SDDC=true \
	SERVE_ORDERS=true \
	SERVE_DPS=true \
	SERVE_API_INTERNAL=true \
	SERVE_API_GHC=false \
	MUTUAL_TLS_ENABLED=true \
	go test -v -count 1 -short $$(go list ./... | grep \\/cmd\\/orders)
else
	@echo "Running acceptance tests for webserver with environment $$TEST_ACC_ENV."
	TEST_ACC_CWD=$(PWD) \
	chamber -r $(CHAMBER_RETRIES) exec orders-$(TEST_ACC_ENV) -- \
	go test -v -count 1 -short $$(go list ./... | grep \\/cmd\\/orders)
endif

#
# ----- END SERVER TARGETS -----
#

#
# ----- START SERVER TEST TARGETS -----
#

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

#
# ----- END SERVER TEST TARGETS -----
#

#
# ----- START DB_DEV TARGETS -----
#

.PHONY: db_dev_destroy
db_dev_destroy: ## Destroy Dev DB
	@echo "Destroying the ${DB_NAME_DEV} database ..."
	/usr/bin/psql --variable "ON_ERROR_STOP=1" postgres://postgres:"${DB_PASSWORD}"@${DB_HOST}:"${DB_PORT}" -c "DROP DATABASE IF EXISTS ${DB_NAME_DEV};" || true

.PHONY: db_dev_create
db_dev_create: ## Create Dev DB
	@echo "Create the ${DB_NAME_DEV} database..."
	DB_NAME=postgres scripts/wait-for-db
	/usr/bin/psql --variable "ON_ERROR_STOP=1" postgres://postgres:"${DB_PASSWORD}"@${DB_HOST}:"${DB_PORT}" -c "CREATE DATABASE ${DB_NAME_DEV}" || true

.PHONY: db_dev_run
db_dev_run: db_dev_create ## Run Dev DB (start and create)

.PHONY: db_dev_reset
db_dev_reset: db_dev_destroy db_dev_run ## Reset Dev DB (destroy and run)

.PHONY: db_dev_migrate_standalone
db_dev_migrate_standalone: bin/orders ## Migrate Dev DB directly
	@echo "Migrating the ${DB_NAME_DEV} database..."
	DB_DEBUG=0 bin/orders migrate -p "file://migrations/${APPLICATION}/secure;file://migrations/${APPLICATION}/schema" -m "migrations/${APPLICATION}/migrations_manifest.txt"

.PHONY: db_dev_migrate
db_dev_migrate: db_dev_migrate_standalone ## Migrate Dev DB

.PHONY: db_dev_psql
db_dev_psql: ## Open PostgreSQL shell for Dev DB
	/usr/bin/psql --variable "ON_ERROR_STOP=1" postgres://postgres:"${DB_PASSWORD}"@${DB_HOST}:"${DB_PORT}"/"${DB_NAME}"

#
# ----- END DB_DEV TARGETS -----
#

#
# ----- START DB_DEPLOYED_MIGRATIONS TARGETS -----
#

.PHONY: db_deployed_migrations_destroy
db_deployed_migrations_destroy: ## Destroy Deployed Migrations DB
	@echo "Destroying the ${DB_NAME_DEPLOYED_MIGRATIONS} database ..."
	/usr/bin/psql --variable "ON_ERROR_STOP=1" postgres://postgres:"${DB_PASSWORD}"@${DB_HOST}:"${DB_PORT}" -c "DROP DATABASE IF EXISTS ${DB_NAME_DEPLOYED_MIGRATIONS};" || true

.PHONY: db_deployed_migrations_create
db_deployed_migrations_create: ## Create Deployed Migrations DB
	@echo "Create the ${DB_NAME_DEPLOYED_MIGRATIONS} database..."
	DB_NAME=postgres scripts/wait-for-db
	/usr/bin/psql --variable "ON_ERROR_STOP=1" postgres://postgres:"${DB_PASSWORD}"@${DB_HOST}:"${DB_PORT}" -c "CREATE DATABASE ${DB_NAME_DEPLOYED_MIGRATIONS}" || true

.PHONY: db_deployed_migrations_run
db_deployed_migrations_run: db_deployed_migrations_create ## Run Deployed Migrations DB (start and create)

.PHONY: db_deployed_migrations_reset
db_deployed_migrations_reset: db_deployed_migrations_destroy db_deployed_migrations_run ## Reset Deployed Migrations DB (destroy and run)

.PHONY: db_deployed_migrations_migrate_standalone
db_deployed_migrations_migrate_standalone: bin/orders ## Migrate Deployed Migrations DB with local secure migrations
	@echo "Migrating the ${DB_NAME_DEPLOYED_MIGRATIONS} database..."
	DB_NAME=${DB_NAME_DEPLOYED_MIGRATIONS} DB_DEBUG=0 bin/orders migrate -p "file://migrations/${APPLICATION}/secure;file://migrations/${APPLICATION}/schema" -m "migrations/${APPLICATION}/migrations_manifest.txt"

.PHONY: db_deployed_migrations_migrate
db_deployed_migrations_migrate: db_deployed_migrations_migrate_standalone ## Migrate Deployed Migrations DB

.PHONY: db_deployed_psql
db_deployed_migrations_psql: ## Open PostgreSQL shell for Deployed Migrations DB
	/usr/bin/psql --variable "ON_ERROR_STOP=1" postgres://postgres:"${DB_PASSWORD}"@${DB_HOST}:"${DB_PORT}"/"${DB_NAME_DEPLOYED_MIGRATIONS}"

#
# ----- END DB_DEPLOYED_MIGRATIONS TARGETS -----
#

#
# ----- START DB_TEST TARGETS -----
#

.PHONY: db_test_destroy
db_test_destroy: ## Destroy Test DB
	@echo "Destroying the ${DB_NAME_TEST} database ..."
	/usr/bin/psql --variable "ON_ERROR_STOP=1" postgres://postgres:"${DB_PASSWORD}"@${DB_HOST}:"${DB_PORT}" -c "DROP DATABASE IF EXISTS ${DB_NAME_TEST};" || true

.PHONY: db_test_create
db_test_create: ## Create Test DB
	@echo "Create the ${DB_NAME_TEST} database..."
	DB_NAME=postgres scripts/wait-for-db
	/usr/bin/psql --variable "ON_ERROR_STOP=1" postgres://postgres:"${DB_PASSWORD}"@${DB_HOST}:"${DB_PORT}" -c "CREATE DATABASE ${DB_NAME_TEST}" || true

.PHONY: db_test_run
db_test_run: db_test_create ## Run Test DB (start and create)

.PHONY: db_test_reset
db_test_reset: db_test_destroy db_test_run ## Reset Test DB (destroy and run)

.PHONY: db_test_migrate_standalone
db_test_migrate_standalone: bin/orders ## Migrate Test DB directly
	@echo "Migrating the ${DB_NAME_TEST} database..."
	DB_NAME=${DB_NAME_TEST} DB_DEBUG=0 bin/orders migrate -p "file://migrations/${APPLICATION}/secure;file://migrations/${APPLICATION}/schema" -m "migrations/${APPLICATION}/migrations_manifest.txt"

.PHONY: db_test_migrate
db_test_migrate: db_test_migrate_standalone ## Migrate Test DB

.PHONY: db_test_psql
db_test_psql: ## Open PostgreSQL shell for Test DB
	/usr/bin/psql --variable "ON_ERROR_STOP=1" postgres://postgres:"${DB_PASSWORD}"@${DB_HOST}:"${DB_PORT}"/"${DB_NAME_TEST}"

#
# ----- END DB_TEST TARGETS -----
#

#
# ----- START RUN DEPLOYED_MIGRATION TARGETS -----
#

.PHONY: run_experimental_migrations
run_experimental_migrations: bin/orders db_deployed_migrations_reset ## Run Experimental migrations against Deployed Migrations DB
	@echo "Migrating the ${DB_NAME_DEPLOYED_MIGRATIONS} database with experimental migrations..."
	MIGRATION_PATH="s3://transcom-ppp-${APPLICATION}-experimental-us-west-2/secure-migrations;file://migrations/$(APPLICATION)/schema" \
	DB_HOST=${DB_HOST} \
	DB_PORT=$(DB_PORT_DEPLOYED_MIGRATIONS) \
	DB_NAME=$(DB_NAME_DEPLOYED_MIGRATIONS) \
	DB_DEBUG=0 \
	bin/orders migrate

#
# ----- END RUN DEPLOYED_MIGRATION TARGETS -----
#

#
# ----- START DOCKER COMPOSE BRANCH TARGETS -----
#

.PHONY: docker_compose_branch_up
docker_compose_branch_up: ## Bring up docker-compose containers for current branch with AWS ECR images
	aws-vault exec "${AWS_PROFILE}" -- docker run -it -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -e AWS_SECURITY_TOKEN -e AWS_SESSION_TOKEN milmove/circleci-docker:milmove-orders aws ecr get-login-password --region "${AWS_DEFAULT_REGION}" | docker login --username AWS --password-stdin "${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com"
	scripts/update-docker-compose
	aws-vault exec "${AWS_PROFILE}" -- docker-compose -f docker-compose.branch.yml up --remove-orphans

.PHONY: docker_compose_branch_down
docker_compose_branch_down: ## Destroy docker-compose containers for current branch
	docker-compose -f docker-compose.branch.yml down
	# Instead of using `--rmi all` which might destroy postgres we just remove the AWS containers
	docker rmi $(shell docker images --filter=reference='*amazonaws*/*:*' --format "{{.ID}}")
	git checkout docker-compose.yml

.PHONY: docker_compose_local_up
docker_compose_local_up: ## Bring up docker-compose containers for current local with AWS ECR images
	aws-vault exec "${AWS_PROFILE}" -- docker-compose -f docker-compose.local.yml up --remove-orphans

.PHONY: docker_compose_local_down
docker_compose_local_down: ## Destroy docker-compose containers for current local
	docker-compose -f docker-compose.local.yml down --rmi local

#
# ----- END DOCKER COMPOSE TARGETS -----
#
