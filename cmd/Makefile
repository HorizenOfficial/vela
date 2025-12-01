BIN_DIR := ../bin
MANAGER_CONF_SRC := ../pkg/manager/manager.conf
MANAGER_CONF_DEST := $(BIN_DIR)/manager.conf
EXECUTOR_CONF_SRC := ../pkg/executor/executor.conf
EXECUTOR_CONF_DEST := $(BIN_DIR)/executor.conf

.PHONY: all manager executor clean check-conf

all: manager executor check-conf

manager: $(BIN_DIR)/manager
	@echo "Manager binary built"

$(BIN_DIR)/manager: ./manager/main.go
	@echo "Building manager..."
	@go build -o $(BIN_DIR)/manager ./manager

executor: $(BIN_DIR)/executor
	@echo "Executor binary built"

$(BIN_DIR)/executor: ./executor/main.go
	@echo "Building executor..."
	@go build -o $(BIN_DIR)/executor ./executor

check-conf:
	@echo "Checking configuration file symlinks..."
	@mkdir -p $(BIN_DIR)
	@if [ ! -L $(MANAGER_CONF_DEST) ] || [ "$(realpath $(MANAGER_CONF_DEST))" != "$(realpath $(MANAGER_CONF_SRC))" ]; then \
		echo "Creating symlink for manager.conf"; \
		rm -f $(MANAGER_CONF_DEST); \
		ln -s $(MANAGER_CONF_SRC) $(MANAGER_CONF_DEST); \
	else \
		echo "manager.conf symlink is correct."; \
	fi
	@if [ ! -L $(EXECUTOR_CONF_DEST) ] || [ "$(realpath $(EXECUTOR_CONF_DEST))" != "$(realpath $(EXECUTOR_CONF_SRC))" ]; then \
		echo "Creating symlink for executor.conf"; \
		rm -f $(EXECUTOR_CONF_DEST); \
		ln -s $(EXECUTOR_CONF_SRC) $(EXECUTOR_CONF_DEST); \
	else \
		echo "executor.conf symlink is correct."; \
	fi

clean:
	@echo "Cleaning up..."
	@rm -f $(BIN_DIR)/manager $(BIN_DIR)/executor
	@rm -f $(MANAGER_CONF_DEST) $(EXECUTOR_CONF_DEST)
