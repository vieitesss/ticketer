"use client";

import { Item } from "@/types/receipt";

interface ReceiptSummaryProps {
  items: Item[];
  discounts: number | null;
}

export default function ReceiptSummary({
  items,
  discounts,
}: ReceiptSummaryProps) {
  const subtotal = items.reduce(
    (sum, item) => sum + item.quantity * item.price,
    0
  );
  const total = subtotal - (discounts || 0);

  return (
    <section className="border-t border-[#2a2d3a] pt-6">
      <article className="bg-[#252936] rounded-lg p-6 border border-[#2a2d3a]">
        <dl className="space-y-3">
          <div className="flex justify-between text-gray-300">
            <dt>Subtotal</dt>
            <dd className="font-medium">€{subtotal.toFixed(2)}</dd>
          </div>
          {discounts && discounts > 0 && (
            <div className="flex justify-between text-green-400">
              <dt>Discounts</dt>
              <dd className="font-medium">-€{discounts.toFixed(2)}</dd>
            </div>
          )}
          <div className="border-t border-[#2a2d3a] pt-3 mt-3">
            <div className="flex justify-between text-2xl font-bold text-white">
              <dt>Total</dt>
              <dd className="text-blue-400">€{total.toFixed(2)}</dd>
            </div>
          </div>
        </dl>
      </article>
    </section>
  );
}
