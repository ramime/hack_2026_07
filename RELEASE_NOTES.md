# RELEASE_NOTES.md — AgencyPulse v0.1.1

**Release Date:** 2026-07-25  
**Milestone:** Pitch Media Boundary Patch (`v0.1.1`)

---

## 🎯 Release Summary

The `v0.1.1` release introduces stable UI testing selectors (`data-testid`), an automated health check API endpoint (`GET /api/health`), and formal alignment with the Pitch Media pipeline boundary.

### Key Highlights in v0.1.1:
- **Health Check API**: `GET /api/health` returning JSON `{"ok": true, "version": "0.1.1"}` for automated orchestration.
- **Stable UI Selectors**: Added `data-testid` attributes to all forms, inputs, buttons, and tables to support Playwright automated screencast recording.
- **Pitch Contract Alignment**: Fully aligned `pitch/manifest.yaml` with shipped application routes and entrypoints.
