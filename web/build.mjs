import { build } from "esbuild";
import { readFileSync, writeFileSync } from "fs";

console.log("[esbuild] start " + new Date().toLocaleTimeString());

const scripts = build({
  entryPoints: ["./src/index.tsx", "./src/**/*.wasm"],
  bundle: true,
  minify: true,
  sourcemap: true,
  outdir: "dist",
  outbase: "src",
  loader: {
    ".wasm": "file",
  },
  external: ["crypto", "fs", "util"],
});

const staticAssets = build({
  entryPoints: ["./src/**/*.html", "./src/wasm-init.js", "./src/assets/**/*.png"],
  minify: false,
  sourcemap: true,
  outdir: "dist",
  outbase: "src",
  loader: {
    ".png": "copy",
    ".html": "copy",
  },
});

const functions = build({
  entryPoints: ["./functions/**/*.ts"],
  bundle: true,
  metafile: true,
  minify: false,
  sourcemap: true,
  outdir: "dist/functions",
  external: ["crypto", "fs", "util", "@cloudflare/*"],
});

await Promise.all([scripts, staticAssets, functions]);

// Hack to force wrangler pages to rebuild
const packageJson = readFileSync("package.json", "utf8");
const pkg = JSON.parse(packageJson);
const version = pkg["version"] ?? "0.0.0";
pkg["version"] = version.replace(/\d+$/, (n) => (parseInt(n) + 1).toString());
writeFileSync("package.json", JSON.stringify(pkg, null, 2));

console.log("[esbuild] done " + new Date().toLocaleTimeString());
