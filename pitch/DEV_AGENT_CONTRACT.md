# Dev-Agent Contract — AgencyPulse Pitch Media Boundary

**Audience:** the agent (or human) building the product app  
**Owner of this file:** Pitch Media pipeline — do not edit  
**Last updated:** 2026-07-25

You build and maintain AgencyPulse features. A separate Pitch Media pipeline
creates the in-app slideshow (`/slides`) and release screencasts. Stay out of
its files; communicate only through the contract below.

---

## DO NOT touch

| Path | Why |
|------|-----|
| `cmd/pitchmedia/**` | Screencast / media CLI |
| `tools/screencast/**` | Playwright recorder |
| `internal/tts/**` | Voiceover providers (later) |
| `pitch/scenes.yaml` | Screencast scenes + narration |
| `pitch/slides.md` | Rendered slide deck |
| `pitch/DEV_AGENT_CONTRACT.md` | This file |
| `artifacts/**` | Generated video/audio (gitignored) |

Product code under `cmd/agencypulse/`, `internal/db/`, `web/`, `locales/`, etc. remains yours.
If you must touch `cmd/agencypulse/main.go`, avoid rewriting the `/slides` and
`/api/health` handlers once the Pitch pipeline has added them — extend elsewhere.

---

## DO maintain: `pitch/manifest.yaml` (you own this exclusively)

Update this file on every meaningful feature merge and **before every release tag**.

### Schema

```yaml
version: "0.2.0-beta.1"          # MUST match version const in cmd/agencypulse/main.go
updated_at: "2026-07-25"         # ISO date of last manifest edit
features:
  - id: time_tracking            # stable snake_case; never rename after shipping
    title_en: "Employee Time Tracking"
    summary_en: "One-line pitch summary."
    slide_bullets_en:            # 2–4 short bullets → become /slides content
      - "Bullet one"
      - "Bullet two"
    route: "/"                   # primary UI path to demo this feature
    status: shipped              # shipped | partial | planned
    demo:
      reset_required: true
      entrypoints:               # CSS selectors; prefer data-testid
        - "[data-testid='time-log-form']"
      notes_en: "What the screencast should show."
```

### Rules

1. Only mark `shipped` / `partial` when the UI is actually reachable.
2. `planned` is fine for roadmap context but must not claim UI that does not exist.
3. Put slide ideas only in `slide_bullets_en` — **never** write `pitch/slides.md`.
4. Do **not** put Playwright step scripts here — only routes, testids, and short notes.
5. Keep English bullets factual and short (hackathon pitch, not marketing fluff).

A starter manifest already exists at `pitch/manifest.yaml`. Keep it in sync as you build.

---

## Stable UI hooks (`data-testid`)

Add these attributes to demo controls. Once shipped, do not rename/remove without
updating `manifest.yaml` → `demo.entrypoints`.

### Required now (current employee time-tracking UI)

| `data-testid` | Element |
|---------------|---------|
| `time-log-form` | Time entry `<form>` |
| `employee-select` | Employee `<select>` |
| `campaign-select` | Campaign `<select>` |
| `hours-input` | Hours `<input>` |
| `description-input` | Description field |
| `submit-time-log` | Submit button |
| `time-log-table` | Recent entries table (or wrapper) |
| `lang-toggle` | Language switch control |
| `theme-toggle` | Theme switch control |
| `reset-demo-button` | Demo reset control in footer |

### Required for v0.2+ Team Lead (`/teamlead`)

| `data-testid` | Element |
|---------------|---------|
| `budget-heatmap` | Heatmap card grid wrapper |
| `ai-estimator-open` | AI Time & Budget Estimator open button |

### Required for v0.3+ Executive (`/executive`)

| `data-testid` | Element |
|---------------|---------|
| `executive-kpis` | KPI strip wrapper |
| `audio-briefing-button` | ElevenLabs briefing play button |

### Required for v0.4 Kiosk (`/tracker`)

| `data-testid` | Element |
|---------------|---------|
| `kiosk-start-btn` | Start timer on a campaign card |
| `kiosk-stop-btn` | Stop active timer |
| `kiosk-emp-select` | Employee selector on kiosk |
| `kiosk-back` | Back to dashboard link |

### Required for v0.5 Client Portal (`/portal/c/<token>`)

| `data-testid` | Element |
|---------------|---------|
| `portal-pin-form` | PIN login form |
| `portal-pin-input` | 4-digit PIN field |
| `portal-submit-pin` | Unlock submit button |
| `portal-content` | Authenticated portal dashboard |
| `portal-asset-card` | Delivered content asset card |
| `portal-logout-btn` | Lock / logout control |

Seed demo tokens must stay stable for screencasts (e.g. Ritter `ritter-sport-8821` / PIN `1234`).

### Add when later features ship

| `data-testid` | Feature |
|---------------|---------|
| `dev-status-panel` | `/dev/status` security dashboard |

---

## Required HTTP APIs (keep stable)

| Method | Path | Status | Contract |
|--------|------|--------|----------|
| `GET` | `/api/health` | **Add if missing** | `200` + JSON `{"ok":true,"version":"<semver>"}` matching app version |
| `POST` | `/api/reset-demo-data` | **Keep** | Restores pitch seed data; callable via `curl` without a browser session; HTMX-compatible response OK |
| `POST` | `/api/set-language` | **Keep** | Accept `lang=en` or `lang=de` (form field); set cookie; English UI must work for screencasts |
| `GET` | `/` and other feature routes | **Keep** | Reachable after reset; no auth wall for local demo capture |

### Example health check

```bash
curl -s http://localhost:8084/api/health
# {"ok":true,"version":"0.2.0-beta.1"}
```

### Example demo reset (Pitch pipeline will call this before capture)

```bash
curl -s -X POST http://localhost:8084/api/reset-demo-data
```

Product APIs such as `POST /api/time-logs` stay yours. The Pitch pipeline drives them
only through UI automation, not as a hard HTTP contract.

---

## `/slides` ownership split

| Concern | Owner |
|---------|--------|
| `pitch/manifest.yaml` content | **Dev-Agent** |
| `pitch/slides.md` | **Pitch pipeline** (generated/derived from manifest) |
| `GET /slides` route + `web/templates/slides.html` | **Pitch pipeline** |
| Slide bullet ideas while coding a feature | Dev → `slide_bullets_en` in manifest |

Do not create competing slide templates or Google-Slides exporters in the app.

---

## Release checklist (Dev-Agent)

Before tagging a release:

1. [ ] Feature works locally (`go run ./cmd/agencypulse`).
2. [ ] `pitch/manifest.yaml` updated (`version` + features).
3. [ ] New UI has the required `data-testid`s; manifest `entrypoints` match.
4. [ ] `POST /api/reset-demo-data` still restores the pitch narrative seed
      (Ritter Sport / Bosch / Porsche story as applicable).
5. [ ] `GET /api/health` returns the same version string as `main.go`.
6. [ ] Do **not** run `cmd/pitchmedia` unless explicitly asked — that is the Pitch side.

---

## What Pitch needs from you (summary)

1. **File:** always-current `pitch/manifest.yaml`
2. **APIs:** `/api/health`, `/api/reset-demo-data`, `/api/set-language`
3. **Selectors:** stable `data-testid`s on demo UI
4. **Hands off:** Pitch tooling and rendered slide/scene files
