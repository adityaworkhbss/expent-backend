// Cloudflare Worker serving the Android app's remote config (UI text + serverUrl).
//
// The JSON lives in Workers KV and is edited from the admin panel (backend PUT /admin/app-config).
// Bindings: KV namespace bound as CONFIG_KV. The key defaults to "remote_config" (same as the
// backend's CF_KV_CONFIG_KEY); override with a CONFIG_KEY variable if you change it there.

const CORS_HEADERS = {
  "Access-Control-Allow-Origin": "*",
  "Access-Control-Allow-Methods": "GET, HEAD, OPTIONS",
  "Access-Control-Allow-Headers": "If-None-Match",
  "Access-Control-Expose-Headers": "ETag",
};

export default {
  async fetch(request, env) {
    if (request.method === "OPTIONS") {
      return new Response(null, { status: 204, headers: CORS_HEADERS });
    }
    if (request.method !== "GET" && request.method !== "HEAD") {
      return json({ error: "Method not allowed" }, 405, { Allow: "GET, HEAD, OPTIONS" });
    }

    const config = await env.CONFIG_KV.get(env.CONFIG_KEY || "remote_config");
    if (config === null) {
      // The app keeps its cached config (or bundled file) when this fails
      return json({ error: "Config not published yet" }, 503);
    }

    // The app sends this back as If-None-Match and gets a body-less 304 when nothing changed
    const etag = `"${await sha256(config)}"`;
    const headers = {
      ...CORS_HEADERS,
      "Content-Type": "application/json; charset=utf-8",
      "Cache-Control": "public, max-age=0, must-revalidate",
      ETag: etag,
    };
    if (matchesEtag(request.headers.get("If-None-Match"), etag)) {
      return new Response(null, { status: 304, headers });
    }
    return new Response(request.method === "HEAD" ? null : config, { status: 200, headers });
  },
};

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
    headers: { ...CORS_HEADERS, "Content-Type": "application/json; charset=utf-8", ...extraHeaders },
  });
}
