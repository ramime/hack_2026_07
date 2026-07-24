# publish-github.ps1 — Squash merge current release to public-release branch & push to public GitHub
param(
    [Parameter(Mandatory=$true, Position=0)]
    [string]$VersionTag
)

$ErrorActionPreference = "Stop"

Write-Host "==> Preparing public release for version: ${VersionTag}..." -ForegroundColor Green

# Ensure working tree is clean
$status = git status --porcelain
if ($status) {
    Write-Error "Working tree is not clean. Commit all changes to private branch first."
    exit 1
}

$currentBranch = (git branch --show-current).Trim()

# Check if public-release branch exists
$branchExists = git rev-parse --verify public-release 2>$null
if (-not $branchExists) {
    Write-Host "==> Creating orphan branch 'public-release'..." -ForegroundColor Green
    git checkout --orphan public-release
    git commit --allow-empty -m "Initial public release branch"
    git checkout $currentBranch
}

Write-Host "==> Merging '${currentBranch}' onto 'public-release' with --squash..." -ForegroundColor Green
git checkout public-release
git merge $currentBranch --squash --allow-unrelated-histories -m "Release ${VersionTag}: AgencyPulse MVP"

Write-Host "==> Setting tag ${VersionTag}..." -ForegroundColor Green
git tag -a -f "${VersionTag}" -m "Public Release ${VersionTag}"

Write-Host "==> Pushing to public GitHub remote 'github'..." -ForegroundColor Green
git push -u github public-release:main --force
git push github "${VersionTag}" --force

Write-Host "==> Returning to branch '${currentBranch}'..." -ForegroundColor Green
git checkout $currentBranch

Write-Host "==> Public GitHub release ${VersionTag} published successfully!" -ForegroundColor Green
