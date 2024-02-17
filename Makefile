# Default make target is to build the wasm file
all: build

# Build the CLI binary
build:
	go build -o ./bin/ ./cmd/cli/main.go

# Build the wasm file
wasm:
	GOOS=js GOARCH=wasm go build -o ./web/test.wasm ./cmd/wasm/main.go
	
test:
	go test -v ./...