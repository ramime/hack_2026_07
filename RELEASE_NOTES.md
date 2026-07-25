# RELEASE_NOTES.md — AgencyPulse v0.7.4

**Release Date:** 2026-07-25  
**Milestone:** Pitch Slides SVG Branding & Chromes (`v0.7.4`)

---

## 🎯 Release Summary

Release `v0.7.4` integrates vector SVG branding assets into the `/slides` presentation viewer and deck chromes.

### Key Highlights
- **SVG Branding in Presentation Deck**: Embedded `/static/logo.svg` directly into the hero slide of `pitch/slides.md`.
- **Slide Renderer Enhancements**: Added brand logo figure detection (`.slide-brand`) in `internal/pitch/slides.go` preventing caption duplication for logo images.
- **Unit Testing**: Added `TestParseSlidesBrandLogo` test suite in `internal/pitch/slides_test.go`.
- **Slide Deck Chrome**: Added `logo-icon.svg` badge to the footer navigation link in `web/templates/slides.html`.
