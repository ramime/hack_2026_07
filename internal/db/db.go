package db

import (
	"database/sql"

	"agencypulse/internal/models"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func InitDB(filepath string) (*DB, error) {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}

	// Enable WAL mode & foreign keys for performance and consistency
	if _, err := db.Exec(`PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;`); err != nil {
		return nil, err
	}

	database := &DB{DB: db}
	if err := database.createTables(); err != nil {
		return nil, err
	}

	if err := database.SeedDataIfEmpty(); err != nil {
		return nil, err
	}

	return database, nil
}

func (db *DB) createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS clients (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS campaigns (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		client_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		target_budget REAL NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS employees (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		role TEXT NOT NULL,
		hourly_rate REAL NOT NULL,
		cost_rate REAL DEFAULT 0,
		billing_rate REAL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS time_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		campaign_id INTEGER NOT NULL,
		employee_id INTEGER NOT NULL,
		hours REAL NOT NULL,
		description TEXT NOT NULL,
		logged_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE,
		FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS active_timer_sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		employee_id INTEGER NOT NULL UNIQUE,
		campaign_id INTEGER NOT NULL,
		task_category TEXT NOT NULL DEFAULT 'Content & Editing',
		started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
		FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS content_assets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		campaign_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		asset_type TEXT NOT NULL,
		status TEXT NOT NULL,
		preview_url TEXT DEFAULT '',
		delivered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS security_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_type TEXT NOT NULL,
		target_token TEXT NOT NULL,
		ip_address TEXT NOT NULL,
		attempt_count INTEGER DEFAULT 1,
		status TEXT NOT NULL,
		details TEXT NOT NULL DEFAULT '',
		logged_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return err
	}

	// Migrations for existing database schemas
	db.Exec("ALTER TABLE employees ADD COLUMN cost_rate REAL DEFAULT 0")
	db.Exec("ALTER TABLE employees ADD COLUMN billing_rate REAL DEFAULT 0")
	db.Exec("ALTER TABLE clients ADD COLUMN portal_token TEXT DEFAULT ''")
	db.Exec("ALTER TABLE clients ADD COLUMN pin_code TEXT DEFAULT '1234'")

	return nil
}

func (db *DB) SeedDataIfEmpty() error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM clients").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	return db.ResetToSeedData()
}

func (db *DB) ResetToSeedData() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear existing tables
	if _, err := tx.Exec("DELETE FROM security_logs; DELETE FROM content_assets; DELETE FROM active_timer_sessions; DELETE FROM time_logs; DELETE FROM campaigns; DELETE FROM clients; DELETE FROM employees;"); err != nil {
		return err
	}

	// Insert Clients with portal tokens & PIN codes
	if _, err := tx.Exec("INSERT INTO clients (id, name, portal_token, pin_code) VALUES (1, 'Ritter Sport', 'ritter-sport-8821', '1234')"); err != nil {
		return err
	}
	tx.Exec("INSERT INTO clients (id, name, portal_token, pin_code) VALUES (2, 'Bosch', 'bosch-4492', '5678')")
	tx.Exec("INSERT INTO clients (id, name, portal_token, pin_code) VALUES (3, 'Porsche', 'porsche-9102', '9900')")

	// Insert Campaigns
	// 1: Ritter Sport - Green (Target 10,000 €)
	tx.Exec("INSERT INTO campaigns (id, client_id, name, target_budget) VALUES (1, 1, 'Summer Special Edition 2026', 10000.0)")
	// 2: Bosch - Yellow (Target 25,000 €)
	tx.Exec("INSERT INTO campaigns (id, client_id, name, target_budget) VALUES (2, 2, 'Smart Home TikTok Launch', 25000.0)")
	// 3: Porsche - Red (Target 50,000 €)
	tx.Exec("INSERT INTO campaigns (id, client_id, name, target_budget) VALUES (3, 3, 'Taycan GT Creator Experience', 50000.0)")

	// Insert Employees with Cost Rate & Billing Rate
	// Sarah: Team Lead (Cost: €60/h, Billing: €120/h)
	tx.Exec("INSERT INTO employees (id, name, role, hourly_rate, cost_rate, billing_rate) VALUES (1, 'Sarah Meyer', 'Team Lead', 120.0, 60.0, 120.0)")
	// Alex: Senior Designer (Cost: €45/h, Billing: €95/h)
	tx.Exec("INSERT INTO employees (id, name, role, hourly_rate, cost_rate, billing_rate) VALUES (2, 'Alex Weber', 'Senior Designer', 95.0, 45.0, 95.0)")
	// Max: Content Creator (Cost: €35/h, Billing: €85/h)
	tx.Exec("INSERT INTO employees (id, name, role, hourly_rate, cost_rate, billing_rate) VALUES (3, 'Max Schmidt', 'Content Creator', 85.0, 35.0, 85.0)")

	// Insert Initial Time Logs
	// Ritter Sport logs (~4,500 € total)
	tx.Exec("INSERT INTO time_logs (campaign_id, employee_id, hours, description, logged_at) VALUES (1, 1, 15.0, 'Concept & Storyboarding', DATETIME('now', '-5 days'))")
	tx.Exec("INSERT INTO time_logs (campaign_id, employee_id, hours, description, logged_at) VALUES (1, 2, 20.0, 'Packaging & Visual Assets', DATETIME('now', '-3 days'))")
	tx.Exec("INSERT INTO time_logs (campaign_id, employee_id, hours, description, logged_at) VALUES (1, 3, 10.0, 'Reels & TikTok Editing', DATETIME('now', '-1 day'))")

	// Bosch logs (~21,000 € total -> 84% of 25,000 € -> Warning Yellow)
	tx.Exec("INSERT INTO time_logs (campaign_id, employee_id, hours, description, logged_at) VALUES (2, 1, 50.0, 'Campaign Strategy & PM', DATETIME('now', '-10 days'))")
	tx.Exec("INSERT INTO time_logs (campaign_id, employee_id, hours, description, logged_at) VALUES (2, 2, 80.0, '3D Motion Graphics & Assets', DATETIME('now', '-7 days'))")
	tx.Exec("INSERT INTO time_logs (campaign_id, employee_id, hours, description, logged_at) VALUES (2, 3, 87.0, 'Influencer Coordination & Video Cuts', DATETIME('now', '-2 days'))")

	// Porsche logs (~57,500 € total -> 115% of 50,000 € -> Danger Red)
	tx.Exec("INSERT INTO time_logs (campaign_id, employee_id, hours, description, logged_at) VALUES (3, 1, 120.0, 'Executive Production & Logistics', DATETIME('now', '-14 days'))")
	tx.Exec("INSERT INTO time_logs (campaign_id, employee_id, hours, description, logged_at) VALUES (3, 2, 220.0, 'High-End Visual Effects & Grading', DATETIME('now', '-8 days'))")
	tx.Exec("INSERT INTO time_logs (campaign_id, employee_id, hours, description, logged_at) VALUES (3, 3, 260.0, 'On-Location Shoot & Audio Master', DATETIME('now', '-1 day'))")

	// Insert Sample Delivered Content Assets for Client Portal
	tx.Exec("INSERT INTO content_assets (campaign_id, title, asset_type, status, preview_url, delivered_at) VALUES (1, 'Summer Special TikTok Teaser Video', 'TikTok Video', 'Delivered', '', DATETIME('now', '-2 days'))")
	tx.Exec("INSERT INTO content_assets (campaign_id, title, asset_type, status, preview_url, delivered_at) VALUES (1, 'Instagram Story Carousel #1', 'Instagram Story', 'Approved', '', DATETIME('now', '-1 day'))")
	tx.Exec("INSERT INTO content_assets (campaign_id, title, asset_type, status, preview_url, delivered_at) VALUES (2, 'Smart Home 3D Render & Key Visual', '3D Motion Graphic', 'In Review', '', DATETIME('now', '-4 days'))")
	tx.Exec("INSERT INTO content_assets (campaign_id, title, asset_type, status, preview_url, delivered_at) VALUES (2, 'Influencer Unboxing & Demo Cut', 'TikTok Video', 'Delivered', '', DATETIME('now', '-2 days'))")
	tx.Exec("INSERT INTO content_assets (campaign_id, title, asset_type, status, preview_url, delivered_at) VALUES (3, 'Taycan GT Track Day 4K Master', 'High-End Video', 'Approved', '', DATETIME('now', '-5 days'))")
	tx.Exec("INSERT INTO content_assets (campaign_id, title, asset_type, status, preview_url, delivered_at) VALUES (3, 'Engine Sound & Spatial Audio Mix', 'Audio Asset', 'Delivered', '', DATETIME('now', '-3 days'))")

	// Insert Sample Security Logs for DevTeam Cockpit
	tx.Exec("INSERT INTO security_logs (event_type, target_token, ip_address, attempt_count, status, details, logged_at) VALUES ('INVALID_PIN', 'ritter-sport-8821', '192.168.178.45', 1, 'WARNING', 'Incorrect PIN attempt (1/3)', DATETIME('now', '-2 hours'))")
	tx.Exec("INSERT INTO security_logs (event_type, target_token, ip_address, attempt_count, status, details, logged_at) VALUES ('INVALID_LINK_SCAN', '/portal/c/unknown-hex-9912', '45.142.120.9', 1, 'WARNING', 'Automated web crawler scanned non-existent portal URL', DATETIME('now', '-1 hour'))")
	tx.Exec("INSERT INTO security_logs (event_type, target_token, ip_address, attempt_count, status, details, logged_at) VALUES ('BUDGET_DRIFT_ALERT', 'Taycan GT Creator Experience', 'system', 1, 'BLOCKED', 'Campaign target budget exceeded (115% usage)', DATETIME('now', '-30 minutes'))")

	return tx.Commit()
}

func (db *DB) GetClientByToken(token string) (*models.Client, error) {
	var c models.Client
	err := db.QueryRow("SELECT id, name, portal_token, pin_code, created_at FROM clients WHERE portal_token = ?", token).
		Scan(&c.ID, &c.Name, &c.PortalToken, &c.PinCode, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (db *DB) GetClientContentAssets(clientID int64) ([]models.ContentAsset, error) {
	rows, err := db.Query(`
		SELECT ca.id, ca.campaign_id, c.name, ca.title, ca.asset_type, ca.status, ca.preview_url, ca.delivered_at
		FROM content_assets ca
		JOIN campaigns c ON ca.campaign_id = c.id
		WHERE c.client_id = ?
		ORDER BY ca.delivered_at DESC
	`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []models.ContentAsset
	for rows.Next() {
		var a models.ContentAsset
		if err := rows.Scan(&a.ID, &a.CampaignID, &a.CampaignName, &a.Title, &a.AssetType, &a.Status, &a.PreviewURL, &a.DeliveredAt); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (db *DB) GetClientCampaignSummaries(clientID int64) ([]models.CampaignBudgetSummary, error) {
	rows, err := db.Query(`
		SELECT 
			c.id as campaign_id,
			cl.name as client_name,
			c.name as campaign_name,
			c.target_budget,
			COALESCE(SUM(tl.hours * COALESCE(NULLIF(e.billing_rate, 0), e.hourly_rate)), 0) as actual_spend,
			COALESCE(SUM(tl.hours), 0) as hours_logged
		FROM campaigns c
		JOIN clients cl ON c.client_id = cl.id
		LEFT JOIN time_logs tl ON c.id = tl.campaign_id
		LEFT JOIN employees e ON tl.employee_id = e.id
		WHERE c.client_id = ?
		GROUP BY c.id, cl.name, c.name, c.target_budget
		ORDER BY c.name ASC
	`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []models.CampaignBudgetSummary
	for rows.Next() {
		var s models.CampaignBudgetSummary
		if err := rows.Scan(&s.CampaignID, &s.ClientName, &s.CampaignName, &s.TargetBudget, &s.ActualSpend, &s.HoursLogged); err != nil {
			return nil, err
		}
		if s.TargetBudget > 0 {
			s.UsagePercent = (s.ActualSpend / s.TargetBudget) * 100.0
		}
		if s.UsagePercent > 100.0 {
			s.Status = "danger"
		} else if s.UsagePercent >= 80.0 {
			s.Status = "warning"
		} else {
			s.Status = "ok"
		}
		summaries = append(summaries, s)
	}
	return summaries, nil
}

func (db *DB) LogSecurityEvent(eventType, targetToken, ipAddress string, attemptCount int, status, details string) (*models.SecurityLog, error) {
	res, err := db.Exec(`
		INSERT INTO security_logs (event_type, target_token, ip_address, attempt_count, status, details, logged_at)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`, eventType, targetToken, ipAddress, attemptCount, status, details)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &models.SecurityLog{
		ID:           id,
		EventType:    eventType,
		TargetToken:  targetToken,
		IPAddress:    ipAddress,
		AttemptCount: attemptCount,
		Status:       status,
		Details:      details,
	}, nil
}

func (db *DB) GetSecurityLogs(limit int) ([]models.SecurityLog, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := db.Query(`
		SELECT id, event_type, target_token, ip_address, attempt_count, status, details, logged_at
		FROM security_logs
		ORDER BY id DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.SecurityLog
	for rows.Next() {
		var l models.SecurityLog
		if err := rows.Scan(&l.ID, &l.EventType, &l.TargetToken, &l.IPAddress, &l.AttemptCount, &l.Status, &l.Details, &l.LoggedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

// Client CRUD
func (db *DB) GetAllClients() ([]models.Client, error) {
	rows, err := db.Query("SELECT id, name, portal_token, pin_code, created_at FROM clients ORDER BY name ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var c models.Client
		if err := rows.Scan(&c.ID, &c.Name, &c.PortalToken, &c.PinCode, &c.CreatedAt); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, nil
}

func (db *DB) CreateClient(name, portalToken, pinCode string) error {
	_, err := db.Exec("INSERT INTO clients (name, portal_token, pin_code) VALUES (?, ?, ?)", name, portalToken, pinCode)
	return err
}

func (db *DB) UpdateClient(id int64, name, portalToken, pinCode string) error {
	_, err := db.Exec("UPDATE clients SET name = ?, portal_token = ?, pin_code = ? WHERE id = ?", name, portalToken, pinCode, id)
	return err
}

func (db *DB) DeleteClient(id int64) error {
	_, err := db.Exec("DELETE FROM clients WHERE id = ?", id)
	return err
}

// Employee CRUD
func (db *DB) GetAllEmployees() ([]models.Employee, error) {
	rows, err := db.Query("SELECT id, name, role, hourly_rate, cost_rate, billing_rate, created_at FROM employees ORDER BY name ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var e models.Employee
		if err := rows.Scan(&e.ID, &e.Name, &e.Role, &e.HourlyRate, &e.CostRate, &e.BillingRate, &e.CreatedAt); err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	return employees, nil
}

func (db *DB) CreateEmployee(name, role string, hourlyRate, costRate, billingRate float64) error {
	_, err := db.Exec("INSERT INTO employees (name, role, hourly_rate, cost_rate, billing_rate) VALUES (?, ?, ?, ?, ?)", name, role, hourlyRate, costRate, billingRate)
	return err
}

func (db *DB) UpdateEmployee(id int64, name, role string, hourlyRate, costRate, billingRate float64) error {
	_, err := db.Exec("UPDATE employees SET name = ?, role = ?, hourly_rate = ?, cost_rate = ?, billing_rate = ? WHERE id = ?", name, role, hourlyRate, costRate, billingRate, id)
	return err
}

func (db *DB) DeleteEmployee(id int64) error {
	_, err := db.Exec("DELETE FROM employees WHERE id = ?", id)
	return err
}

// Campaign CRUD
func (db *DB) GetAllCampaigns() ([]models.Campaign, error) {
	rows, err := db.Query(`
		SELECT c.id, c.client_id, cl.name as client_name, c.name, c.target_budget, c.created_at
		FROM campaigns c
		JOIN clients cl ON c.client_id = cl.id
		ORDER BY c.name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []models.Campaign
	for rows.Next() {
		var cmp models.Campaign
		if err := rows.Scan(&cmp.ID, &cmp.ClientID, &cmp.ClientName, &cmp.Name, &cmp.TargetBudget, &cmp.CreatedAt); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, cmp)
	}
	return campaigns, nil
}

func (db *DB) CreateCampaign(clientID int64, name string, targetBudget float64) error {
	_, err := db.Exec("INSERT INTO campaigns (client_id, name, target_budget) VALUES (?, ?, ?)", clientID, name, targetBudget)
	return err
}

func (db *DB) UpdateCampaign(id int64, clientID int64, name string, targetBudget float64) error {
	_, err := db.Exec("UPDATE campaigns SET client_id = ?, name = ?, target_budget = ? WHERE id = ?", clientID, name, targetBudget, id)
	return err
}

func (db *DB) DeleteCampaign(id int64) error {
	_, err := db.Exec("DELETE FROM campaigns WHERE id = ?", id)
	return err
}



