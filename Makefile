# global options
Q=@

# go options
PKGS       := ./cmd/...
TESTFLAGS  :=
GOFLAGS    :=
SERVER_CMD_PATH := ./cmd/blog/blog.go
CLI_CMD_PATH := ./cmd/cli/cli.go

# server options
SERVER_ADDR :=

# db options
DB_DSN := ./blog.db

# ------------------------------------------------------------------------------
#  run
.PHONY: run
run:
	$QDB_DSN=$(DB_DSN) go run $(SERVER_CMD_PATH) $(SERVER_ADDR)

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
