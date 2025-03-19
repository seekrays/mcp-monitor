all: build

# Version and application information
VERSION := 1.0.0
APPNAME := mcp-monitor
BUILDDIR := ./bin

# Build for current platform
build:
	mkdir -p $(BUILDDIR)
	go build -o $(BUILDDIR)/$(APPNAME) main.go

# Run the project
run: build
	$(BUILDDIR)/$(APPNAME)

# Clean the project
clean:
	rm -rf $(BUILDDIR)

# Create output directory
.PHONY: init
init:
	mkdir -p $(BUILDDIR)

# Cross-platform build - Windows
.PHONY: build-windows-amd64
build-windows-amd64: init
	GOOS=windows GOARCH=amd64 go build -o $(BUILDDIR)/$(APPNAME)_windows_amd64.exe main.go

# Cross-platform build - macOS (Intel)
.PHONY: build-darwin-amd64
build-darwin-amd64: init
	GOOS=darwin GOARCH=amd64 go build -o $(BUILDDIR)/$(APPNAME)_darwin_amd64 main.go

# Cross-platform build - macOS (Apple Silicon)
.PHONY: build-darwin-arm64
build-darwin-arm64: init
	GOOS=darwin GOARCH=arm64 go build -o $(BUILDDIR)/$(APPNAME)_darwin_arm64 main.go

# Cross-platform build - All platforms
.PHONY: build-all
build-all: build-windows-amd64 build-darwin-amd64 build-darwin-arm64
	@echo "All platforms built successfully"
	@ls -la $(BUILDDIR)

# Help information
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make               - Build for current platform"
	@echo "  make run           - Build and run"
	@echo "  make clean         - Clean build outputs"
	@echo "  make build-windows-amd64 - Build for Windows (amd64)"
	@echo "  make build-darwin-amd64  - Build for macOS (Intel)"
	@echo "  make build-darwin-arm64  - Build for macOS (Apple Silicon)"
	@echo "  make build-all     - Build for all platforms"
	@echo "  make help          - Show this help information"
