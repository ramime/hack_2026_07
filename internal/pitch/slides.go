package pitch

import (
	"html/template"
	"os"
	"strings"
)

// Slide is one deck page parsed from pitch/slides.md.
type Slide struct {
	Title string
	Body  string
	HTML  template.HTML // simple rendered body (paragraphs + ul)
}

// LoadSlides reads a markdown file split on --- separators.
func LoadSlides(path string) ([]Slide, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseSlides(string(raw)), nil
}

// ParseSlides splits markdown into slides.
func ParseSlides(raw string) []Slide {
	// Strip HTML comments
	for {
		start := strings.Index(raw, "<!--")
		if start < 0 {
			break
		}
		end := strings.Index(raw[start:], "-->")
		if end < 0 {
			raw = raw[:start]
			break
		}
		raw = raw[:start] + raw[start+end+3:]
	}

	parts := strings.Split(raw, "\n---\n")
	var slides []Slide
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		lines := strings.Split(part, "\n")
		title := ""
		var bodyLines []string
		for i, line := range lines {
			trim := strings.TrimSpace(line)
			if i == 0 && strings.HasPrefix(trim, "# ") {
				title = strings.TrimSpace(trim[2:])
				continue
			}
			bodyLines = append(bodyLines, line)
		}
		body := strings.TrimSpace(strings.Join(bodyLines, "\n"))
		slides = append(slides, Slide{
			Title: title,
			Body:  body,
			HTML:  template.HTML(renderBody(body)),
		})
	}
	return slides
}

func renderBody(body string) string {
	if body == "" {
		return ""
	}
	var b strings.Builder
	lines := strings.Split(body, "\n")
	inList := false
	flushList := func() {
		if inList {
			b.WriteString("</ul>")
			inList = false
		}
	}
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" {
			flushList()
			continue
		}
		if alt, src, ok := parseMarkdownImage(trim); ok {
			flushList()
			b.WriteString(`<figure class="slide-shot">`)
			b.WriteString(`<img src="`)
			b.WriteString(escapeHTML(src))
			b.WriteString(`" alt="`)
			b.WriteString(escapeHTML(alt))
			b.WriteString(`" loading="lazy">`)
			if alt != "" {
				b.WriteString(`<figcaption>`)
				b.WriteString(escapeHTML(alt))
				b.WriteString(`</figcaption>`)
			}
			b.WriteString(`</figure>`)
			continue
		}
		if strings.HasPrefix(trim, "- ") {
			if !inList {
				b.WriteString("<ul>")
				inList = true
			}
			b.WriteString("<li>")
			b.WriteString(inlineMarkdown(strings.TrimSpace(trim[2:])))
			b.WriteString("</li>")
			continue
		}
		flushList()
		if strings.HasPrefix(trim, "**") && strings.HasSuffix(trim, "**") && len(trim) > 4 && !strings.Contains(trim[2:len(trim)-2], "**") {
			b.WriteString("<p class=\"slide-label\"><strong>")
			b.WriteString(escapeHTML(trim[2 : len(trim)-2]))
			b.WriteString("</strong></p>")
			continue
		}
		b.WriteString("<p>")
		b.WriteString(inlineMarkdown(trim))
		b.WriteString("</p>")
	}
	flushList()
	return b.String()
}

// parseMarkdownImage parses ![alt](src).
func parseMarkdownImage(line string) (alt, src string, ok bool) {
	if !strings.HasPrefix(line, "![") {
		return "", "", false
	}
	rest := line[2:]
	closeAlt := strings.Index(rest, "](")
	if closeAlt < 0 {
		return "", "", false
	}
	alt = rest[:closeAlt]
	rest = rest[closeAlt+2:]
	if !strings.HasSuffix(rest, ")") {
		return "", "", false
	}
	src = strings.TrimSuffix(rest, ")")
	if src == "" {
		return "", "", false
	}
	return alt, src, true
}

func escapeHTML(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
	)
	return replacer.Replace(s)
}

// inlineMarkdown supports a minimal **bold** subset after HTML escaping.
func inlineMarkdown(s string) string {
	escaped := escapeHTML(s)
	var b strings.Builder
	for {
		start := strings.Index(escaped, "**")
		if start < 0 {
			b.WriteString(escaped)
			break
		}
		rest := escaped[start+2:]
		end := strings.Index(rest, "**")
		if end < 0 {
			b.WriteString(escaped)
			break
		}
		b.WriteString(escaped[:start])
		b.WriteString("<strong>")
		b.WriteString(rest[:end])
		b.WriteString("</strong>")
		escaped = rest[end+2:]
	}
	return b.String()
}
