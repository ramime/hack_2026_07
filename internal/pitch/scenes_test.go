package pitch

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadScenes(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))
	doc, err := LoadScenes(filepath.Join(root, "pitch/scenes.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Scenes) < 1 {
		t.Fatal("expected scenes")
	}
	if doc.Scenes[0].NarrationEN == "" {
		t.Fatal("expected narration")
	}
}
