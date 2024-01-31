$ErrorActionPreference = "Stop"

$BINARY_NAME = "you.exe"
$OUTPUT_DIR = "build"
$MIN_GO_VERSION = "1.18.0"

# Check if Go is installed
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "Go could not be found, please install it."
    exit
}

# Check if the correct Go version is installed
$CURRENT_GO_VERSION = ((go version) -replace "go version go", "") -replace " windows/amd64", ""
if ([version]$CURRENT_GO_VERSION -le [version]$MIN_GO_VERSION) {
    Write-Host "Go version must be greater than $MIN_GO_VERSION"
    exit
}

# Build the application
Write-Host "Building..."
go build -o $OUTPUT_DIR/$BINARY_NAME main.go

Write-Host "Build complete"