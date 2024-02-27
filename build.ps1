$ErrorActionPreference = "Stop"

$BINARY_NAME = "youchat.exe"
$OUTPUT_DIR = "build"
$MIN_GO_VERSION = "1.18.0"

try {
    # Check if Go is installed
    if (-not (Get-Command go -ErrorAction Stop)) {
        throw "Go could not be found, please install it."
    }

    # Check if the correct Go version is installed
    $CURRENT_GO_VERSION = ((go version) -replace "go version go", "") -replace " windows/amd64", ""
    if ([version]$CURRENT_GO_VERSION -le [version]$MIN_GO_VERSION) {
        throw "Go version must be greater than $MIN_GO_VERSION"
    }

    # Create output directory if not exists
    if (-not (Test-Path $OUTPUT_DIR)) {
        New-Item -ItemType Directory -Path $OUTPUT_DIR | Out-Null
    }

    # Build the application
    Write-Host "Building..."
    go build -o "$OUTPUT_DIR\$BINARY_NAME" main.go

    Write-Host "Build complete"
}
catch {
    Write-Host "Error: $_"
    exit 1
}