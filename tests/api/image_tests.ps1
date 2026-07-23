Write-Host "`n########## IMAGE ##########" -ForegroundColor Magenta

if ($TestUpload) {
    $tempDir = if ($env:TEMP) { $env:TEMP } elseif ($env:TMPDIR) { $env:TMPDIR } else { [System.IO.Path]::GetTempPath() }
    $dummyImagePath = Join-Path $tempDir "saq_dummy_upload_$suffix.png"
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
    $uploadedImagePath = "images/dummy-$suffix.png"
    Write-Host "Lewati upload file asli (pakai -TestUpload buat tes POST /images/upload)." -ForegroundColor DarkYellow
}

$imageBody = @{
    item_id    = $itemId
    image_path = $uploadedImagePath
    is_primary = $true
} | ConvertTo-Json -Compress
$imageResp = Invoke-Api -Method "POST" -Path "/images" -Body $imageBody -StepName "Create Image"
$script:imageId   = $imageResp.data.id

Invoke-Api -Method "GET" -Path "/images" -StepName "List Images"
Invoke-Api -Method "GET" -Path "/images/$imageId" -StepName "Get Image by ID"

$imageUpdateBody = @{ is_primary = $false } | ConvertTo-Json -Compress
Invoke-Api -Method "PUT" -Path "/images/$imageId" -Body $imageUpdateBody -StepName "Update Image"

$invalidImageBody = @{ image_path = "images/invalid-$suffix.png" } | ConvertTo-Json -Compress
Invoke-Api -Method "POST" -Path "/images" -Body $invalidImageBody -StepName "Create Image tanpa owner (harus gagal)" -ExpectFail
