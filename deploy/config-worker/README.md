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
5. **Switch the worker.** The current `expent-config` Worker only serves a static file, and the dashboard
   doesn't allow bindings on those, so deploy over it with Wrangler (same name, same URL):
   1. Put the namespace ID from step 1 into [`wrangler.jsonc`](wrangler.jsonc) (`REPLACE_WITH_NAMESPACE_ID`).
   2. From this folder run `npx wrangler login` (opens the browser once), then `npx wrangler deploy`.

   Wrangler may warn that the deployed Worker differs from your local config (it had static assets); confirm.
   Only do this after step 4, or the worker will answer `503 Config not published yet`.
6. **Check it:** `curl -i https://expent-config.adityaworkhbss.workers.dev/` should return `200` with the JSON
   and an `ETag` header; repeating with `-H 'If-None-Match: <that etag>'` should return `304`.

## Shared key

The worker only serves the config to requests carrying the right `X-Config-Key` header. The app sends it
on every config fetch; the value is built into the app from `CONFIG_API_KEY` in `local.properties`
(gitignored). The key ships inside the APK, so this keeps casual visitors and scrapers out; it isn't
a secret from someone who unpacks the app. Never put secrets in the config itself.

The same value lives in three places:

| Where | How to set it |
|---|---|
| App (`expent/local.properties`) | `CONFIG_API_KEY=...` (CI: environment variable `CONFIG_API_KEY`) |
| Worker | `npx wrangler secret put CONFIG_API_KEY` (paste the value when asked) |
| Backend (Render env + local `.env`) | `CONFIG_API_KEY=...`, used only by the panel's empty-KV fallback |

Order for turning it on: set the worker secret and the backend env, install an app build that has the key,
then `npx wrangler deploy`. App builds without the key get `401` from then on and keep using their cached
config. To give older installs time to update, deploy with `"REQUIRE_CONFIG_KEY": "false"` first (requests
without the key are only logged; see `npx wrangler tail`), then switch it to `"true"` and deploy again.

To rotate the key: change it in all three places, ship the app update, and redeploy the worker.

## Notes

- Every publish first copies the current value to `remote_config_previous` in the same namespace. To roll
  back, copy that value over `remote_config` in the KV dashboard (or paste it into the panel's JSON tab).
- KV is eventually consistent: a publish can take up to ~60 seconds to reach every Cloudflare location.
  Apps then pick it up on next launch, when they reconnect, or on their 15-minute sync.
- The app also ships `app/src/main/assets/remote_config.json`, used only if its very first launch is
  offline. Keep it roughly in sync when you make big changes.
