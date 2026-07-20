Write-Host "Installing tools..."

# Install Go jika belum ada
if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Go..."
    winget install Go.Go --accept-source-agreements --accept-package-agreements

    # Reload PATH dari environment terbaru
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" +
                [System.Environment]::GetEnvironmentVariable("Path","User")
}

# Tambahkan Go bin
$env:Path += ";$env:USERPROFILE\go\bin"

# Install Goose jika belum ada
if (!(Get-Command goose -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Goose..."
    go install github.com/pressly/goose/v3/cmd/goose@v3.24.1
}

# Install SQLite CLI jika belum ada
if (!(Get-Command sqlite3 -ErrorAction SilentlyContinue)) {
    Write-Host "Installing SQLite..."
    winget install SQLite.SQLite --accept-source-agreements --accept-package-agreements
}

Write-Host ""
Write-Host "Downloading Go dependencies..."
go mod download

Write-Host ""
Write-Host "Installed:"
Write-Host "----------------"

Write-Host "Go    : " -NoNewLine
go version

Write-Host "Goose : " -NoNewLine
goose -version

Write-Host "SQLite: " -NoNewLine
sqlite3 --version

Write-Host ""
Write-Host "Done"