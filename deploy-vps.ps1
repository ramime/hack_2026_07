# deploy-vps.ps1 — Local cross-compile and SSH deploy to Kubuntu VPS from PowerShell
$ErrorActionPreference = "Stop"

$VpsHost = $env:VPS_HOST
if (-not $VpsHost) { $VpsHost = "192.168.178.100" }
$VpsUser = $env:VPS_USER
if (-not $VpsUser) { $VpsUser = "ralf" }
$VpsPath = "/opt/agencypulse"

Write-Host "==> Building Linux binary (GOOS=linux GOARCH=amd64)..." -ForegroundColor Green
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o build/agencypulse ./cmd/agencypulse
Remove-Item Env:\GOOS
Remove-Item Env:\GOARCH

Write-Host "==> Creating remote directory structure on VPS..." -ForegroundColor Green
ssh "${VpsUser}@${VpsHost}" "mkdir -p ${VpsPath}/bin ${VpsPath}/web ${VpsPath}/locales"

Write-Host "==> Stopping service on VPS..." -ForegroundColor Green
ssh "${VpsUser}@${VpsHost}" "sudo systemctl stop agencypulse 2>/dev/null || true"

Write-Host "==> Deploying binary and assets to VPS via SCP..." -ForegroundColor Green
scp build/agencypulse "${VpsUser}@${VpsHost}:${VpsPath}/bin/agencypulse"
if (Test-Path "web") { scp -r web "${VpsUser}@${VpsHost}:${VpsPath}/" }
if (Test-Path "locales") { scp -r locales "${VpsUser}@${VpsHost}:${VpsPath}/" }

Write-Host "==> Restarting service on VPS..." -ForegroundColor Green
ssh "${VpsUser}@${VpsHost}" "sudo systemctl start agencypulse 2>/dev/null || sudo systemctl restart agencypulse 2>/dev/null || echo 'Service start skipped'"

Write-Host "==> VPS Deployment complete!" -ForegroundColor Green
