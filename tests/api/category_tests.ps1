Write-Host "`n########## CATEGORY ##########" -ForegroundColor Magenta

$categoryBody = @{
    name        = "Laptop-$suffix"
    description = "Kategori laptop kantor"
} | ConvertTo-Json -Compress
$categoryResp = Invoke-Api -Method "POST" -Path "/categories" -Body $categoryBody -StepName "Create Category"
$script:categoryId   = $categoryResp.data.id

Invoke-Api -Method "GET" -Path "/categories" -StepName "List Categories"
Invoke-Api -Method "GET" -Path "/categories/$categoryId" -StepName "Get Category by ID"

$categoryUpdateBody = @{ description = "Kategori laptop & notebook kantor" } | ConvertTo-Json -Compress
Invoke-Api -Method "PUT" -Path "/categories/$categoryId" -Body $categoryUpdateBody -StepName "Update Category"
