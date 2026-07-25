# RELEASE_NOTES.md — AgencyPulse v0.6.0

**Release Date:** 2026-07-25  
**Milestone:** Security, Fraud Detection & n8n Automation (`v0.6.0`)

---

## 🎯 Release Summary

Release `v0.6.0` introduces the **DevTeam Security & Fraud Cockpit** (`/dev/status`), automated brute-force protection for the Client Portal, invalid link scan auditing, and an async n8n Webhook Dispatcher for security incidents and budget drift notifications.

### Key Highlights
- **DevTeam Security Cockpit (`/dev/status`)**: Real-time system health dashboard featuring SQLite WAL mode monitoring, live security incident log stream, editable n8n Webhook Target URL input, and interactive incident simulation controls.
- **PIN Brute-Force Lockout**: Automatically tracks failed PIN attempts on client portal links; locks portal access after 3 consecutive failures for 15 minutes and records a `BLOCKED` security incident.
- **Invalid Link Scanner Logging**: Detects and logs accesses to non-existent portal tokens (`INVALID_LINK_SCAN`) to catch web crawlers or unauthorized scanners.
- **n8n Webhook Dispatcher**: Async HTTP JSON dispatcher sending alerts for security incidents and budget drift events (>100% budget) to an n8n webhook endpoint with fallback simulated logging.
- **Interactive Simulations**: 1-click simulation buttons on `/dev/status` to test PIN brute force, invalid link scan, and budget drift webhook dispatches live during pitches.
- **i18n & Versioning**: Full German/English translations and version updated to `v0.6.0-beta.1`.
