# Pitch Media

**Current pitch target:** `v0.5.2`  
**Primary demo base:** `http://localhost:8084`

## Slides (few, domain-focused)

`pitch/slides.md` — short deck with UI screenshots from `web/static/pitch/`.

Refresh screenshots (app must be running):

```bash
cd tools/screencast && BASE_URL=http://localhost:8084 node shots.mjs
```

View: http://localhost:8084/slides

## Screencast

```bash
go run ./cmd/agencypulse
go run ./cmd/pitchmedia
# → artifacts/agencypulse-0.5.2.mp4
```

## Portal demo

| Client | URL | PIN |
|--------|-----|-----|
| Ritter Sport | `/portal/c/ritter-sport-8821` | `1234` |
