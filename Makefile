BINARY_NAME=you
OUTPUT_DIR=build
MIN_GO_VERSION=1.18.0

.PHONY: build clean test run

build:
 @echo "Building..."
 @if [ $(go version | cut -d " " -f 3 | tr -d 'go') \< $(MIN_GO_VERSION) ]; then \
    echo "Go version must be greater than $(MIN_GO_VERSION)"; \
    exit 1; \
 fi
 go build -o $(OUTPUT_DIR)/$(BINARY_NAME) -v

clean:
 @echo "Cleaning..."
 rm -rf $(OUTPUT_DIR)

test:
 @echo "Testing..."
 go test -v ./...

run: build
 @echo "Running..."
 ./$(OUTPUT_DIR)/$(BINARY_NAME)