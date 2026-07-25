# RELEASE_NOTES.md — AgencyPulse v0.7.1

**Release Date:** 2026-07-25  
**Milestone:** Table Formatting & Kiosk Improvements (`v0.7.1`)

---

## 🎯 Release Summary

Release `v0.7.1` is a patch release fixing table column width wrapping for localized German currency amounts and updating the 800x480 Kiosk header UI.

### Key Highlights
- **Table Column Formatting**: Applied global `white-space: nowrap;` for data table headers (`th`) and cells (`td`), ensuring currency values (e.g. `57.400,00 €`) and table headings stay clean on a single line.
- **Kiosk Header UI**: Updated Kiosk header title text to `AgencyPulse Kiosk` without icon, alongside language toggle (`🌐 EN`/`🌐 DE`) and theme switcher (`🌙 Dark`/`☀️ Light`).
