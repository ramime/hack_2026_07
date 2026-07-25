package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type Manager struct {
	mu           sync.RWMutex
	translations map[string]map[string]string
}

func NewManager(localesDir string) (*Manager, error) {
	m := &Manager{
		translations: make(map[string]map[string]string),
	}

	langs := []string{"de", "en"}
	for _, lang := range langs {
		filePath := fmt.Sprintf("%s/%s.json", localesDir, lang)
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read locale file %s: %w", filePath, err)
		}

		var dict map[string]string
		if err := json.Unmarshal(data, &dict); err != nil {
			return nil, fmt.Errorf("failed to parse locale file %s: %w", filePath, err)
		}

		m.translations[lang] = dict
	}

	return m, nil
}

func (m *Manager) T(lang, key string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if lang != "de" && lang != "en" {
		lang = "de"
	}

	if dict, exists := m.translations[lang]; exists {
		if val, found := dict[key]; found {
			return val
		}
	}

	// Fallback to German or key itself
	if dict, exists := m.translations["de"]; exists {
		if val, found := dict[key]; found {
			return val
		}
	}

	return key
}

// FormatNumber formats numbers with thousand separators and decimal places according to locale.
// German ("de"): 1.234,56
// English ("en"): 1,234.56
func FormatNumber(lang string, val interface{}, decimals int) string {
	var f float64
	switch v := val.(type) {
	case float64:
		f = v
	case float32:
		f = float64(v)
	case int:
		f = float64(v)
	case int64:
		f = float64(v)
	case int32:
		f = float64(v)
	case uint:
		f = float64(v)
	case uint64:
		f = float64(v)
	default:
		return fmt.Sprintf("%v", val)
	}

	neg := f < 0
	if neg {
		f = -f
	}

	formatStr := fmt.Sprintf("%%.%df", decimals)
	str := fmt.Sprintf(formatStr, f)

	parts := strings.Split(str, ".")
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = parts[1]
	}

	var result []rune
	n := len(intPart)
	for i, r := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			if lang == "de" {
				result = append(result, '.')
			} else {
				result = append(result, ',')
			}
		}
		result = append(result, r)
	}

	resStr := string(result)
	if neg {
		resStr = "-" + resStr
	}

	if decimals > 0 {
		if lang == "de" {
			resStr = resStr + "," + decPart
		} else {
			resStr = resStr + "." + decPart
		}
	}
	return resStr
}

// FormatAmount formats monetary amounts with thousand separators and currency symbol according to locale.
// German ("de"): 1.234,56 €
// English ("en"): €1,234.56
func FormatAmount(lang string, amount float64, decimals int) string {
	numStr := FormatNumber(lang, amount, decimals)
	if lang == "de" {
		return numStr + " €"
	}
	if strings.HasPrefix(numStr, "-") {
		return "-€" + strings.TrimPrefix(numStr, "-")
	}
	return "€" + numStr
}

// FormatPercent formats percentages according to locale.
// German ("de"): 85,5 %
// English ("en"): 85.5%
func FormatPercent(lang string, val float64, decimals int) string {
	numStr := FormatNumber(lang, val, decimals)
	if lang == "de" {
		return numStr + " %"
	}
	return numStr + "%"
}

// FormatDate formats dates according to locale.
// German ("de"): 25.07.2026
// English ("en"): Jul 25, 2026
func FormatDate(lang string, t time.Time) string {
	if lang == "de" {
		return t.Format("02.01.2006")
	}
	return t.Format("Jan 02, 2006")
}

// FormatDateTime formats date and time according to locale.
// German ("de"): 25.07.2026 14:01
// English ("en"): Jan 02, 2006 14:01
func FormatDateTime(lang string, t time.Time) string {
	if lang == "de" {
		return t.Format("02.01.2006 15:04")
	}
	return t.Format("Jan 02, 2006 15:04")
}

