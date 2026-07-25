package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"agencypulse/internal/db"
	"agencypulse/internal/i18n"
	"agencypulse/internal/models"
	"agencypulse/internal/pitch"
	"agencypulse/internal/tts"
)

const version = "0.3.0"


type PageData struct {
	Version         string
	Lang            string
	ActiveNav       string
	Employees       []models.Employee
	Campaigns       []models.Campaign
	TimeLogs        []models.TimeLog
	BudgetSummaries []models.CampaignBudgetSummary
	CountHealthy    int
	CountWarning    int
	CountDanger     int
	SuccessMsg      string
	ExecKPIs        models.ExecutiveKPIs
	ClientProfits   []models.ClientProfitability
	EmpEfficiencies []models.EmployeeEfficiency
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

	funcMap := template.FuncMap{
		"t": func(lang, key string) string {
			return i18nMgr.T(lang, key)
		},
	}

	empTmpl, err := template.New("layout.html").Funcs(funcMap).ParseFiles("web/templates/layout.html", "web/templates/employee.html")
	if err != nil {
		log.Fatalf("Failed to parse employee templates: %v", err)
	}

	teamleadTmpl, err := template.New("layout.html").Funcs(funcMap).ParseFiles("web/templates/layout.html", "web/templates/teamlead.html")
	if err != nil {
		log.Fatalf("Failed to parse teamlead templates: %v", err)
	}

	executiveTmpl, err := template.New("layout.html").Funcs(funcMap).ParseFiles("web/templates/layout.html", "web/templates/executive.html")
	if err != nil {
		log.Fatalf("Failed to parse executive templates: %v", err)
	}

	slidesTmpl, err := template.ParseFiles("web/templates/slides.html")
	if err != nil {
		log.Fatalf("Failed to parse slides template: %v", err)
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
		empRows, err := database.Query("SELECT id, name, role, hourly_rate, cost_rate, billing_rate, created_at FROM employees ORDER BY name ASC")
		if err != nil {
			return nil, err
		}
		defer empRows.Close()

		var employees []models.Employee
		for empRows.Next() {
			var emp models.Employee
			if err := empRows.Scan(&emp.ID, &emp.Name, &emp.Role, &emp.HourlyRate, &emp.CostRate, &emp.BillingRate, &emp.CreatedAt); err != nil {
				return nil, err
			}
			if emp.BillingRate == 0 {
				emp.BillingRate = emp.HourlyRate
			}
			if emp.CostRate == 0 {
				emp.CostRate = emp.BillingRate * 0.5
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

		// Fetch budget summaries for Team Lead view
		summaryRows, err := database.Query(`
			SELECT 
				c.id, cl.name, c.name, c.target_budget,
				COALESCE(SUM(tl.hours * COALESCE(NULLIF(e.billing_rate, 0), e.hourly_rate)), 0) AS actual_spend,
				COALESCE(SUM(tl.hours), 0) AS hours_logged
			FROM campaigns c
			JOIN clients cl ON c.client_id = cl.id
			LEFT JOIN time_logs tl ON c.id = tl.campaign_id
			LEFT JOIN employees e ON tl.employee_id = e.id
			GROUP BY c.id, cl.name, c.name, c.target_budget
			ORDER BY cl.name ASC, c.name ASC
		`)
		if err != nil {
			return nil, err
		}
		defer summaryRows.Close()

		var summaries []models.CampaignBudgetSummary
		var countHealthy, countWarning, countDanger int

		for summaryRows.Next() {
			var s models.CampaignBudgetSummary
			if err := summaryRows.Scan(&s.CampaignID, &s.ClientName, &s.CampaignName, &s.TargetBudget, &s.ActualSpend, &s.HoursLogged); err != nil {
				return nil, err
			}
			if s.TargetBudget > 0 {
				s.UsagePercent = (s.ActualSpend / s.TargetBudget) * 100.0
			}
			if s.UsagePercent > 100.0 {
				s.Status = "danger"
				countDanger++
			} else if s.UsagePercent >= 80.0 {
				s.Status = "warning"
				countWarning++
			} else {
				s.Status = "ok"
				countHealthy++
			}
			summaries = append(summaries, s)
		}

		// Calculate Client Profitability Breakdown
		clientRows, err := database.Query(`
			SELECT 
				cl.id, cl.name,
				COUNT(DISTINCT c.id) as campaign_count,
				COALESCE(SUM(tl.hours), 0) as total_hours,
				COALESCE(SUM(tl.hours * COALESCE(NULLIF(e.billing_rate, 0), e.hourly_rate)), 0) as billed_revenue,
				COALESCE(SUM(tl.hours * COALESCE(NULLIF(e.cost_rate, 0), e.hourly_rate * 0.5)), 0) as labor_cost
			FROM clients cl
			LEFT JOIN campaigns c ON cl.id = c.client_id
			LEFT JOIN time_logs tl ON c.id = tl.campaign_id
			LEFT JOIN employees e ON tl.employee_id = e.id
			GROUP BY cl.id, cl.name
			ORDER BY cl.name ASC
		`)
		if err != nil {
			return nil, err
		}
		defer clientRows.Close()

		var clientProfits []models.ClientProfitability
		var totalRev, totalCost float64

		for clientRows.Next() {
			var cp models.ClientProfitability
			if err := clientRows.Scan(&cp.ClientID, &cp.ClientName, &cp.CampaignCount, &cp.TotalHours, &cp.BilledRevenue, &cp.LaborCost); err != nil {
				return nil, err
			}
			cp.NetProfit = cp.BilledRevenue - cp.LaborCost
			if cp.BilledRevenue > 0 {
				cp.MarginPercent = (cp.NetProfit / cp.BilledRevenue) * 100.0
			}
			totalRev += cp.BilledRevenue
			totalCost += cp.LaborCost
			clientProfits = append(clientProfits, cp)
		}

		// Calculate Employee Efficiency Breakdown
		effRows, err := database.Query(`
			SELECT 
				e.id, e.name, e.role,
				COALESCE(NULLIF(e.cost_rate, 0), e.hourly_rate * 0.5) as cost_rate,
				COALESCE(NULLIF(e.billing_rate, 0), e.hourly_rate) as billing_rate,
				COALESCE(SUM(tl.hours), 0) as hours_logged
			FROM employees e
			LEFT JOIN time_logs tl ON e.id = tl.employee_id
			GROUP BY e.id, e.name, e.role, e.cost_rate, e.billing_rate, e.hourly_rate
			ORDER BY e.name ASC
		`)
		if err != nil {
			return nil, err
		}
		defer effRows.Close()

		var empEfficiencies []models.EmployeeEfficiency
		for effRows.Next() {
			var ee models.EmployeeEfficiency
			if err := effRows.Scan(&ee.EmployeeID, &ee.EmployeeName, &ee.Role, &ee.CostRate, &ee.BillingRate, &ee.HoursLogged); err != nil {
				return nil, err
			}
			ee.BilledRevenue = ee.HoursLogged * ee.BillingRate
			ee.LaborCost = ee.HoursLogged * ee.CostRate
			ee.NetContribution = ee.BilledRevenue - ee.LaborCost
			if ee.BilledRevenue > 0 {
				ee.MarginPercent = (ee.NetContribution / ee.BilledRevenue) * 100.0
			}
			empEfficiencies = append(empEfficiencies, ee)
		}

		netProfit := totalRev - totalCost
		var agencyMargin float64
		if totalRev > 0 {
			agencyMargin = (netProfit / totalRev) * 100.0
		}

		execKPIs := models.ExecutiveKPIs{
			TotalRevenue:   totalRev,
			TotalLaborCost: totalCost,
			NetProfit:      netProfit,
			AgencyMargin:   agencyMargin,
			AtRiskCount:    countWarning + countDanger,
		}

		return &PageData{
			Version:         version,
			Lang:            lang,
			Employees:       employees,
			Campaigns:       campaigns,
			TimeLogs:        logs,
			BudgetSummaries: summaries,
			CountHealthy:    countHealthy,
			CountWarning:    countWarning,
			CountDanger:     countDanger,
			SuccessMsg:      successMsg,
			ExecKPIs:        execKPIs,
			ClientProfits:   clientProfits,
			EmpEfficiencies: empEfficiencies,
		}, nil
	}

	// Route: GET /slides (pitch deck — content from pitch/slides.md)
	http.HandleFunc("/slides", func(w http.ResponseWriter, r *http.Request) {
		slides, err := pitch.LoadSlides("pitch/slides.md")
		if err != nil {
			http.Error(w, "slides unavailable: "+err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			Version string
			Slides  []pitch.Slide
		}{Version: version, Slides: slides}
		if err := slidesTmpl.ExecuteTemplate(w, "slides.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Route: GET / (Employee View)
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
		data.ActiveNav = "employee"
		if err := empTmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Route: GET /teamlead (Team Lead Cockpit View)
	http.HandleFunc("/teamlead", func(w http.ResponseWriter, r *http.Request) {
		lang := getLang(r)
		data, err := fetchPageData(lang, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.ActiveNav = "teamlead"
		if err := teamleadTmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Route: GET /executive (Executive View)
	http.HandleFunc("/executive", func(w http.ResponseWriter, r *http.Request) {
		lang := getLang(r)
		data, err := fetchPageData(lang, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.ActiveNav = "executive"
		if err := executiveTmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Route: GET /api/audio-briefing (ElevenLabs Audio Executive Briefing)
	http.HandleFunc("/api/audio-briefing", func(w http.ResponseWriter, r *http.Request) {
		lang := getLang(r)
		data, err := fetchPageData(lang, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		briefingText := fmt.Sprintf(
			"AgencyPulse Executive Briefing: Total revenue is %d euros with net profit of %d euros and agency margin of %.1f percent. There are %d campaigns currently in warning or danger status.",
			int(data.ExecKPIs.TotalRevenue),
			int(data.ExecKPIs.NetProfit),
			data.ExecKPIs.AgencyMargin,
			data.ExecKPIs.AtRiskCount,
		)

		audioBytes, err := tts.GenerateExecutiveAudio(briefingText)
		if err != nil {
			http.Error(w, "Failed to generate audio briefing: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "audio/wav")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		w.Write(audioBytes)
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
		data.ActiveNav = "employee"

		// HTMX request: render "content" block directly
		if r.Header.Get("HX-Request") == "true" {
			if err := empTmpl.ExecuteTemplate(w, "content", data); err != nil {
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
		data.ActiveNav = "employee"

		if r.Header.Get("HX-Request") == "true" {
			if err := empTmpl.ExecuteTemplate(w, "content", data); err != nil {
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
