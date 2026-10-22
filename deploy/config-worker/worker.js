// Cloudflare Worker serving the Android app's remote config (UI text + serverUrl).
//
// The JSON lives in Workers KV and is edited from the admin panel (backend PUT /admin/app-config).
// Bindings / vars / secrets (see wrangler.jsonc):
//   CONFIG_KV           KV namespace with the config
//   CONFIG_KEY          KV key, default "remote_config" (same as the backend's CF_KV_CONFIG_KEY)
//   CONFIG_API_KEY      secret: the shared key the app sends in X-Config-Key (`npx wrangler secret put`)
//   REQUIRE_CONFIG_KEY  "true" (default): requests without the right key get 401
//                       "false": wrong/missing keys are only logged (for rolling out to old app versions)
//
// The key ships inside the app, so this stops casual access to the URL, not a determined attacker.

const CONFIG_KEY_HEADER = "X-Config-Key";

export default {
  async fetch(request, env) {
    if (request.method !== "GET" && request.method !== "HEAD") {
      return json({ error: "Method not allowed" }, 405, { Allow: "GET, HEAD" });
    }

    if (!env.CONFIG_API_KEY) {
      // Fail closed: a missing secret shouldn't silently make the config public
      console.log(JSON.stringify({ configKey: "secret CONFIG_API_KEY is not set" }));
      return json({ error: "Server misconfigured" }, 500);
    }
    const provided = request.headers.get(CONFIG_KEY_HEADER);
    if (!provided || !timingSafeEqual(provided, env.CONFIG_API_KEY)) {
      const enforce = (env.REQUIRE_CONFIG_KEY || "true").toLowerCase() !== "false";
      console.log(JSON.stringify({ configKey: provided ? "wrong key" : "missing key", enforced: enforce }));
      if (enforce) return json({ error: "Unauthorized" }, 401);
    }

    const config = await env.CONFIG_KV.get(env.CONFIG_KEY || "remote_config");
    if (config === null) {
      // The app keeps its cached config (or bundled file) when this fails
      return json({ error: "Config not published yet" }, 503);
    }

    // The app sends this back as If-None-Match and gets a body-less 304 when nothing changed
    const etag = `"${await sha256(config)}"`;
    const headers = {
      "Content-Type": "application/json; charset=utf-8",
      // Only for key holders, so keep it out of shared caches
      "Cache-Control": "private, no-cache",
      ETag: etag,
    };
    if (matchesEtag(request.headers.get("If-None-Match"), etag)) {
      return new Response(null, { status: 304, headers });
    }
    return new Response(request.method === "HEAD" ? null : config, { status: 200, headers });
  },
};

// Compares without stopping at the first difference, so response timing doesn't reveal the key
function timingSafeEqual(a, b) {
  const x = new TextEncoder().encode(a);
  const y = new TextEncoder().encode(b);
  let diff = x.length ^ y.length;
  for (let i = 0; i < Math.max(x.length, y.length); i++) {
    diff |= (x[i] ?? 0) ^ (y[i] ?? 0);
  }
  return diff === 0;
}

function matchesEtag(ifNoneMatch, etag) {
  if (!ifNoneMatch) return false;
  return ifNoneMatch.split(",").some((tag) => tag.trim().replace(/^W\//, "") === etag);
}

async function sha256(text) {
  const digest = await crypto.subtle.digest("SHA-256", new TextEncoder().encode(text));
  return [...new Uint8Array(digest)].map((b) => b.toString(16).padStart(2, "0")).join("");
}

function json(body, status, extraHeaders = {}) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json; charset=utf-8", "Cache-Control": "no-store", ...extraHeaders },
  });
}
