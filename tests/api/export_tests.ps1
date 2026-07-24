Write-Host "`n########## EXPORT ##########" -ForegroundColor Magenta

# 1. Export CSV (ZIP)
$urlCsv = "$BASE_URL/exports/csv"
Write-Host "==> [Export CSV] GET $urlCsv" -ForegroundColor Cyan

$curlArgsCsv = @(
    "-s", "-i",
    "$urlCsv"
)

$rawCsv = & $script:CurlCommand @curlArgsCsv
$linesCsv = $rawCsv -split "\r?\n"

$statusLineCsv = $linesCsv | Where-Object { $_ -like "HTTP/*" } | Select-Object -First 1
$contentTypeLineCsv = $linesCsv | Where-Object { $_ -like "Content-Type:*" } | Select-Object -First 1
$contentDispLineCsv = $linesCsv | Where-Object { $_ -like "Content-Disposition:*" } | Select-Object -First 1

if ($statusLineCsv -match "200") {
    Write-Host "HTTP Status: 200 OK" -ForegroundColor Yellow
} else {
    $script:FailCount++
    $script:Failures += "[Export CSV] HTTP status bukan 200: $statusLineCsv"
}

if ($contentTypeLineCsv -like "*application/zip*") {
    Write-Host "Header Content-Type: application/zip OK" -ForegroundColor Green
} else {
    $script:FailCount++
    $script:Failures += "[Export CSV] Content-Type bukan application/zip: $contentTypeLineCsv"
}

if ($contentDispLineCsv -like "*attachment; filename=exports.zip*") {
    Write-Host "Header Content-Disposition OK" -ForegroundColor Green
    $script:PassCount++
} else {
    $script:FailCount++
    $script:Failures += "[Export CSV] Content-Disposition bukan attachment; filename=exports.zip: $contentDispLineCsv"
}

# 2. Export XLSX
$urlXlsx = "$BASE_URL/exports/xlsx"
Write-Host "==> [Export XLSX] GET $urlXlsx" -ForegroundColor Cyan

$curlArgsXlsx = @(
    "-s", "-i",
    "$urlXlsx"
)

$rawXlsx = & $script:CurlCommand @curlArgsXlsx
$linesXlsx = $rawXlsx -split "\r?\n"

$statusLineXlsx = $linesXlsx | Where-Object { $_ -like "HTTP/*" } | Select-Object -First 1
$contentTypeLineXlsx = $linesXlsx | Where-Object { $_ -like "Content-Type:*" } | Select-Object -First 1
$contentDispLineXlsx = $linesXlsx | Where-Object { $_ -like "Content-Disposition:*" } | Select-Object -First 1

if ($statusLineXlsx -match "200") {
    Write-Host "HTTP Status: 200 OK" -ForegroundColor Yellow
} else {
    $script:FailCount++
    $script:Failures += "[Export XLSX] HTTP status bukan 200: $statusLineXlsx"
}

if ($contentTypeLineXlsx -like "*application/vnd.openxmlformats-officedocument.spreadsheetml.sheet*") {
    Write-Host "Header Content-Type XLSX OK" -ForegroundColor Green
} else {
    $script:FailCount++
    $script:Failures += "[Export XLSX] Content-Type bukan spreadsheetml.sheet: $contentTypeLineXlsx"
}

if ($contentDispLineXlsx -like "*attachment; filename=exports.xlsx*") {
    Write-Host "Header Content-Disposition OK" -ForegroundColor Green
    $script:PassCount++
} else {
    $script:FailCount++
    $script:Failures += "[Export XLSX] Content-Disposition bukan attachment; filename=exports.xlsx: $contentDispLineXlsx"
}

