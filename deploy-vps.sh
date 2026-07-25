#!/usr/bin/env bash
# deploy-vps.sh — Local cross-compile and SSH deploy to VPS
set -e

VPS_TARGET="${VPS_TARGET:-ralf-metzing-vps}"
VPS_PATH="${VPS_PATH:-/opt/agencypulse}"
SERVICE_NAME="agencypulse"

echo "==> Building Linux binary (GOOS=linux GOARCH=amd64)..."
GOOS=linux GOARCH=amd64 go build -o build/agencypulse ./cmd/agencypulse

echo "==> Creating remote directory structure if needed..."
ssh "${VPS_TARGET}" "mkdir -p ${VPS_PATH}/bin ${VPS_PATH}/web ${VPS_PATH}/locales ${VPS_PATH}/pitch"

echo "==> Stopping service on VPS..."
ssh "${VPS_TARGET}" "systemctl stop ${SERVICE_NAME} 2>/dev/null || true"

echo "==> Deploying binary and assets to VPS..."
scp build/agencypulse "${VPS_TARGET}:${VPS_PATH}/bin/agencypulse"
if [ -d "web" ]; then scp -r web "${VPS_TARGET}:${VPS_PATH}/"; fi
if [ -d "locales" ]; then scp -r locales "${VPS_TARGET}:${VPS_PATH}/"; fi
# Pitch deck markdown for GET /slides (pitch/slides.md)
if [ -d "pitch" ]; then scp -r pitch "${VPS_TARGET}:${VPS_PATH}/"; fi

echo "==> Restarting service on VPS..."
ssh "${VPS_TARGET}" "systemctl start ${SERVICE_NAME} 2>/dev/null || systemctl restart ${SERVICE_NAME} 2>/dev/null || echo 'Service start skipped'"

echo "==> VPS Deployment complete!"
