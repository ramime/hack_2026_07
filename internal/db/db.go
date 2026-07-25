package db

import (
	"database/sql"

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
	`
	if _, err := db.Exec(schema); err != nil {
		return err
	}

	// Migrations for existing database schemas
	db.Exec("ALTER TABLE employees ADD COLUMN cost_rate REAL DEFAULT 0")
	db.Exec("ALTER TABLE employees ADD COLUMN billing_rate REAL DEFAULT 0")

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
	if _, err := tx.Exec("DELETE FROM active_timer_sessions; DELETE FROM time_logs; DELETE FROM campaigns; DELETE FROM clients; DELETE FROM employees;"); err != nil {
		return err
	}

	// Insert Clients
	resRS, err := tx.Exec("INSERT INTO clients (id, name) VALUES (1, 'Ritter Sport')")
	if err != nil {
		return err
	}
	_ = resRS

	tx.Exec("INSERT INTO clients (id, name) VALUES (2, 'Bosch')")
	tx.Exec("INSERT INTO clients (id, name) VALUES (3, 'Porsche')")

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

	return tx.Commit()
}
