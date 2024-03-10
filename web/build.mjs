import { build } from "esbuild";
import { readFileSync, writeFileSync } from "fs";

console.log("[esbuild] start " + new Date().toLocaleTimeString());

const scripts = build({
  entryPoints: ["./src/index.tsx", "./src/**/*.wasm"],
  entryNames: "[dir]/[name]-[hash]",
  bundle: true,
  minify: true,
  metafile: true,
  sourcemap: false,
  outdir: "dist",
  outbase: "src",
  loader: {
    ".wasm": "file",
  },
  external: ["crypto", "fs", "util"],
});

const staticAssets = build({
  entryPoints: ["./src/**/*.html", "./src/wasm-init.js", "./src/assets/**/*.png", "./src/assets/**/*.jpg"],
  minify: false,
  sourcemap: false,
  outdir: "dist",
  outbase: "src",
  loader: {
    ".png": "copy",
    ".jpg": "copy",
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

const output = await Promise.all([scripts, staticAssets, functions]);

// Adjust index.html so it points to the newly hashed script name
const distPath = Object.keys(output[0].metafile.outputs).find((key) => output[0].metafile.outputs[key].entryPoint === "src/index.tsx").replace("dist/", "");
const importDistScript = `<script src="/${distPath}"></script>`;

const html = readFileSync("dist/index.html", "utf8");
const newHtml = html.replace(/<script src=".\/index.js"><\/script>/, importDistScript);
writeFileSync("dist/index.html", newHtml);

// Hack to force wrangler pages to rebuild
const packageJson = readFileSync("package.json", "utf8");
const pkg = JSON.parse(packageJson);
const version = pkg["version"] ?? "0.0.0";
pkg["version"] = version.replace(/\d+$/, (n) => (parseInt(n) + 1).toString());
writeFileSync("package.json", JSON.stringify(pkg, null, 2));

console.log("[esbuild] done " + new Date().toLocaleTimeString());
