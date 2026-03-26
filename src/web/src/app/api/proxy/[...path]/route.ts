import { NextRequest, NextResponse } from "next/server";

const BACKEND_URL =
  process.env.BACKEND_URL || "https://api.claway.cc/api/v1";

// Allowed path prefixes to prevent SSRF.
// All prefixes must end with "/" or be followed by "/" or query params at runtime.
const ALLOWED_PREFIXES = [
  "public/",
  "auth/",
  "ideas/",
  "ideas",    // exact match for /ideas (list endpoint)
  "contributions/",
  "me/",
  "me",       // exact match for /me
  "draft/",
];

export async function GET(
  req: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  return proxy(req, await params);
}

export async function POST(
  req: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  return proxy(req, await params);
}

export async function PUT(
  req: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  return proxy(req, await params);
}

export async function DELETE(
  req: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  return proxy(req, await params);
}

async function proxy(req: NextRequest, params: { path: string[] }) {
  const path = params.path.join("/");

  // Validate path: reject traversal, encoded traversal, and non-whitelisted prefixes
  if (
    path.includes("..") ||
    path.includes("//") ||
    path.includes("%2e") ||
    path.includes("%2E") ||
    !ALLOWED_PREFIXES.some((p) =>
      path === p || path.startsWith(p.endsWith("/") ? p : p + "/") || path === p
    )
  ) {
    return NextResponse.json({ error: "forbidden" }, { status: 403 });
  }

  const url = new URL(`${BACKEND_URL}/${path}`);

  // Ensure resolved URL still points to backend (prevent host override)
  const backendOrigin = new URL(BACKEND_URL).origin;
  if (url.origin !== backendOrigin) {
    return NextResponse.json({ error: "forbidden" }, { status: 403 });
  }

  // Forward query params
  req.nextUrl.searchParams.forEach((value, key) => {
    url.searchParams.set(key, value);
  });

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  // Forward auth header
  const auth = req.headers.get("authorization");
  if (auth) {
    headers["Authorization"] = auth;
  }

  const fetchOptions: RequestInit = {
    method: req.method,
    headers,
  };

  // Forward body for non-GET requests
  if (req.method !== "GET" && req.method !== "HEAD") {
    try {
      fetchOptions.body = await req.text();
    } catch {
      // no body
    }
  }

  try {
    const res = await fetch(url.toString(), fetchOptions);
    const data = await res.text();

    return new NextResponse(data, {
      status: res.status,
      headers: {
        "Content-Type": res.headers.get("Content-Type") || "application/json",
      },
    });
  } catch {
    return NextResponse.json(
      { error: "proxy error" },
      { status: 502 }
    );
  }
}
