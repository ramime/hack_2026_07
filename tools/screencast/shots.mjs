#!/usr/bin/env node
/**
 * Capture domain screenshots for /slides into web/static/pitch/
 * Usage: BASE_URL=http://localhost:8084 node shots.mjs
 */
import { chromium } from "playwright";
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const baseURL = process.env.BASE_URL || "http://localhost:8084";
const __dirname = path.dirname(fileURLToPath(import.meta.url));
const outDir = path.resolve(__dirname, "../../web/static/pitch");

fs.mkdirSync(outDir, { recursive: true });

const shots = [
  { name: "heatmap.png", path: "/teamlead", prep: null },
  { name: "executive.png", path: "/executive", prep: null },
  { name: "kiosk.png", path: "/tracker", prep: null },
  {
    name: "portal.png",
    path: "/portal/c/ritter-sport-8821",
    prep: async (page) => {
      const form = page.locator("[data-testid='portal-pin-form']");
      if (await form.count()) {
        await page.fill("[data-testid='portal-pin-input']", "1234");
        await Promise.all([
          page.waitForLoadState("networkidle").catch(() => {}),
          page.click("[data-testid='portal-submit-pin']"),
        ]);
        await page.locator("[data-testid='portal-content']").waitFor({ state: "visible", timeout: 10000 });
      }
    },
  },
  {
    name: "security.png",
    path: "/dev/status",
    prep: async (page) => {
      const btn = page.locator("[data-testid='sim-brute-force-button']");
      if (await btn.count()) {
        await btn.click();
        await page.waitForTimeout(800);
      }
    },
  },
  {
    name: "masterdata.png",
    path: "/masterdata?tab=portal",
    prep: async (page) => {
      await page.locator("[data-testid='portal-views-section']").waitFor({ state: "visible", timeout: 10000 }).catch(() => {});
      await page.waitForTimeout(300);
    },
  },
];

const browser = await chromium.launch({ headless: true });
const context = await browser.newContext({
  viewport: { width: 1280, height: 720 },
  locale: "en-US",
});
await context.addCookies([{ name: "lang", value: "en", url: baseURL }]);
const page = await context.newPage();

// Ensure seed data for portal tokens / heatmap narrative
await page.request.post(`${baseURL}/api/reset-demo-data`).catch(() => {});

for (const shot of shots) {
  const url = new URL(shot.path, baseURL).toString();
  console.log(`→ ${shot.name} ${url}`);
  await page.goto(url, { waitUntil: "networkidle" });
  if (shot.prep) await shot.prep(page);
  await page.waitForTimeout(400);
  const dest = path.join(outDir, shot.name);
  await page.screenshot({ path: dest, type: "png" });
  console.log(`  wrote ${dest}`);
}

await browser.close();
console.log("Done.");
