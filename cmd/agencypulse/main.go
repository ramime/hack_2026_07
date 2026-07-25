package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"agencypulse/internal/db"
	"agencypulse/internal/i18n"
	"agencypulse/internal/models"
	"agencypulse/internal/n8n"
	"agencypulse/internal/pitch"
	"agencypulse/internal/tts"
)

const version = "0.7.4"



type KioskPageData struct {
	Version     string
	Lang        string
	ActiveEmpID int64
	Employees   []models.Employee
	Cards       []models.KioskCard
}


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

var (
	currentN8NWebhookURL = os.Getenv("N8N_WEBHOOK_URL")
	failedPINAttempts    = make(map[string]int)
	tokenLockouts        = make(map[string]time.Time)
)

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
		"formatAmount":   i18n.FormatAmount,
		"formatNumber":   i18n.FormatNumber,
		"formatPercent":  i18n.FormatPercent,
		"formatDate":     i18n.FormatDate,
		"formatDateTime": i18n.FormatDateTime,
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

	kioskTmpl, err := template.New("kiosk.html").Funcs(funcMap).ParseFiles("web/templates/kiosk.html")
	if err != nil {
		log.Fatalf("Failed to parse kiosk templates: %v", err)
	}

	slidesTmpl, err := template.ParseFiles("web/templates/slides.html")
	if err != nil {
		log.Fatalf("Failed to parse slides template: %v", err)
	}

	portalTmpl, err := template.New("portal.html").Funcs(funcMap).ParseFiles("web/templates/portal.html")
	if err != nil {
		log.Fatalf("Failed to parse portal templates: %v", err)
	}

	devStatusTmpl, err := template.New("layout.html").Funcs(funcMap).ParseFiles("web/templates/layout.html", "web/templates/dev_status.html")
	if err != nil {
		log.Fatalf("Failed to parse dev_status templates: %v", err)
	}

	masterdataTmpl, err := template.New("layout.html").Funcs(funcMap).ParseFiles("web/templates/layout.html", "web/templates/masterdata.html")
	if err != nil {
		log.Fatalf("Failed to parse masterdata templates: %v", err)
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

	// Helper to fetch kiosk page data
	fetchKioskData := func(lang string, activeEmpID int64) (*KioskPageData, error) {
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
			employees = append(employees, emp)
		}

		if activeEmpID <= 0 && len(employees) > 0 {
			activeEmpID = employees[0].ID
		}

		// Fetch active timer session for active employee
		activeSessionMap := make(map[int64]models.ActiveTimerSession)
		sessionRows, err := database.Query("SELECT id, employee_id, campaign_id, task_category, started_at FROM active_timer_sessions WHERE employee_id = ?", activeEmpID)
		if err == nil {
			defer sessionRows.Close()
			for sessionRows.Next() {
				var sess models.ActiveTimerSession
				if err := sessionRows.Scan(&sess.ID, &sess.EmployeeID, &sess.CampaignID, &sess.TaskCategory, &sess.StartedAt); err == nil {
					activeSessionMap[sess.CampaignID] = sess
				}
			}
		}

		// Fetch budget summaries for campaigns
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

		var cards []models.KioskCard
		for summaryRows.Next() {
			var card models.KioskCard
			if err := summaryRows.Scan(&card.CampaignID, &card.ClientName, &card.CampaignName, &card.TargetBudget, &card.ActualSpend, &card.HoursLogged); err != nil {
				return nil, err
			}
			if card.TargetBudget > 0 {
				card.UsagePercent = (card.ActualSpend / card.TargetBudget) * 100.0
			}
			if card.UsagePercent > 100.0 {
				card.Status = "danger"
			} else if card.UsagePercent >= 80.0 {
				card.Status = "warning"
			} else {
				card.Status = "ok"
			}

			if sess, exists := activeSessionMap[card.CampaignID]; exists {
				card.IsActive = true
				card.ActiveEmpID = activeEmpID
				card.StartedAtUnix = sess.StartedAt.Unix()
				card.TaskCategory = sess.TaskCategory
			} else {
				card.TaskCategory = "Content & Editing"
			}

			cards = append(cards, card)
		}

		return &KioskPageData{
			Version:     version,
			Lang:        lang,
			ActiveEmpID: activeEmpID,
			Employees:   employees,
			Cards:       cards,
		}, nil
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

	// Route: GET /tracker (800x480 Hardware Kiosk Touch Quick-Tracker)
	http.HandleFunc("/tracker", func(w http.ResponseWriter, r *http.Request) {
		lang := getLang(r)
		empID, _ := strconv.ParseInt(r.FormValue("employee_id"), 10, 64)

		data, err := fetchKioskData(lang, empID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Header.Get("HX-Request") == "true" {
			if err := kioskTmpl.ExecuteTemplate(w, "kiosk_cards", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if err := kioskTmpl.ExecuteTemplate(w, "kiosk.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Helper for rendering /dev/status page
	renderDevStatusPage := func(w http.ResponseWriter, r *http.Request, successMsg, errorMsg string) {
		lang := getLang(r)
		logs, _ := database.GetSecurityLogs(50)
		data := models.DevStatusData{
			Version:        version,
			Lang:           lang,
			DBJournalMode:  "WAL",
			SecurityStatus: "Active",
			N8NWebhookURL:  currentN8NWebhookURL,
			N8NActive:      strings.TrimSpace(currentN8NWebhookURL) != "",
			TotalAlerts:    len(logs),
			Logs:           logs,
			SuccessMsg:     successMsg,
			ErrorMsg:       errorMsg,
		}

		if r.Header.Get("HX-Request") == "true" {
			if err := devStatusTmpl.ExecuteTemplate(w, "content", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		type DevPageData struct {
			models.DevStatusData
			ActiveNav string
		}
		pageData := DevPageData{
			DevStatusData: data,
			ActiveNav:     "dev_status",
		}

		if err := devStatusTmpl.ExecuteTemplate(w, "layout.html", pageData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	// Route: GET /dev/status
	http.HandleFunc("/dev/status", func(w http.ResponseWriter, r *http.Request) {
		renderDevStatusPage(w, r, "", "")
	})

	// Route: POST /api/dev/update-webhook-url
	http.HandleFunc("/api/dev/update-webhook-url", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			currentN8NWebhookURL = strings.TrimSpace(r.FormValue("webhook_url"))
			renderDevStatusPage(w, r, "n8n Webhook Target URL updated successfully!", "")
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Route: POST /api/dev/simulate-brute-force
	http.HandleFunc("/api/dev/simulate-brute-force", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			database.LogSecurityEvent("INVALID_PIN", "ritter-sport-8821", "192.168.178.99", 3, "BLOCKED", "Simulated 3x PIN brute force attack on Ritter Sport portal")
			n8n.DispatchWebhook(currentN8NWebhookURL, n8n.WebhookPayload{
				Event:     "INVALID_PIN",
				Token:     "ritter-sport-8821",
				IP:        "192.168.178.99",
				Status:    "BLOCKED",
				Details:   "Simulated 3x PIN brute force attack",
				Timestamp: time.Now(),
			})
			renderDevStatusPage(w, r, "Simulated 3x PIN Brute-Force attack logged and dispatched to n8n!", "")
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Route: POST /api/dev/simulate-link-scan
	http.HandleFunc("/api/dev/simulate-link-scan", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			database.LogSecurityEvent("INVALID_LINK_SCAN", "/portal/c/suspicious-hex-scanner", "45.142.120.55", 1, "WARNING", "Simulated automated web crawler scanning invalid portal links")
			n8n.DispatchWebhook(currentN8NWebhookURL, n8n.WebhookPayload{
				Event:     "INVALID_LINK_SCAN",
				Token:     "/portal/c/suspicious-hex-scanner",
				IP:        "45.142.120.55",
				Status:    "WARNING",
				Details:   "Simulated invalid link scan",
				Timestamp: time.Now(),
			})
			renderDevStatusPage(w, r, "Simulated Invalid Link Scan logged and dispatched to n8n!", "")
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Route: POST /api/dev/simulate-budget-drift
	http.HandleFunc("/api/dev/simulate-budget-drift", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			database.LogSecurityEvent("BUDGET_DRIFT_ALERT", "Taycan GT Creator Experience", "system", 1, "BLOCKED", "Campaign target budget exceeded (115% usage)")
			n8n.DispatchWebhook(currentN8NWebhookURL, n8n.WebhookPayload{
				Event:     "BUDGET_DRIFT_ALERT",
				Token:     "Taycan GT Creator Experience",
				IP:        "system",
				Status:    "BLOCKED",
				Details:   "Campaign budget limit exceeded (115% usage)",
				Timestamp: time.Now(),
			})
			renderDevStatusPage(w, r, "Simulated Budget Drift alert logged and dispatched to n8n!", "")
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Helper for Master Data Page Data
	fetchMasterDataPageData := func(lang, tab, successMsg, errorMsg string) (*models.MasterDataPageData, error) {
		clients, err := database.GetAllClients()
		if err != nil {
			return nil, err
		}
		employees, err := database.GetAllEmployees()
		if err != nil {
			return nil, err
		}
		campaigns, err := database.GetAllCampaigns()
		if err != nil {
			return nil, err
		}
		if tab == "" {
			tab = "portal"
		}
		return &models.MasterDataPageData{
			Version:    version,
			Lang:       lang,
			ActiveNav:  "masterdata",
			ActiveTab:  tab,
			Clients:    clients,
			Employees:  employees,
			Campaigns:  campaigns,
			SuccessMsg: successMsg,
			ErrorMsg:   errorMsg,
		}, nil
	}

	renderMasterDataPage := func(w http.ResponseWriter, r *http.Request, tab, successMsg, errorMsg string) {
		lang := getLang(r)
		data, err := fetchMasterDataPageData(lang, tab, successMsg, errorMsg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := masterdataTmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	// Route: GET /masterdata
	http.HandleFunc("/masterdata", func(w http.ResponseWriter, r *http.Request) {
		tab := r.URL.Query().Get("tab")
		renderMasterDataPage(w, r, tab, "", "")
	})

	// Route: POST /api/masterdata/employee/save
	http.HandleFunc("/api/masterdata/employee/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		name := strings.TrimSpace(r.FormValue("name"))
		role := strings.TrimSpace(r.FormValue("role"))
		costRate, _ := strconv.ParseFloat(r.FormValue("cost_rate"), 64)
		billingRate, _ := strconv.ParseFloat(r.FormValue("billing_rate"), 64)
		tab := r.FormValue("active_tab")
		if tab == "" {
			tab = "employees"
		}

		if name == "" || role == "" {
			renderMasterDataPage(w, r, tab, "", "Name und Rolle dürfen nicht leer sein.")
			return
		}

		var err error
		if id > 0 {
			err = database.UpdateEmployee(id, name, role, billingRate, costRate, billingRate)
		} else {
			err = database.CreateEmployee(name, role, billingRate, costRate, billingRate)
		}
		if err != nil {
			renderMasterDataPage(w, r, tab, "", "Fehler beim Speichern: "+err.Error())
			return
		}
		renderMasterDataPage(w, r, tab, "Mitarbeiter erfolgreich gespeichert!", "")
	})

	// Route: POST /api/masterdata/employee/delete
	http.HandleFunc("/api/masterdata/employee/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if id > 0 {
			if err := database.DeleteEmployee(id); err != nil {
				renderMasterDataPage(w, r, "employees", "", "Fehler beim Löschen: "+err.Error())
				return
			}
		}
		renderMasterDataPage(w, r, "employees", "Mitarbeiter gelöscht!", "")
	})

	// Route: POST /api/masterdata/client/save
	http.HandleFunc("/api/masterdata/client/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		name := strings.TrimSpace(r.FormValue("name"))
		portalToken := strings.TrimSpace(r.FormValue("portal_token"))
		pinCode := strings.TrimSpace(r.FormValue("pin_code"))
		tab := r.FormValue("active_tab")
		if tab == "" {
			tab = "clients"
		}

		if name == "" || portalToken == "" || len(pinCode) != 4 {
			renderMasterDataPage(w, r, tab, "", "Kundenname, Portal-Token und ein 4-stelliger PIN-Code sind erforderlich.")
			return
		}

		var err error
		if id > 0 {
			err = database.UpdateClient(id, name, portalToken, pinCode)
		} else {
			err = database.CreateClient(name, portalToken, pinCode)
		}
		if err != nil {
			renderMasterDataPage(w, r, tab, "", "Fehler beim Speichern: "+err.Error())
			return
		}
		renderMasterDataPage(w, r, tab, "Kunde & PIN erfolgreich gespeichert!", "")
	})

	// Route: POST /api/masterdata/client/delete
	http.HandleFunc("/api/masterdata/client/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if id > 0 {
			if err := database.DeleteClient(id); err != nil {
				renderMasterDataPage(w, r, "clients", "", "Fehler beim Löschen: "+err.Error())
				return
			}
		}
		renderMasterDataPage(w, r, "clients", "Kunde gelöscht!", "")
	})

	// Route: POST /api/masterdata/campaign/save
	http.HandleFunc("/api/masterdata/campaign/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		clientID, _ := strconv.ParseInt(r.FormValue("client_id"), 10, 64)
		name := strings.TrimSpace(r.FormValue("name"))
		targetBudget, _ := strconv.ParseFloat(r.FormValue("target_budget"), 64)
		tab := r.FormValue("active_tab")
		if tab == "" {
			tab = "campaigns"
		}

		if clientID <= 0 || name == "" || targetBudget <= 0 {
			renderMasterDataPage(w, r, tab, "", "Kunde, Name und ein gültiges Zielbudget sind erforderlich.")
			return
		}

		var err error
		if id > 0 {
			err = database.UpdateCampaign(id, clientID, name, targetBudget)
		} else {
			err = database.CreateCampaign(clientID, name, targetBudget)
		}
		if err != nil {
			renderMasterDataPage(w, r, tab, "", "Fehler beim Speichern: "+err.Error())
			return
		}
		renderMasterDataPage(w, r, tab, "Kampagne erfolgreich gespeichert!", "")
	})

	// Route: POST /api/masterdata/campaign/delete
	http.HandleFunc("/api/masterdata/campaign/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if id > 0 {
			if err := database.DeleteCampaign(id); err != nil {
				renderMasterDataPage(w, r, "campaigns", "", "Fehler beim Löschen: "+err.Error())
				return
			}
		}
		renderMasterDataPage(w, r, "campaigns", "Kampagne gelöscht!", "")
	})

	// Route: Client Portal (/portal/c/{token})

	http.HandleFunc("/portal/c/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/portal/c/")
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) == 0 || parts[0] == "" {
			http.NotFound(w, r)
			return
		}

		token := parts[0]
		action := ""
		if len(parts) > 1 {
			action = parts[1]
		}

		client, err := database.GetClientByToken(token)
		if err != nil || client == nil {
			database.LogSecurityEvent("INVALID_LINK_SCAN", r.URL.Path, r.RemoteAddr, 1, "WARNING", "Automated web crawler or user scanned non-existent portal URL")
			n8n.DispatchWebhook(currentN8NWebhookURL, n8n.WebhookPayload{
				Event:     "INVALID_LINK_SCAN",
				Token:     r.URL.Path,
				IP:        r.RemoteAddr,
				Status:    "WARNING",
				Details:   "Scan of non-existent portal token URL",
				Timestamp: time.Now(),
			})
			http.NotFound(w, r)
			return
		}

		// Handle Logout
		if action == "logout" && r.Method == http.MethodPost {
			http.SetCookie(w, &http.Cookie{
				Name:     "portal_auth_" + token,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
			http.Redirect(w, r, "/portal/c/"+token, http.StatusSeeOther)
			return
		}

		lang := r.URL.Query().Get("lang")
		if lang == "" {
			lang = getLang(r)
		}

		// Check Lockout
		if until, locked := tokenLockouts[token]; locked && time.Now().Before(until) {
			data := models.ClientPortalData{
				Version:       version,
				Lang:          lang,
				Client:        *client,
				Authenticated: false,
				ErrorMsg:      i18nMgr.T(lang, "portal_locked_out"),
			}
			w.WriteHeader(http.StatusForbidden)
			portalTmpl.ExecuteTemplate(w, "portal.html", data)
			return
		}

		cookie, err := r.Cookie("portal_auth_" + token)
		authenticated := (err == nil && cookie.Value == "true")

		// Handle PIN Authentication
		if action == "auth" && r.Method == http.MethodPost {
			inputPIN := r.FormValue("pin")
			if inputPIN == client.PinCode {
				delete(failedPINAttempts, token)
				delete(tokenLockouts, token)
				authenticated = true
				http.SetCookie(w, &http.Cookie{
					Name:     "portal_auth_" + token,
					Value:    "true",
					Path:     "/",
					Expires:  time.Now().Add(24 * time.Hour),
					HttpOnly: true,
				})
				http.Redirect(w, r, "/portal/c/"+token, http.StatusSeeOther)
				return
			}

			// Invalid PIN
			failedPINAttempts[token]++
			attempts := failedPINAttempts[token]
			status := "WARNING"
			errMsg := i18nMgr.T(lang, "portal_pin_invalid")
			if attempts >= 3 {
				tokenLockouts[token] = time.Now().Add(15 * time.Minute)
				status = "BLOCKED"
				errMsg = i18nMgr.T(lang, "portal_locked_out")
			}

			database.LogSecurityEvent("INVALID_PIN", token, r.RemoteAddr, attempts, status, fmt.Sprintf("Incorrect PIN attempt (%d/3)", attempts))
			n8n.DispatchWebhook(currentN8NWebhookURL, n8n.WebhookPayload{
				Event:     "INVALID_PIN",
				Token:     token,
				IP:        r.RemoteAddr,
				Status:    status,
				Details:   fmt.Sprintf("Failed PIN attempt (%d/3)", attempts),
				Timestamp: time.Now(),
			})

			data := models.ClientPortalData{
				Version:       version,
				Lang:          lang,
				Client:        *client,
				Authenticated: false,
				ErrorMsg:      errMsg,
			}
			if attempts >= 3 {
				w.WriteHeader(http.StatusForbidden)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
			portalTmpl.ExecuteTemplate(w, "portal.html", data)
			return
		}


		var campaigns []models.CampaignBudgetSummary
		var assets []models.ContentAsset
		if authenticated {
			campaigns, _ = database.GetClientCampaignSummaries(client.ID)
			assets, _ = database.GetClientContentAssets(client.ID)
		}

		data := models.ClientPortalData{
			Version:       version,
			Lang:          lang,
			Client:        *client,
			Campaigns:     campaigns,
			ContentAssets: assets,
			Authenticated: authenticated,
		}

		if err := portalTmpl.ExecuteTemplate(w, "portal.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Helper for stopping active timer and logging 15-min billing rounded time
	stopActiveTimerSession := func(empID int64) error {
		var sessionID, campID int64
		var category string
		var startedAt time.Time

		err := database.QueryRow(
			"SELECT id, campaign_id, task_category, started_at FROM active_timer_sessions WHERE employee_id = ?",
			empID,
		).Scan(&sessionID, &campID, &category, &startedAt)

		if err == nil && sessionID > 0 {
			durationMinutes := time.Since(startedAt).Minutes()
			// 15-minute agency billing interval rounding (min 0.25h)
			roundedHours := math.Max(0.25, math.Ceil(durationMinutes/15.0)*0.25)
			description := category + " (Kiosk Quick-Track)"

			tx, err := database.Begin()
			if err != nil {
				return err
			}
			defer tx.Rollback()

			if _, err := tx.Exec(
				"INSERT INTO time_logs (campaign_id, employee_id, hours, description, logged_at) VALUES (?, ?, ?, ?, ?)",
				campID, empID, roundedHours, description, time.Now(),
			); err != nil {
				return err
			}

			if _, err := tx.Exec("DELETE FROM active_timer_sessions WHERE id = ?", sessionID); err != nil {
				return err
			}

			return tx.Commit()
		}
		return nil
	}

	// Route: POST /api/kiosk/start
	http.HandleFunc("/api/kiosk/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		empID, _ := strconv.ParseInt(r.FormValue("employee_id"), 10, 64)
		campID, _ := strconv.ParseInt(r.FormValue("campaign_id"), 10, 64)
		category := r.FormValue("task_category")
		if category == "" {
			category = "Content & Editing"
		}

		if empID > 0 && campID > 0 {
			// Auto-stop any existing active timer for this employee
			_ = stopActiveTimerSession(empID)

			// Start new active timer session
			_, err := database.Exec(
				"INSERT INTO active_timer_sessions (employee_id, campaign_id, task_category, started_at) VALUES (?, ?, ?, ?)",
				empID, campID, category, time.Now(),
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		lang := getLang(r)
		data, err := fetchKioskData(lang, empID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Header.Get("HX-Request") == "true" {
			if err := kioskTmpl.ExecuteTemplate(w, "kiosk_cards", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		http.Redirect(w, r, "/tracker?employee_id="+strconv.FormatInt(empID, 10), http.StatusSeeOther)
	})

	// Route: POST /api/kiosk/stop
	http.HandleFunc("/api/kiosk/stop", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		empID, _ := strconv.ParseInt(r.FormValue("employee_id"), 10, 64)

		if empID > 0 {
			if err := stopActiveTimerSession(empID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		lang := getLang(r)
		data, err := fetchKioskData(lang, empID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Header.Get("HX-Request") == "true" {
			if err := kioskTmpl.ExecuteTemplate(w, "kiosk_cards", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		http.Redirect(w, r, "/tracker?employee_id="+strconv.FormatInt(empID, 10), http.StatusSeeOther)
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
