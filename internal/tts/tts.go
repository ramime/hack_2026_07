package tts

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Provider names supported by the pitch media CLI.
const (
	ProviderSkip       = "skip"
	ProviderEdge       = "edge"
	ProviderElevenLabs = "elevenlabs"
)

// Synthesize writes narration audio for a scene to outPath.
// ProviderSkip is a no-op and returns nil without creating a file.
func Synthesize(provider, text, outPath string) error {
	switch provider {
	case "", ProviderSkip:
		return nil
	case ProviderEdge:
		return synthesizeEdge(text, outPath)
	case ProviderElevenLabs:
		return synthesizeElevenLabs(text, outPath)
	default:
		return fmt.Errorf("unknown tts provider %q (use skip|edge|elevenlabs)", provider)
	}
}

func synthesizeEdge(text, outPath string) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	// edge-tts CLI: pip install edge-tts
	cmd := exec.Command("edge-tts", "--text", text, "--write-media", outPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("edge-tts failed (is edge-tts installed?): %w", err)
	}
	return nil
}

func synthesizeElevenLabs(text, outPath string) error {
	key := os.Getenv("ELEVENLABS_API_KEY")
	if key == "" {
		return fmt.Errorf("ELEVENLABS_API_KEY is not set")
	}
	_ = text
	_ = outPath
	// Placeholder for a later ElevenLabs HTTP client shared with in-app briefing.
	return fmt.Errorf("elevenlabs TTS is not implemented yet; use -tts skip or -tts edge")
}
