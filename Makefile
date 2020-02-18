

#
# ----- START CHECK TARGETS -----
#


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
#
# ----- START BIN TARGETS -----
#

### Go Tool Targets

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
