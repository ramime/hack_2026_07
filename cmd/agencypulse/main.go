package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"agencypulse/internal/db"
	"agencypulse/internal/i18n"
	"agencypulse/internal/models"
)

const version = "0.1.1"

type PageData struct {
	Version    string
	Lang       string
	Employees  []models.Employee
	Campaigns  []models.Campaign
	TimeLogs   []models.TimeLog
	SuccessMsg string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "agencypulse.db"
	}

	database, err := db.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	i18nMgr, err := i18n.NewManager("locales")
	if err != nil {
		log.Fatalf("Failed to load i18n manager: %v", err)
	}

	tmpl, err := template.New("base").Funcs(template.FuncMap{
		"t": func(lang, key string) string {
			return i18nMgr.T(lang, key)
		},
	}).ParseFiles("web/templates/layout.html", "web/templates/employee.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	// Static file handler
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Get current language from cookie or default to "de"
	getLang := func(r *http.Request) string {
		cookie, err := r.Cookie("lang")
		if err == nil && (cookie.Value == "de" || cookie.Value == "en") {
			return cookie.Value
		}
		return "de"
	}

	// Helper to fetch page data
	fetchPageData := func(lang, successMsg string) (*PageData, error) {
		// Fetch employees
		empRows, err := database.Query("SELECT id, name, role, hourly_rate, created_at FROM employees ORDER BY name ASC")
		if err != nil {
			return nil, err
		}
		defer empRows.Close()

		var employees []models.Employee
		for empRows.Next() {
			var emp models.Employee
			if err := empRows.Scan(&emp.ID, &emp.Name, &emp.Role, &emp.HourlyRate, &emp.CreatedAt); err != nil {
				return nil, err
			}
			employees = append(employees, emp)
		}

		// Fetch campaigns with client names
		campRows, err := database.Query(`
			SELECT c.id, c.client_id, cl.name, c.name, c.target_budget, c.created_at
			FROM campaigns c
			JOIN clients cl ON c.client_id = cl.id
			ORDER BY cl.name ASC, c.name ASC
		`)
		if err != nil {
			return nil, err
		}
		defer campRows.Close()

		var campaigns []models.Campaign
		for campRows.Next() {
			var camp models.Campaign
			if err := campRows.Scan(&camp.ID, &camp.ClientID, &camp.ClientName, &camp.Name, &camp.TargetBudget, &camp.CreatedAt); err != nil {
				return nil, err
			}
			campaigns = append(campaigns, camp)
		}

		// Fetch recent time logs
		logRows, err := database.Query(`
			SELECT tl.id, tl.campaign_id, c.name, tl.employee_id, e.name, tl.hours, tl.description, tl.logged_at
			FROM time_logs tl
			JOIN campaigns c ON tl.campaign_id = c.id
			JOIN employees e ON tl.employee_id = e.id
			ORDER BY tl.logged_at DESC
			LIMIT 20
		`)
		if err != nil {
			return nil, err
		}
		defer logRows.Close()

		var logs []models.TimeLog
		for logRows.Next() {
			var tl models.TimeLog
			if err := logRows.Scan(&tl.ID, &tl.CampaignID, &tl.CampaignName, &tl.EmployeeID, &tl.EmployeeName, &tl.Hours, &tl.Description, &tl.LoggedAt); err != nil {
				return nil, err
			}
			logs = append(logs, tl)
		}

		return &PageData{
			Version:    version,
			Lang:       lang,
			Employees:  employees,
			Campaigns:  campaigns,
			TimeLogs:   logs,
			SuccessMsg: successMsg,
		}, nil
	}

	// Route: GET /
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		lang := getLang(r)
		data, err := fetchPageData(lang, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Route: POST /api/time-logs
	http.HandleFunc("/api/time-logs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		empID, _ := strconv.ParseInt(r.FormValue("employee_id"), 10, 64)
		campID, _ := strconv.ParseInt(r.FormValue("campaign_id"), 10, 64)
		hours, _ := strconv.ParseFloat(r.FormValue("hours"), 64)
		desc := r.FormValue("description")

		if empID > 0 && campID > 0 && hours > 0 && desc != "" {
			_, err := database.Exec(
				"INSERT INTO time_logs (campaign_id, employee_id, hours, description, logged_at) VALUES (?, ?, ?, ?, ?)",
				campID, empID, hours, desc, time.Now(),
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		lang := getLang(r)
		data, err := fetchPageData(lang, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// HTMX request: render "content" block directly
		if r.Header.Get("HX-Request") == "true" {
			if err := tmpl.ExecuteTemplate(w, "content", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Route: POST /api/reset-demo-data
	http.HandleFunc("/api/reset-demo-data", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := database.ResetToSeedData(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lang := getLang(r)
		msg := i18nMgr.T(lang, "demo_reset_success")
		data, err := fetchPageData(lang, msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Header.Get("HX-Request") == "true" {
			if err := tmpl.ExecuteTemplate(w, "content", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Route: GET /api/health
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"ok":      true,
			"version": version,
		})
	})

	// Route: POST /api/set-language
	http.HandleFunc("/api/set-language", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		newLang := r.FormValue("lang")
		if newLang != "de" && newLang != "en" {
			newLang = "de"
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "lang",
			Value:    newLang,
			Path:     "/",
			Expires:  time.Now().Add(365 * 24 * time.Hour),
			HttpOnly: true,
		})

		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Refresh", "true")
			w.WriteHeader(http.StatusOK)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Printf("AgencyPulse v%s starting on http://localhost:%s...", version, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
