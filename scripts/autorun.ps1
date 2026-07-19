$ROOT = Resolve-Path (Join-Path $PSScriptRoot "..")

$TOOL_DIR      = Join-Path $ROOT "tools"
$MIGRATION_DIR = Join-Path $ROOT "migrations"
$APP_DIR       = Join-Path $ROOT "cmd/server"
$ENV_FILE      = Join-Path $ROOT ".env"

# Load library
. (Join-Path $TOOL_DIR "zuu-powershell-dotenv/Import-DotEnv.ps1")

# Load .env
Import-DotEnv -Path $ENV_FILE

# MySQL DSN
$DSN = "$($env:DB_USER):$($env:DB_PASS)@tcp($($env:DB_HOST):$($env:DB_PORT))/$($env:DB_NAME)?parseTime=true"

# Run migrations
goose `
    -dir $MIGRATION_DIR `
    mysql $DSN `
    up

if ($LASTEXITCODE -ne 0) {
    throw "Migration failed."
}

# Run application
go run $APP_DIR