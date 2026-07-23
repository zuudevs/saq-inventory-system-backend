<#
.SYNOPSIS
    Test script CRUD lengkap untuk semua endpoint SAQ Inventory System Backend.

.DESCRIPTION
    Memakai native curl (BUKAN Invoke-WebRequest/Invoke-RestMethod) supaya
    perilaku request/response tetap konsisten seperti panggilan luar
    PowerShell (mis. dari Postman/curl asli), dan tetap bisa jalan di
    Windows maupun GitHub Actions Ubuntu.

    Urutan test (dependency-aware, dibersihkan lagi di akhir):
    1. Brand            -> POST, GET all, GET by id, PUT, DELETE
    2. Location          -> POST, GET all, GET by id, PUT, DELETE
    3. Category           -> POST, GET all, GET by id, PUT
    4. Metadata Structure -> POST /categories/{id}/metadata-structure, GET, PUT (Update)
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
    pastikan server sudah nyala dan .env sudah diisi PORT-nya.

    Contoh:
        ./tests/api_test.ps1
        ./tests/api_test.ps1 -TestUpload
        ./tests/api_test.ps1 -SkipCleanup
#>

param(
    [switch] $TestUpload,
    [switch] $SkipCleanup
)

$suffix = Get-Date -Format "yyyyMMddHHmmss"

# Load helpers and configure environment
. (Join-Path $PSScriptRoot "api/helpers.ps1")

# Execute test parts in order (using dot sourcing to share scopes)
. (Join-Path $PSScriptRoot "api/brand_tests.ps1")
. (Join-Path $PSScriptRoot "api/location_tests.ps1")
. (Join-Path $PSScriptRoot "api/category_tests.ps1")
. (Join-Path $PSScriptRoot "api/metadata_structure_tests.ps1")
. (Join-Path $PSScriptRoot "api/item_tests.ps1")
. (Join-Path $PSScriptRoot "api/image_tests.ps1")
. (Join-Path $PSScriptRoot "api/cleanup_tests.ps1")

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