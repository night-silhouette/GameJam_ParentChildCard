$root = Get-Location

Write-Host "Scanning orphan .import files..." -ForegroundColor Cyan

Get-ChildItem -Path $root -Recurse -Filter *.import | ForEach-Object {

    $importFile = $_.FullName
    $originalFile = $importFile -replace '\.import$', ''

    if (-not (Test-Path $originalFile)) {

        Write-Host "Deleting: $importFile" -ForegroundColor Yellow
        Remove-Item $importFile -Force
    }
}

Write-Host "Done." -ForegroundColor Green