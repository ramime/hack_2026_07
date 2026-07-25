package pitch

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ScenesFile is the pitch/scenes.yaml document.
type ScenesFile struct {
	Version  string   `yaml:"version" json:"version"`
	Lang     string   `yaml:"lang" json:"lang"`
	Viewport Viewport `yaml:"viewport" json:"viewport"`
	Scenes   []Scene  `yaml:"scenes" json:"scenes"`
}

type Viewport struct {
	Width  int `yaml:"width" json:"width"`
	Height int `yaml:"height" json:"height"`
}

type Scene struct {
	ID          string `yaml:"id" json:"id"`
	Path        string `yaml:"path" json:"path"`
	MaxSeconds  int    `yaml:"max_seconds" json:"max_seconds"`
	NarrationEN string `yaml:"narration_en" json:"narration_en"`
	Steps       []Step `yaml:"steps" json:"steps"`
}

type Step struct {
	Action        string `yaml:"action" json:"action"`
	TestID        string `yaml:"testid" json:"testid,omitempty"`
	LabelContains string `yaml:"label_contains" json:"label_contains,omitempty"`
	Value         string `yaml:"value" json:"value,omitempty"`
	MS            int    `yaml:"ms" json:"ms,omitempty"`
}

// LoadScenes reads pitch/scenes.yaml.
func LoadScenes(path string) (*ScenesFile, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc ScenesFile
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	if doc.Viewport.Width == 0 {
		doc.Viewport.Width = 1280
	}
	if doc.Viewport.Height == 0 {
		doc.Viewport.Height = 720
	}
	return &doc, nil
}
