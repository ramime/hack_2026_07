package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"agencypulse/internal/pitch"
	"agencypulse/internal/tts"
)

func main() {
	base := flag.String("base", "http://localhost:8084", "AgencyPulse base URL (app must already be running)")
	out := flag.String("out", "artifacts", "Output directory for video/audio")
	ttsProvider := flag.String("tts", tts.ProviderSkip, "TTS provider: skip|edge|elevenlabs (default skip)")
	scenesPath := flag.String("scenes", "pitch/scenes.yaml", "Scenes file")
	dryRun := flag.Bool("dry-run", false, "Validate scenes and exit without capture")
	skipReset := flag.Bool("skip-reset", false, "Skip POST /api/reset-demo-data")
	maxSeconds := flag.Int("max-seconds", 120, "Trim final video to this many seconds")
	flag.Parse()

	doc, err := pitch.LoadScenes(*scenesPath)
	if err != nil {
		fatalf("load scenes: %v", err)
	}
	fmt.Printf("Loaded %d scenes from %s (version %s)\n", len(doc.Scenes), *scenesPath, doc.Version)

	if *dryRun {
		for _, s := range doc.Scenes {
			fmt.Printf("  - %s path=%s narration=%q steps=%d\n", s.ID, s.Path, trunc(s.NarrationEN, 60), len(s.Steps))
		}
		fmt.Println("Dry run OK.")
		return
	}

	if err := os.MkdirAll(*out, 0o755); err != nil {
		fatalf("mkdir out: %v", err)
	}

	version, err := healthCheck(*base)
	if err != nil {
		fatalf("health check failed (%s): %v\nStart the app first: go run ./cmd/agencypulse", *base, err)
	}
	fmt.Printf("Health OK — version %s\n", version)

	if !*skipReset {
		if err := postForm(*base+"/api/reset-demo-data", nil); err != nil {
			fatalf("reset demo data: %v", err)
		}
		fmt.Println("Demo data reset.")
	}

	if err := postForm(*base+"/api/set-language", url.Values{"lang": {"en"}}); err != nil {
		fmt.Printf("Warning: set-language failed: %v (continuing; capture sets lang cookie)\n", err)
	}

	runConfigPath := filepath.Join(*out, "run-config.json")
	if err := writeJSON(runConfigPath, doc); err != nil {
		fatalf("write run config: %v", err)
	}

	rawVideo := filepath.Join(*out, "raw.webm")
	if err := runCapture(*base, runConfigPath, *out, "raw.webm"); err != nil {
		fatalf("capture: %v", err)
	}

	audioDir := filepath.Join(*out, "audio")
	var audioFiles []string
	if *ttsProvider != tts.ProviderSkip && *ttsProvider != "" {
		if err := os.MkdirAll(audioDir, 0o755); err != nil {
			fatalf("mkdir audio: %v", err)
		}
		for _, scene := range doc.Scenes {
			if strings.TrimSpace(scene.NarrationEN) == "" {
				continue
			}
			ap := filepath.Join(audioDir, scene.ID+".mp3")
			if err := tts.Synthesize(*ttsProvider, scene.NarrationEN, ap); err != nil {
				fatalf("tts %s: %v", scene.ID, err)
			}
			if _, err := os.Stat(ap); err == nil {
				audioFiles = append(audioFiles, ap)
			}
		}
	} else {
		fmt.Println("TTS skipped (default) — captions are burned into the screencast.")
	}

	safeVer := strings.ReplaceAll(version, "/", "-")
	finalPath := filepath.Join(*out, fmt.Sprintf("agencypulse-%s.mp4", safeVer))
	if err := muxVideo(rawVideo, audioFiles, finalPath, *maxSeconds); err != nil {
		fatalf("ffmpeg: %v", err)
	}

	fmt.Println()
	fmt.Println("Done.")
	fmt.Printf("  Video:  %s\n", finalPath)
	fmt.Printf("  Slides: %s/slides\n", strings.TrimRight(*base, "/"))
	fmt.Println("Next: open the MP4, tweak pitch/scenes.yaml if needed, re-run.")
	fmt.Println("Optional: upload to Loom/Drive and set the link in SUBMISSION.md")
}

func healthCheck(base string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(strings.TrimRight(base, "/") + "/api/health")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}
	var body struct {
		OK      bool   `json:"ok"`
		Version string `json:"version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	if !body.OK {
		return "", fmt.Errorf("ok=false")
	}
	if body.Version == "" {
		body.Version = "unknown"
	}
	return body.Version, nil
}

func postForm(endpoint string, values url.Values) error {
	client := &http.Client{Timeout: 15 * time.Second}
	var body io.Reader
	if values != nil {
		body = strings.NewReader(values.Encode())
	}
	req, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return err
	}
	if values != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return nil
}

func writeJSON(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func runCapture(base, runConfig, outDir, videoName string) error {
	toolDir := "tools/screencast"
	if err := ensurePlaywright(toolDir); err != nil {
		return err
	}

	absConfig, err := filepath.Abs(runConfig)
	if err != nil {
		return err
	}
	absOut, err := filepath.Abs(outDir)
	if err != nil {
		return err
	}

	cmd := exec.Command("node", "capture.mjs")
	cmd.Dir = toolDir
	cmd.Env = append(os.Environ(),
		"BASE_URL="+base,
		"RUN_CONFIG="+absConfig,
		"OUT_DIR="+absOut,
		"VIDEO_NAME="+videoName,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ensurePlaywright(toolDir string) error {
	pwMod := filepath.Join(toolDir, "node_modules", "playwright")
	if _, err := os.Stat(pwMod); err != nil {
		fmt.Println("Installing Playwright npm package...")
		install := exec.Command("npm", "install")
		install.Dir = toolDir
		install.Stdout = os.Stdout
		install.Stderr = os.Stderr
		if err := install.Run(); err != nil {
			return fmt.Errorf("npm install: %w", err)
		}
	}

	// Always ensure the browser revision for this Playwright version exists
	// in the default user cache (~/.cache/ms-playwright). Harmless if already present.
	fmt.Println("Ensuring Playwright Chromium browser is installed...")
	browsers := exec.Command("npx", "playwright", "install", "chromium")
	browsers.Dir = toolDir
	browsers.Stdout = os.Stdout
	browsers.Stderr = os.Stderr
	// Do not inherit a custom PLAYWRIGHT_BROWSERS_PATH from agent/CI sandboxes
	// when the user runs locally — strip it so browsers land in the default cache.
	browsers.Env = filterEnv(os.Environ(), "PLAYWRIGHT_BROWSERS_PATH")
	if err := browsers.Run(); err != nil {
		return fmt.Errorf("playwright install chromium: %w\nTry manually: cd tools/screencast && npx playwright install chromium", err)
	}
	return nil
}

func filterEnv(env []string, dropKey string) []string {
	prefix := dropKey + "="
	out := make([]string, 0, len(env))
	for _, e := range env {
		if strings.HasPrefix(e, prefix) {
			continue
		}
		out = append(out, e)
	}
	return out
}

func muxVideo(rawWebm string, audioFiles []string, outMP4 string, maxSeconds int) error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH")
	}

	args := []string{"-y", "-i", rawWebm}
	// v1: captions are burned in; optional single concatenated audio later.
	// For skip TTS we just remux/transcode and trim.
	if len(audioFiles) == 1 {
		args = append(args, "-i", audioFiles[0], "-c:v", "libx264", "-pix_fmt", "yuv420p", "-c:a", "aac", "-shortest")
	} else if len(audioFiles) > 1 {
		// Concatenate audio with filter, keep simple: use first clip only for v1
		fmt.Println("Warning: multiple TTS clips — muxing first clip only in v1")
		args = append(args, "-i", audioFiles[0], "-c:v", "libx264", "-pix_fmt", "yuv420p", "-c:a", "aac", "-shortest")
	} else {
		args = append(args, "-c:v", "libx264", "-pix_fmt", "yuv420p", "-an")
	}
	if maxSeconds > 0 {
		args = append(args, "-t", fmt.Sprintf("%d", maxSeconds))
	}
	args = append(args, outMP4)

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "pitchmedia: "+format+"\n", args...)
	os.Exit(1)
}
