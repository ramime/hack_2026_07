# RELEASE_NOTES.md — AgencyPulse v0.1.0

**Release Date:** 2026-07-25  
**Milestone:** Foundation & Time Tracking (`v0.1.0`)

---

## 🎯 Release Summary

The `v0.1.0` release establishes the database foundation, multi-language support (i18n), pitch storytelling seed data, glassmorphic dark mode styling, and HTMX reactive time tracking.

### Key Highlights in v0.1.0:
- **Pure Go SQLite Integration**: `modernc.org/sqlite` database layer with schema migrations and seed data initialization.
- **Pitch Storytelling Seed Data**: Pre-loaded clients and campaigns (Ritter Sport = Green, Bosch = Yellow, Porsche = Red) with historical time logs.
- **Dark Mode Glassmorphic UI**: High-contrast, sleek design system (`web/static/styles.css`) using custom CSS tokens.
- **Reactive HTMX Time Entry**: Form for logging hours per employee & campaign with dynamic table updates without page reloads.
- **Language Switcher (i18n)**: Seamless German (`DE`) & English (`EN`) UI translation toggling.
- **1-Click Demo Reset**: Endpoint (`POST /api/reset-demo-data`) to instantly restore seed state during live pitch presentations.
