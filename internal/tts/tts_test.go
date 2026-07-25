package tts

import "testing"

func TestSynthesizeSkip(t *testing.T) {
	if err := Synthesize(ProviderSkip, "hello", "/tmp/should-not-exist-pitch.mp3"); err != nil {
		t.Fatal(err)
	}
	if err := Synthesize("", "hello", "/tmp/should-not-exist-pitch.mp3"); err != nil {
		t.Fatal(err)
	}
}

func TestSynthesizeUnknown(t *testing.T) {
	if err := Synthesize("nope", "x", "y.mp3"); err == nil {
		t.Fatal("expected error")
	}
}
