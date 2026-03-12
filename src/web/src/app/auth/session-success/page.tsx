import Link from "next/link";

export default function SessionSuccessPage() {
  return (
    <main
      className="flex min-h-screen items-center justify-center px-6"
      style={{ background: "var(--surface)" }}
    >
      <div
        className="w-full max-w-md rounded-2xl p-8 text-center"
        style={{ background: "var(--surface)", border: "1px solid var(--line)" }}
      >
        <div className="mb-4 text-4xl">&#10003;</div>
        <h1
          className="mb-3 text-xl font-bold"
          style={{ color: "var(--ink)" }}
        >
          授权成功
        </h1>
        <p
          className="mb-6 text-sm leading-relaxed"
          style={{ color: "var(--ink-soft)" }}
        >
          你可以关闭此页面，返回 IM 继续操作。
        </p>
        <Link
          href="/"
          className="inline-block rounded-xl px-5 py-2.5 text-sm font-medium text-white"
          style={{
            background: "linear-gradient(135deg, var(--accent), var(--accent-deep))",
          }}
        >
          返回首页
        </Link>
      </div>
    </main>
  );
}
