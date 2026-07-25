#!/usr/bin/env node
/**
 * AgencyPulse screencast capture with burned-in caption overlay.
 * Invoked by cmd/pitchmedia with env:
 *   BASE_URL, RUN_CONFIG (path to JSON), OUT_DIR, VIDEO_NAME
 */
import { chromium } from "playwright";
import fs from "node:fs";
import path from "node:path";

const baseURL = process.env.BASE_URL || "http://localhost:8084";
const runConfigPath = process.env.RUN_CONFIG;
const outDir = process.env.OUT_DIR || "artifacts";
const videoName = process.env.VIDEO_NAME || "raw.webm";

if (!runConfigPath) {
  console.error("RUN_CONFIG env is required (JSON scenes file)");
  process.exit(1);
}

const config = JSON.parse(fs.readFileSync(runConfigPath, "utf8"));
const viewport = config.viewport || { width: 1280, height: 720 };
const scenes = config.scenes || [];

fs.mkdirSync(outDir, { recursive: true });
const videoDir = path.join(outDir, "pw-video");
fs.rmSync(videoDir, { recursive: true, force: true });
fs.mkdirSync(videoDir, { recursive: true });

// Solid caption plate so text stays readable on light/busy UI.
const CAPTION_STYLE = `
#ap-pitch-caption {
  position: fixed;
  left: 50%;
  bottom: 20px;
  transform: translateX(-50%);
  z-index: 2147483647;
  pointer-events: none;
  box-sizing: border-box;
  max-width: min(1100px, calc(100vw - 32px));
  padding: 12px 22px;
  background: rgba(0, 0, 0, 0.92);
  border: 1px solid rgba(255, 255, 255, 0.22);
  border-radius: 12px;
  box-shadow: 0 10px 28px rgba(0, 0, 0, 0.55);
  color: #ffffff;
  font: 600 20px/1.35 "DM Sans", system-ui, sans-serif;
  text-align: center;
  letter-spacing: 0.01em;
  text-shadow: none;
}
#ap-pitch-caption[data-empty="1"] { display: none; }
`;

async function ensureCaption(page, text) {
  await page.addStyleTag({ content: CAPTION_STYLE }).catch(() => {});
  await page.evaluate((caption) => {
    let el = document.getElementById("ap-pitch-caption");
    if (!el) {
      el = document.createElement("div");
      el.id = "ap-pitch-caption";
      document.documentElement.appendChild(el);
    }
    // Re-apply inline plate styles in case page CSS fights the injected sheet.
    el.style.cssText = [
      "position:fixed",
      "left:50%",
      "bottom:20px",
      "transform:translateX(-50%)",
      "z-index:2147483647",
      "pointer-events:none",
      "box-sizing:border-box",
      "max-width:min(1100px, calc(100vw - 32px))",
      "padding:12px 22px",
      "background:rgba(0,0,0,0.92)",
      "border:1px solid rgba(255,255,255,0.22)",
      "border-radius:12px",
      "box-shadow:0 10px 28px rgba(0,0,0,0.55)",
      "color:#ffffff",
      "font:600 20px/1.35 DM Sans, system-ui, sans-serif",
      "text-align:center",
      "letter-spacing:0.01em",
    ].join(";");
    const t = (caption || "").trim();
    el.textContent = t;
    el.setAttribute("data-empty", t ? "0" : "1");
    el.style.display = t ? "block" : "none";
  }, text || "");
}

async function runStep(page, step, scene) {
  const action = step.action;
  if (action === "set_caption") {
    await ensureCaption(page, scene.narration_en);
    return;
  }
  if (action === "wait") {
    if (step.testid) {
      const sel = `[data-testid="${step.testid}"]`;
      await page.locator(sel).first().waitFor({ state: "visible", timeout: 10000 });
    }
    await page.waitForTimeout(step.ms || 1000);
    return;
  }
  if (!step.testid) {
    console.warn(`skip step without testid: ${JSON.stringify(step)}`);
    return;
  }
  const sel = `[data-testid="${step.testid}"]`;
  const loc = page.locator(sel).first();
  await loc.waitFor({ state: "visible", timeout: 10000 });

  if (action === "click") {
    await loc.click();
    // Allow form navigations (e.g. portal PIN auth) to settle.
    await page.waitForLoadState("networkidle").catch(() => {});
    return;
  }
  if (action === "fill") {
    await loc.fill(step.value || "");
    return;
  }
  if (action === "select") {
    if (step.label_contains) {
      const options = await loc.locator("option").allTextContents();
      const match = options.find((o) => o.includes(step.label_contains));
      if (!match) {
        throw new Error(`no option containing ${JSON.stringify(step.label_contains)} in ${sel}`);
      }
      await loc.selectOption({ label: match });
      return;
    }
    await loc.selectOption(step.value || "");
    return;
  }
  console.warn(`unknown action: ${action}`);
}

async function main() {
  console.log(`Capturing ${scenes.length} scenes from ${baseURL}`);
  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext({
    viewport,
    recordVideo: {
      dir: videoDir,
      size: viewport,
    },
    locale: "en-US",
  });
  await context.addCookies([
    {
      name: "lang",
      value: "en",
      url: baseURL,
    },
  ]);

  const page = await context.newPage();

  for (const scene of scenes) {
    const url = new URL(scene.path || "/", baseURL).toString();
    console.log(`→ scene ${scene.id}: ${url}`);
    await page.goto(url, { waitUntil: "networkidle" });
    await ensureCaption(page, scene.narration_en);

    const steps = scene.steps || [];
    if (steps.length === 0) {
      const hold = Math.max(2000, (scene.max_seconds || 5) * 1000);
      await page.waitForTimeout(Math.min(hold, 8000));
      continue;
    }

    for (const step of steps) {
      await runStep(page, step, scene);
      await ensureCaption(page, scene.narration_en);
    }

    const pad = Math.max(800, Math.min((scene.max_seconds || 5) * 200, 3000));
    await page.waitForTimeout(pad);
  }

  const video = page.video();
  await page.close();
  await context.close();
  await browser.close();

  if (!video) {
    throw new Error("Playwright did not produce a video");
  }
  const tmpPath = await video.path();
  const dest = path.join(outDir, videoName);
  fs.renameSync(tmpPath, dest);
  fs.rmSync(videoDir, { recursive: true, force: true });
  console.log(`Wrote ${dest}`);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
