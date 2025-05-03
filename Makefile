lint:
	golines -w . && golint ./...

test:
	go test ./...

build-mac:
	GOOS='darwin' GOARCH='arm64' go build -o ./bin/dash-darwin-arm64

build-linux:
	GOOS='linux' GOARCH='amd64' go build -o ./bin/dash-linux-amd64

build-windows:
	GOOS='windows' GOARCH='amd64' go build -o ./bin/dash-windows-amd64.exe

build: build-mac build-linux build-windows
