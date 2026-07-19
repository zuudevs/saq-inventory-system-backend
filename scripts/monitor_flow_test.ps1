<#
.SYNOPSIS
    Test script end-to-end: buat kategori "Monitor" -> definisikan metadata
    structure-nya -> buat item dengan metadata sesuai structure tsb.

.DESCRIPTION
    Memakai curl.exe (BUKAN Invoke-WebRequest/Invoke-RestMethod) supaya
    perilaku request/response persis seperti dipanggil dari luar PowerShell
    (mis. dari Postman/curl asli), dan supaya mudah dibaca kalau mau
    di-debug manual di terminal.

    Alur:
    1. POST /categories                              -> buat kategori Monitor
    2. POST /categories/{id}/metadata-structure       -> definisikan field metadata Monitor
    3. POST /items                                    -> buat item dgn metadata sesuai field di atas
    4. GET  /items/{id}                                -> verifikasi metadata ikut kebawa balik

.NOTES
    Jalankan dari root project (atau lewat ./scripts/monitor_flow_test.ps1),
    pastikan server sudah nyala (lihat scripts/autorun.ps1) dan .env sudah
    diisi PORT-nya.
#>

$ErrorActionPreference = "Stop"

$ROOT     = Resolve-Path (Join-Path $PSScriptRoot "..")
$TOOL_DIR = Join-Path $ROOT "tools"
$ENV_FILE = Join-Path $ROOT ".env"

# Load .env (butuh $env:PORT)
. (Join-Path $TOOL_DIR "zuu-powershell-dotenv/Import-DotEnv.ps1")
Import-DotEnv -Path $ENV_FILE

$PORT     = if ($env:PORT) { $env:PORT } else { "8080" }
$BASE_URL = "http://localhost:$PORT"

# ---------------------------------------------------------------------------
# Helper: panggil curl.exe, parse JSON response, dan stop kalau HTTP-nya gagal
# atau body-nya "success": false.
# ---------------------------------------------------------------------------
function Invoke-Api {
    param(
        [Parameter(Mandatory = $true)][string] $Method,
        [Parameter(Mandatory = $true)][string] $Path,
        [string] $Body = $null,
        [string] $StepName = ""
    )

    $url = "$BASE_URL$Path"

    Write-Host ""
    Write-Host "==> [$StepName] $Method $url" -ForegroundColor Cyan
    if ($Body) {
        Write-Host $Body -ForegroundColor DarkGray
    }

    # -s: silent, -w: sisipkan HTTP status code di baris terakhir output
    # supaya kita bisa pisahkan body vs status tanpa parsing header manual.
    $curlArgs = @(
        "-s", "-X", $Method,
        "$url",
        "-H", "Content-Type: application/json",
        "-w", "`nHTTP_STATUS:%{http_code}"
    )

    if ($Body) {
        $curlArgs += @("-d", $Body)
    }

    $raw = & curl.exe @curlArgs
    $lines = $raw -split "`n"

    $statusLine = $lines | Where-Object { $_ -like "HTTP_STATUS:*" }
    $statusCode = [int]($statusLine -replace "HTTP_STATUS:", "")
    $jsonText   = ($lines | Where-Object { $_ -notlike "HTTP_STATUS:*" }) -join "`n"

    Write-Host "<== HTTP $statusCode" -ForegroundColor Yellow
    Write-Host $jsonText

    $parsed = $null
    try {
        $parsed = $jsonText | ConvertFrom-Json
    } catch {
        throw "[$StepName] gagal parse response JSON: $jsonText"
    }

    if ($statusCode -ge 400 -or $parsed.success -eq $false) {
        throw "[$StepName] gagal (HTTP $statusCode): $($parsed.message)"
    }

    return $parsed
}

# ---------------------------------------------------------------------------
# 1. Buat kategori "Monitor"
# ---------------------------------------------------------------------------
$categoryBody = @{
    name        = "Monitor"
    description = "Monitor / layar display untuk komputer kantor"
} | ConvertTo-Json -Compress

$categoryResp = Invoke-Api -Method "POST" -Path "/categories" -Body $categoryBody -StepName "Create Category"
$categoryId   = $categoryResp.data.id

Write-Host "Category ID: $categoryId" -ForegroundColor Green

# ---------------------------------------------------------------------------
# 2. Definisikan metadata structure untuk kategori Monitor
#    Field: ukuran layar (float), resolusi (enum), refresh rate (int),
#    panel type (enum), ada speaker (bool), tanggal garansi habis (date)
# ---------------------------------------------------------------------------
$metadataStructureBody = @{
    fields = @(
        @{
            name     = "screen_size_inch"
            label    = "Ukuran Layar (inch)"
            type     = "float"
            precision = 4
            scale     = 1
            nullable = $false
        },
        @{
            name     = "resolution"
            label    = "Resolusi"
            type     = "enum"
            options  = @("1920x1080", "2560x1440", "3840x2160")
            nullable = $false
        },
        @{
            name     = "panel_type"
            label    = "Tipe Panel"
            type     = "enum"
            options  = @("IPS", "VA", "TN", "OLED")
            nullable = $false
        },
        @{
            name     = "refresh_rate_hz"
            label    = "Refresh Rate (Hz)"
            type     = "int"
            nullable = $false
        },
        @{
            name     = "has_speaker"
            label    = "Punya Speaker Built-in"
            type     = "bool"
            nullable = $false
            default  = "0"
        },
        @{
            name     = "warranty_expiry"
            label    = "Garansi Berlaku Sampai"
            type     = "date"
            nullable = $true
        }
    )
} | ConvertTo-Json -Depth 5 -Compress

$structureResp = Invoke-Api `
    -Method "POST" `
    -Path "/categories/$categoryId/metadata-structure" `
    -Body $metadataStructureBody `
    -StepName "Create Metadata Structure (Monitor)"

Write-Host "Metadata structure dibuat untuk category_id: $($structureResp.data.category_id)" -ForegroundColor Green

# ---------------------------------------------------------------------------
# 3. Buat item Monitor lengkap dengan payload metadata sesuai structure di atas
# ---------------------------------------------------------------------------
$itemBody = @{
    category_id    = $categoryId
    asset_code     = "MON-0001"
    name           = "Dell UltraSharp U2723QE"
    item_condition = "good"
    item_status    = "active"
    notes          = "Monitor meja resepsionis lantai 2"
    metadata       = @{
        screen_size_inch = 27.0
        resolution       = "3840x2160"
        panel_type       = "IPS"
        refresh_rate_hz  = 60
        has_speaker      = $false
        warranty_expiry  = "2027-07-19"
    }
} | ConvertTo-Json -Depth 5 -Compress

$itemResp = Invoke-Api -Method "POST" -Path "/items" -Body $itemBody -StepName "Create Item (Monitor)"
$itemId   = $itemResp.data.id

Write-Host "Item ID: $itemId" -ForegroundColor Green

# ---------------------------------------------------------------------------
# 4. Verifikasi: GET /items/{id} harus ikut membawa balik metadata dari DB
#    (bukan cuma echo payload seperti response POST)
# ---------------------------------------------------------------------------
$verifyResp = Invoke-Api -Method "GET" -Path "/items/$itemId" -StepName "Verify Item + Metadata"

if ($null -eq $verifyResp.data.metadata) {
    throw "GAGAL: metadata kosong/null waktu GET /items/$itemId — cek ItemService.FindByID"
}

Write-Host ""
Write-Host "=== SUMMARY ===" -ForegroundColor Magenta
Write-Host "Category ID : $categoryId"
Write-Host "Item ID     : $itemId"
Write-Host "Metadata    : $($verifyResp.data.metadata | ConvertTo-Json -Compress)"
Write-Host ""
Write-Host "OK - alur category -> metadata structure -> item + metadata sukses semua." -ForegroundColor Green
