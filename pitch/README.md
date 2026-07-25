# Pitch Media

**Current pitch target:** `v0.7.0-beta` — Master Data & Client PIN overview  
**Primary demo base:** `http://localhost:8084`

## Slides

Short domain deck with screenshots → http://localhost:8084/slides

```bash
cd tools/screencast && BASE_URL=http://localhost:8084 node shots.mjs
```

## Screencast

Captions use a solid black plate for readability.

```bash
go run ./cmd/pitchmedia
# → artifacts/agencypulse-0.7.0-beta.1.mp4
```

## Key demo URLs

| View | URL |
|------|-----|
| Master data (portal PINs) | `/masterdata?tab=portal` |
| Client portal (Ritter) | `/portal/c/ritter-sport-8821` PIN `1234` |
| Team Lead heatmap | `/teamlead` |
