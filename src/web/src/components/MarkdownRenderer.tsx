"use client";

import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

// Security: Explicitly disallow dangerous HTML elements.
// react-markdown does not render raw HTML by default (safe), but we add
// explicit overrides for defense-in-depth against future plugin additions.
const disallowedElements = ["script", "iframe", "object", "embed", "form", "input", "style", "link"];

export default function MarkdownRenderer({ content }: { content: string }) {
  return (
    <div
      className="prose-claway max-w-none rounded-[12px] p-5 text-sm leading-relaxed"
      style={{
        background: "var(--surface-muted)",
        border: "1px solid var(--line)",
      }}
    >
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        disallowedElements={disallowedElements}
        unwrapDisallowed={true}
      >
        {content}
      </ReactMarkdown>
    </div>
  );
}
