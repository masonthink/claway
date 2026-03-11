export default function MarkdownRenderer({ content }: { content: string }) {
  return (
    <div className="max-w-none">
      <pre
        className="whitespace-pre-wrap rounded-[12px] p-5 font-body text-sm leading-relaxed"
        style={{
          background: "var(--surface-muted)",
          border: "1px solid var(--line)",
        }}
      >
        {content}
      </pre>
    </div>
  );
}
