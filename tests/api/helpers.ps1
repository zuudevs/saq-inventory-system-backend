$ErrorActionPreference = "Stop"

$ROOT     = Resolve-Path (Join-Path $PSScriptRoot "..\..")
$ENV_FILE = Join-Path $ROOT ".env"

function Import-DotEnvFile {
    param(
        [Parameter(Mandatory = $true)][string] $Path
    )

    if (-not (Test-Path $Path)) {
        Write-Host "Warning: .env not found at $Path" -ForegroundColor Yellow
        return
    }

    foreach ($lineRaw in Get-Content $Path) {
        $line = $lineRaw.Trim()

        if (-not $line -or $line.StartsWith("#")) {
            continue
        }

        if ($line.StartsWith("export ")) {
            $line = $line.Substring(7).Trim()
        }

        if ($line -notmatch '^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.*)\s*$') {
            continue
        }

        $key = $matches[1]
        $value = $matches[2].Trim()

        if (($value.StartsWith('"') -and $value.EndsWith('"')) -or
            ($value.StartsWith("'") -and $value.EndsWith("'"))) {
            $value = $value.Substring(1, $value.Length - 2)
        }

        [System.Environment]::SetEnvironmentVariable(
            $key,
            $value,
            [System.EnvironmentVariableTarget]::Process
        )
    }
}

function Get-NativeCurlCommand {
    if ($IsWindows) {
        $cmd = Get-Command curl.exe -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
        if ($cmd) { return $cmd.Source }
    } else {
        $cmd = Get-Command curl -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
        if ($cmd) { return $cmd.Source }
    }

    throw "Native curl binary not found. Please install curl."
}

Import-DotEnvFile -Path $ENV_FILE

$script:CurlCommand = Get-NativeCurlCommand

$PORT     = if ($env:PORT) { $env:PORT } else { "8080" }
$BASE_URL = "http://localhost:$PORT"

$script:PassCount = 0
$script:FailCount = 0
$script:Failures  = @()

function Invoke-Api {
    param(
        [Parameter(Mandatory = $true)][string] $Method,
        [Parameter(Mandatory = $true)][string] $Path,
        [string] $Body = $null,
        [string] $StepName = "",
        [switch] $ExpectFail,
        [switch] $ContinueOnError
    )

    $url = "$BASE_URL$Path"

    Write-Host ""
    Write-Host "==> [$StepName] $Method $url" -ForegroundColor Cyan
    if ($Body) {
        Write-Host $Body -ForegroundColor DarkGray
    }

    $curlArgs = @(
        "-s", "-X", $Method,
        "$url",
        "-H", "Content-Type: application/json",
        "-w", "`nHTTP_STATUS:%{http_code}"
    )

    $tempFile = $null
    if ($Body) {
        $tempFile = [System.IO.Path]::GetTempFileName()
        [System.IO.File]::WriteAllText($tempFile, $Body)
        $curlArgs += @("-d", "@$tempFile")
    }

    $raw = & $script:CurlCommand @curlArgs
    if ($tempFile) {
        Remove-Item $tempFile -ErrorAction SilentlyContinue
    }
    $lines = $raw -split "`n"

    $statusLine = $lines | Where-Object { $_ -like "HTTP_STATUS:*" }
    $statusCode = [int]($statusLine -replace "HTTP_STATUS:", "")
    $jsonText   = ($lines | Where-Object { $_ -notlike "HTTP_STATUS:*" }) -join "`n"

    Write-Host "<== HTTP $statusCode" -ForegroundColor Yellow
    Write-Host $jsonText

    $parsed = $null
    try {
        $parsed = $jsonText | ConvertFrom-Json
    } catch {
        $script:FailCount++
        $script:Failures += "[$StepName] gagal parse response JSON: $jsonText"
        if ($ContinueOnError) { return $null }
        throw "[$StepName] gagal parse response JSON: $jsonText"
    }

    $isError = ($statusCode -ge 400 -or $parsed.success -eq $false)

    if ($ExpectFail) {
        if ($isError) {
            Write-Host "OK (expected fail): $($parsed.message)" -ForegroundColor Green
            $script:PassCount++
        } else {
            Write-Host "TIDAK SESUAI HARAPAN: request ini seharusnya gagal tapi sukses" -ForegroundColor Red
            $script:FailCount++
            $script:Failures += "[$StepName] diharapkan gagal tapi malah sukses"
        }
        return $parsed
    }

    if ($isError) {
        $script:FailCount++
        $script:Failures += "[$StepName] gagal (HTTP $statusCode): $($parsed.message)"
        if ($ContinueOnError) { return $parsed }
        throw "[$StepName] gagal (HTTP $statusCode): $($parsed.message)"
    }

    $script:PassCount++
    return $parsed
}

function Invoke-ApiRaw {
    param(
        [Parameter(Mandatory = $true)][string[]] $CurlArgs,
        [Parameter(Mandatory = $true)][string] $StepName
    )

    Write-Host ""
    Write-Host "==> [$StepName] (raw curl)" -ForegroundColor Cyan

    $raw = & $script:CurlCommand @CurlArgs
    $lines = $raw -split "`n"
    $statusLine = $lines | Where-Object { $_ -like "HTTP_STATUS:*" }
    $statusCode = [int]($statusLine -replace "HTTP_STATUS:", "")
    $jsonText   = ($lines | Where-Object { $_ -notlike "HTTP_STATUS:*" }) -join "`n"

    Write-Host "<== HTTP $statusCode" -ForegroundColor Yellow
    Write-Host $jsonText

    $parsed = $null
    try { $parsed = $jsonText | ConvertFrom-Json } catch { }

    if ($statusCode -ge 400 -or ($parsed -and $parsed.success -eq $false)) {
        $script:FailCount++
        $script:Failures += "[$StepName] gagal (HTTP $statusCode)"
    } else {
        $script:PassCount++
    }

    return $parsed
}
