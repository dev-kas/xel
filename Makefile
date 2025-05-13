VERSION ?= 0.2.1
ENGINE_VERSION ?= 1.0.0

test:
	go test ./...

build-mac:
	GOOS='darwin' GOARCH='amd64' go build -ldflags="-X xel/shared.RuntimeVersion=$(VERSION) -X xel/shared.EngineVersion=$(ENGINE_VERSION)" -o ./bin/xel-darwin-amd64
	GOOS='darwin' GOARCH='arm64' go build -ldflags="-X xel/shared.RuntimeVersion=$(VERSION) -X xel/shared.EngineVersion=$(ENGINE_VERSION)" -o ./bin/xel-darwin-arm64

build-linux:
	GOOS='linux' GOARCH='amd64' go build -ldflags="-X xel/shared.RuntimeVersion=$(VERSION) -X xel/shared.EngineVersion=$(ENGINE_VERSION)" -o ./bin/xel-linux-amd64
	GOOS='linux' GOARCH='arm64' go build -ldflags="-X xel/shared.RuntimeVersion=$(VERSION) -X xel/shared.EngineVersion=$(ENGINE_VERSION)" -o ./bin/xel-linux-arm64

build-windows:
	GOOS='windows' GOARCH='amd64' go build -ldflags="-X xel/shared.RuntimeVersion=$(VERSION) -X xel/shared.EngineVersion=$(ENGINE_VERSION)" -o ./bin/xel-windows-amd64.exe
	GOOS='windows' GOARCH='arm64' go build -ldflags="-X xel/shared.RuntimeVersion=$(VERSION) -X xel/shared.EngineVersion=$(ENGINE_VERSION)" -o ./bin/xel-windows-arm64.exe

build: build-mac build-linux build-windows

build-all: build
