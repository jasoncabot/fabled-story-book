# Default make target is to build the wasm file
all: build

# Build the CLI binary
build:
	go build -o ./bin/ ./cmd/cli/main.go

# Build the wasm file
wasm:
	GOOS=js GOARCH=wasm tinygo build -o ./web/src/jabl.wasm -target wasm -no-debug ./cmd/wasm/main.go

# Run tests and generate coverage report
test:
	go test -v -coverprofile=coverage.out -covermode=atomic  ./...
	go tool cover -html=coverage.out -o coverage.html
