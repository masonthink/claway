// Claway plugin entry point for OpenClaw Gateway.

import { ClawayClient } from "./client";
import { registerTools } from "./tools";
import { registerContextInjection } from "./context";

export default {
  id: "claway",
  name: "Claway",
  configSchema: {
    type: "object",
    properties: {
      platformUrl: { type: "string" },
      apiKey: { type: "string" },
    },
    required: ["platformUrl", "apiKey"],
  },

  register(api: any) {
    const config = api.config || {};
    const platformUrl = config.platformUrl;
    const apiKey = config.apiKey;

    if (!platformUrl || !apiKey) {
      api.logger?.warn?.(
        "Claway plugin: platformUrl and apiKey are required. Plugin will not register tools."
      );
      return;
    }

    const client = new ClawayClient(platformUrl, apiKey);

    // Register all agent tools
    registerTools(api, client);

    // Register context injection for active tasks
    registerContextInjection(api, client);

    api.logger?.info?.("Claway plugin loaded successfully.");
  },
};
