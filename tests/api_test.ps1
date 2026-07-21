<#
.SYNOPSIS
    Test script CRUD lengkap untuk semua endpoint SAQ Inventory System Backend.

.DESCRIPTION
    Memakai native curl (BUKAN Invoke-WebRequest/Invoke-RestMethod) supaya
    perilaku request/response tetap konsisten seperti panggilan luar
    PowerShell (mis. dari Postman/curl asli), dan tetap bisa jalan di
    Windows maupun GitHub Actions Ubuntu.
    Konvensinya sama seperti scripts/monitor_flow_test.ps1 di project ini.

    Urutan test (dependency-aware, dibersihkan lagi di akhir):
    1. Brand            -> POST, GET all, GET by id, PUT, DELETE
    2. Location          -> POST, GET all, GET by id, PUT, DELETE
    3. Category           -> POST, GET all, GET by id, PUT
    4. Metadata Structure -> POST /categories/{id}/metadata-structure, GET
    5. Item (pakai category+brand+location baru + metadata)
                          -> POST, GET all, GET by id, PUT
    6. Image (nempel ke item di atas)
                          -> POST (metadata record), GET all, GET by id, PUT,
                             POST /images/upload (opsional, lihat -TestUpload)
    7. Cleanup: DELETE image -> item -> brand baru2 -> location baru2
                -> category (metadata structure ikut ke-cascade delete)

.PARAMETER TestUpload
    Kalau di-set, script akan coba upload file dummy ke POST /images/upload
    (butuh file $DummyImagePath ada / akan dibuat otomatis kalau belum ada).

.PARAMETER SkipCleanup
    Kalau di-set, data hasil test TIDAK dihapus di akhir (berguna buat cek
    manual isi database/storage setelah run).

.NOTES
    Jalankan dari root project (atau lewat ./tests/api_test.ps1),
    pastikan server sudah nyala (lihat scripts/autorun.ps1) dan .env sudah
    diisi PORT-nya.

    Contoh:
        ./tests/api_test.ps1
        ./tests/api_test.ps1 -TestUpload
        ./tests/api_test.ps1 -SkipCleanup
#>

param(
    [switch] $TestUpload,
    [switch] $SkipCleanup
)

$ErrorActionPreference = "Stop"

$ROOT     = Resolve-Path (Join-Path $PSScriptRoot "..")
$ENV_FILE = Join-Path $ROOT ".env"

function Import-DotEnvFile {
    param(
        [Parameter(Mandatory = $true)][string] $Path
    )

    if (-not (Test-Path $Path)) {
        Write-Host "Warning: .env not found at $Path" -ForegroundColor Yellow
        return
    }

    foreach ($lineRaw in Get-Content $Path) {
        $line = $lineRaw.Trim()

        if (-not $line -or $line.StartsWith("#")) {
            continue
        }

        if ($line.StartsWith("export ")) {
            $line = $line.Substring(7).Trim()
        }

        if ($line -notmatch '^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.*)\s*$') {
            continue
        }

        $key = $matches[1]
        $value = $matches[2].Trim()

        if (($value.StartsWith('"') -and $value.EndsWith('"')) -or
            ($value.StartsWith("'") -and $value.EndsWith("'"))) {
            $value = $value.Substring(1, $value.Length - 2)
        }

        [System.Environment]::SetEnvironmentVariable(
            $key,
            $value,
            [System.EnvironmentVariableTarget]::Process
        )
    }
}

function Get-NativeCurlCommand {
    if ($IsWindows) {
        $cmd = Get-Command curl.exe -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
        if ($cmd) { return $cmd.Source }
    } else {
        $cmd = Get-Command curl -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
        if ($cmd) { return $cmd.Source }
    }

    throw "Native curl binary not found. Please install curl."
}

Import-DotEnvFile -Path $ENV_FILE

$script:CurlCommand = Get-NativeCurlCommand

$PORT     = if ($env:PORT) { $env:PORT } else { "8080" }
$BASE_URL = "http://localhost:$PORT"

$script:PassCount = 0
$script:FailCount = 0
$script:Failures  = @()

# ---------------------------------------------------------------------------
# Helper: panggil native curl, parse JSON response.
# $ExpectFail = $true dipakai buat test case yang MEMANG diharapkan gagal
# (mis. GET by id yang sudah dihapus) supaya tidak menghentikan script.
# ---------------------------------------------------------------------------
function Invoke-Api {
    param(
        [Parameter(Mandatory = $true)][string] $Method,
        [Parameter(Mandatory = $true)][string] $Path,
        [string] $Body = $null,
        [string] $StepName = "",
        [switch] $ExpectFail,
        [switch] $ContinueOnError
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

    $raw = & $script:CurlCommand @curlArgs
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
        $script:FailCount++
        $script:Failures += "[$StepName] gagal parse response JSON: $jsonText"
        if ($ContinueOnError) { return $null }
        throw "[$StepName] gagal parse response JSON: $jsonText"
    }

    $isError = ($statusCode -ge 400 -or $parsed.success -eq $false)

    if ($ExpectFail) {
        if ($isError) {
            Write-Host "OK (expected fail): $($parsed.message)" -ForegroundColor Green
            $script:PassCount++
        } else {
            Write-Host "TIDAK SESUAI HARAPAN: request ini seharusnya gagal tapi sukses" -ForegroundColor Red
            $script:FailCount++
            $script:Failures += "[$StepName] diharapkan gagal tapi malah sukses"
        }
        return $parsed
    }

    if ($isError) {
        $script:FailCount++
        $script:Failures += "[$StepName] gagal (HTTP $statusCode): $($parsed.message)"
        if ($ContinueOnError) { return $parsed }
        throw "[$StepName] gagal (HTTP $statusCode): $($parsed.message)"
    }

    $script:PassCount++
    return $parsed
}

function Invoke-ApiRaw {
    # Buat request non-JSON (mis. multipart upload) langsung lewat native curl,
    # tetap hitung pass/fail dari status HTTP.
    param(
        [Parameter(Mandatory = $true)][string[]] $CurlArgs,
        [Parameter(Mandatory = $true)][string] $StepName
    )

    Write-Host ""
    Write-Host "==> [$StepName] (raw curl)" -ForegroundColor Cyan

    $raw = & $script:CurlCommand @CurlArgs
    $lines = $raw -split "`n"
    $statusLine = $lines | Where-Object { $_ -like "HTTP_STATUS:*" }
    $statusCode = [int]($statusLine -replace "HTTP_STATUS:", "")
    $jsonText   = ($lines | Where-Object { $_ -notlike "HTTP_STATUS:*" }) -join "`n"

    Write-Host "<== HTTP $statusCode" -ForegroundColor Yellow
    Write-Host $jsonText

    $parsed = $null
    try { $parsed = $jsonText | ConvertFrom-Json } catch { }

    if ($statusCode -ge 400 -or ($parsed -and $parsed.success -eq $false)) {
        $script:FailCount++
        $script:Failures += "[$StepName] gagal (HTTP $statusCode)"
    } else {
        $script:PassCount++
    }

    return $parsed
}

$suffix = Get-Date -Format "yyyyMMddHHmmss"

# ===========================================================================
# 1. BRAND CRUD
# ===========================================================================
Write-Host "`n########## BRAND ##########" -ForegroundColor Magenta

$brandBody = @{ name = "Dell-$suffix" } | ConvertTo-Json -Compress
$brandResp = Invoke-Api -Method "POST" -Path "/brands" -Body $brandBody -StepName "Create Brand"
$brandId   = $brandResp.data.id

Invoke-Api -Method "GET" -Path "/brands" -StepName "List Brands"
Invoke-Api -Method "GET" -Path "/brands/$brandId" -StepName "Get Brand by ID"

$brandUpdateBody = @{ name = "Dell Technologies-$suffix" } | ConvertTo-Json -Compress
Invoke-Api -Method "PUT" -Path "/brands/$brandId" -Body $brandUpdateBody -StepName "Update Brand"

Invoke-Api -Method "GET" -Path "/brands/999999" -StepName "Get Brand by ID (not found)" -ExpectFail

# ===========================================================================
# 2. LOCATION CRUD
# ===========================================================================
Write-Host "`n########## LOCATION ##########" -ForegroundColor Magenta

$locationBody = @{
    name        = "Ruang Server-$suffix"
    room_code   = "SRV-$suffix"
    description = "Ruang server lantai 3"
} | ConvertTo-Json -Compress
$locationResp = Invoke-Api -Method "POST" -Path "/locations" -Body $locationBody -StepName "Create Location"
$locationId   = $locationResp.data.id

Invoke-Api -Method "GET" -Path "/locations" -StepName "List Locations"
Invoke-Api -Method "GET" -Path "/locations/$locationId" -StepName "Get Location by ID"

$locationUpdateBody = @{ description = "Ruang server lantai 3, sayap timur" } | ConvertTo-Json -Compress
Invoke-Api -Method "PUT" -Path "/locations/$locationId" -Body $locationUpdateBody -StepName "Update Location"

# ===========================================================================
# 3. CATEGORY CRUD
# ===========================================================================
Write-Host "`n########## CATEGORY ##########" -ForegroundColor Magenta

$categoryBody = @{
    name        = "Laptop-$suffix"
    description = "Kategori laptop kantor"
} | ConvertTo-Json -Compress
$categoryResp = Invoke-Api -Method "POST" -Path "/categories" -Body $categoryBody -StepName "Create Category"
$categoryId   = $categoryResp.data.id

Invoke-Api -Method "GET" -Path "/categories" -StepName "List Categories"
Invoke-Api -Method "GET" -Path "/categories/$categoryId" -StepName "Get Category by ID"

$categoryUpdateBody = @{ description = "Kategori laptop & notebook kantor" } | ConvertTo-Json -Compress
Invoke-Api -Method "PUT" -Path "/categories/$categoryId" -Body $categoryUpdateBody -StepName "Update Category"

# ===========================================================================
# 4. METADATA STRUCTURE (per category)
# ===========================================================================
Write-Host "`n########## METADATA STRUCTURE ##########" -ForegroundColor Magenta

$metadataStructureBody = @{
    fields = @(
        @{ name = "ram_gb";   label = "RAM (GB)";  type = "int";   nullable = $false },
        @{ name = "storage_type"; label = "Tipe Storage"; type = "enum"; options = @("HDD", "SSD", "NVMe"); nullable = $false },
        @{ name = "warranty_expiry"; label = "Garansi Sampai"; type = "date"; nullable = $true }
    )
} | ConvertTo-Json -Depth 5 -Compress

Invoke-Api -Method "POST" -Path "/categories/$categoryId/metadata-structure" -Body $metadataStructureBody -StepName "Create Metadata Structure"
Invoke-Api -Method "GET" -Path "/categories/$categoryId/metadata-structure" -StepName "Get Metadata Structure"

# ===========================================================================
# 5. ITEM CRUD
# ===========================================================================
Write-Host "`n########## ITEM ##########" -ForegroundColor Magenta

$itemBody = @{
    brand_id       = $brandId
    category_id    = $categoryId
    location_id    = $locationId
    asset_code     = "LAP-$suffix"
    name           = "Dell Latitude 5440"
    item_condition = "good"
    item_status    = "active"
    notes          = "Laptop testing dari api_test_full.ps1"
    metadata       = @{
        ram_gb          = 16
        storage_type    = "NVMe"
        warranty_expiry = "2027-01-01"
    }
} | ConvertTo-Json -Depth 5 -Compress
$itemResp = Invoke-Api -Method "POST" -Path "/items" -Body $itemBody -StepName "Create Item"
$itemId   = $itemResp.data.id

Invoke-Api -Method "GET" -Path "/items" -StepName "List Items"
$itemGetResp = Invoke-Api -Method "GET" -Path "/items/$itemId" -StepName "Get Item by ID"

if ($null -eq $itemGetResp.data.metadata) {
    $script:FailCount++
    $script:Failures += "[Get Item by ID] metadata null, seharusnya ikut kebawa balik"
}

$itemUpdateBody = @{ item_status = "maintenance"; notes = "Lagi diservice" } | ConvertTo-Json -Compress
Invoke-Api -Method "PUT" -Path "/items/$itemId" -Body $itemUpdateBody -StepName "Update Item"

# ===========================================================================
# 6. IMAGE CRUD (+ optional upload)
# ===========================================================================
Write-Host "`n########## IMAGE ##########" -ForegroundColor Magenta

if ($TestUpload) {
    $tempDir = if ($env:TEMP) { $env:TEMP } elseif ($env:TMPDIR) { $env:TMPDIR } else { [System.IO.Path]::GetTempPath() }
    $dummyImagePath = Join-Path $tempDir "saq_dummy_upload_$suffix.png"
    # PNG 1x1 transparan, base64, biar tidak perlu file asli buat test upload.
    $pngBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="
    [IO.File]::WriteAllBytes($dummyImagePath, [Convert]::FromBase64String($pngBase64))

    $uploadArgs = @(
        "-s", "-X", "POST",
        "$BASE_URL/images/upload",
        "-F", "file=@$dummyImagePath;type=image/png",
        "-w", "`nHTTP_STATUS:%{http_code}"
    )
    $uploadResp = Invoke-ApiRaw -CurlArgs $uploadArgs -StepName "Upload Image File"
    $uploadedImagePath = $uploadResp.data.image_path

    Remove-Item $dummyImagePath -ErrorAction SilentlyContinue
} else {
    # Tanpa upload beneran, pakai path dummy langsung supaya CRUD /images
    # tetap bisa dites (image_path cuma string di DB, tidak divalidasi ada
    # filenya atau tidak di endpoint Create/Update).
    $uploadedImagePath = "images/dummy-$suffix.png"
    Write-Host "Lewati upload file asli (pakai -TestUpload buat tes POST /images/upload)." -ForegroundColor DarkYellow
}

$imageBody = @{
    item_id    = $itemId
    image_path = $uploadedImagePath
    is_primary = $true
} | ConvertTo-Json -Compress
$imageResp = Invoke-Api -Method "POST" -Path "/images" -Body $imageBody -StepName "Create Image"
$imageId   = $imageResp.data.id

Invoke-Api -Method "GET" -Path "/images" -StepName "List Images"
Invoke-Api -Method "GET" -Path "/images/$imageId" -StepName "Get Image by ID"

$imageUpdateBody = @{ is_primary = $false } | ConvertTo-Json -Compress
Invoke-Api -Method "PUT" -Path "/images/$imageId" -Body $imageUpdateBody -StepName "Update Image"

# Validasi bisnis: create image tanpa location_id maupun item_id harus ditolak
$invalidImageBody = @{ image_path = "images/invalid-$suffix.png" } | ConvertTo-Json -Compress
Invoke-Api -Method "POST" -Path "/images" -Body $invalidImageBody -StepName "Create Image tanpa owner (harus gagal)" -ExpectFail

# ===========================================================================
# 7. CLEANUP (urutan mengikuti FK: image -> item -> brand/location -> category)
# ===========================================================================
if (-not $SkipCleanup) {
    Write-Host "`n########## CLEANUP ##########" -ForegroundColor Magenta

    Invoke-Api -Method "DELETE" -Path "/images/$imageId" -StepName "Delete Image" -ContinueOnError
    Invoke-Api -Method "DELETE" -Path "/items/$itemId" -StepName "Delete Item" -ContinueOnError
    Invoke-Api -Method "DELETE" -Path "/brands/$brandId" -StepName "Delete Brand" -ContinueOnError
    Invoke-Api -Method "DELETE" -Path "/locations/$locationId" -StepName "Delete Location" -ContinueOnError
    # Metadata structure ikut ter-cascade delete waktu category dihapus.
    Invoke-Api -Method "DELETE" -Path "/categories/$categoryId" -StepName "Delete Category" -ContinueOnError

    Invoke-Api -Method "GET" -Path "/items/$itemId" -StepName "Verify Item Deleted" -ExpectFail
} else {
    Write-Host "`nSkip cleanup (-SkipCleanup). Data yang dibuat:" -ForegroundColor DarkYellow
    Write-Host "  brand_id=$brandId location_id=$locationId category_id=$categoryId item_id=$itemId image_id=$imageId"
}

# ===========================================================================
# SUMMARY
# ===========================================================================
Write-Host ""
Write-Host "=== SUMMARY ===" -ForegroundColor Magenta
Write-Host "Pass : $script:PassCount" -ForegroundColor Green
Write-Host "Fail : $script:FailCount" -ForegroundColor $(if ($script:FailCount -gt 0) { "Red" } else { "Green" })

if ($script:FailCount -gt 0) {
    Write-Host ""
    Write-Host "Detail kegagalan:" -ForegroundColor Red
    $script:Failures | ForEach-Object { Write-Host " - $_" -ForegroundColor Red }
    exit 1
}

Write-Host ""
Write-Host "OK - semua endpoint CRUD (brand, location, category, metadata-structure, item, image) lolos test." -ForegroundColor Green
exit 0