export default function Logo({ className = "h-6 w-6" }: { className?: string }) {
  return (
    <svg
      viewBox="0 0 28 28"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={className}
      aria-hidden="true"
    >
      {/* Two overlapping rounded squares — collaboration / convergence */}
      <rect x="1" y="5" width="18" height="18" rx="5" fill="var(--accent, #e65c46)" opacity="0.85" />
      <rect x="9" y="1" width="18" height="18" rx="5" fill="var(--accent-deep, #bf3f30)" />
    </svg>
  );
}
