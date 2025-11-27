BINARY_NAME=bin/mop
SRC=./cmd/mop/mop.go

# Build for current platform
build: build-windows

# Build for all platforms
build-all: build-linux build-windows build-mac

# Build for Linux (amd64)
build-linux:
	set GOOS=linux&& set GOARCH=amd64&& go build -o ${BINARY_NAME}-linux-amd64 ${SRC}

# Build for Linux (arm64)
build-linux-arm:
	set GOOS=linux&& set GOARCH=arm64&& go build -o ${BINARY_NAME}-linux-arm64 ${SRC}

# Build for Windows (amd64)
build-windows:
	set GOOS=windows&& set GOARCH=amd64&& go build -o ${BINARY_NAME}-windows-amd64.exe ${SRC}

# Build for macOS (amd64)
build-mac:
	set GOOS=darwin&& set GOARCH=amd64&& go build -o ${BINARY_NAME}-darwin-amd64 ${SRC}

# Build for macOS (arm64 - Apple Silicon)
build-mac-arm:
	set GOOS=darwin&& set GOARCH=arm64&& go build -o ${BINARY_NAME}-darwin-arm64 ${SRC}

# Clean build artifacts
clean:
	del /Q ${BINARY_NAME}.exe ${BINARY_NAME}-* 2>nul || exit /b 0

.PHONY: build build-all build-linux build-linux-arm build-windows build-mac build-mac-arm clean
