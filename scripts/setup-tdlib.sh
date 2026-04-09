#!/usr/bin/env bash
set -euo pipefail

TDLIB_DIR="${TDLIB_DIR:-$HOME/td}"
TDLIB_PREFIX="${TDLIB_DIR}/tdlib"
NPROC=$(nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4)

BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()  { echo -e "${BLUE}[*]${NC} $1"; }
ok()    { echo -e "${GREEN}[✓]${NC} $1"; }
warn()  { echo -e "${YELLOW}[!]${NC} $1"; }

# Check if TDLib is already installed
check_tdlib() {
    for path in "$TDLIB_PREFIX" /usr/local /usr /opt/homebrew/opt/tdlib /opt/homebrew; do
        if [ -f "$path/lib/libtdjson.so" ] || [ -f "$path/lib/libtdjson.dylib" ]; then
            ok "TDLib found at $path"
            return 0
        fi
    done
    return 1
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "macos" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *)       echo "unknown" ;;
    esac
}

# Install deps on Linux
install_deps_linux() {
    info "Installing build dependencies..."
    if command -v apt-get &>/dev/null; then
        sudo apt-get update -qq
        sudo apt-get install -y -qq build-essential cmake gperf zlib1g-dev libssl-dev git
    elif command -v dnf &>/dev/null; then
        sudo dnf install -y gcc-c++ cmake gperf zlib-devel openssl-devel git
    elif command -v pacman &>/dev/null; then
        sudo pacman -S --noconfirm base-devel cmake gperf openssl zlib git
    else
        warn "Unknown package manager. Please install: cmake, gperf, zlib, openssl, git"
        exit 1
    fi
    ok "Dependencies installed"
}

# Install on macOS via brew
install_macos() {
    if ! command -v brew &>/dev/null; then
        warn "Homebrew not found. Install from https://brew.sh"
        exit 1
    fi
    info "Installing TDLib via Homebrew..."
    brew install tdlib
    ok "TDLib installed via brew"
}

# Build TDLib from source
build_tdlib() {
    info "Building TDLib from source (this takes a few minutes)..."

    if [ -d "${TDLIB_DIR}/td-src" ]; then
        info "Source already exists, updating..."
        cd "${TDLIB_DIR}/td-src" && git pull --depth 1 2>/dev/null || true
    else
        mkdir -p "$TDLIB_DIR"
        git clone --depth 1 https://github.com/tdlib/td.git "${TDLIB_DIR}/td-src"
    fi

    cd "${TDLIB_DIR}/td-src"
    rm -rf build && mkdir build && cd build

    info "Configuring CMake..."
    cmake -DCMAKE_BUILD_TYPE=Release \
          -DCMAKE_INSTALL_PREFIX="$TDLIB_PREFIX" \
          .. 2>&1 | tail -5

    info "Compiling with ${NPROC} cores..."
    cmake --build . -j"$NPROC" 2>&1 | tail -3

    info "Installing to ${TDLIB_PREFIX}..."
    cmake --install .

    ok "TDLib built and installed to ${TDLIB_PREFIX}"

    # Add to ldconfig on Linux
    if [ "$(detect_os)" = "linux" ]; then
        echo "${TDLIB_PREFIX}/lib" | sudo tee /etc/ld.so.conf.d/tdlib.conf >/dev/null 2>&1 || true
        sudo ldconfig 2>/dev/null || true
        ok "Library path registered with ldconfig"
    fi
}

main() {
    echo ""
    echo "  ╭──────────────────────────────────╮"
    echo "  │  Telegram CLI - TDLib Setup      │"
    echo "  ╰──────────────────────────────────╯"
    echo ""

    if check_tdlib; then
        ok "TDLib is already installed. Nothing to do."
        return 0
    fi

    OS=$(detect_os)
    info "Detected OS: ${OS}"

    case "$OS" in
        macos)
            install_macos
            ;;
        linux)
            install_deps_linux
            build_tdlib
            ;;
        windows)
            warn "On Windows, use MSYS2:"
            warn "  pacman -S mingw-w64-x86_64-cmake mingw-w64-x86_64-gperf"
            warn "  Then build TDLib manually or use vcpkg"
            exit 1
            ;;
        *)
            warn "Unsupported OS. Build TDLib manually:"
            warn "  https://github.com/tdlib/td#building"
            exit 1
            ;;
    esac
}

main "$@"
