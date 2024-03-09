```
npm run watch:css
npm run watch:web
```

Adding a new section to the KV backed store

```
npx wrangler kv:key put "{source}:{section}" --path="./section.jabl" --binding="SECTIONS" --local
```