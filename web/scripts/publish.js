// Usage: node ./scripts/publish.ts --path=<dir>
// e.g: node scripts/publish.ts ./../assets/example01

const fs = require("fs");

const args = process.argv.slice(2);
const directory = args[0];

if (!fs.existsSync(directory)) {
  console.error(`Directory '${directory}' does not exist.`);
  process.exit(1);
}

const readAll = (key, path) => {
  const jabls = [];
  const filesAndFolders = fs.readdirSync(path);
  for (const item of filesAndFolders) {
    const itemPath = `${path}/${item}`;
    const stats = fs.statSync(itemPath);
    if (stats.isDirectory()) {
      jabls.push(...readAll(key, itemPath));
    } else if (itemPath.endsWith(".jabl")) {
      const trimmed = itemPath.replace(directory + "/", "");
      jabls.push({ key, path: itemPath, value: trimmed });
    }
  }
  return jabls;
};

const source = directory.split("/").pop();

console.log(`Reading all .jabl files in '${directory}' and uploading to '${source}'...`);

const toUpload = readAll(source, directory);

const localOnly = process.argv[3] !== "--remote";
toUpload.forEach(({ key, path, value }) => {
  const cmd = `npx wrangler kv:key put "${key}:${value}" --path="${path}" --binding="SECTIONS"${localOnly ? " --local" : ""}`;
  console.log(cmd);
  const result = require("child_process").execSync(cmd).toString();
  console.log(result);
});
