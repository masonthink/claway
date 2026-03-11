import { Coins } from "lucide-react";

export default function CreditBalance({ balance }: { balance: number }) {
  return (
    <div
      className="inline-flex items-center gap-2 rounded-[10px] px-3.5 py-2"
      style={{ background: "rgba(43,198,164,0.12)" }}
    >
      <Coins className="h-4 w-4 text-seafoam" />
      <span className="text-sm font-medium text-ink-soft">Credits</span>
      <span className="text-base font-bold" style={{ color: "rgb(26,107,91)" }}>
        {balance.toFixed(0)}
      </span>
    </div>
  );
}
