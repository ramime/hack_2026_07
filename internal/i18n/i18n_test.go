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
