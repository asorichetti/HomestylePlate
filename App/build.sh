#!/bin/bash
# Build script for Homestyle Plate
# Creates a distributable application bundle

set -e

echo "🔨 Building Homestyle Plate..."

# Clean previous build
rm -rf build

# Detect platform
PLATFORM=$(uname -s)
BUILD_DIR=""

if [ "$PLATFORM" = "Darwin" ]; then
    APP_NAME="Homestyle Plate.app"
    BUILD_DIR="build/$APP_NAME/Contents/MacOS"
    mkdir -p "$BUILD_DIR"
else
    BUILD_DIR="build"
    mkdir -p "$BUILD_DIR"
fi

# Copy frontend files to build directory
echo "📦 Copying frontend files..."
if [ "$PLATFORM" = "Darwin" ]; then
    cp -R "$PWD/../Frontend" "build/$APP_NAME/Contents/"
else
    cp -R "$PWD/../Frontend" "$BUILD_DIR/"
fi

# Change to App directory for building
cd "$PWD"

echo "🔧 Installing dependencies..."
go mod tidy

echo "🏗️  Building Go binary..."
export CGO_ENABLED=1

# Build
go build -o "$BUILD_DIR/HomestylePlate" -ldflags="-s -w" server.go

# Platform-specific setup
if [ "$PLATFORM" = "Darwin" ]; then
    echo "🍎 Creating macOS app bundle..."
    
    # Create Info.plist
    cat > "build/$APP_NAME/Contents/Info.plist" << 'EOFPLIST'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>Homestyle Plate</string>
    <key>CFBundleExecutable</key>
    <string>HomestylePlate</string>
    <key>CFBundleIdentifier</key>
    <string>com.homestyleplate.app</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleSignature</key>
    <string>????</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.15</string>
    <key>NSHighResolutionCapable</key>
    <true/>
</dict>
</plist>
EOFPLIST
    
    # Make executable
    chmod +x "$BUILD_DIR/HomestylePlate"
    
    echo "✅ macOS app bundle created: build/$APP_NAME"
    echo ""
    echo "📦 To create a .dmg:"
    echo "   cd build && hdiutil create -volname 'Homestyle Plate' -srcfolder '$APP_NAME' -ov -format UDZO 'Homestyle Plate.dmg'"
    
else
    echo "⚙️  Setting up launcher..."
    chmod +x "$BUILD_DIR/HomestylePlate"
    
    # Create launcher script
    cat > "$BUILD_DIR/launch.sh" << 'EOFSCRIPT'
#!/bin/bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
"$SCRIPT_DIR/HomestylePlate" &
BGPID=$!
sleep 2
xdg-open "http://localhost:8080" 2>/dev/null || xdg-open "http://localhost:8080"
echo "✅ Homestyle Plate is running on http://localhost:8080"
echo "Press Ctrl+C to stop"
wait $BGPID
EOFSCRIPT
    chmod +x "$BUILD_DIR/launch.sh"
    
    echo "✅ Linux binary created: build/HomestylePlate"
    echo ""
    echo "📦 To create a tarball:"
    echo "   cd build && tar -czf HomestylePlate-linux.tar.gz *"
fi

echo ""
echo "📦 Application ready in: $(pwd)/build"
echo ""
echo "💡 The application includes:"
echo "   • Self-contained SQLite database"
echo "   • All frontend files (HTML/CSS/JS)"
echo "   • No Docker required"
echo "   • No external dependencies needed"
echo "   • Just one binary + data files"
