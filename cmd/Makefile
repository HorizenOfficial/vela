BIN_DIR := ../bin
FLAG_FILE := $(BIN_DIR)/.built
GIT_VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -X github.com/horizen-pes/pkg/version.Version=$(GIT_VERSION)

# Configuration files
MANAGER_CONF_SRC := ../pkg/manager/manager.conf
MANAGER_CONF_DEST := $(BIN_DIR)/manager.conf
EXECUTOR_CONF_SRC := ../pkg/executor/executor.conf
EXECUTOR_CONF_DEST := $(BIN_DIR)/executor.conf

# Automatically discover all Go source files for proper dependency tracking
MANAGER_SRC_FILES := $(shell go list -deps -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./manager)
EXECUTOR_SRC_FILES := $(shell go list -deps -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./executor)
ADMINCLI_SRC_FILES := $(shell go list -deps -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./admincli)

.PHONY: all manager executor admincli check-conf clean list-deps clean-flag print-nothing

all: clean-flag manager executor admincli check-conf print-nothing

# Ensure binary directory exists
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# ----------------------
# Manager
# ----------------------
manager: $(BIN_DIR)/manager

$(BIN_DIR)/manager: $(BIN_DIR) ../go.mod ../go.sum $(MANAGER_SRC_FILES)
	@echo "Building manager (version: $(GIT_VERSION))..."
	@go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/manager ./manager
	@echo "Manager binary built"
	@touch $(FLAG_FILE)

# ----------------------
# Executor
# ----------------------
executor: $(BIN_DIR)/executor

$(BIN_DIR)/executor: $(BIN_DIR) ../go.mod ../go.sum $(EXECUTOR_SRC_FILES)
	@echo "Building executor (version: $(GIT_VERSION))..."
	@go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/executor ./executor
	@echo "Executor binary built"
	@touch $(FLAG_FILE)

# ----------------------
# Admin CLI
# ----------------------
admincli: $(BIN_DIR)/admincli

$(BIN_DIR)/admincli: $(BIN_DIR) ../go.mod ../go.sum $(ADMINCLI_SRC_FILES)
	@echo "Building admincli..."
	@go build -o $(BIN_DIR)/admincli ./admincli
	@echo "Admin CLI binary built"
	@touch $(FLAG_FILE)

# ----------------------
# Config symlinks
# ----------------------
check-conf: $(BIN_DIR)
	@if [ ! -e $(MANAGER_CONF_DEST) ]; then \
		echo "Creating symlink: $(MANAGER_CONF_DEST) -> $(MANAGER_CONF_SRC)"; \
		ln -s $(MANAGER_CONF_SRC) $(MANAGER_CONF_DEST); \
		touch $(FLAG_FILE); \
	fi
	@if [ ! -e $(EXECUTOR_CONF_DEST) ]; then \
		echo "Creating symlink: $(EXECUTOR_CONF_DEST) -> $(EXECUTOR_CONF_SRC)"; \
		ln -s $(EXECUTOR_CONF_SRC) $(EXECUTOR_CONF_DEST); \
		touch $(FLAG_FILE); \
	fi

# ----------------------
# Nothing to do message
# ----------------------
clean-flag:
	@rm -f $(FLAG_FILE)

print-nothing:
	@if [ ! -f $(FLAG_FILE) ]; then \
		echo "Nothing to do"; \
	fi
	@rm -f $(FLAG_FILE)

# ----------------------
# Dependencies
# ----------------------
list-deps:
	@echo "Manager Source Dependencies:"
	@echo $(MANAGER_SRC_FILES) | tr ' ' '\n'
	@echo "\nExecutor Source Dependencies:"
	@echo $(EXECUTOR_SRC_FILES) | tr ' ' '\n'
	@echo "\nAdmin CLI Source Dependencies:"
	@echo $(ADMINCLI_SRC_FILES) | tr ' ' '\n'

# ----------------------
# Clean
# ----------------------
clean:
	@echo "Cleaning up..."
	@rm -f $(BIN_DIR)/manager $(BIN_DIR)/executor $(BIN_DIR)/admincli
	@rm -f $(MANAGER_CONF_DEST) $(EXECUTOR_CONF_DEST)
