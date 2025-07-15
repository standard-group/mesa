Write-Host "mesa client build tool" -ForegroundColor Cyan
Write-Host "select build type (1-5)" -ForegroundColor Cyan
Write-Host "1. stable"
Write-Host "2. beta"
Write-Host "3. nightly"
Write-Host "4. debug"
Write-Host "5. dev"
Write-Host ""

$choice = Read-Host "Enter your choice (1-5)"

switch ($choice) {
    "1" {
        $env:MESA_BUILD_TYPE = "release"
        Write-Host "Building Stable..." -ForegroundColor Green
    }
    "2" {
        $env:MESA_BUILD_TYPE = "beta"
        Write-Host "Building Beta..." -ForegroundColor Yellow
    }
    "3" {
        $env:MESA_BUILD_TYPE = "nightly"
        Write-Host "Building Nightly..." -ForegroundColor Magenta
    }
    "4" {
        $env:MESA_BUILD_TYPE = "debug"
        Write-Host "Building Debug..." -ForegroundColor Red
    }
    "5" {
        $env:MESA_BUILD_TYPE = "internal"
        Write-Host "Building Dev Build..." -ForegroundColor Blue
    }
    Default {
        Write-Host "Invalid choice. Exiting." -ForegroundColor Red
        exit 1
    }
}

Write-Host ""
Write-Host "MESA_BUILD_TYPE = $env:MESA_BUILD_TYPE"

$distDir = "dist"
if (!(Test-Path $distDir)) {
    New-Item -ItemType Directory -Path $distDir | Out-Null
}

$targets = @(
    @{ GOOS = "linux"; GOARCH = "amd64"; SUFFIX = "linux-amd64" },
    @{ GOOS = "linux"; GOARCH = "arm64"; SUFFIX = "linux-arm64" },
    @{ GOOS = "linux"; GOARCH = "arm"; GOARM = "7"; SUFFIX = "linux-armv7" },
    @{ GOOS = "windows"; GOARCH = "amd64"; SUFFIX = "windows-amd64.exe" },
    @{ GOOS = "windows"; GOARCH = "arm64"; SUFFIX = "windows-arm64.exe" }
)

foreach ($target in $targets) {
    $env:GOOS = $target.GOOS
    $env:GOARCH = $target.GOARCH
    if ($target.ContainsKey("GOARM")) {
        $env:GOARM = $target.GOARM
    } else {
        Remove-Item Env:\GOARM -ErrorAction SilentlyContinue
    }
    $output = "$distDir/mesa-$($target.SUFFIX)"
    Write-Host "Building $output..." -ForegroundColor Cyan
    go build -o $output ./cmd/mesa/main.go
}

Remove-Item Env:\GOOS
Remove-Item Env:\GOARCH
Remove-Item Env:\GOARM -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "SHA256 checksums:" -ForegroundColor Cyan
Get-ChildItem $distDir | ForEach-Object {
    $hash = Get-FileHash $_.FullName -Algorithm SHA256
    Write-Host "$($hash.Hash)  $($_.Name)"
}

Remove-Item Env:\MESA_BUILD_TYPE 