# SCOPE.md — AgencyPulse MVP Scope

Project milestones and feature scope per release step.

---

## 🎯 In-Scope (MVP Release Milestones)

| Version | Milestone | Key Features |
| :--- | :--- | :--- |
| `v0.0.1` | **Bootstrap** | Minimal Go HTTP webserver returning `"Hello World!"` (Status 200), `deploy-vps.sh/ps1`, `publish-github.sh/ps1`, VS Code tasks. |
| `v0.1.0` | **Foundation** | Go webserver, SQLite database, storytelling seed data (Ritter Sport, Bosch, Porsche), i18n middleware, Dark Mode UI, HTMX employee time tracking form, 1-click Demo Reset endpoint. |
| `v0.2.0` | **Team Lead Cockpit** | Budget heatmap, traffic-light alerts (>80% warning / >100% danger), AI Time & Budget Estimator modal. |
| `v0.3.0` | **Executive View** | Profitability KPIs (Revenue, Labor Cost, Profit, Margins %), billing rate breakdowns, ElevenLabs Audio Briefing. |
| `v0.4.0` | **Hardware Kiosk** | Touch-optimized 800x480 Kiosk view (`/tracker`) with favorite campaign cards, live digital stopwatches, strict Start/Stop logic. |
| `v0.5.0` | **Client Portal** | Cryptic URLs (`/portal/c/<uuid>`) + 4-digit PIN authentication gate, transparent content asset view. |
| `v0.6.0` | **Security & Fraud** | Live DevTeam status monitor (`/dev/status`), PIN brute-force detection, invalid link scanner, n8n webhook alerts. |
| `v1.0.0` | **Pitch Release** | Final polishing, fal.ai thumbnail integration, Firecrawl competitor scraping, completed slide deck `PITCH_SLIDES.md`. |

---

## ⛔ Out-of-Scope (Future / Post-Hackathon)

- Multi-tenant cloud database with complex user authentication / OAuth2.
- Invoice PDF generation & direct payment gateway integration.
- Full mobile native apps (iOS / Android).
