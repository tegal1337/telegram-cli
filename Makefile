.PHONY: build run test lint clean install-tdlib

BINARY_NAME=tele-tui
BUILD_DIR=bin
TDLIB_DIR=$(HOME)/td

# TDLib build flags
CGO_CFLAGS=-I$(TDLIB_DIR)/tdlib/include
CGO_LDFLAGS=-L$(TDLIB_DIR)/tdlib/lib -ltdjson

build:
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" \
		go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/teletui

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

test:
	go test ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BUILD_DIR)
	rm -rf .tdlib

install-tdlib:
	@echo "Installing TDLib..."
	@echo "See: https://github.com/zelenin/go-tdlib#installation"
	@echo ""
	@echo "Ubuntu/Debian:"
	@echo "  sudo apt install -y build-essential cmake gperf zlib1g-dev libssl-dev"
	@echo "  git clone --depth 1 https://github.com/tdlib/td.git $(TDLIB_DIR)/td-src"
	@echo "  cd $(TDLIB_DIR)/td-src && mkdir build && cd build && cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=$(TDLIB_DIR)/tdlib .. && cmake --build . -j$$(nproc) && cmake --install ."
