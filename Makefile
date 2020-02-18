

# Convenience for LDFLAGS
WEBSERVER_LDFLAGS=-X main.gitBranch=$(shell git branch | grep \* | cut -d ' ' -f2) -X main.gitCommit=$(shell git rev-list -1 HEAD)
GC_FLAGS=-trimpath=$(GOPATH)


.PHONY: help
help:  ## Print the help documentation
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


#
# ----- END PREAMBLE -----
#

#
# ----- START CHECK TARGETS -----
#


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
pkg/gen/: $(shell find swagger -type f -name *.yaml)
	scripts/gen-server

.PHONY: server_build
server_build: bin/orders ## Build the server

# This command is for running the server by itself, it will serve the compiled frontend on its own
# Note: Don't double wrap with aws-vault because the pkg/cli/vault.go will handle it
server_run_standalone: check_log_dir server_build db_dev_run
	DEBUG_LOGGING=true ./bin/milmove serve 2>&1 | tee -a log/dev.log

# This command will rebuild the swagger go code and rerun server on any changes
server_run:
	find ./swagger -type f -name "*.yaml" | entr -c -r make server_run_default
# This command runs the server behind gin, a hot-reload server
# Note: Gin is not being used as a proxy so assigning odd port and laddr to keep in IPv4 space.
# Note: The INTERFACE envar is set to configure the gin build, milmove_gin, local IP4 space with default port 8080.
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
