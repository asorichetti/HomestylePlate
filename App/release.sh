#!/bin/bash
# ============================================================
# Homestyle Plate - Release Script
# Builds distributable releases for macOS, Linux, and Windows
# ============================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
VERSION="$(cat "$SCRIPT_DIR/VERSION" | tr -d '[:space:]')"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

# Check Docker is running (needed for cross-compilation from macOS)
check_docker() {
    if command -v docker &>/dev/null && docker info &>/dev/null; then
        return 0
    fi
    return 1
}

log_info()    { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[OK]${NC}   $1"; }
log_warn()    { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error()   { echo -e "${RED}[ERR]${NC}  $1"; }
log_step()    { echo -e "${CYAN}→${NC}    $1"; }

# Parse args
TARGET_PLATFORM="${1:-all}"
APP_NAME="Homestyle Plate"
BINARY_NAME="HomestylePlate"

# ============================================================
# Build binary for target OS/Arch
# ============================================================
build_binary() {
    local target_os=$1
    local target_arch=$2
    local output=$3
    local is_cross=$4  # "yes" if cross-compiling from macOS
    
    log_step "Building binary for $target_os/$target_arch..."
    
    if [[ "$is_cross" == "yes" ]]; then
        # On macOS, cross-compiling requires either a cross-compiler or Docker
        if [[ "$(uname -s)" != "Darwin" ]]; then
            # Not on macOS, native cross-compilation should work
            : # fall through to native build
        elif [[ "$target_os" == "windows" ]]; then
            if command -v x86_64-w64-mingw32-gcc &>/dev/null; then
                export CC=x86_64-w64-mingw32-gcc
                log_info "Using mingw cross-compiler (brew install mingw-w64)"
            elif check_docker; then
                log_info "Using Docker for Windows cross-compile..."
                cross_compile_docker "windows" "amd64" "$output"
                return $?
            else
                log_error "Cannot cross-compile to Windows from macOS."
                log_error "Install either:"
                log_error "  1. mingw-w64:  brew install mingw-w64"
                log_error "  2. Docker:     brew install --cask docker && docker start"
                return 1
            fi
        elif [[ "$target_os" == "linux" ]]; then
            if check_docker; then
                log_info "Using Docker for Linux cross-compile..."
                cross_compile_docker "linux" "$target_arch" "$output"
                return $?
            else
                log_error "Cannot cross-compile to Linux from macOS."
                log_error "Install either:"
                log_error "  1. Linux machine with gcc"
                log_error "  2. Docker:     brew install --cask docker && docker start"
                return 1
            fi
        fi
        # Cross-compilation: need appropriate toolchain
        if [[ "$target_os" == "windows" ]]; then
            # Check for mingw or use docker
            if command -v x86_64-w64-mingw32-gcc &>/dev/null; then
                export CC=x86_64-w64-mingw32-gcc
                log_info "Using mingw cross-compiler"
            elif command -v docker &>/dev/null; then
                log_info "Using Docker for Windows cross-compile..."
                cross_compile_docker "windows" "amd64" "$output"
                export CC=""
                return $?
            else
                log_error "Cross-compiling to Windows requires either:"
                log_error "  1. mingw-w64: brew install mingw-w64"
                log_error "  2. Docker: docker pull golang:1.23"
                return 1
            fi
        elif [[ "$target_os" == "linux" ]]; then
            if command -v gcc &>/dev/null && [[ "$(uname -s)" == "Linux" ]]; then
                # Native Linux - gcc works
                log_info "Using native gcc"
            elif command -v docker &>/dev/null; then
                log_info "Using Docker for Linux cross-compile..."
                cross_compile_docker "linux" "$target_arch" "$output"
                export CC=""
                return $?
            else
                log_error "Cross-compiling to Linux requires either:"
                log_error "  1. A Linux machine with gcc, or"
                log_error "  2. Docker: brew install --cask docker"
                return 1
            fi
        fi
    fi
    
    export CGO_ENABLED=1
    export GOOS="$target_os"
    export GOARCH="$target_arch"
    
    go build -ldflags="-s -w" -o "$output" "$SCRIPT_DIR/server.go"
    
    unset CGO_ENABLED GOOS GOARCH CC CXX
    return 0
}

# Docker-based cross compilation
cross_compile_docker() {
    local target_os=$1
    local target_arch=$2
    local output=$3
    
    log_step "Pulling Go build image..."
    
    # On Apple Silicon, we need linux/amd64 platform for x86_64 builds
    local platform=""
    if [[ "$(uname -m)" == "arm64" ]]; then
        platform="--platform linux/amd64"
    fi
    
    # Use a build container - gcc is pre-installed in golang image
    if [[ "$target_os" == "windows" ]]; then
        log_step "Installing cross-compiler and compiling Windows binary..."
        docker run --rm $platform \
            -v "$SCRIPT_DIR:/src" \
            -v "${SCRIPT_DIR}/build_cross:/build" \
            -w /src \
            golang:1.23-bookworm \
            bash -c "
                apt-get clean > /dev/null 2>&1 && \
                apt-get update -o Acquire::AllowInsecureRepositories=true 2>&1 | tail -1 > /dev/null && \
                apt-get install -y --allow-unauthenticated --no-install-recommends gcc-mingw-w64-x86-64 > /dev/null 2>&1 && \
                export CGO_ENABLED=1 && \
                export GOOS=windows && \
                export GOARCH=amd64 && \
                export CC=x86_64-w64-mingw32-gcc && \
                export CXX=x86_64-w64-mingw32-g++ && \
                go mod download > /dev/null 2>&1 && \
                go build -ldflags='-s -w' -o /build/HomestylePlate.exe server.go && \
                echo 'Build succeeded' || echo 'Build failed'
            "
        if [[ -f "${SCRIPT_DIR}/build_cross/HomestylePlate.exe" ]]; then
            # Keep .exe extension for Windows
            mv "${SCRIPT_DIR}/build_cross/HomestylePlate.exe" "${output}.exe"
        fi
        
    elif [[ "$target_os" == "linux" && "$target_arch" == "amd64" ]]; then
        # Linux x86_64 - gcc is already in the image
        log_step "Compiling Linux binary..."
        docker run --rm $platform \
            -v "$SCRIPT_DIR:/src" \
            -v "${SCRIPT_DIR}/build_cross:/build" \
            -w /src \
            golang:1.23-bookworm \
            bash -c "
                export CGO_ENABLED=1 && \
                export GOOS=linux && \
                export GOARCH=amd64 && \
                go mod download > /dev/null 2>&1 && \
                go build -ldflags='-s -w' -o /build/HomestylePlate server.go && \
                echo 'Build succeeded' || echo 'Build failed'
            "
        if [[ -f "${SCRIPT_DIR}/build_cross/HomestylePlate" ]]; then
            mv "${SCRIPT_DIR}/build_cross/HomestylePlate" "$output"
        fi
        
    elif [[ "$target_os" == "linux" && "$target_arch" == "arm64" ]]; then
        log_step "Installing ARM cross-compiler and compiling..."
        docker run --rm \
            -v "$SCRIPT_DIR:/src" \
            -v "${SCRIPT_DIR}/build_cross:/build" \
            -w /src \
            golang:1.23-bookworm \
            bash -c "
                apt-get clean > /dev/null 2>&1 && \
                apt-get update -o Acquire::AllowInsecureRepositories=true 2>&1 | tail -1 > /dev/null && \
                apt-get install -y --allow-unauthenticated --no-install-recommends gcc-aarch64-linux-gnu > /dev/null 2>&1 && \
                export CGO_ENABLED=1 && \
                export GOOS=linux && \
                export GOARCH=arm64 && \
                export CC=aarch64-linux-gnu-gcc && \
                export CXX=aarch64-linux-gnu-g++ && \
                go mod download > /dev/null 2>&1 && \
                go build -ldflags='-s -w' -o /build/HomestylePlate server.go && \
                echo 'Build succeeded' || echo 'Build failed'
            "
        if [[ -f "${SCRIPT_DIR}/build_cross/HomestylePlate" ]]; then
            mv "${SCRIPT_DIR}/build_cross/HomestylePlate" "$output"
        fi
    fi
    
    # Cleanup
    rm -rf "${SCRIPT_DIR}/build_cross"
    return 0
}

# ============================================================
# macOS DMG Build
# ============================================================
build_macos() {
    log_info "══════════════════════════════════════════"
    log_info "  Building macOS Release v$VERSION"
    log_info "══════════════════════════════════════════"
    
    local build_dir="$SCRIPT_DIR/build_macos"
    local app_dir="$build_dir/$APP_NAME.app"
    local output_dir="$REPO_ROOT/releases"
    
    rm -rf "$build_dir"
    mkdir -p "$output_dir" "$app_dir/Contents/MacOS" "$app_dir/Contents/Frontend"
    
    # Copy frontend
    log_step "Copying frontend assets..."
    cp -R "$REPO_ROOT/Frontend/"* "$app_dir/Contents/Frontend/"
    
    # Build binary
    build_binary "darwin" "arm64" "$app_dir/Contents/MacOS/$BINARY_NAME" "no"
    chmod +x "$app_dir/Contents/MacOS/$BINARY_NAME"
    
    # Info.plist
    cat > "$app_dir/Contents/Info.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key><string>$APP_NAME</string>
    <key>CFBundleExecutable</key><string>$BINARY_NAME</string>
    <key>CFBundleIdentifier</key><string>com.homestyleplate.app</string>
    <key>CFBundleVersion</key><string>$VERSION</string>
    <key>CFBundleShortVersionString</key><string>$VERSION</string>
    <key>CFBundlePackageType</key><string>APPL</string>
    <key>LSMinimumSystemVersion</key><string>11.0</string>
    <key>NSHighResolutionCapable</key><true/>
    <key>LSApplicationCategoryType</key><string>public.app-category.utilities</string>
</dict>
</plist>
EOF
    
    # README
    cat > "$app_dir/Contents/README.txt" << EOF
Homestyle Plate - Meal Randomizer v$VERSION

Installation:
  1. Drag to Applications folder
  2. Double-click to launch
  3. Browser opens at http://localhost:8080

Data is stored locally via SQLite. First run seeds meal data.
EOF
    
    # Create DMG
    log_step "Creating DMG..."
    local dmg="$output_dir/${APP_NAME}_${VERSION}_macos.dmg"
    hdiutil create \
        -volname "$APP_NAME" \
        -srcfolder "$app_dir" \
        -ov -format UDZO -fs APFS \
        "$dmg"
    
    rm -rf "$build_dir"
    log_success "macOS DMG: $(basename "$dmg") ($(du -h "$dmg" | cut -f1))"
}

# ============================================================
# Linux Tarball Build
# ============================================================
build_linux() {
    local arch=$1
    local arch_label="amd64"
    [[ "$arch" == "arm64" ]] && arch_label="arm64"
    
    log_info "══════════════════════════════════════════"
    log_info "  Building Linux v$VERSION ($arch_label)"
    log_info "══════════════════════════════════════════"
    
    local build_dir="$SCRIPT_DIR/build_linux"
    local output_dir="$REPO_ROOT/releases"
    local is_cross="no"
    [[ "$(uname -s)" == "Darwin" ]] && is_cross="yes"
    
    rm -rf "$build_dir"
    mkdir -p "$output_dir" "$build_dir/$APP_NAME/Frontend"
    
    log_step "Copying frontend assets..."
    cp -R "$REPO_ROOT/Frontend/"* "$build_dir/$APP_NAME/Frontend/"
    
    log_step "Building binary..."
    if ! build_binary "linux" "$arch" "$build_dir/$APP_NAME/$BINARY_NAME" "$is_cross"; then
        return 1
    fi
    chmod +x "$build_dir/$APP_NAME/$BINARY_NAME"
    
    # Launcher
    cat > "$build_dir/$APP_NAME/launch.sh" << 'EOF'
#!/bin/bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
"$SCRIPT_DIR/HomestylePlate" &
BGPID=$!
sleep 2
xdg-open "http://localhost:8080" 2>/dev/null || chromium-browser "http://localhost:8080" 2>/dev/null || firefox "http://localhost:8080" 2>/dev/null
echo "Running at http://localhost:8080"
wait $BGPID
EOF
    chmod +x "$build_dir/$APP_NAME/launch.sh"
    
    # README
    cat > "$build_dir/$APP_NAME/README.txt" << EOF
Homestyle Plate v$VERSION (Linux $arch_label)

Quick Start:
  ./launch.sh          # Start + open browser
  ./HomestylePlate     # Start server only
  # Then visit http://localhost:8080

Self-contained binary, no dependencies needed.
EOF

    log_step "Creating tarball..."
    local tarball="$output_dir/${APP_NAME}_${VERSION}_linux-${arch_label}.tar.gz"
    tar -czf "$tarball" -C "$build_dir" "$APP_NAME"
    rm -rf "$build_dir"
    log_success "Linux tarball: $(basename "$tarball") ($(du -h "$tarball" | cut -f1))"
}

# ============================================================
# Windows ZIP Build
# ============================================================
build_windows() {
    log_info "══════════════════════════════════════════"
    log_info "  Building Windows v$VERSION (x64)"
    log_info "══════════════════════════════════════════"
    
    local build_dir="$SCRIPT_DIR/build_windows"
    local output_dir="$REPO_ROOT/releases"
    local is_cross="no"
    [[ "$(uname -s)" == "Darwin" ]] && is_cross="yes"
    
    rm -rf "$build_dir"
    mkdir -p "$output_dir" "$build_dir/$APP_NAME/Frontend"
    
    log_step "Copying frontend assets..."
    cp -R "$REPO_ROOT/Frontend/"* "$build_dir/$APP_NAME/Frontend/"
    
    log_step "Building binary..."
    if ! build_binary "windows" "amd64" "$build_dir/$APP_NAME/$BINARY_NAME" "$is_cross"; then
        return 1
    fi
    
    # Launcher
    cat > "$build_dir/$APP_NAME/launch.bat" << 'EOF'
@echo off
start "" "HomestylePlate.exe"
timeout /t 3 /nobreak >nul
start "" "http://localhost:8080"
echo Running at http://localhost:8080
HomestylePlate.exe
EOF
    
    # README
    cat > "$build_dir/$APP_NAME/README.txt" << EOF
Homestyle Plate v$VERSION (Windows x64)

Quick Start:
  launch.bat           # Start + open browser
  HomestylePlate.exe   # Start server only
  # Then visit http://localhost:8080

Self-contained binary, no dependencies needed.
EOF

    log_step "Creating ZIP..."
    local zip_name="${APP_NAME}_${VERSION}_windows-amd64.zip"
    local zip="$output_dir/$zip_name"
    
    if command -v zip &>/dev/null; then
        cd "$build_dir/$APP_NAME"
        zip -r "$output_dir/$zip_name" .
        cd "$SCRIPT_DIR"
    elif command -v 7z &>/dev/null; then
        7z a "$zip" "$build_dir/$APP_NAME"/*
    else
        log_warn "No zip/7z found - creating folder package"
        mkdir -p "$output_dir/$APP_NAME"
        cp -R "$build_dir/$APP_NAME"/* "$output_dir/$APP_NAME/"
        rm -rf "$build_dir"
        log_success "Files extracted to: $APP_NAME/"
        return 0
    fi
    
    rm -rf "$build_dir"
    log_success "Windows ZIP: $(basename "$zip") ($(du -h "$zip" | cut -f1))"
}

# ============================================================
# Release Manifest
# ============================================================
generate_manifest() {
    local manifest="$REPO_ROOT/releases/RELEASE_MANIFEST.txt"
    
    {
        echo "=================================================="
        echo "Homestyle Plate v$VERSION Release Manifest"
        echo "Generated: $(date -u '+%Y-%m-%d %H:%M:%S UTC')"
        echo "=================================================="
        echo ""
        
        for f in "$REPO_ROOT/releases"/*; do
            [[ -f "$f" ]] || continue
            [[ "$f" == *"/RELEASE_MANIFEST.txt" ]] && continue
            
            local size=$(du -h "$f" | cut -f1)
            local name=$(basename "$f")
            local md5="N/A"
            if command -v md5 &>/dev/null; then
                md5=$(md5 -q "$f" 2>/dev/null || echo "N/A")
            elif command -v md5sum &>/dev/null; then
                md5=$(md5sum "$f" | cut -d' ' -f1)
            fi
            
            echo "  File: $name"
            echo "  Size: $size"
            echo "  MD5:  $md5"
            echo ""
        done
        
        echo "Platform Requirements:"
        echo "  macOS DMG  - macOS 11.0 (Big Sur) or later"
        echo "  Linux      - x86_64 or ARM64, self-contained"
        echo "  Windows    - Windows 10 x64 or later, self-contained"
        echo ""
        echo "All platforms: SQLite local storage, auto-seeds meal data."
        echo "=================================================="
    } > "$manifest"
    
    log_success "Manifest: $(basename "$manifest")"
}

# ============================================================
# Main
# ============================================================
main() {
    log_info "══════════════════════════════════════════"
    log_info "  Homestyle Plate Release Builder"
    log_info "  Version: $VERSION | $(date '+%Y-%m-%d %H:%M:%S')"
    log_info "  Host: $(uname -s) $(uname -m)"
    log_info "  Go: $(go version 2>/dev/null | awk '{print $3}')"
    log_info "══════════════════════════════════════════"
    
    [[ -z "$VERSION" ]] && { log_error "Empty version. Check App/VERSION"; exit 1; }
    
    mkdir -p "$REPO_ROOT/releases"
    
    case "$TARGET_PLATFORM" in
        macos)
            build_macos
            ;;
        linux-amd64)
            build_linux "amd64"
            ;;
        linux-arm64)
            build_linux "arm64"
            ;;
        windows-amd64)
            build_windows
            ;;
        linux)
            log_warn "Defaulting to linux-amd64. Use linux-arm64 for ARM."
            build_linux "amd64"
            ;;
        windows)
            build_windows
            ;;
        all)
            build_macos
            echo ""
            if ! build_linux "amd64"; then
                log_warn "Linux build skipped (requires Docker on macOS)"
            fi
            echo ""
            if ! build_windows; then
                log_warn "Windows build skipped (requires Docker on macOS)"
            fi
            ;;
        *)
            log_error "Unknown platform: $TARGET_PLATFORM"
            echo "Usage: $0 [all|macos|linux-amd64|linux-arm64|windows-amd64]"
            exit 1
            ;;
    esac
    
    generate_manifest
    
    log_info "══════════════════════════════════════════"
    log_info "  Release Complete! 🎉"
    log_info "══════════════════════════════════════════"
    log_info "Artifacts in: $REPO_ROOT/releases/"
    for f in "$REPO_ROOT/releases"/*; do
        [[ -f "$f" ]] || continue
        [[ "$f" == *"/RELEASE_MANIFEST.txt" ]] && continue
        log_success "  $(basename "$f") ($(du -h "$f" | cut -f1))"
    done
}

main "$@"
