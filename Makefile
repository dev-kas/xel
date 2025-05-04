test:
	go test ./...

build-mac:
	GOOS='darwin' GOARCH='arm64' go build -o ./bin/xel-darwin-arm64

build-linux:
	GOOS='linux' GOARCH='amd64' go build -o ./bin/xel-linux-amd64

build-windows:
	GOOS='windows' GOARCH='amd64' go build -o ./bin/xel-windows-amd64.exe

build: build-mac build-linux build-windows
