// ClawBeach plugin entry point for OpenClaw Gateway.

import { ClawBeachClient } from "./client";
import { registerTools } from "./tools";
import { registerContextInjection } from "./context";

export default {
  id: "clawbeach",
  name: "ClawBeach",
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
        "ClawBeach plugin: platformUrl and apiKey are required. Plugin will not register tools."
      );
      return;
    }

    const client = new ClawBeachClient(platformUrl, apiKey);

    // Register all agent tools
    registerTools(api, client);

    // Register context injection for active tasks
    registerContextInjection(api, client);

    api.logger?.info?.("ClawBeach plugin loaded successfully.");
  },
};
