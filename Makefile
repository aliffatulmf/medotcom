BINARY_NAME=medotcom
OUTPUT_DIR=build

build:
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME) -v