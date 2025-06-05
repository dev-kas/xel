VERSION ?= 0.2.1
ENGINE_VERSION ?= 1.0.0

test:
	go test ./...

# macOS builds
build-darwin-amd64:
	# Note: macOS builds with CGO_ENABLED=1 should be done on macOS machines
	# For local testing on non-macOS, we'll use CGO_ENABLED=0
	GOOS='darwin' GOARCH='amd64' CGO_ENABLED=0 go build -ldflags="-X xel/shared.RuntimeVersion=$(VERSION) -X xel/shared.EngineVersion=$(ENGINE_VERSION)" -o ./bin/xel-darwin-amd64

build-darwin-arm64:
	# Note: macOS builds with CGO_ENABLED=1 should be done on macOS machines
	# For local testing on non-macOS, we'll use CGO_ENABLED=0
	GOOS='darwin' GOARCH='arm64' CGO_ENABLED=0 go build -ldflags="-X xel/shared.RuntimeVersion=$(VERSION) -X xel/shared.EngineVersion=$(ENGINE_VERSION)" -o ./bin/xel-darwin-arm64

build-mac: build-darwin-amd64 build-darwin-arm64

# Linux builds
build-linux-amd64:
	GOOS='linux' GOARCH='amd64' CGO_ENABLED=1 go build -ldflags="-X xel/shared.RuntimeVersion=$(VERSION) -X xel/shared.EngineVersion=$(ENGINE_VERSION)" -o ./bin/xel-linux-amd64

build-linux-arm64:
	CC=aarch64-linux-gnu-gcc GOOS='linux' GOARCH='arm64' CGO_ENABLED=1 go build -ldflags="-X xel/shared.RuntimeVersion=$(VERSION) -X xel/shared.EngineVersion=$(ENGINE_VERSION)" -o ./bin/xel-linux-arm64

build-linux: build-linux-amd64 build-linux-arm64

# Windows builds
build-windows-amd64:
	CC=x86_64-w64-mingw32-gcc GOOS='windows' GOARCH='amd64' CGO_ENABLED=1 go build -ldflags="-X xel/shared.RuntimeVersion=$(VERSION) -X xel/shared.EngineVersion=$(ENGINE_VERSION)" -o ./bin/xel-windows-amd64.exe

build-windows: build-windows-amd64

# All builds
build-all: build-mac build-linux build-windows
build: build-all
