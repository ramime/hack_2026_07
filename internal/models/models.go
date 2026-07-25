package models

import "time"

type Client struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	PortalToken string    `json:"portal_token"`
	PinCode     string    `json:"pin_code"`
	CreatedAt   time.Time `json:"created_at"`
}

type Campaign struct {
	ID           int64     `json:"id"`
	ClientID     int64     `json:"client_id"`
	ClientName   string    `json:"client_name,omitempty"`
	Name         string    `json:"name"`
	TargetBudget float64   `json:"target_budget"` // in EUR
	ActualBudget float64   `json:"actual_budget"` // in EUR (calculated or cached)
	StatusColor  string    `json:"status_color"`  // green, yellow, red
	CreatedAt    time.Time `json:"created_at"`
}

type Employee struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Role        string    `json:"role"`        // dev, design, PM, lead
	HourlyRate  float64   `json:"hourly_rate"` // legacy field/billing rate
	CostRate    float64   `json:"cost_rate"`   // internal labor cost rate in EUR/h
	BillingRate float64   `json:"billing_rate"`// client billing rate in EUR/h
	CreatedAt   time.Time `json:"created_at"`
}

type TimeLog struct {
	ID           int64     `json:"id"`
	CampaignID   int64     `json:"campaign_id"`
	CampaignName string    `json:"campaign_name,omitempty"`
	EmployeeID   int64     `json:"employee_id"`
	EmployeeName string    `json:"employee_name,omitempty"`
	Hours        float64   `json:"hours"`
	Description  string    `json:"description"`
	LoggedAt     time.Time `json:"logged_at"`
}

type CampaignBudgetSummary struct {
	CampaignID   int64   `json:"campaign_id"`
	ClientName   string  `json:"client_name"`
	CampaignName string  `json:"campaign_name"`
	TargetBudget float64 `json:"target_budget"`
	ActualSpend  float64 `json:"actual_spend"`
	HoursLogged  float64 `json:"hours_logged"`
	UsagePercent float64 `json:"usage_percent"`
	Status       string  `json:"status"` // ok, warning, danger
}

type ExecutiveKPIs struct {
	TotalRevenue   float64 `json:"total_revenue"`
	TotalLaborCost float64 `json:"total_labor_cost"`
	NetProfit      float64 `json:"net_profit"`
	AgencyMargin   float64 `json:"agency_margin"`
	AtRiskCount    int     `json:"at_risk_count"`
}

type ClientProfitability struct {
	ClientID      int64   `json:"client_id"`
	ClientName    string  `json:"client_name"`
	CampaignCount int     `json:"campaign_count"`
	TotalHours    float64 `json:"total_hours"`
	BilledRevenue float64 `json:"billed_revenue"`
	LaborCost     float64 `json:"labor_cost"`
	NetProfit     float64 `json:"net_profit"`
	MarginPercent float64 `json:"margin_percent"`
}

type EmployeeEfficiency struct {
	EmployeeID      int64   `json:"employee_id"`
	EmployeeName    string  `json:"employee_name"`
	Role            string  `json:"role"`
	CostRate        float64 `json:"cost_rate"`
	BillingRate     float64 `json:"billing_rate"`
	HoursLogged     float64 `json:"hours_logged"`
	BilledRevenue   float64 `json:"billed_revenue"`
	LaborCost       float64 `json:"labor_cost"`
	NetContribution float64 `json:"net_contribution"`
	MarginPercent   float64 `json:"margin_percent"`
}

type ActiveTimerSession struct {
	ID           int64     `json:"id"`
	EmployeeID   int64     `json:"employee_id"`
	CampaignID   int64     `json:"campaign_id"`
	TaskCategory string    `json:"task_category"`
	StartedAt    time.Time `json:"started_at"`
}

type KioskCard struct {
	CampaignID    int64     `json:"campaign_id"`
	ClientName    string    `json:"client_name"`
	CampaignName  string    `json:"campaign_name"`
	TargetBudget  float64   `json:"target_budget"`
	ActualSpend   float64   `json:"actual_spend"`
	HoursLogged   float64   `json:"hours_logged"`
	UsagePercent  float64   `json:"usage_percent"`
	Status        string    `json:"status"` // ok, warning, danger
	IsActive      bool      `json:"is_active"`
	ActiveEmpID   int64     `json:"active_emp_id"`
	StartedAtUnix int64     `json:"started_at_unix"`
	TaskCategory  string    `json:"task_category"`
}

type ContentAsset struct {
	ID           int64     `json:"id"`
	CampaignID   int64     `json:"campaign_id"`
	CampaignName string    `json:"campaign_name,omitempty"`
	Title        string    `json:"title"`
	AssetType    string    `json:"asset_type"` // TikTok Video, Instagram Reel, 3D Motion Graphic, etc.
	Status       string    `json:"status"`     // Delivered, Approved, In Review
	PreviewURL   string    `json:"preview_url"`
	DeliveredAt  time.Time `json:"delivered_at"`
}

type ClientPortalData struct {
	Version       string                  `json:"version"`
	Lang          string                  `json:"lang"`
	Client        Client                  `json:"client"`
	Campaigns     []CampaignBudgetSummary `json:"campaigns"`
	ContentAssets []ContentAsset          `json:"content_assets"`
	Authenticated bool                    `json:"authenticated"`
	ErrorMsg      string                  `json:"error_msg,omitempty"`
}




