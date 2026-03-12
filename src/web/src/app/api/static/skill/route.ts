import { readFileSync } from "fs";
import { join } from "path";

export function GET() {
  const content = readFileSync(join(process.cwd(), "public", "skill.md"), "utf-8");
  return new Response(content, {
    headers: { "Content-Type": "text/markdown; charset=utf-8" },
  });
}
