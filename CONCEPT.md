# CONCEPT: AgencyPulse (Social Media Budget & Time Tracker)

**Project:** Hackathon MVP for Social Media Agencies  
**Tech Stack:** Go (Golang), HTMX, SQLite, Vanilla CSS, i18n (Multilingual DE & EN)  
**Focus:** Budget Drift Detection ("out of control"), Profitability Analysis, Fast Time Tracking, Protected Client Portal, Raspberry Pi Quick-Tracker, Fraud Detection, Pitch Demo Reset & Dual-AI Workflow.

---

## 🎯 1. Problem Statement & Goals

Social Media Agencies often manage dozens of parallel campaigns for various clients with different content formats (TikTok, Instagram Reels, Carousels, Stories).  
Common challenges:
- **Budget Drift:** Campaigns quietly exceed their allocated hours or euro budgets.
- **Lack of Transparency:** Team leads only see which team members/teams are overbooked at the end of the month.
- **Unclear Margins:** The difference between internal hourly rates (labor costs) and billing hourly rates (client rates) is rarely tracked transparently per project.
- **Client Communication:** Clients demand transparent real-time insights into campaign progress without seeing internal agency margins or individual developer/creator rates.
- **Frequent Context Switching:** Creators & managers constantly switch between tasks (15 min call, 30 min email, 2h video editing). Traditional time tracking is too cumbersome.
- **Security & Fraud:** Public client links carry risks of brute-force attempts (incorrect PINs) or third-party link scanning.
- **Pitch Reliability:** During a live demo, data gets manipulated; a 1-click reset back to the perfect demo state is essential.

**AgencyPulse** solves this through a real-time dashboard featuring traffic-light alert systems, drift warnings, role-specific views, a **Raspberry Pi Touchscreen Quick-Tracker**, **Storytelling Demo Data**, **1-Click Demo Reset**, a **DevTeam Security & Fraud Status Page**, and a **parallel Dual-AI Presentation Workflow**.

---

## 👥 2. Roles, Data Model & Pitch Demo Data

### Entities & Database Schema
1. **Clients:** Name, Industry, Monthly Budget, **`portal_token` (Cryptic URL UUID)**, **`portal_pin` (Access PIN)**.
2. **Campaigns:** Client, Name, Total Budget (€), Status (Green, Yellow, Red), **`is_favorite` (Boolean for Quick-Tracker)**.
3. **Teams & Team Leads:** Team Name, Team Lead Name.
4. **Employees:** Name, Role, Team, **Internal Hourly Cost (€/h)**, **Client Billing Rate (€/h)**.
5. **Time Logs:** Campaign, Employee, Hours Spent (minute-accurate + 15-min billing intervals), Content Type, Task Description, Timestamp.
6. **Active Timer Sessions:** Employee ID, Campaign ID, Start Time (`started_at`).
7. **Security & Fraud Logs:** Event Type (`INVALID_PIN`, `INVALID_LINK_SCAN`), Target Token, Attempt Count, Timestamp, Status (`BLOCKED`, `WARNING`).

---

### 🎭 2.1 Demo Scenarios for Live Presentation (Storytelling)
The demo data is the core of the pitch and represents 3 realistic agency campaign stages:

| Client | Campaign | Status | Budget Consumption | Pitch Storyline |
| :--- | :--- | :---: | :---: | :--- |
| **Ritter Sport** | *Summer Variety Carousel Push* | 🟢 **Green** | ~35% (€1,225 / €3,500) | **Fresh Campaign:** High margin, on schedule. Demonstrates the ideal state in the dashboard. |
| **Bosch Smart Home** | *Smart Lock Reels & Stories* | 🟡 **Yellow** | ~82% (€3,690 / €4,500) | **Warning Phase:** Campaign approaching limit. Team lead gets an early alert in the dashboard. |
| **Porsche Stuttgart** | *Taycan TikTok Launch* | 🔴 **Red** | ~108% (€6,480 / €6,000) | **"Out of Control":** Overbooked due to 3D rendering costs. Triggers n8n alert & red heatmap cards! |

---

### 🔄 2.2 Hidden 1-Click Demo Reset
- **Function:** A subtle button in the footer / dev area (`🔄 Reset Demo Data` or shortcut) calls the HTMX endpoint `POST /api/reset-demo-data`.
- **Impact:** Clears the SQLite database in milliseconds and reinstates the exact starting state of the 3 pitch scenarios. If data is altered during live testing, 1 click restores the baseline demo state instantly.

---

## 🌐 3. Multilingual Support & i18n (DE | EN)

- **Language Toggle in Header & Portals:** A switcher (**🇩🇪 DE | 🇬🇧 EN**) in the main header as well as in the client portal & kiosk dynamically updates all UI labels, table headers, buttons, and alert messages.
- **Technical Approach:** Go-based i18n Dictionary System (`locales/de.json` & `locales/en.json`). Selected language stored in session/cookie and rendered via HTMX.
- **Demo Data:** The underlying master data in the SQLite database remains in German for realism during local German agency pitches.

---

## 🖥️ 4. Six Core Views (Role Switcher, Touch Kiosk, Client Portal & DevTeam Monitor)

### 👑 Executive / Agency View (Profitability & Margen)
- **Top KPIs:** Total Revenue (€), Net Profit (€), Agency Margin (%), At-Risk Campaigns.
- **Client Profitability:** Billing Revenue vs. Actual Personnel Cost, Contribution Margin.
- **Employee Efficiency:** Billed Hours x Client Rate vs. Internal Salary Cost.

### 🛡️ Team Lead View (Budget Drift & Alerts)
- **Budget Heatmap:** Visual progress bars for every campaign (`Target` vs. `Actual`).
- **Danger Alerts ("Out of Control"):** Automatic highlighting of all campaigns reaching >80% or >100% budget usage.
- **AI Time & Budget Estimator:** Calculator for pre-estimating effort based on content deliverables (e.g., 3x TikTok Videos = ~9h effort).

### ⏱️ Employee View (Standard Fast Time Tracking)
- Minimalist form to log hours on campaigns and content types.
- Immediate HTMX update of campaign budgets without page reload.

### 📱 Hardware Kiosk & Quick-Tracker View (800x480 Raspberry Pi & Mobile Display)
- **Target Device:** Specifically optimized for 7" Raspberry Pi displays (800x480 pixels), while seamlessly responsive on smartphones, tablets & desktops (`/tracker`).
- **UI Layout:** Card grid featuring live digital stopwatch timers (`00:14:52`) & large touch-friendly **▶ START** / **⏹ STOP** buttons.
- **Strict Stop/Start Logic:** Active timers must be explicitly stopped (logging minute-accurate to SQLite) before another timer can start.

### 🤝 Client View (Protected Portal via Cryptic Link + PIN)
- **Access:** Anonymized URL path (e.g., `/portal/c/a8f9c2d1-4e7b-4a39-9c5d-88f12a34b5c6`) + 4-digit PIN gate.
- **Client Perspective:** Transparent progress & delivered content assets without revealing internal hourly rates or agency margins.
- **i18n Support:** Dedicated language switch icon (DE | EN) in the portal header.

### 🚨 DevTeam & Security Status View (Fraud Detection & Incident Monitor)
- **Real-time Security Cockpit:** Live feed of all platform security events (`/dev/status`).
- **Fraud Detection:** 3x failed PIN attempts or invalid URL scans trigger instant alerts (n8n webhooks to Slack/Telegram).

---

## ⚡ 5. Sponsor Integrations (Seamless & Practical)

| Sponsor | Function in Concept | Seamless Use Case (MVP) |
| :--- | :--- | :--- |
| **n8n** | **Budget Alert & Fraud Automation** | **1. Budget Drift:** Webhook triggers Slack alert at >100% budget.<br>**2. Security Fraud:** Fires instantly on brute-force PIN attempts or link scans. |
| **ElevenLabs** | **Executive Voice Briefing & Audio Alerts** | Button **"🔊 Play Audio Briefing"** in Executive Dashboard & Kiosk. ElevenLabs generates a concise 15-second audio executive summary. |
| **fal.ai** | **Visual Content Preview in AI Estimator** | Generates real-time AI moodboard thumbnails for estimated campaign deliverables inside the calculator. |
| **Firecrawl** | **Competitor Content Scan & Auto-Budgeting** | Scrapes social media profiles of new clients to recommend realistic monthly hourly budgets. |
| **Cursor** | **Dev Tooling / Building** | Rapid development of the Go + HTMX + SQLite MVP during the hackathon day. |

---

## 🏗️ 6. Architecture & Directory Structure

```
hack_2026_07/
├── KONZEPT.md                  # German Concept Specification
├── CONCEPT.md                  # English Concept Specification (This Document)
├── PITCH_SLIDES.md             # Pitch Presentation Deck (Generated automatically by Pitch Agent)
├── deploy-vps.sh               # Local cross-compile & SSH deploy to Kubuntu VPS (Linux)
├── deploy-vps.ps1              # Local cross-compile & SSH deploy to Kubuntu VPS (Windows PowerShell)
├── publish-github.sh           # Squash release commit & tag sync to public GitHub (Linux)
├── publish-github.ps1          # Squash release commit & tag sync to public GitHub (Windows PowerShell)
├── go.mod                      # Go Module Definition
├── main.go                     # Web server Entrypoint
├── locales/
│   ├── de.json                 # German Translations
│   └── en.json                 # English Translations
├── database/
│   └── db.go                   # SQLite Setup, Schema & Pitch Demo Data (with Reset capability)
├── templates/
│   ├── index.html              # Main Shell Layout with HTMX, Header, i18n Toggle & Demo Reset Footer Button
│   └── partials/
│       ├── executive_view.html # Executive Dashboard Partial
│       ├── teamlead_view.html  # Team Lead Dashboard Partial
│       ├── employee_view.html  # Employee Standard Time Tracking Partial
│       ├── kiosk_tracker.html  # 800x480 Touch Kiosk Quick-Tracker (Start/Stop Cards)
│       ├── client_login.html   # PIN Entry Mask for Client Portal
│       ├── client_portal.html  # Protected Client Dashboard
│       ├── dev_status.html     # DevTeam Security & Fraud Status Dashboard
│       └── ai_modal.html       # AI Estimator Partial
└── static/
    └── css/
        └── style.css           # Modern Vanilla CSS (Dark Mode & 800x480 Kiosk Responsive Breakpoints)
```

---

## 🤖 7. Dual-AI Agent Workflow (Dev Agent & Pitch Agent)

During the hackathon, two AI agents work in tandem:

```
[ Developer + Dev-AI Agent ] ──> Builds & verifies Step (e.g. v0.0.1, v0.1.0) ──> Git Commit & Git Tag
                                                                                               │
                                                                                               ▼
[ Pitch-AI Agent ] <── Reads completed code & spec from Tag v0.0.1 ── Updates PITCH_SLIDES.md
```

1. **Dev-AI Agent (IDE Agent):** Develops source code for each step (`v0.0.1`, `v0.1.0`, `v0.2.0`, etc.), launches the server, runs functional tests, and applies annotated Git tags.
2. **Pitch-AI Agent (Documentation & Presentation Agent):** As soon as a release tag is set, the Pitch Agent reads the finished code and updates the presentation slides in **`PITCH_SLIDES.md`**.
3. **Outcome:** Upon reaching version `v1.0.0`, both the working software and the complete pitch deck are 100% finished and synchronized.

---

## 🚀 8. Development Roadmap, Quality Gate & Versioning

### 📌 Quality Gate, Dual Remote & Release Scripts:
1. **Runtime Verification:** After completing each step, start the Go server (`go run main.go`), open the browser, and test the functionality live.
2. **Git Commit & Private Push:** Create a clean Git commit, set an annotated Git tag (e.g. `git tag -a v0.0.1`), and push to the private Git repository (`git push private main --tags`).
3. **VPS Deployment (`deploy-vps.sh` / `deploy-vps.ps1`):** Cross-compiles Linux binary (`GOOS=linux`), transfers files via SSH to Kubuntu VPS, and restarts `agencypulse` systemd service.
4. **Public GitHub Release (`publish-github.sh` / `publish-github.ps1`):** Accepts release tag parameter (e.g. `./publish-github.sh v0.0.1`), squash-merges `main` to `public-release` branch, sets public tag `v0.0.1`, and pushes a clean single release commit + tag to GitHub (`origin`).
5. **Patch Versions (`v0.X.1`):** Reserved exclusively for bug fixes in existing features.

---

### 🔹 Step 0 — `v0.0.1`: Minimal Webserver & Deployment Setup (Bootstrap)
- **Runnable App:** Minimal Go webserver on port 8080 returning unformatted text `"Hello World!"` with HTTP Status `200 OK`.
- **Deployment Scripts:** `deploy-vps.sh` / `deploy-vps.ps1` for VPS deployment and `publish-github.sh` / `publish-github.ps1` for squashed public GitHub releases.
- **Tooling:** `.vscode/tasks.json` for local compilation and execution (`Ctrl+Shift+B`).
- **Pitch Agent Task:** Initializes `PITCH_SLIDES.md` deck structure.

### 🔹 Step 1 — `v0.1.0`: Foundation, Storyline Demo Data, Reset Button & Standard Time Tracking
- **Runnable App:** Go web server with i18n middleware, SQLite database, pitch demo data (Ritter Sport = Green, Bosch = Yellow, Porsche = Red) & **Demo Reset Endpoint (`POST /api/reset-demo-data`)**.
- **UI:** Base Dark Mode Design System (Vanilla CSS) with header, navigation, **Language Toggle (DE | EN)**, and **discrete Demo Reset Button**.
- **Feature:** Standard Employee View with HTMX form for manual time tracking.
- **Pitch Agent Task:** Creates slides 1–3 in `PITCH_SLIDES.md` (Problem, Vision, Core Time Tracking UI).

### 🔹 Step 2 — `v0.2.0`: Team Lead Dashboard & Budget Drift Alerts
- **Runnable App:** Team Lead cockpit with budget heatmap (`Target` vs. `Actual` in %) and i18n labels.
- **Feature:** Automatic traffic light categorization (`ok`, `warning`, `danger` at >80% / >100% budget).
- **Feature:** AI Time & Budget Estimator Modal for pre-estimating content deliverables.
- **Pitch Agent Task:** Appends slides 4–5 in `PITCH_SLIDES.md` (Budget Drift "Out of Control" & AI Estimator).

### 🔹 Step 3 — `v0.3.0`: Executive Cockpit & Profitability Analysis
- **Runnable App:** Executive view with financial KPIs (Total Revenue, Labor Cost, Net Profit, Agency Margin %).
- **Feature:** Breakdown of client profitability and employee hourly rates (Internal vs. Client Billing).
- **Feature:** Audio Briefing Button (ElevenLabs Voice Summary integration).
- **Pitch Agent Task:** Appends slide 6 in `PITCH_SLIDES.md` (Profitability & ElevenLabs Voice Briefing).

### 🔹 Step 4 — `v0.4.0`: Hardware Kiosk & Quick-Tracker (800x480 Raspberry Pi Display)
- **Runnable App:** Dedicated touch view (`/tracker`) optimized for 800x480 pixels (also responsive for mobile/desktop) with i18n support.
- **Feature:** Favorite campaign cards with digital live stopwatch timers and touch-friendly **Start / Stop** buttons.
- **Feature:** Strict manual-stop logic with minute-accurate SQLite logging & 15-minute billing calculations.
- **Pitch Agent Task:** Appends slide 7 in `PITCH_SLIDES.md` (Hardware Kiosk & 800x480 Touch UX).

### 🔹 Step 5 — `v0.5.0`: Protected Client Portal (Cryptic Link + PIN)
- **Runnable App:** Anonymized/cryptic URLs for clients (e.g. `/portal/c/<uuid>`).
- **Feature:** HTMX PIN entry gate (4-digit PIN) guarding portal access.
- **Feature:** Client view showing campaign progress and delivered content assets with DE/EN switcher.
- **Pitch Agent Task:** Appends slide 8 in `PITCH_SLIDES.md` (Transparent Client Portal & Security Gate).

### 🔹 Step 6 — `v0.6.0`: DevTeam Security, Fraud Detection & n8n Integration
- **Runnable App:** Live DevTeam Security Cockpit (`/dev/status`).
- **Feature:** Fraud detection logs for 3x failed PIN attempts (brute force) & invalid portal URL scans.
- **Feature:** n8n Webhook Dispatcher (sends security alerts & budget drift notifications to Slack/Telegram).
- **Pitch Agent Task:** Appends slide 9 in `PITCH_SLIDES.md` (Security Fraud Detection & n8n Automation).

### 🏆 Step 7 — `v1.0.0`: Hackathon Pitch Release
- **Runnable App:** Overall system polish, integration of fal.ai AI thumbnails, and Firecrawl competitor scraping.
- **Outcome:** Production-ready, fully integrated live demonstration for the hackathon pitch & completed slide deck `PITCH_SLIDES.md`.

---

## 📢 9. Teaser & Matchmaking Pitch (for Discord & Networking)

Here is the ready-to-use pitch text for Discord or team matchmaking sessions:

> **⚡ AgencyPulse – "Stop Social Media Campaign Budgets From Running Out of Control!"**
>
> **The Problem:** Social Media agencies burn thousands of Euros every day because campaign budgets get quietly overbooked, hourly rates remain opaque, and creators lose track of time while context-switching between TikTok, Reels & client emails.
>
> **Our Solution:** **AgencyPulse** – an ultra-fast, visually stunning agency cockpit & hardware touch kiosk.
> - 🛡️ **Budget Heatmap & Live Alerts:** Instant warnings when campaigns hit >80% or >100% of their budget.
> - 👑 **Executive Profitability Cockpit:** Real-time net profit & margins per client & employee.
> - 📱 **800x480 Hardware Kiosk (Raspberry Pi):** Touch stopwatch with Start/Stop buttons for lightning-fast task switching at your desk.
> - 🤝 **Protected Client Portal:** PIN-secured cryptic links for clients (hiding internal agency margins).
> - 🤖 **AI & Automation Power:** ElevenLabs Voice Briefings, fal.ai Preview Generator & n8n Slack Webhooks!
>
> **Tech Stack:** Go (Golang) + HTMX + SQLite + Vanilla CSS Dark Mode (High Speed, Zero Overhead).
>
> 👥 **Looking for Teammates in:**
> - 🎨 **UI/CSS Fine-Tuning & Styling** (Glassmorphism & Touch UX)
> - 🤖 **AI & Automation Integration** (n8n Webhooks, ElevenLabs Audio, fal.ai Thumbnails)
> - 🎬 **Pitch & Storytelling** (Live demo orchestration & storytelling)
>
> *A clear step-by-step roadmap (v0.0.1 to v1.0.0) and an automated Dual-AI workflow are ready to go! Want to build a winning project with real-world impact? Join us!* 🚀
