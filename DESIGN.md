# DESIGN.md — AgencyPulse Design System

Visual design tokens, color palette, glassmorphism specs, and 800x480 Kiosk responsive breakpoints.

---

## 🎨 Color Palette & Design Tokens

```css
:root {
  --bg-dark: #0f172a;               /* Deep slate dark background */
  --bg-card: rgba(30, 41, 59, 0.7);  /* Translucent card background */
  --bg-card-hover: rgba(51, 65, 85, 0.8);
  --border-color: rgba(255, 255, 255, 0.1);
  --text-main: #f8fafc;             /* Primary text */
  --text-muted: #94a3b8;            /* Secondary text */
  
  --accent-primary: #6366f1;        /* Indigo accent */
  --accent-gradient: linear-gradient(135deg, #6366f1 0%, #a855f7 100%);
  --accent-glow: rgba(99, 102, 241, 0.25);
  
  --status-ok: #10b981;             /* Green (In Budget) */
  --status-warning: #f59e0b;        /* Yellow (Warning >80%) */
  --status-danger: #ef4444;         /* Red (Out of Control >100%) */

  --font-family: 'Plus Jakarta Sans', system-ui, sans-serif;
  --radius-lg: 16px;
  --radius-md: 12px;
  --radius-sm: 8px;
}
```

---

## 📱 800x480 Kiosk Breakpoints

- Target resolution: 800 x 480 pixels (7" Raspberry Pi Touchscreen Display).
- Responsive breakpoint `@media (max-width: 800px) and (max-height: 480px)`:
  - Header height minimized to `48px`.
  - Grid layout: 2x2 or 3x2 large touch cards.
  - Digital stopwatch font size: `2.2rem` bold.
  - Touch button min-size: `56px x 56px` for touch accuracy.
