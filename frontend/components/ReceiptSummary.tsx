"use client";

interface ReceiptSummaryProps {
  subtotal: number;
  discounts: number;
  totalAmount: number;
}

export default function ReceiptSummary({
  subtotal,
  discounts,
  totalAmount,
}: ReceiptSummaryProps) {
  return (
    <section className="border-t border-[#2a2d3a] pt-6">
      <article className="bg-[#252936] rounded-lg p-6 border border-[#2a2d3a]">
        <dl className="space-y-3">
          <div className="flex justify-between text-gray-300">
            <dt>Subtotal</dt>
            <dd className="font-medium">€{subtotal.toFixed(2)}</dd>
          </div>
          {discounts > 0 && (
            <div className="flex justify-between text-green-400">
              <dt>Discounts</dt>
              <dd className="font-medium">-€{discounts.toFixed(2)}</dd>
            </div>
          )}
          <div className="border-t border-[#2a2d3a] pt-3 mt-3">
            <div className="flex justify-between text-2xl font-bold text-white">
              <dt>Total</dt>
              <dd className="text-blue-400">€{totalAmount.toFixed(2)}</dd>
            </div>
          </div>
        </dl>
      </article>
    </section>
  );
}
