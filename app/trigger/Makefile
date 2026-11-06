APP := trigger_app
OUT := build/$(APP).wasm
PROD_OUT := production_build/$(APP).wasm
WAT := $(OUT:.wasm=.wat)
WAT_PROD := $(PROD_OUT:.wasm=.wat)

SRC := $(shell find . -name '*.go')

.PHONY: all build production_build clean wat production_wat

build: $(OUT)

production_build: $(PROD_OUT)

all: build production_build wat production_wat

clean:
	rm -rf build production_build

$(OUT): $(SRC)
	@mkdir -p $(@D)
	@echo "Building development WASM module..."
	tinygo build -o $@ -target=wasip1 .
	@ls -lh $@ | awk '{print "Binary size: " $$5}'

$(PROD_OUT): $(SRC)
	@mkdir -p $(@D)
	@echo "Building production WASM module..."
	tinygo build -o $@ -opt=s -no-debug -target=wasip1 .
	@ls -lh $@ | awk '{print "Binary size: " $$5}'

wat: $(WAT)

production_wat: $(WAT_PROD)

%.wat: %.wasm
	@echo "Converting WASM to WAT..."
	wasm2wat $< -o $@
