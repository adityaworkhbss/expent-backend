# App config worker

Serves the Android app's remote config (UI text + `serverUrl`) at
`https://expent-config.adityaworkhbss.workers.dev/`. The JSON is stored in Workers KV and edited from the
admin panel's **App Config** page, which saves through the backend (`PUT /api/v1/admin/app-config`).

```
Admin panel ──PUT──> Backend ──Cloudflare API──> KV ("remote_config") <──GET── Worker <──GET── Android app
```

## One-time setup (no downtime)

Do the steps in this order: the old worker keeps serving the app until the last step.

1. **Create the KV namespace.** Cloudflare dashboard → Storage & Databases → KV → Create, e.g. `expent-config`.
   Copy its **Namespace ID**.
2. **Create an API token.** My Profile → API Tokens → Create Token → Custom token with the permission
   **Account · Workers KV Storage · Edit**, limited to your account.
3. **Configure the backend** (Render → Environment, and your local `.env`):
   ```
   CF_ACCOUNT_ID=<your account id, shown in the dashboard sidebar/URL>
   CF_API_TOKEN=<token from step 2>
   CF_KV_NAMESPACE_ID=<namespace id from step 1>
   CF_KV_CONFIG_KEY=remote_config
   CONFIG_PUBLIC_URL=https://expent-config.adityaworkhbss.workers.dev/
   ```
   Redeploy the backend.
4. **Seed KV from the panel.** Open the panel → App Config. While KV is empty it shows the config the app
   gets today, with a blue "isn't in the new storage yet" note. Click **Review & publish** → **Publish**.
5. **Switch the worker.** Workers & Pages → `expent-config` → Settings → Bindings → add a **KV namespace**
   binding named `CONFIG_KV` pointing at the namespace from step 1. Then Edit code, replace everything with
   [`worker.js`](worker.js), and Deploy.
6. **Check it:** `curl -i https://expent-config.adityaworkhbss.workers.dev/` should return `200` with the JSON
   and an `ETag` header; repeating with `-H 'If-None-Match: <that etag>'` should return `304`.

## Notes

- Every publish first copies the current value to `remote_config_previous` in the same namespace. To roll
  back, copy that value over `remote_config` in the KV dashboard (or paste it into the panel's JSON tab).
- KV is eventually consistent: a publish can take up to ~60 seconds to reach every Cloudflare location.
  Apps then pick it up on next launch, when they reconnect, or on their 15-minute sync.
- The app also ships `app/src/main/assets/remote_config.json`, used only if its very first launch is
  offline. Keep it roughly in sync when you make big changes.
