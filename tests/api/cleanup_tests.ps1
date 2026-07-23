if (-not $SkipCleanup) {
    Write-Host "`n########## CLEANUP ##########" -ForegroundColor Magenta

    Invoke-Api -Method "DELETE" -Path "/images/$imageId" -StepName "Delete Image" -ContinueOnError
    Invoke-Api -Method "DELETE" -Path "/items/$itemId" -StepName "Delete Item" -ContinueOnError
    Invoke-Api -Method "DELETE" -Path "/brands/$brandId" -StepName "Delete Brand" -ContinueOnError
    Invoke-Api -Method "DELETE" -Path "/locations/$locationId" -StepName "Delete Location" -ContinueOnError
    Invoke-Api -Method "DELETE" -Path "/categories/$categoryId" -StepName "Delete Category" -ContinueOnError

    Invoke-Api -Method "GET" -Path "/items/$itemId" -StepName "Verify Item Deleted" -ExpectFail
} else {
    Write-Host "`nSkip cleanup (-SkipCleanup). Data yang dibuat:" -ForegroundColor DarkYellow
    Write-Host "  brand_id=$brandId location_id=$locationId category_id=$categoryId item_id=$itemId image_id=$imageId"
}
