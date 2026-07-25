# CHANGELOG.md — Version History

All notable changes to AgencyPulse will be documented in this file.

## [v0.1.4] — 2026-07-25 (Pitch Media & UI Asset Patch)
### Added
- Maintenance and pitch media pipeline alignment release.

## [v0.1.3] — 2026-07-25 (Pitch Deck Deployment Patch)
### Added
- Deployment script support (`deploy-vps.sh` and `deploy-vps.ps1`) for copying `pitch/` directory assets to VPS (`/opt/agencypulse/pitch`).
- Updated `DEPLOYMENT.md` documentation.

## [v0.1.2] — 2026-07-25 (Slides Footer Link Patch)
### Added
- Subtle `/slides` navigation link in application footer (`layout.html`) with `data-testid="slides-link"`.
- i18n keys for slides navigation (`locales/de.json` & `locales/en.json`).

## [v0.1.1] — 2026-07-25 (Pitch Media Boundary & Health API Patch)
### Added
- Health check HTTP endpoint (`GET /api/health`) returning JSON status and SemVer version.
- Stable `data-testid` attributes across HTML templates (`layout.html`, `employee.html`) for Playwright automated pitch recording.
- Pitch contract documentation and manifest synchronization (`pitch/manifest.yaml`).
- Unit test suite for API endpoints (`cmd/agencypulse/main_test.go`).

## [v0.1.0] — 2026-07-25 (Foundation & Time Tracking Release)
### Added
- Integrated pure Go SQLite database (`modernc.org/sqlite`) with automatic migrations.
- Pitch storytelling seed data (Ritter Sport, Bosch, Porsche) with initial campaign time logs.
- Bilingual i18n support (`locales/de.json` & `locales/en.json`) with header language toggle (`DE`/`EN`).
- Dark Mode Glassmorphic Design System (`web/static/styles.css`).
- HTMX employee view for reactive time tracking entry and recent log list.
- 1-click Demo Reset endpoint (`POST /api/reset-demo-data`).
- Unit test suite for database and i18n modules.

## [v0.0.1] — 2026-07-24 (Bootstrap Release)
### Added
- Minimal Go HTTP webserver returning `"Hello World!"` with HTTP 200 OK.
- VS Code Build & Run task (`.vscode/tasks.json`) mapped to `Ctrl+Shift+B`.
- Cross-platform VPS deployment scripts (`deploy-vps.sh` & `deploy-vps.ps1`).
- Public GitHub squashed release scripts (`publish-github.sh` & `publish-github.ps1`).
- Project documentation templates and English language guidelines (`AGENTS.md`).
