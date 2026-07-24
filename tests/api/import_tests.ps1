Write-Host "`n########## IMPORT ##########" -ForegroundColor Magenta

# 1. First download exported XLSX to temp file
$exportUrl = "$BASE_URL/exports/xlsx"
$tempFile = [System.IO.Path]::GetTempFileName()
$tmpXlsx = "$tempFile.xlsx"
Remove-Item -Force $tempFile -ErrorAction SilentlyContinue

Write-Host "==> [Export XLSX for Import Test] GET $exportUrl -> $tmpXlsx" -ForegroundColor Cyan

$curlArgsExport = @(
    "-s",
    "-o", "$tmpXlsx",
    "$exportUrl"
)
& $script:CurlCommand @curlArgsExport

if (-not (Test-Path $tmpXlsx)) {
    $script:FailCount++
    $script:Failures += "[Import XLSX] Failed to download temporary XLSX file for import test"
} else {
    Write-Host "Downloaded exported workbook successfully" -ForegroundColor Green

    # 2. Upload to POST /imports/xlsx
    $importUrl = "$BASE_URL/imports/xlsx"
    Write-Host "==> [Import XLSX] POST $importUrl" -ForegroundColor Cyan

    $curlArgsImport = @(
        "-s", "-i",
        "-F", "file=@$tmpXlsx",
        "$importUrl"
    )

    $rawImport = & $script:CurlCommand @curlArgsImport
    $linesImport = $rawImport -split "\r?\n"

    $statusLineImport = $linesImport | Where-Object { $_ -like "HTTP/*" } | Select-Object -First 1

    if ($statusLineImport -match "200") {
        Write-Host "HTTP Status: 200 OK" -ForegroundColor Yellow
        $script:PassCount++
    } else {
        $script:FailCount++
        $script:Failures += "[Import XLSX] HTTP status bukan 200: $statusLineImport"
    }

    # Clean up temp file
    Remove-Item -Force $tmpXlsx -ErrorAction SilentlyContinue
}
