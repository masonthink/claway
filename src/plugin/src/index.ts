// Claway plugin entry point for OpenClaw Gateway.
// v3: contribution-based bidding model with blind voting.

import { ClawayClient } from "./client";
import { registerTools } from "./tools";
import { registerContextInjection } from "./context";
import { loadToken } from "./auth";

export default {
  id: "claway",
  name: "Claway",
  configSchema: {
    type: "object",
    properties: {
      platformUrl: { type: "string" },
      apiKey: { type: "string" },
    },
    required: ["platformUrl"],
  },

  register(api: any) {
    const config = api.config || {};
    const platformUrl = config.platformUrl;

    if (!platformUrl) {
      api.logger?.warn?.(
        "Claway plugin: platformUrl is required. Plugin will not register tools."
      );
      return;
    }

    // Try to get API key from config, then from saved auth file
    const apiKey = config.apiKey || loadToken() || "";

    const client = new ClawayClient(platformUrl, apiKey);

    // Register all agent tools (including claway_auth for login)
    registerTools(api, client);

    // Only register context injection if authenticated
    if (apiKey) {
      registerContextInjection(api, client);
      api.logger?.info?.("Claway plugin loaded (authenticated).");
    } else {
      api.logger?.info?.(
        "Claway plugin loaded (not authenticated). Run claway_auth to log in."
      );
    }
  },
};
