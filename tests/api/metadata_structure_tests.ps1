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

# Test the new PUT (update) endpoint
$metadataStructureUpdateBody = @{
    fields = @(
        @{ name = "ram_gb";   label = "RAM (GB)";  type = "int";   nullable = $false },
        @{ name = "storage_type"; label = "Tipe Storage"; type = "enum"; options = @("HDD", "SSD", "NVMe"); nullable = $false },
        @{ name = "brand_gpu"; label = "Brand GPU"; type = "string"; nullable = $true }
    )
} | ConvertTo-Json -Depth 5 -Compress

Invoke-Api -Method "PUT" -Path "/categories/$categoryId/metadata-structure" -Body $metadataStructureUpdateBody -StepName "Update Metadata Structure (PUT)"
