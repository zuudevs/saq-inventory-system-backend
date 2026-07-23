Write-Host "`n########## BRAND ##########" -ForegroundColor Magenta

$brandBody = @{ name = "Dell-$suffix" } | ConvertTo-Json -Compress
$brandResp = Invoke-Api -Method "POST" -Path "/brands" -Body $brandBody -StepName "Create Brand"
$script:brandId   = $brandResp.data.id

Invoke-Api -Method "GET" -Path "/brands" -StepName "List Brands"
Invoke-Api -Method "GET" -Path "/brands/$brandId" -StepName "Get Brand by ID"

$brandUpdateBody = @{ name = "Dell Technologies-$suffix" } | ConvertTo-Json -Compress
Invoke-Api -Method "PUT" -Path "/brands/$brandId" -Body $brandUpdateBody -StepName "Update Brand"

Invoke-Api -Method "GET" -Path "/brands/999999" -StepName "Get Brand by ID (not found)" -ExpectFail
