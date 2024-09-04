TARGET_DIR=./cmd/app/
OUTPUT_BINARY=console-api

build:
	cd $(TARGET_DIR) && go build -o $(OUTPUT_BINARY)

build-v:
	cd $(TARGET_DIR) && go build -o $(OUTPUT_BINARY) -v ./...

test:
	go test ./... -coverprofile=./coverage.out && go tool cover -func=coverage.out | grep total

test-verbose:
	go test -v ./... -coverprofile=./coverage.out
	
lint:
	golangci-lint run --config ./.golangci.yaml

fmt:
	find . -name '*.go' -exec gofmt -s -w {} \;