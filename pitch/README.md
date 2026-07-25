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

**Pitch pipeline:** consumes the manifest, owns slides/scenes/capture CLI (coming next).

## Quick links for the Dev-Agent

1. Open [DEV_AGENT_CONTRACT.md](DEV_AGENT_CONTRACT.md)
2. Update [manifest.yaml](manifest.yaml) when features ship
3. Add required `data-testid`s listed in the contract
4. Implement `GET /api/health` if missing
