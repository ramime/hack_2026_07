# Pitch Media

**Current pitch target:** `v0.5.0-beta` — Client Portal + prior features  
**Primary demo base:** `http://localhost:8084`

## Demo portal credentials (seed)

| Client | URL | PIN |
|--------|-----|-----|
| Ritter Sport | `/portal/c/ritter-sport-8821` | `1234` |
| Bosch | `/portal/c/bosch-4492` | `5678` |
| Porsche | `/portal/c/porsche-9102` | `9900` |

## v0.5 screencast story

1. `/` — time entry  
2. `/teamlead` — budget heatmap  
3. `/executive` — profitability KPIs  
4. `/tracker` — kiosk start/stop  
5. `/portal/c/ritter-sport-8821` — PIN unlock + client view  
6. Reset demo  

## Runbook

```bash
go run ./cmd/agencypulse          # health → 0.5.0-beta.x
# http://localhost:8084/slides
# http://localhost:8084/portal/c/ritter-sport-8821

go run ./cmd/pitchmedia -dry-run
go run ./cmd/pitchmedia
```

Output: `artifacts/agencypulse-0.5.0-beta.1.mp4`

## Agent boundary

| File | Owner |
|------|--------|
| [DEV_AGENT_CONTRACT.md](DEV_AGENT_CONTRACT.md) | Pitch |
| [manifest.yaml](manifest.yaml) | **Dev-Agent only** |
| [slides.md](slides.md) | Pitch |
| [scenes.yaml](scenes.yaml) | Pitch |
