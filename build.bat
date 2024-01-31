@echo off
SETLOCAL

SET BINARY_NAME=you.exe
SET OUTPUT_DIR=build
SET MIN_GO_VERSION=1.18.0

REM Check if Go is installed
where /q go
IF ERRORLEVEL 1 (
    echo Go could not be found, please install it.
    EXIT /B
)

REM Check if the correct Go version is installed
FOR /F "tokens=3" %%i IN ('go version') DO (
    FOR /F "tokens=1,2 delims=." %%j IN ("%%i") DO (
        IF %%j LEQ 1 IF %%k LEQ 18 (
            echo Go version must be greater than %MIN_GO_VERSION%
            EXIT /B
        )
    )
)

REM Build the application
echo Building...
go build -o %OUTPUT_DIR%/%BINARY_NAME% main.go

echo Build complete
ENDLOCAL