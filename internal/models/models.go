package models

import "time"

type Client struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
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
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`        // dev, design, PM, lead
	HourlyRate float64   `json:"hourly_rate"` // in EUR
	CreatedAt  time.Time `json:"created_at"`
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

