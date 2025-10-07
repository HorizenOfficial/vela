APP := simple_app
OUT := build/$(APP).wasm
PROD_OUT := production_build/$(APP).wasm

SRC := $(shell find . -name '*.go')

.PHONY: all build production_build clean

build: $(OUT)

production_build: $(PROD_OUT)

all: build production_build

clean:
	rm -rf build production_build

$(OUT): $(SRC)
	@mkdir -p $(@D)
	@echo "Building development WASM module..."
	tinygo build -o $@ -target=wasi .

$(PROD_OUT): $(SRC)
	@mkdir -p $(@D)
	@echo "Building production WASM module..."
	tinygo build -o $@ -opt=2 -target=wasi .
