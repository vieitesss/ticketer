import { Card, CardBody, Spinner } from "@heroui/react";
import { ReceiptListItem } from "@/types/receipt";

interface Props {
  receipts: ReceiptListItem[];
  selectedId: string | null;
  isLoading: boolean;
  error: string | null;
  onReceiptClick: (id: string) => void;
}

export default function ReceiptListSidebar({
  receipts,
  selectedId,
  isLoading,
  error,
  onReceiptClick,
}: Props) {
  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Spinner color="primary" size="lg" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-4">
        <Card className="bg-red-900/20 border border-red-500">
          <CardBody>
            <p className="text-red-400 text-sm">{error}</p>
          </CardBody>
        </Card>
      </div>
    );
  }

  if (receipts.length === 0) {
    return (
      <div className="flex items-center justify-center h-full p-4">
        <p className="text-gray-500 text-center">
          No receipts yet. Upload your first receipt to get started!
        </p>
      </div>
    );
  }

  return (
    <div className="flex-1 overflow-y-auto p-4 space-y-2">
      {receipts.map((receipt) => (
        <Card
          key={receipt.id}
          isPressable
          onPress={() => onReceiptClick(receipt.id)}
          classNames={{
            base: `cursor-pointer transition-all w-full ${
              selectedId === receipt.id
                ? "bg-[#1e2330] border-2 border-[#3b82f6]"
                : "bg-[#16181f] hover:!bg-[#1e2330] border-2 border-[#2a2d3a]"
            }`,
            body: "p-3",
          }}
        >
          <CardBody>
            <div className="flex justify-between items-start mb-2">
              <h3 className="font-semibold text-white text-sm">
                {receipt.store_name}
              </h3>
              <span className="text-xs text-gray-400">
                {receipt.item_count} {receipt.item_count === 1 ? "item" : "items"}
              </span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-xs text-gray-400">
                {formatDate(receipt.bought_date)}
              </span>
              <span className="text-sm font-bold text-[#3b82f6]">
                ${receipt.total_amount.toFixed(2)}
              </span>
            </div>
          </CardBody>
        </Card>
      ))}
    </div>
  );
}

function formatDate(dateString: string): string {
  try {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  } catch {
    return dateString;
  }
}
