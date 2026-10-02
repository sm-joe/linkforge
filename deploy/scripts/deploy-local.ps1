$ErrorActionPreference = "Stop"

$RepoRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")
$ComposeFile = Join-Path $RepoRoot "deploy\compose\compose.yaml"
$ImageFile = Join-Path $RepoRoot "deployment\images.env"

Write-Host "========================================"
Write-Host "LinkForge Local Deployment"
Write-Host "========================================"

if (-not (Test-Path $ComposeFile)) {
    throw "Compose file not found: $ComposeFile"
}

if (-not (Test-Path $ImageFile)) {
    throw "Deployment image references not found: $ImageFile"
}

Write-Host ""
Write-Host "Loading immutable image references..."

Get-Content $ImageFile | ForEach-Object {
    if ($_ -match '^\s*([^#=]+)=(.+)$') {
        [System.Environment]::SetEnvironmentVariable(
            $matches[1].Trim(),
            $matches[2].Trim()
        )
    }
}

$ApiImage = [System.Environment]::GetEnvironmentVariable("LINKFORGE_API_IMAGE")
$WebImage = [System.Environment]::GetEnvironmentVariable("LINKFORGE_WEB_IMAGE")

if ([string]::IsNullOrWhiteSpace($ApiImage)) {
    throw "LINKFORGE_API_IMAGE is not defined."
}

if ([string]::IsNullOrWhiteSpace($WebImage)) {
    throw "LINKFORGE_WEB_IMAGE is not defined."
}

Write-Host ""
Write-Host "API image:"
Write-Host $ApiImage

Write-Host ""
Write-Host "Web image:"
Write-Host $WebImage

Write-Host ""
Write-Host "Pulling immutable images..."

docker pull $ApiImage
if ($LASTEXITCODE -ne 0) {
    throw "Failed to pull API image."
}

docker pull $WebImage
if ($LASTEXITCODE -ne 0) {
    throw "Failed to pull Web image."
}

Write-Host ""
Write-Host "Starting LinkForge..."

docker compose `
    -f $ComposeFile `
    up -d

if ($LASTEXITCODE -ne 0) {
    throw "Docker Compose deployment failed."
}

Write-Host ""
Write-Host "Waiting for API health..."

$healthy = $false

for ($i = 1; $i -le 30; $i++) {
    try {
        $response = Invoke-WebRequest `
            -Uri "http://localhost:8080/healthz" `
            -UseBasicParsing `
            -TimeoutSec 3

        if ($response.StatusCode -eq 200) {
            $healthy = $true
            break
        }
    }
    catch {
        Start-Sleep -Seconds 2
    }
}

if (-not $healthy) {
    docker compose -f $ComposeFile ps
    docker compose -f $ComposeFile logs --tail=100 api
    throw "API health check failed."
}

Write-Host ""
Write-Host "========================================"
Write-Host "LinkForge deployment successful"
Write-Host "========================================"

docker compose -f $ComposeFile ps