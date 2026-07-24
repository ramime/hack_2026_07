# KONZEPT: AgencyPulse (Social Media Budget & Time Tracker)

**Projekt:** Hackathon MVP für Social Media Agenturen  
**Tech-Stack:** Go (Golang), HTMX, SQLite, Vanilla CSS, i18n (Mehrsprachig DE & EN)  
**Fokus:** Budget-Drift-Erkennung ("aus dem Ruder gelaufen"), Profitabilitäts-Analyse, schnelle Zeiterfassung, geschütztes Kunden-Portal, Raspberry-Pi Quick-Tracker, Fraud-Detection, Pitch Demo Reset & Dual-AI-Workflow.

---

## 🎯 1. Problemstellung & Zielsetzung

Social Media Agenturen managen oft Dutzende parallele Kampagnen für verschiedene Kunden mit unterschiedlichen Content-Formaten (TikTok, Instagram Reels, Carousels, Stories).  
Häufige Probleme:
- **Budget-Drift:** Kampagnen überschreiten unbemerkt ihr Stunden- oder Euro-Budget.
- **Fehlende Transparenz:** Teamleads sehen erst am Monatsende, welche Mitarbeiter/Teams überbucht sind.
- **Unklare Margen:** Der Unterschied zwischen internem Stundensatz (Personalkosten) und Abrechnungs-Stundensatz (Kundensatz) wird pro Projekt oft nicht transparent nachverfolgt.
- **Kunden-Kommunikation:** Kunden verlangen transparente Echtzeit-Einblicke in ihren Kampagnenfortschritt, ohne jedoch interne Agentur-Margen oder Stundensätze sehen zu dürfen.
- **Kontext-Wechsel im Alltag:** Creator & Manager wechseln ständig zwischen Aufgaben (15 Min Anruf, 30 Min E-Mail, 2h Videoschnitt). Klassische Zeiterfassung ist zu träge.
- **Security & Fraud:** Öffentliche Kunden-Links bergen das Risiko von Brute-Force-Versuchen (Falsche PINs) oder Abscannen ungültiger Links durch Dritte.
- **Pitch-Sicherheit:** Während der Live-Vorführung werden Daten verändert; ein schneller Reset auf den perfekten Vorführ-Zustand ist unerlässlich.

**AgencyPulse** löst dies durch ein Echtzeit-Dashboard mit Ampelsystem, Ausreißer-Alerts, Rollen-spezifischen Sichten, einem **Raspberry-Pi Touchscreen Quick-Tracker**, **Storytelling-Demodaten**, **1-Klick Demo-Reset**, einer **DevTeam Security & Fraud Statusseite** sowie einem **parallelen Dual-AI Präsentations-Workflow**.

---

## 👥 2. Rollen, Datenmodell & Pitch-Demodaten

### Entitäten & Datenbankschema
1. **Kunden (Clients):** Name, Branche, Monatsbudget, **`portal_token` (Kryptische URL UUID)**, **`portal_pin` (PIN für Zugang)**.
2. **Kampagnen (Campaigns):** Kunde, Name, Gesamtbudget (€), Status (Grün, Gelb, Rot), **`is_favorite` (Boolean für Quick-Tracker)**.
3. **Teams & Teamleads:** Teamname, Teamleiter.
4. **Mitarbeiter (Employees):** Name, Rolle, Team, **Stundensatz Intern (€/h)**, **Stundensatz Kunde (€/h)**.
5. **Zeiterfassung (Time Logs):** Kampagne, Mitarbeiter, Arbeitsstunden (minutengenau + 15-Min-Abrechnungs-Intervall), Content-Typ, Beschreibung, Zeitstempel.
6. **Active Timer Sessions:** Mitarbeiter-ID, Kampagnen-ID, Startzeitpunkt (`started_at`).
7. **Security & Fraud Logs:** Event-Typ (`INVALID_PIN`, `INVALID_LINK_SCAN`), Ziel-Token, Versuchs-Anzahl, Timestamp, Status (`BLOCKED`, `WARNING`).

---

### 🎭 2.1 Demo-Szenarien für die Live-Präsentation (Storytelling)
Die Demodaten sind das Herzstück des Pitches und bilden 3 realistische Phasen von Agentur-Kampagnen ab:

| Kunde | Kampagne | Status | Budget-Verbrauch | Story für den Pitch |
| :--- | :--- | :---: | :---: | :--- |
| **Ritter Sport** | *Sommer-Sorte Carousel Push* | 🟢 **Grün** | ~35% (1.225 € / 3.500 €) | **Frische Kampagne:** Hohe Marge, pünktlich im Zeitplan. Zeigt den Idealzustand im Dashboard. |
| **Bosch Smart Home** | *Smart Lock Reels & Stories* | 🟡 **Gelb** | ~82% (3.690 € / 4.500 €) | **Warnphase:** Kampagne nähert sich dem Limit. Teamlead erhält Vorwarnung im Dashboard. |
| **Porsche Stuttgart** | *Taycan TikTok Launch* | 🔴 **Rot** | ~108% (6.480 € / 6.000 €) | **"Aus dem Ruder gelaufen":** Überbuchung durch 3D-Renderings. Löst n8n Alert & rote Heatmap-Karten aus! |

---

### 🔄 2.2 Versteckter 1-Klick Demo-Reset
- **Funktion:** Ein diskreter Button im Footer / Dev-Bereich (`🔄 Reset Demo Data` oder Shortcut) ruft den HTMX-Endpoint `POST /api/reset-demo-data` auf.
- **Wirkung:** Leert in Bruchteilen einer Sekunde die SQLite-Datenbank und baut den perfekten Zustand der 3 Pitch-Szenarien neu auf. Falls die Daten beim Zeigen "zerspielt" werden, bringt 1 Klick die App sofort wieder auf den Ausgangszustand zurück.

---

## 🌐 3. Mehrsprachigkeit & i18n (DE | EN)

- **Sprach-Toggle im Header & Portalen:** Ein Umschalter (**🇩🇪 DE | 🇬🇧 EN**) im Header sowie im Kunden-Portal & Kiosk schaltet die Sprache aller UI-Texte, Tabellen-Header, Buttons und Alert-Meldungen live um.
- **Technischer Ansatz:** Go-basiertes i18n Dictionary System (`locales/de.json` & `locales/en.json`).
- **Demo-Daten:** Die in der SQLite-Datenbank hinterlegten Stammdaten verbleiben in deutscher Sprache.

---

## 🖥️ 4. Sechs Kern-Perspektiven (Role Switcher, Touch-Kiosk, Kunden-Portal & DevTeam Monitor)

### 👑 Executive / Agentur View (Profitabilität & Margen)
- **Top KPIs:** Gesamtumsatz (€), Reingewinn (€), Agentur-Marge (%), Gefährdete Kampagnen.
- **Profitabilität pro Kunde:** Einnahmen vs. echte Personalkosten, Deckungsbeitrag.
- **Mitarbeiter-Effizienz:** Geleistete Stunden x Kundensatz vs. Gehaltskosten.

### 🛡️ Teamlead View (Budget-Drift & Warnungen)
- **Budget Heatmap:** Optische Fortschrittsbalken je Kampagne (`Soll` vs. `Ist`).
- **Gefahren-Alerts ("Aus dem Ruder gelaufen"):** Automatische Hervorhebung aller Kampagnen mit >80% oder >100% Budgetverbrauch.
- **AI Time & Budget Estimator:** Rechner zur Vorab-Kalkulation von Zeitschätzungen basierend auf Content-Formaten (z.B. 3x TikTok Video = ~9h Aufwand).

### ⏱️ Mitarbeiter View (Schnell-Zeiterfassung Standard)
- Minimalistisches Formular zur Buchung von Stunden auf Kampagnen und Content-Typen.
- Sofortige HTMX-Aktualisierung des Kampagnen-Budgets ohne Reload.

### 📱 Hardware-Kiosk & Quick-Tracker View (800x480 Raspberry Pi & Mobile Display)
- **Zielgruppe & Device:** Speziell optimiert für 7" Raspberry Pi Displays (800x480 Pixel), aber auch stufenlos skalierbar auf Smartphone, Tablet & Desktop (`/tracker`).
- **UI-Layout:** Kachel-Grid mit Live-Digital-Stoppuhr (`00:14:52`) & großen Touch-Buttons **▶ START** / **⏹ STOP**.
- **Strikte Stop/Start-Logik:** Der aktive Timer muss explizit gestoppt werden, bevor ein anderer startet.

### 🤝 Kunden View (Geschütztes Portal via Kryptischer Link + PIN)
- **Zugang:** Anonymisierte URL (z.B. `/portal/c/a8f9c2d1-4e7b-4a39-9c5d-88f12a34b5c6`) + 4-stellige PIN.
- **Kunden-Sicht:** Transparenter Fortschritt & gebuchte Zeiten ohne Einblick in interne Stundensätze.
- **i18n Support:** Eigenes Sprach-Icon im Portal-Header (DE | EN).

### 🚨 DevTeam & Security Status View (Fraud-Detection & Incident Monitor)
- **Echtzeit-Sicherheits-Cockpit:** Live-Feed aller Security-Events (`/dev/status`).
- **Fraud Detection:** 3x falsche PIN-Eingabe oder Aufruf ungültiger Portal-Links schlägt Alarm (n8n Webhook an Slack/Telegram).

---

## ⚡ 5. Sponsoren-Integration (Nahtlos & Praxisnah)

| Sponsor | Funktion im Konzept | Nahtloser Usecase (MVP) |
| :--- | :--- | :--- |
| **n8n** | **Budget-Alert & Fraud-Automation** | **1. Budget Drift:** Webhook löst Slack-Alert bei >100% Budget aus.<br>**2. Security Fraud:** bei Brute-Force PIN-Versuchen oder Link-Scan. |
| **ElevenLabs** | **Executive Voice-Briefing & Audio-Alerts** | Im Executive Dashboard & Kiosk gibt es ein **"🔊 Audio Briefing / Status"**. ElevenLabs generiert ein prägnantes Audio-Summary. |
| **fal.ai** | **Visual Content-Preview im AI Estimator** | Generiert in Echtzeit KI-Beispiel-Thumbnails für geplante Kampagnen-Assets im Kalkulator. |
| **Firecrawl** | **Competitor Content-Scan & Auto-Budgeting** | Scrawlt Social-Media-Profile von Neukunden & empfiehlt monatliche Stunden-Budgets. |
| **Cursor** | **Dev-Tooling / Building** | Rapid Development des Go + HTMX + SQLite MVPs während des Hackathons. |

---

## 🏗️ 6. Architektur & Ordnerstruktur

```
hack_2026_07/
├── KONZEPT.md                  # Dieses Dokument
├── CONCEPT.md                  # Englisches Konzept-Dokument
├── PITCH_SLIDES.md             # Automatisch vom Pitch-Agent erstelltes Präsentations-Deck (Marp/Markdown Format)
├── deploy-vps.sh               # Local cross-compile & SSH deploy to Kubuntu VPS (Linux)
├── deploy-vps.ps1              # Local cross-compile & SSH deploy to Kubuntu VPS (Windows PowerShell)
├── publish-github.sh           # Squash release commit & tag sync to public GitHub (Linux)
├── publish-github.ps1          # Squash release commit & tag sync to public GitHub (Windows PowerShell)
├── go.mod                      # Go Modul-Definition
├── main.go                     # Webserver Entrypoint
├── locales/
│   ├── de.json                 # Übersetzungen Deutsch
│   └── en.json                 # Übersetzungen Englisch
├── database/
│   └── db.go                   # SQLite Setup, Schema & Pitch-Demodaten (mit Reset-Funktion)
├── templates/
│   ├── index.html              # Haupt-Layout mit HTMX, Header, i18n Toggle & Demo-Reset Footer Button
│   └── partials/
│       ├── executive_view.html # Executive Dashboard Partial
│       ├── teamlead_view.html  # Teamlead Dashboard Partial
│       ├── employee_view.html  # Mitarbeiter Zeiterfassung Standard Partial
│       ├── kiosk_tracker.html  # 800x480 Touch-Kiosk Quick-Tracker (Start/Stop Cards)
│       ├── client_login.html   # PIN-Eingabemasken für Kunden-Portal
│       ├── client_portal.html  # Geschütztes Kunden-Dashboard
│       ├── dev_status.html     # DevTeam Security & Fraud Status Dashboard
│       └── ai_modal.html       # AI Estimator Partial
└── static/
    └── css/
        └── style.css           # Modernes Vanilla CSS (Dark Mode & 800x480 Kiosk Responsive Breakpoints)
```

---

## 🤖 7. Dual-AI-Agenten Workflow (Dev-Agent & Pitch-Agent)

Während des Hackathons arbeiten zwei KI-Agenten parallel Hand in Hand:

```
[ Developer + Dev-AI Agent ] ──> Entwickelt & verifiziert Step (z.B. v0.0.1, v0.1.0) ──> Git Commit & Git Tag
                                                                                               │
                                                                                               ▼
[ Pitch-AI Agent ] <── Liest fertige Dokumentation & Code ab Tag v0.0.1 ── Aktualisiert PITCH_SLIDES.md
```

1. **Dev-AI Agent (IDE Agent):** Entwickelt den Quellcode für den jeweiligen Schritt (`v0.0.1`, `v0.1.0`, `v0.2.0`, etc.), startet den Server, führt Funktionstests durch und setzt den annotierten Git-Tag.
2. **Pitch-AI Agent (Dokumentations & Präsentations-Agent):** Sobald der Dev-Agent ein Release markiert hat, greift der Pitch-Agent auf die fertige Codebasis zu und dokumentiert den genauen Stand in der Datei **`PITCH_SLIDES.md`**.
3. **Ergebnis:** Bei Erreichen von Version `v1.0.0` ist nicht nur die Software fertig, sondern auch die Präsentation für die Hackathon-Jury auf Punkt genau fertiggestellt und mit den realen Features synchronisiert.

---

## 🚀 8. Entwicklungs-Roadmap, Quality Gate & Versionierung

### 📌 Quality Gate, Dual-Remote & Release-Skripte:
1. **Lauffähigkeit prüfen:** Nach Fertigstellung jedes Schrittes wird der Go-Server gestartet (`go run main.go`), im Browser aufgerufen und die Funktionalität live getestet.
2. **Git-Commit & Private Push:** Nach erfolgreicher Prüfung wird ein sauberer Git-Commit erstellt, ein annotierter Git-Tag gesetzt (z.B. `git tag -a v0.0.1`) und per `git push private main --tags` auf das private Git (`ssh://git@192.168.178.100/ralf/hack_2026_07.git`) hochgeladen.
3. **VPS Deployment (`deploy-vps.sh` / `deploy-vps.ps1`):** Kompiliert lokal die Linux-Binärdatei (`GOOS=linux`), überträgt die Artefakte per SSH auf den Kubuntu VPS und startet den `agencypulse` systemd-Dienst neu.
4. **Öffentlicher GitHub Release (`publish-github.sh` / `publish-github.ps1`):** Nimmt die Versionsnummer entgegen (z.B. `./publish-github.sh v0.0.1`), führt ein `git merge main --squash` auf den `public-release` Branch durch, bereinigt temporäre Dateien, setzt den öffentlichen Tag `v0.0.1` und pusht den sauberen Einzelcommit + Tag nach GitHub (`origin`).
5. **Patch-Versionen (`v0.X.1`):** Werden ausschließlich verwendet, falls im Nachgang ein Fehler/Bug in einem bereits bestehenden Release behoben wird.

---

### 🔹 Step 0 — `v0.0.1`: Minimal Webserver & Deployment Setup (Bootstrap)
- **Lauffähige App:** Minimaler Go-Webserver auf Port 8080, der unformatierten Text `"Hello World!"` mit HTTP Status `200 OK` zurückgibt.
- **Deployment Skripte:** `deploy-vps.sh` / `deploy-vps.ps1` für VPS Deployment sowie `publish-github.sh` / `publish-github.ps1` für gequetschte GitHub-Releases.
- **Tooling:** `.vscode/tasks.json` für lokales Kompilieren und Starten (`Strg+Umschalt+B`).
- **Pitch Agent Task:** Initialisiert `PITCH_SLIDES.md` Deck-Struktur.

### 🔹 Step 1 — `v0.1.0`: Foundation, Storyline-Demodaten, Reset-Button & Standard-Zeiterfassung
- **Lauffähige App:** Go-Webserver mit i18n-Middleware, SQLite-Datenbank, Pitch-Demodaten (Ritter Sport = Grün, Bosch = Gelb, Porsche = Rot) & **Demo-Reset Endpoint (`POST /api/reset-demo-data`)**.
- **UI:** Basis Dark-Mode Design System (Vanilla CSS) mit Header, Navigation, **Sprach-Toggle (DE | EN)** und **diskretem Demo-Reset Button**.
- **Feature:** Standard Mitarbeiter-View mit HTMX-Formular zur manuellen Zeiterfassung.
- **Pitch Agent Task:** Erstellt Folien 1–3 in `PITCH_SLIDES.md` (Problem, Vision, Core Time Tracking UI).

### 🔹 Step 2 — `v0.2.0`: Teamlead Dashboard & Budget-Drift Alerts
- **Lauffähige App:** Teamlead-Cockpit mit Budget-Heatmap (`Soll` vs. `Ist` in %) und i18n-Labels.
- **Feature:** Automatische Ampel-Einstufung (`ok`, `warning`, `danger` bei >80% / >100% Budget).
- **Feature:** AI Time & Budget Estimator Modal zur Vorab-Kalkulation von Content-Paketen.
- **Pitch Agent Task:** Ergänzt Folien 4–5 in `PITCH_SLIDES.md` (Budget-Drift "Aus dem Ruder gelaufen" & AI Estimator).

### 🔹 Step 3 — `v0.3.0`: Executive Cockpit & Profitabilitäts-Analysen
- **Lauffähige App:** Executive View mit Finanz-KPIs (Gesamtumsatz, Personalkosten, Reingewinn, Agentur-Marge %).
- **Feature:** Aufschlüsselung der Kunden-Profitabilität und Mitarbeiter-Stundensätze.
- **Feature:** Audio Briefing Button (Integration von ElevenLabs Audio Summary).
- **Pitch Agent Task:** Ergänzt Folie 6 in `PITCH_SLIDES.md` (Profitabilität & ElevenLabs Voice Briefing).

### 🔹 Step 4 — `v0.4.0`: Hardware-Kiosk & Quick-Tracker (800x480 Raspberry Pi Display)
- **Lauffähige App:** Dedizierte Touch-Ansicht (`/tracker`) optimiert für 800x480 Pixel (auch responsive für Mobil/Desktop) mit i18n-Support.
- **Feature:** Favoriten-Kampagnen Kacheln mit riesiger Live-Digital-Stoppuhr und touch-freundlichen **Start / Stop** Buttons.
- **Feature:** Strikte Manual-Stop-Logik mit minutengenauer SQLite-Verbuchung & 15-Minuten-Abrechnungs-Kalkulation.
- **Pitch Agent Task:** Ergänzt Folie 7 in `PITCH_SLIDES.md` (Hardware Kiosk & 800x480 Touch UX).

### 🔹 Step 5 — `v0.5.0`: Geschütztes Kunden-Portal (Kryptischer Link + PIN)
- **Lauffähige App:** Anonymisierte/kryptische URLs für Kunden (z.B. `/portal/c/<uuid>`).
- **Feature:** HTMX-PIN-Eingabegate (4-stellige PIN) vor der Portal-Freischaltung.
- **Feature:** Kunden-Ansicht mit Budgetfortschritt und ausgelieferten Content-Assets mit DE/EN Umschalter.
- **Pitch Agent Task:** Ergänzt Folie 8 in `PITCH_SLIDES.md` (Transparentes Kunden-Portal & Security Gate).

### 🔹 Step 6 — `v0.6.0`: DevTeam Security, Fraud-Detection & n8n Integration
- **Lauffähige App:** Live DevTeam Security Cockpit (`/dev/status`).
- **Feature:** Fraud Detection Log für 3x falsche PINs (Brute Force) & Zugriffe auf ungültige Portal-URLs.
- **Feature:** n8n Webhook Dispatcher (sendet Security-Alerts & Budget-Drift Notifications an Slack/Telegram).
- **Pitch Agent Task:** Ergänzt Folie 9 in `PITCH_SLIDES.md` (Security Fraud Detection & n8n Automation).

### 🏆 Step 7 — `v1.0.0`: Hackathon Pitch Release
- **Lauffähige App:** Gesamtsystem-Polishing, Integration von fal.ai KI-Thumbnails und Firecrawl Competitor Scrape.
- **Ergebnis:** Schlüsselfertige, voll-integrierte Live-Demonstration für den Hackathon Pitch & fertiges Präsentationsdeck `PITCH_SLIDES.md`.

---

## 📢 9. Teaser & Mitstreiter-Pitch (für Discord & Matchmaking)

Hier ist der fertige Pitch-Text, mit dem du auf Discord oder bei der Vorstellungsrunde Teammitglieder begeistern kannst:

> **⚡ AgencyPulse – "Schluss mit Budgets, die aus dem Ruder laufen!"**
>
> **Das Problem:** Social Media Agenturen verbrennen täglich Tausende Euro, weil Kampagnen-Budgets unbemerkt überbucht werden, Stundensätze unklar sind und Creator zwischen TikTok, Reels & E-Mails den Überblick bei der Zeiterfassung verlieren.
>
> **Unsere Lösung:** **AgencyPulse** – ein ultra-schnelles, visuell beeindruckendes Agentur-Cockpit & Touch-Kiosk.
> - 🛡️ **Budget Heatmap & Live Alerts:** Sofortige Warnung, wenn Kampagnen >80% oder >100% erreichen.
> - 👑 **Executive Profitabilitäts-Cockpit:** Reingewinn & Margen pro Kunde & Mitarbeiter in Echtzeit.
> - 📱 **800x480 Hardware-Kiosk (Raspberry Pi):** Touch-Stoppuhr mit Start/Stop-Knöpfen für blitzschnellen Aufgabenwechsel am Arbeitsplatz.
> - 🤝 **Geschütztes Kunden-Portal:** PIN-gesicherte Kryptolinks für Kunden (ohne interne Margen-Einsicht).
> - 🤖 **AI & Automation Power:** ElevenLabs Voice Briefings, fal.ai Preview-Generator & n8n Slack Webhooks!
>
> **Tech Stack:** Go (Golang) + HTMX + SQLite + Vanilla CSS Dark Mode (High Speed, Zero Overhead).
>
> 👥 **Gesucht werden Mitstreiter für:**
> - 🎨 **UI/CSS Fine-Tuning & Styling** (Glassmorphism & Touch UX)
> - 🤖 **AI & Automation Integration** (n8n Webhooks, ElevenLabs Audio, fal.ai Thumbnails)
> - 🎬 **Pitch & Storytelling** (Vorführung & Demo-Kalkulation)
>
> *Ein klarer Schritt-für-Schritt-Plan (v0.0.1 bis v1.0.0) und ein automatischer Dual-AI-Workflow stehen bereit! Hast du Bock auf ein Gewinner-Projekt mit echtem Praxisnutzen? Melde dich!* 🚀
