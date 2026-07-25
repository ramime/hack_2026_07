package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
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
