# Pitch Media

Release slideshow + screencast pipeline for AgencyPulse.

## Agent boundary

| File | Owner |
|------|--------|
| [DEV_AGENT_CONTRACT.md](DEV_AGENT_CONTRACT.md) | Pitch (instructions for Dev-Agent) |
| [manifest.yaml](manifest.yaml) | **Dev-Agent only** |
| [slides.md](slides.md) | Pitch only |
| [scenes.yaml](scenes.yaml) | Pitch only |

**Dev-Agent:** read `DEV_AGENT_CONTRACT.md`, maintain `manifest.yaml`, add `data-testid`s and `/api/health`. Do not edit the other files in this folder.

**Pitch pipeline:** owns slides/scenes/capture CLI under `cmd/pitchmedia` and `tools/screencast`.

## Runbook (learn the flow)

Prerequisites: Go, Node/npm, ffmpeg, Chromium via Playwright (installed on first `pitchmedia` run).

### 1. Start the app

```bash
go run ./cmd/agencypulse
# http://localhost:8084
```

### 2. Check slides

Open [http://localhost:8084/slides](http://localhost:8084/slides)  
Navigate with `←` `→` / Space. Content comes from `pitch/slides.md`.

### 3. Validate scenes (no browser)

```bash
go run ./cmd/pitchmedia -dry-run
```

### 4. Capture screencast (default: silent + burned-in captions)

```bash
go run ./cmd/pitchmedia
# equivalent: -tts skip
```

What it does:

1. `GET /api/health`
2. `POST /api/reset-demo-data`
3. Playwright walkthrough from `pitch/scenes.yaml` with caption overlay
4. `ffmpeg` → `artifacts/agencypulse-<version>.mp4` (≤ 120s)

### 5. Review and iterate

```bash
ls -la artifacts/
# tweak narration / steps in pitch/scenes.yaml, then re-run
```

### 6. Optional voiceover (later)

```bash
go run ./cmd/pitchmedia -tts edge        # requires edge-tts
go run ./cmd/pitchmedia -tts elevenlabs  # stub until API client lands
```

Captions stay burned-in either way.

### 7. Submission

Upload the MP4 to Loom/Drive manually and set the link in `SUBMISSION.md`.

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
2. Update [manifest.yaml](manifest.yaml) when features ship
3. Keep `data-testid`s stable
4. Keep `GET /api/health` in sync with the app version
