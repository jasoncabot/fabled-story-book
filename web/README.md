```
npm run watch:css
npm run watch:web
```

Adding a new section to the KV backed store

```
node ./scripts/publish.js ./src/example-4 --remote
npx wrangler kv:key put "{source}:{section}" --path="./section.jabl" --binding="SECTIONS" --local
```
