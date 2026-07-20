$ROOT = Resolve-Path (Join-Path $PSScriptRoot "..")

$TOOL_DIR      = Join-Path $ROOT "tools"
$MIGRATION_DIR = Join-Path $ROOT "migrations"
$APP_DIR       = Join-Path $ROOT "cmd/server"
$ENV_FILE      = Join-Path $ROOT ".env"

# Load library
. (Join-Path $TOOL_DIR "zuu-powershell-dotenv/Import-DotEnv.ps1")

# Load .env
Import-DotEnv -Path $ENV_FILE

$DB_PATH = Join-Path $ROOT $env:DB_PATH
$DB_DIR = Split-Path $DB_PATH -Parent

if (-not (Test-Path $DB_DIR)) {
    New-Item -ItemType Directory -Path $DB_DIR -Force | Out-Null
}

# Run migrations
goose `
    -dir $MIGRATION_DIR `
    sqlite $DB_PATH `
    up

if ($LASTEXITCODE -ne 0) {
    throw "Migration failed."
}

# Run application
go run $APP_DIR