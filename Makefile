# global options
Q=@

# go options
PKGS       := ./cmd/...
TESTFLAGS  :=
GOFLAGS    :=

# server options
SERVER_ADDR :=

# ------------------------------------------------------------------------------
#  run
.PHONY: run
run:
	$QDB_DSN="./blog.db" go run ./cmd/blog/blog.go $(SERVER_ADDR)

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
