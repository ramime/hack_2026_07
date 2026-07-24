# AgencyPulse ⚡

**Real-Time Social Media Agency Budget & Time Tracking Cockpit**

AgencyPulse gives social media agencies full control over campaign budgets, team profitability, employee billing rates, and instant context-switching time tracking.

---

## 🚀 Features

- **🛡️ Budget Drift Alerts:** Heatmap dashboard highlighting campaigns exceeding 80% or 100% of their budget.
- **👑 Executive Profitability Cockpit:** Real-time net profit and margin analysis per client and employee.
- **📱 800x480 Hardware Kiosk:** Touch-optimized stopwatch timer for 7" Raspberry Pi displays and mobile devices (`/tracker`).
- **🤝 Protected Client Portal:** PIN-secured cryptic URLs (`/portal/c/<uuid>`) for clients hiding internal agency margins.
- **🚨 Security & Fraud Detection:** Live DevTeam security feed (`/dev/status`) detecting brute-force PIN attempts and link scanning.
- **🤖 AI & Automation Power:** Integrated n8n webhooks, ElevenLabs voice briefings, and fal.ai content previews.
- **🌐 Multilingual (i18n):** German & English UI support.

---

## 🛠️ Tech Stack

- **Backend:** Go (Golang 1.22+)
- **Frontend:** HTMX + Vanilla CSS (Dark Mode Design System)
- **Database:** SQLite (`modernc.org/sqlite` pure Go)
- **Deployment:** Linux VPS (Kubuntu) + Cross-platform PowerShell & Bash scripts

---

## 💻 Local Development

### Prerequisites
- Go 1.22 or higher

### Build & Run
```bash
# Run locally
go run ./cmd/agencypulse

# Build executable
go build -o build/agencypulse ./cmd/agencypulse
```

### VS Code
Press `Ctrl+Shift+B` to compile and run the local server.

---

## 📜 Documentation

- [CONCEPT.md](CONCEPT.md) — Full English project specification & roadmap.
- [KONZEPT.md](KONZEPT.md) — Master specification in German.
- [SCOPE.md](SCOPE.md) — Release milestones & scope.
- [DESIGN.md](DESIGN.md) — UI design tokens & 800x480 Kiosk specs.
- [DECISIONS.md](DECISIONS.md) — Architectural decision records (ADRs).
- [DEPLOYMENT.md](DEPLOYMENT.md) — VPS deployment instructions.
