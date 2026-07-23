Write-Host "`n########## LOCATION ##########" -ForegroundColor Magenta

$locationBody = @{
    name        = "Ruang Server-$suffix"
    room_code   = "SRV-$suffix"
    description = "Ruang server lantai 3"
} | ConvertTo-Json -Compress
$locationResp = Invoke-Api -Method "POST" -Path "/locations" -Body $locationBody -StepName "Create Location"
$script:locationId   = $locationResp.data.id

Invoke-Api -Method "GET" -Path "/locations" -StepName "List Locations"
Invoke-Api -Method "GET" -Path "/locations/$locationId" -StepName "Get Location by ID"

$locationUpdateBody = @{ description = "Ruang server lantai 3, sayap timur" } | ConvertTo-Json -Compress
Invoke-Api -Method "PUT" -Path "/locations/$locationId" -Body $locationUpdateBody -StepName "Update Location"
