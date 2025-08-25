$ErrorActionPreference = "Stop"

try {
    $goVersion = & go version
} catch {
    Write-Host "Go not found. Please install Go 1.22.2 or newer."
    exit 1
}

if ($goVersion -match "go([0-9]+\.[0-9]+\.[0-9]+)") {
    $installedVersion = $Matches[1]
    if ([version]$installedVersion -lt [version]"1.22.2") {
        Write-Host "Go version found ($installedVersion) is lower than 1.22.2. Please update Go."
        exit 1
    }
}

if (-not (Test-Path -Path "./build")) {
    New-Item -ItemType Directory -Path "./build" | Out-Null
}

Write-Host "Installing Go dependencies..."
& go mod tidy

Write-Host "Building for Windows..."
& go build -ldflags="-s -w" -o "build/TheDarmMark_win.exe" ./src/main.go

Write-Host "Build completed. Executable is in the build/ directory."
