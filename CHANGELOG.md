# CHANGELOG.md — Version History

All notable changes to AgencyPulse will be documented in this file.

## [v0.7.4] — 2026-07-25 (Pitch Slides SVG Branding & Chromes Release)
### Added
- Integrated SVG Vector Brand Logos (`/static/logo.svg` & `/static/logo-icon.svg`) directly into `/slides` presentation deck.
- Added custom slide logo CSS styles (`.slide-brand`, `.slides-chrome-logo`) and unit test suite `TestParseSlidesBrandLogo`.

## [v0.7.3] — 2026-07-25 (Vector SVG Logo Integration Release)
### Added
- Vector SVG App Logos (`web/static/logo-icon.svg` & `web/static/logo.svg`) featuring a modern geometric pulse-wave letter 'A' badge with indigo-purple gradient aura.
- Favicon integration across all views (`layout.html`, `kiosk.html`, `portal.html`).

## [v0.7.2] — 2026-07-25 (Full Origin Portal URL Resolution Patch Release)
### Fixed
- Portal Link input fields in `/masterdata?tab=portal` now automatically populate with complete, full origin URLs (e.g. `http://localhost:8084/portal/c/<token>` or `https://<domain>/portal/c/<token>`).
- Direct copying from input fields or clicking the copy button yields a fully clickable absolute URL.

## [v0.7.1] — 2026-07-25 (Table Currency Column Width Patch Release)
### Fixed
- Added global `white-space: nowrap;` table styles in `web/static/styles.css` preventing currency amounts (`57.400,00 €`) and table headers from wrapping onto multiple lines.
- Updated 800x480 Kiosk header title text to `AgencyPulse Kiosk` without icon.

## [v0.7.0] — 2026-07-25 (Master Data Management & Client PIN Overview Release)
### Added
- Unified Master Data Dashboard (`/masterdata`) with navigation header link `📁 Stammdaten`.
- Sub-section for **Client Views & PINs**: List of all clients with portal links (`/portal/c/{token}`), PIN codes (`1234`), quick copy button, direct portal opening, and instant PIN update form.
- Full CRUD management for **Employees (Mitarbeiter)**: Name, Role, Cost Rate (€/h), and Billing Rate (€/h).
- Full CRUD management for **Clients (Kunden)**: Client Name, Portal Token, and Security PIN Code.
- Full CRUD management for **Campaigns (Kampagnen)**: Client Assignment, Campaign Title, and Target Budget (€).
- Database helper methods (`GetAllClients`, `CreateClient`, `UpdateClient`, `DeleteClient`, `GetAllEmployees`, `CreateEmployee`, `UpdateEmployee`, `DeleteEmployee`, `GetAllCampaigns`, `CreateCampaign`, `UpdateCampaign`, `DeleteCampaign`).
- German and English localization keys in `locales/de.json` and `locales/en.json`.
- Updated `pitch/manifest.yaml` with feature entry `masterdata_management`.

## [v0.6.0] — 2026-07-25 (Security, Fraud Detection & n8n Automation Release)
### Added
- DevTeam Security Cockpit (`/dev/status`) featuring system health cards, live security audit log table, and configurable n8n Webhook Target URL input.
- SQLite `security_logs` table & database methods (`LogSecurityEvent`, `GetSecurityLogs`) for persistent event logging.
- Client Portal PIN brute-force protection: locks portal access after 3 failed attempts, logs `BLOCKED` incident, and dispatches webhook.
- Invalid link scan auditing: logs accesses to non-existent portal tokens (`INVALID_LINK_SCAN`).
- n8n Webhook Dispatcher module (`internal/n8n/webhook.go`) supporting async POST alerts with fallback simulated logging.
- Interactive attack simulation buttons on `/dev/status` for PIN brute force, invalid link scan, and campaign budget drift.
- Navigation header tab for `/dev/status` and full German & English i18n support.

## [v0.5.2] — 2026-07-25 (Number & Date Formatting Localization Release)

### Added
- Comprehensive i18n number, currency amount, percentage, and date/time formatting helpers (`FormatNumber`, `FormatAmount`, `FormatPercent`, `FormatDate`, `FormatDateTime`).
- German locale support: thousand dot separators (`1.234,56 €`), decimal commas, formatted dates (`25.07.2026 14:01`).
- English locale support: thousand comma separators (`€1,234.56`), decimal points, formatted dates (`Jul 25, 2026 14:01`).
- Updated template displays across Employee, Team Lead, Executive, 800x480 Kiosk, and Client Portal views.
- Dynamic client-side locale JS formatting for AI Estimator calculation modal.

## [v0.5.1] — 2026-07-25 (Protected Client Portal & Pitch Slides Release)
### Added
- Anonymized/cryptic client portal URLs (`/portal/c/<token>`) for Ritter Sport (`ritter-sport-8821`), Bosch (`bosch-4492`), and Porsche (`porsche-9102`).
- 4-digit PIN security authentication gate (`data-testid="portal-pin-form"`) with invalid PIN error feedback.
- Client Portal dashboard displaying transparent campaign budget usage and delivered content assets.
- Confidentiality shield hiding internal employee cost rates, margins, and individual employee names.
- Direct Client Portal launch links added to Executive Cockpit and Team Lead budget cards.
- Integrated pitch presentation slide for Client Portal in `pitch/slides.md` and slide viewer (`/slides`).
- i18n support for client portal in English and German.
- Updated `pitch/manifest.yaml` marking `client_portal` as shipped for `v0.5.1`.

## [v0.2.0] — 2026-07-25 (Team Lead Cockpit & Budget Drift Alerts Release)
### Added
- Team Lead Cockpit (`/teamlead`) with real-time budget heatmap (`Target Budget` vs `Actual Spend` in € and %).
- Top-level KPI summary cards for **Healthy** (<80%), **Warning** (80–100%), and **Danger** (>100%) campaign budget states.
- Color-coded progress bars and status badges (Green, Yellow, Red) for client campaigns (Ritter Sport, Bosch, Porsche).
- Interactive **AI Time & Budget Estimator Modal** (`data-testid="ai-estimator-open"`) for forecasting deliverable hours and budget cost.
- Navigation header links for role switching between Employee view (`/`) and Team Lead Cockpit (`/teamlead`).
- Full bilingual i18n support (`locales/de.json` & `locales/en.json`).
- Updated `pitch/manifest.yaml` with stable `data-testid` hooks (`budget-heatmap`, `ai-estimator-open`).

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
