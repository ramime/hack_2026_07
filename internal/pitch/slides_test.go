package pitch

import (
	"strings"
	"testing"
)

func TestParseSlides(t *testing.T) {
	raw := `<!-- comment -->

# Title One

Hello world.

---

# Title Two

- Bullet A
- Bullet B
`
	slides := ParseSlides(raw)
	if len(slides) != 2 {
		t.Fatalf("got %d slides, want 2", len(slides))
	}
	if slides[0].Title != "Title One" {
		t.Errorf("title0 = %q", slides[0].Title)
	}
	if !strings.Contains(string(slides[0].HTML), "<p>Hello world.</p>") {
		t.Errorf("html0 = %q", slides[0].HTML)
	}
	if slides[1].Title != "Title Two" {
		t.Errorf("title1 = %q", slides[1].Title)
	}
	if !strings.Contains(string(slides[1].HTML), "<li>Bullet A</li>") {
		t.Errorf("html1 = %q", slides[1].HTML)
	}
}
