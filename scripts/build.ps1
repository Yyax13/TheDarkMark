# build.ps1

$ErrorActionPreference = "Stop"

if (-not (Test-Path -Path "./build")) {
    New-Item -ItemType Directory -Path "./build" | Out-Null
}

Write-Host "Building for Windows..."
go build -ldflags="-s -w" -o "build/onTop-C2_win.exe" ./src/main.go

Write-Host "Build completed. Executable is in the build/ directory."
