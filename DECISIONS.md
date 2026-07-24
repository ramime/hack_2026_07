# DECISIONS.md — Architectural Decision Records (ADRs)

Chronological record of key technical and architectural decisions for AgencyPulse.

---

## ADR 001: Go + HTMX + SQLite Tech Stack
- **Date:** 2026-07-24
- **Status:** Accepted
- **Context:** Fast, reliable 1-day hackathon MVP development requiring zero C-compiler overhead and instant reactivity.
- **Decision:** Use Go (Golang) 1.22+ standard library `net/http`, `modernc.org/sqlite` (pure Go SQLite driver), and HTMX for frontend reactivity with Vanilla CSS.

## ADR 002: Dual Remote & Release Script Strategy
- **Date:** 2026-07-24
- **Status:** Accepted
- **Context:** Detailed development history must remain private on `git.rmhl.de`, while public GitHub should receive clean, squashed release commits per version tag.
- **Decision:** Use two remotes (`private` and `origin`). Maintain full commit tree on `private`. Provide `deploy-vps.sh/ps1` for local-to-VPS deployment and `publish-github.sh/ps1` for squashed public GitHub releases.

## ADR 003: 800x480 Raspberry Pi Touch Kiosk
- **Date:** 2026-07-24
- **Status:** Accepted
- **Context:** Agency creators switch tasks rapidly (phone calls, emails, editing).
- **Decision:** Build a dedicated `/tracker` view optimized for 800x480 touchscreen displays with touch-friendly Start/Stop cards and minute-accurate logging.

## ADR 004: Multilingual Support (i18n DE / EN)
- **Date:** 2026-07-24
- **Status:** Accepted
- **Context:** The app targets international teams, while `KONZEPT.md` remains the sole German master document.
- **Decision:** Build a Go i18n dictionary system (`locales/de.json` & `locales/en.json`) with a header language toggle. Master seed data remains in German.
