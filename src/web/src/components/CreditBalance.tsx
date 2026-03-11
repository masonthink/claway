import { Coins } from "lucide-react";

export default function CreditBalance({ balance }: { balance: number }) {
  return (
    <div className="flex items-center gap-2 rounded-lg bg-indigo-50 px-4 py-2">
      <Coins className="h-5 w-5 text-indigo-600" />
      <span className="text-sm font-medium text-gray-600">Credits</span>
      <span className="text-lg font-bold text-indigo-700">
        {balance.toFixed(0)}
      </span>
    </div>
  );
}
