.PHONY: build run test lint clean setup

BINARY_NAME=tele-tui
BUILD_DIR=bin
TDLIB_DIR=$(HOME)/td

# Auto-detect TDLib location
TDLIB_PATHS = $(TDLIB_DIR)/tdlib /usr/local /usr /opt/homebrew/opt/tdlib /opt/homebrew
TDLIB_PREFIX := $(shell for p in $(TDLIB_PATHS); do \
	if [ -f "$$p/lib/libtdjson.so" ] || [ -f "$$p/lib/libtdjson.dylib" ] || [ -f "$$p/bin/tdjson.dll" ]; then \
		echo "$$p"; break; \
	fi; \
done)

ifeq ($(TDLIB_PREFIX),)
$(warning TDLib not found. Run: make setup)
TDLIB_PREFIX=$(TDLIB_DIR)/tdlib
endif

CGO_CFLAGS=-I$(TDLIB_PREFIX)/include
CGO_LDFLAGS=-L$(TDLIB_PREFIX)/lib -ltdjson

build:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" \
		go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/teletui
	@echo "Built: $(BUILD_DIR)/$(BINARY_NAME)"

run: build
	LD_LIBRARY_PATH=$(TDLIB_PREFIX)/lib:$$LD_LIBRARY_PATH ./$(BUILD_DIR)/$(BINARY_NAME)

test:
	go test ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BUILD_DIR)

# Auto-install everything: TDLib + deps + build
setup:
	@./scripts/setup-tdlib.sh
	@echo ""
	@echo "Done! Now run: make build && make run"
