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
        brand_gpu       = "Nvidia"
    }
} | ConvertTo-Json -Depth 5 -Compress
$itemResp = Invoke-Api -Method "POST" -Path "/items" -Body $itemBody -StepName "Create Item"
$script:itemId   = $itemResp.data.id

Invoke-Api -Method "GET" -Path "/items" -StepName "List Items"
$itemGetResp = Invoke-Api -Method "GET" -Path "/items/$itemId" -StepName "Get Item by ID"

if ($null -eq $itemGetResp.data.metadata) {
    $script:FailCount++
    $script:Failures += "[Get Item by ID] metadata null, seharusnya ikut kebawa balik"
}

$itemUpdateBody = @{ item_status = "maintenance"; notes = "Lagi diservice" } | ConvertTo-Json -Compress
Invoke-Api -Method "PUT" -Path "/items/$itemId" -Body $itemUpdateBody -StepName "Update Item"
