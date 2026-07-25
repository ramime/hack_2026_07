# Pitch Media

Release slideshow + screencast pipeline for AgencyPulse.

**Current pitch target:** `v0.2.0-beta` — Employee time tracking + Team Lead cockpit  
**Primary demo base:** `http://localhost:8084` (Live may lag until the next VPS deploy)

## Agent boundary

| File | Owner |
|------|--------|
| [DEV_AGENT_CONTRACT.md](DEV_AGENT_CONTRACT.md) | Pitch (instructions for Dev-Agent) |
| [manifest.yaml](manifest.yaml) | **Dev-Agent only** |
| [slides.md](slides.md) | Pitch only |
| [scenes.yaml](scenes.yaml) | Pitch only |

**Dev-Agent:** read `DEV_AGENT_CONTRACT.md`, maintain `manifest.yaml`, add `data-testid`s and `/api/health`. Do not edit the other files in this folder.

**Pitch pipeline:** owns slides/scenes/capture CLI under `cmd/pitchmedia` and `tools/screencast`.

## v0.2 screencast story

1. `/` — intro + time entry (HTMX)
2. `/teamlead` — budget heatmap (Ritter / Bosch / Porsche traffic lights)
3. `/teamlead` — open AI Time & Budget Estimator modal
4. Reset demo data for the next pitch

Slides mirror that story in [slides.md](slides.md) (`GET /slides`).

## Runbook

Prerequisites: Go, Node/npm, ffmpeg. Playwright Chromium installs on first `pitchmedia`
run into `~/.cache/ms-playwright`. Manual fallback:

```bash
cd tools/screencast && npx playwright install chromium
```

### 1. Start the v0.2 app

```bash
go run ./cmd/agencypulse
# http://localhost:8084  → health should report 0.2.0-beta.x
```

### 2. Check slides & cockpit

- Slides: [http://localhost:8084/slides](http://localhost:8084/slides)
- Team Lead: [http://localhost:8084/teamlead](http://localhost:8084/teamlead)

### 3. Validate scenes

```bash
go run ./cmd/pitchmedia -dry-run
```

### 4. Capture screencast (silent + burned-in captions)

```bash
# against local v0.2 (default)
go run ./cmd/pitchmedia

# against Live only after v0.2 is deployed there
go run ./cmd/pitchmedia -base https://trw5wo98w.ralf-metzing.de
```

Pipeline: health → reset-demo → Playwright (`scenes.yaml` + captions) → `ffmpeg` →  
`artifacts/agencypulse-<version>.mp4` (≤ 120s). Default `-tts skip`.

### 5. Review and iterate

```bash
ls -la artifacts/
# edit pitch/scenes.yaml or pitch/slides.md, then re-run
```

### 6. Optional voiceover (later)

```bash
go run ./cmd/pitchmedia -tts edge
go run ./cmd/pitchmedia -tts elevenlabs
```

### 7. Submission

Upload the MP4 to Loom/Drive; set the link in `SUBMISSION.md`.

## Useful flags

| Flag | Default | Meaning |
|------|---------|---------|
| `-base` | `http://localhost:8084` | Running app URL |
| `-out` | `artifacts` | Output directory |
| `-tts` | `skip` | `skip` \| `edge` \| `elevenlabs` |
| `-scenes` | `pitch/scenes.yaml` | Scene script |
| `-dry-run` | false | Parse/print scenes only |
| `-skip-reset` | false | Do not call reset-demo-data |
| `-max-seconds` | `120` | Trim final MP4 |

## Quick links for the Dev-Agent

1. Open [DEV_AGENT_CONTRACT.md](DEV_AGENT_CONTRACT.md)
2. Keep [manifest.yaml](manifest.yaml) in sync (`budget_heatmap` / `ai-estimator-open` shipped for v0.2)
3. Keep `data-testid`s stable on `/` and `/teamlead`
4. Keep `GET /api/health` aligned with the app version
