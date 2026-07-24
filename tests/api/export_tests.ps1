Write-Host "`n########## EXPORT ##########" -ForegroundColor Magenta

$url = "$BASE_URL/exports/items"
Write-Host "==> [Export Items CSV] GET $url" -ForegroundColor Cyan

$curlArgs = @(
    "-s", "-i",
    "$url"
)

$raw = & $script:CurlCommand @curlArgs
$lines = $raw -split "\r?\n"

$statusLine = $lines | Where-Object { $_ -like "HTTP/*" } | Select-Object -First 1
$contentTypeLine = $lines | Where-Object { $_ -like "Content-Type:*" } | Select-Object -First 1
$contentDispLine = $lines | Where-Object { $_ -like "Content-Disposition:*" } | Select-Object -First 1

if ($statusLine -match "200") {
    Write-Host "HTTP Status: 200 OK" -ForegroundColor Yellow
} else {
    $script:FailCount++
    $script:Failures += "[Export Items CSV] HTTP status bukan 200: $statusLine"
}

if ($contentTypeLine -like "*text/csv*") {
    Write-Host "Header Content-Type: text/csv OK" -ForegroundColor Green
} else {
    $script:FailCount++
    $script:Failures += "[Export Items CSV] Content-Type bukan text/csv: $contentTypeLine"
}

if ($contentDispLine -like "*attachment; filename=items.csv*") {
    Write-Host "Header Content-Disposition OK" -ForegroundColor Green
} else {
    $script:FailCount++
    $script:Failures += "[Export Items CSV] Content-Disposition bukan attachment: $contentDispLine"
}

$expectedHeader = "Brand ID,ID,Category ID,Location ID,Asset Code,Name,Item Condition,Item Status,Notes,Created At,Updated At"
$rawText = $raw -join "`n"
if ($rawText -like "*$expectedHeader*") {
    Write-Host "CSV Headers OK" -ForegroundColor Green
    $script:PassCount++
} else {
    $script:FailCount++
    $script:Failures += "[Export Items CSV] Response body tidak mengandung CSV headers yang diharapkan"
}
