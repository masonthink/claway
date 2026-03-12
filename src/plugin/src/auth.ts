// CLI OAuth authentication flow for Claway.
// Starts a localhost HTTP server, opens browser for X OAuth, receives JWT callback.

import * as http from "http";
import * as fs from "fs";
import * as path from "path";
import * as os from "os";

const CLI_PORT = 19876;
const AUTH_TIMEOUT_MS = 120_000; // 2 minutes

interface AuthResult {
  token: string;
}

/**
 * Get the path to the Claway auth config file.
 */
function getAuthFilePath(): string {
  const configDir = path.join(os.homedir(), ".config", "claway");
  return path.join(configDir, "auth.json");
}

/**
 * Save JWT token to ~/.config/claway/auth.json
 */
function saveToken(token: string): void {
  const authFile = getAuthFilePath();
  const configDir = path.dirname(authFile);

  // Create directory if it doesn't exist
  if (!fs.existsSync(configDir)) {
    fs.mkdirSync(configDir, { recursive: true, mode: 0o700 });
  }

  const data = JSON.stringify({ token, saved_at: new Date().toISOString() }, null, 2);
  fs.writeFileSync(authFile, data, { mode: 0o600 });
}

/**
 * Load saved JWT token from ~/.config/claway/auth.json
 */
export function loadToken(): string | null {
  const authFile = getAuthFilePath();
  try {
    const data = JSON.parse(fs.readFileSync(authFile, "utf-8"));
    return data.token || null;
  } catch {
    return null;
  }
}

/**
 * Run the OAuth login flow:
 * 1. Start localhost HTTP server on CLI_PORT
 * 2. Return the auth URL for user to open in browser
 * 3. Wait for callback with JWT token
 * 4. Save token and return
 */
export function runAuthFlow(platformUrl: string): Promise<AuthResult> {
  return new Promise((resolve, reject) => {
    const authUrl = `${platformUrl}/api/v1/auth/x?cli_port=${CLI_PORT}`;

    const server = http.createServer((req, res) => {
      const url = new URL(req.url || "/", `http://127.0.0.1:${CLI_PORT}`);

      if (url.pathname === "/callback") {
        const token = url.searchParams.get("token");

        if (token) {
          // Save token
          saveToken(token);

          // Send success page
          res.writeHead(200, { "Content-Type": "text/html; charset=utf-8" });
          res.end(successHTML());

          // Close server and resolve
          server.close();
          resolve({ token });
        } else {
          res.writeHead(400, { "Content-Type": "text/html; charset=utf-8" });
          res.end(errorHTML("Missing token in callback"));
        }
      } else {
        res.writeHead(404);
        res.end("Not Found");
      }
    });

    // Timeout
    const timeout = setTimeout(() => {
      server.close();
      reject(new Error("Authentication timed out (2 minutes). Please try again."));
    }, AUTH_TIMEOUT_MS);

    server.on("close", () => clearTimeout(timeout));

    server.on("error", (err: NodeJS.ErrnoException) => {
      clearTimeout(timeout);
      if (err.code === "EADDRINUSE") {
        reject(new Error(`Port ${CLI_PORT} is already in use. Close the other process and try again.`));
      } else {
        reject(err);
      }
    });

    server.listen(CLI_PORT, "127.0.0.1", () => {
      // Server is ready - the tool will return the URL to the user
    });

    // Store authUrl and server reference so the tool can access them
    (runAuthFlow as any)._pending = { authUrl, server };
  });
}

function successHTML(): string {
  return `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>Claway - Login Success</title>
<style>body{font-family:system-ui;display:flex;justify-content:center;align-items:center;min-height:100vh;margin:0;background:#0a0a0a;color:#e5e5e5}
.card{text-align:center;padding:2rem 3rem;border-radius:12px;border:1px solid #333}
.check{font-size:3rem;margin-bottom:1rem}
p{color:#999;margin-top:0.5rem}</style></head>
<body><div class="card"><div class="check">&#10003;</div><h2>Login Successful</h2><p>You can close this tab and return to your terminal.</p></div></body></html>`;
}

function errorHTML(msg: string): string {
  return `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>Claway - Auth Error</title>
<style>body{font-family:system-ui;display:flex;justify-content:center;align-items:center;min-height:100vh;margin:0;background:#0a0a0a;color:#e5e5e5}
.card{text-align:center;padding:2rem 3rem;border-radius:12px;border:1px solid #333}
p{color:#ff6b6b}</style></head>
<body><div class="card"><h2>Authentication Error</h2><p>${msg}</p></div></body></html>`;
}
