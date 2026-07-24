# DEPLOYMENT.md — VPS Deployment Guide

Instructions for deploying AgencyPulse to a Linux (Kubuntu) VPS.

---

## 💻 Target Environment

- **OS:** Linux (Kubuntu / Ubuntu / Debian)
- **Target Path:** `/opt/agencypulse`
- **Default Port:** `8080`
- **Service Manager:** systemd (`agencypulse.service`)

---

## 🚀 Deployment Automation

### From Linux / macOS / Bash:
```bash
./deploy-vps.sh
```

### From Windows PowerShell:
```powershell
.\deploy-vps.ps1
```

---

## ⚙️ systemd Service Configuration

Create `/etc/systemd/system/agencypulse.service` on the VPS:

```ini
[Unit]
Description=AgencyPulse Social Media Budget & Time Tracker
After=network.target

[Service]
Type=simple
User=ralf
WorkingDirectory=/opt/agencypulse
ExecStart=/opt/agencypulse/bin/agencypulse
Restart=always
RestartSec=5
Environment=PORT=8080

[Install]
WantedBy=multi-user.target
```

Enable and start service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable agencypulse
sudo systemctl start agencypulse
```
