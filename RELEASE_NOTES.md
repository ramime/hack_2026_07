# RELEASE_NOTES.md — AgencyPulse v0.1.3

**Release Date:** 2026-07-25  
**Milestone:** Pitch Deck Deployment Patch (`v0.1.3`)

---

## 🎯 Release Summary

The `v0.1.3` patch release updates the deployment scripts (`deploy-vps.sh` & `deploy-vps.ps1`) to automatically deploy the `pitch/` directory assets (including `pitch/slides.md`) to the VPS environment (`/opt/agencypulse/pitch`).

### Key Highlights in v0.1.3:
- **Pitch Asset Synchronization**: Automated copying of the `pitch/` folder during deployment for `GET /slides` serving.
- **Documentation Updates**: Updated `DEPLOYMENT.md` to document the remote directory structure.
