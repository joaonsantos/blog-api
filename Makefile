# global options
Q=@

# go options
PKGS       := ./cmd/...
TESTFLAGS  :=
GOFLAGS    :=
SERVER_CMD_PATH := ./cmd/blog/blog.go
CLI_CMD_PATH := ./cmd/cli/cli.go

# server options
SERVER_ADDR := :8080

# db options
DB_DSN := file:blog.db?cache=shared

# ------------------------------------------------------------------------------
#  run
.PHONY: run
run:
	$Qgo run $(SERVER_CMD_PATH) -addr $(SERVER_ADDR) -db $(DB_DSN)

# ------------------------------------------------------------------------------
#  init
.PHONY: init
init:
	$Qgo run $(CLI_CMD_PATH) -db $(DB_DSN)

# ------------------------------------------------------------------------------
#  test
.PHONY: test
test: test-unit

.PHONY: test-unit
test-unit:
	@echo
	@echo "==> Running unit tests <=="
	@echo
	$Qgo test $(GOFLAGS) $(PKGS) $(TESTFLAGS)
