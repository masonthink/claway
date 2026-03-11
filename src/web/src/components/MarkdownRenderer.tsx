// Simple markdown renderer for MVP
// Renders markdown as preformatted text with basic formatting

export default function MarkdownRenderer({ content }: { content: string }) {
  return (
    <div className="prose prose-indigo max-w-none">
      <pre className="whitespace-pre-wrap rounded-lg bg-gray-50 p-6 font-sans text-sm leading-relaxed text-gray-800">
        {content}
      </pre>
    </div>
  );
}
