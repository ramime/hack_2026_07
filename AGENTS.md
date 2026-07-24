# AGENTS.md — Language & Development Conventions

Guidelines for AI Coding Assistants working on AgencyPulse.

---

## 🌐 Language Conventions (Variant B — English Domain & Documentation)

- **Master Specification:** `KONZEPT.md` is the **ONLY German document** in this repository (kept for hackathon pitch guidelines).
- **Project Documentation:** All other documentation (`README.md`, `CONCEPT.md`, `SCOPE.md`, `DESIGN.md`, `DECISIONS.md`, `DEPLOYMENT.md`, `CHANGELOG.md`, `RELEASE_NOTES.md`, `TODO.md`) **MUST BE WRITTEN IN ENGLISH**.
- **Domain Language & Code:** All domain terms, variable names, struct fields, database tables/columns, functions, and inline comments **MUST BE IN ENGLISH**.
  - `Client`, `Campaign`, `Employee`, `TimeLog`, `HourlyRate`, `BillingRate`, `BudgetDrift`, `FraudLog`.
- **UI User-Facing Texts:** Handled via the i18n system (`locales/de.json` and `locales/en.json`).

---

## 🏗️ Architecture Principles

- **Backend:** Go (Golang 1.26+) using standard library `net/http` for fast, lightweight routing.
- **Frontend:** HTMX for dynamic reactivity + Vanilla CSS (Dark Mode Design System) without heavyweight JS frameworks.
- **Database:** SQLite (`modernc.org/sqlite` pure Go driver for zero-CGO C-compiler dependencies).
- **Deployment:**
  - Full git commit history pushed to private remote (`ssh://git@192.168.178.100/ralf/hack_2026_07.git`).
  - VPS deployment via `deploy-vps.sh` / `deploy-vps.ps1`.
  - Public squashed releases pushed to GitHub (`origin`) via `publish-github.sh` / `publish-github.ps1`.

---

## 🚀 Release & Versioning Policy

- SemVer tags: `v0.0.1` (Bootstrap), `v0.1.0` (Foundation & Time Tracking), `v0.2.0` (Team Lead & Alerts), `v0.3.0` (Executive View), `v0.4.0` (800x480 Kiosk), `v0.5.0` (Client Portal), `v0.6.0` (Fraud & n8n), `v1.0.0` (Pitch Release).
- Patch tags `v0.X.1` reserved strictly for bug fixes.
- Every release MUST leave a fully runnable application tested locally before tagging.
