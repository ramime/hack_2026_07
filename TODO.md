# TODO.md — AgencyPulse Milestone Checklist

---

## 🟢 v0.0.1 — Bootstrap (COMPLETED)
- [x] Minimal Go webserver (`cmd/agencypulse/main.go`) returning `"Hello World!"`.
- [x] `.vscode/tasks.json` configured for `Ctrl+Shift+B`.
- [x] `deploy-vps.sh` & `deploy-vps.ps1` scripts created.
- [x] `publish-github.sh` & `publish-github.ps1` scripts created.

---

## 🟢 v0.1.0 — Foundation & Time Tracking (COMPLETED)
- [x] SQLite database setup (`modernc.org/sqlite`).
- [x] Pitch storytelling seed data (Ritter Sport = Green, Bosch = Yellow, Porsche = Red).
- [x] i18n middleware and dictionary loader (`locales/de.json` & `locales/en.json`).
- [x] Dark mode UI shell & header with Role Switcher & Language Toggle.
- [x] Employee view with HTMX time tracking form.
- [x] 1-click Demo Reset endpoint (`POST /api/reset-demo-data`).

---

## ⏳ v0.2.0 — Team Lead Cockpit & Alerts
- [ ] Budget heatmap (`Target` vs. `Actual`).
- [ ] Traffic-light alert cards (>80% / >100%).
- [ ] AI Time & Budget Estimator modal.

---

## ⏳ v0.3.0 — Executive View & Profitability
- [ ] Financial KPIs (Revenue, Cost, Net Profit, Margins %).
- [ ] ElevenLabs Audio Briefing integration.

---

## ⏳ v0.4.0 — Hardware Kiosk (800x480)
- [ ] Touch-optimized `/tracker` view with stopwatch cards and Start/Stop buttons.

---

## ⏳ v0.5.0 — Protected Client Portal
- [ ] Cryptic URL routing & 4-digit PIN authentication gate.

---

## ⏳ v0.6.0 — Security & Fraud Detection
- [ ] DevTeam live security dashboard (`/dev/status`) & n8n webhooks.

---

## 🏆 v1.0.0 — Pitch Release
- [ ] fal.ai thumbnails & Firecrawl competitor scraping.
