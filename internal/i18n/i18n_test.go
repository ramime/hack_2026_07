package i18n

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestI18nTranslations(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	localesDir := filepath.Join(filepath.Dir(filename), "../../locales")

	mgr, err := NewManager(localesDir)
	if err != nil {
		t.Fatalf("Failed to initialize i18n manager: %v", err)
	}

	deTitle := mgr.T("de", "app_title")
	if deTitle != "AgencyPulse" {
		t.Errorf("Expected AgencyPulse, got %s", deTitle)
	}

	deHeader := mgr.T("de", "log_time_heading")
	if deHeader != "Arbeitszeit erfassen" {
		t.Errorf("Expected 'Arbeitszeit erfassen', got %s", deHeader)
	}

	enHeader := mgr.T("en", "log_time_heading")
	if enHeader != "Log Working Hours" {
		t.Errorf("Expected 'Log Working Hours', got %s", enHeader)
	}
}

func TestFormatNumberAndAmount(t *testing.T) {
	if got := FormatNumber("de", 1234567.89, 2); got != "1.234.567,89" {
		t.Errorf("FormatNumber de failed: got %s, want 1.234.567,89", got)
	}
	if got := FormatNumber("en", 1234567.89, 2); got != "1,234,567.89" {
		t.Errorf("FormatNumber en failed: got %s, want 1,234,567.89", got)
	}

	if got := FormatAmount("de", 2470.0, 0); got != "2.470 €" {
		t.Errorf("FormatAmount de failed: got %s, want 2.470 €", got)
	}
	if got := FormatAmount("en", 2470.0, 0); got != "€2,470" {
		t.Errorf("FormatAmount en failed: got %s, want €2,470", got)
	}

	if got := FormatPercent("de", 85.5, 1); got != "85,5 %" {
		t.Errorf("FormatPercent de failed: got %s, want 85,5 %%", got)
	}
	if got := FormatPercent("en", 85.5, 1); got != "85.5%" {
		t.Errorf("FormatPercent en failed: got %s, want 85.5%%", got)
	}
}

