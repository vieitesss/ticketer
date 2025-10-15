import { useState } from "react";
import { Button, Spinner, Input } from "@heroui/react";
import { Receipt } from "@/types/receipt";
import { api } from "@/lib/api";

interface Props {
  receipt: Receipt | null;
  isLoading: boolean;
  onDelete: (id: string) => void;
  onUpdate: (id: string) => void;
}

export default function ReceiptDetailView({
  receipt,
  isLoading,
  onDelete,
  onUpdate,
}: Props) {
  const [editingItemId, setEditingItemId] = useState<string | null>(null);
  const [editQuantity, setEditQuantity] = useState<string>("");
  const [editPrice, setEditPrice] = useState<string>("");
  const [isSaving, setIsSaving] = useState(false);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Spinner size="lg" className="text-[#3b82f6]" />
      </div>
    );
  }

  if (!receipt) {
    return (
      <div className="flex items-center justify-center h-full">
        <p className="text-gray-500 text-lg">
          Select a receipt to view details
        </p>
      </div>
    );
  }

  const handleEditClick = (itemId: string, quantity: number, price: number) => {
    setEditingItemId(itemId);
    setEditQuantity(quantity.toString());
    setEditPrice(price.toString());
  };

  const handleSaveEdit = async (itemId: string) => {
    const quantity = parseFloat(editQuantity);
    const pricePaid = parseFloat(editPrice);

    if (isNaN(quantity) || quantity <= 0) {
      alert("Quantity must be a positive number");
      return;
    }

    if (isNaN(pricePaid) || pricePaid < 0) {
      alert("Price must be a non-negative number");
      return;
    }

    try {
      setIsSaving(true);
      await api.updateItem(itemId, {
        quantity,
        price_paid: pricePaid,
      });
      setEditingItemId(null);
      onUpdate(receipt.id);
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to update item");
    } finally {
      setIsSaving(false);
    }
  };

  const handleCancelEdit = () => {
    setEditingItemId(null);
    setEditQuantity("");
    setEditPrice("");
  };

  return (
    <div className="flex items-start justify-center p-8 min-h-full">
      <div className="w-full max-w-md">
        <div className="bg-[#16181f] text-gray-100 p-8 shadow-2xl border border-[#2a2d3a] font-mono text-sm">
          <div className="text-center border-b-2 border-dashed border-[#3a3d4a] pb-4 mb-4">
            <h2 className="text-xl font-bold mb-1 text-white">{receipt.store.name}</h2>
            <p className="text-xs text-gray-400">{formatDate(receipt.bought_date)}</p>
          </div>

          <div className="space-y-1 mb-4">
            {receipt.items.map((item) => (
              <div key={item.id} className="border-b border-[#2a2d3a] pb-2 mb-2">
                <div className="flex justify-between items-start">
                  <span className="font-semibold flex-1 pr-2 text-gray-200">{item.product_name}</span>
                  {editingItemId !== item.id && (
                    <button
                      onClick={() =>
                        handleEditClick(item.id, item.quantity, item.price_paid)
                      }
                      className="text-xs text-[#3b82f6] hover:text-[#60a5fa] ml-2"
                    >
                      [edit]
                    </button>
                  )}
                </div>
                {editingItemId === item.id ? (
                  <div className="mt-2 space-y-2 bg-[#1e2330] p-2 rounded border border-[#2a2d3a]">
                    <div className="flex gap-2 items-center">
                      <label className="text-xs w-16 text-gray-400">Qty:</label>
                      <Input
                        type="number"
                        value={editQuantity}
                        onChange={(e) => setEditQuantity(e.target.value)}
                        size="sm"
                        className="flex-1"
                        min="0"
                        step="0.001"
                        classNames={{
                          input: "text-white",
                          inputWrapper: "bg-[#0f1117] border border-[#2a2d3a]",
                        }}
                      />
                    </div>
                    <div className="flex gap-2 items-center">
                      <label className="text-xs w-16 text-gray-400">Price:</label>
                      <Input
                        type="number"
                        value={editPrice}
                        onChange={(e) => setEditPrice(e.target.value)}
                        size="sm"
                        className="flex-1"
                        min="0"
                        step="0.01"
                        startContent={<span className="text-gray-400">$</span>}
                        classNames={{
                          input: "text-white",
                          inputWrapper: "bg-[#0f1117] border border-[#2a2d3a]",
                        }}
                      />
                    </div>
                    <div className="flex gap-2 justify-end">
                      <button
                        onClick={() => handleSaveEdit(item.id)}
                        disabled={isSaving}
                        className="px-3 py-1 bg-[#3b82f6] text-white text-xs rounded hover:bg-[#2563eb] disabled:opacity-50"
                      >
                        {isSaving ? "..." : "Save"}
                      </button>
                      <button
                        onClick={handleCancelEdit}
                        disabled={isSaving}
                        className="px-3 py-1 bg-[#2a2d3a] text-gray-300 text-xs rounded hover:bg-[#3a3d4a] disabled:opacity-50"
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                ) : (
                  <div className="flex justify-between text-xs mt-1">
                    <span className="text-gray-400">
                      {item.quantity.toFixed(3)} x ${item.price_paid.toFixed(2)}
                    </span>
                    <span className="font-semibold text-gray-200">${item.subtotal.toFixed(2)}</span>
                  </div>
                )}
              </div>
            ))}
          </div>

          <div className="border-t-2 border-dashed border-[#3a3d4a] pt-3 space-y-1">
            <div className="flex justify-between text-sm text-gray-300">
              <span>SUBTOTAL:</span>
              <span>${receipt.subtotal.toFixed(2)}</span>
            </div>
            {receipt.discounts > 0 && (
              <div className="flex justify-between text-sm text-green-400">
                <span>DISCOUNT:</span>
                <span>-${receipt.discounts.toFixed(2)}</span>
              </div>
            )}
            <div className="flex justify-between text-lg font-bold text-white border-t border-[#3a3d4a] pt-2 mt-2">
              <span>TOTAL:</span>
              <span>${receipt.total_amount.toFixed(2)}</span>
            </div>
          </div>

          <div className="text-center text-xs text-gray-500 mt-6 pt-4 border-t border-[#2a2d3a]">
            <p>THANK YOU FOR YOUR PURCHASE</p>
          </div>
        </div>

        <div className="mt-6 text-center">
          <Button
            size="sm"
            onPress={() => onDelete(receipt.id)}
            className="bg-red-600 text-white hover:bg-red-700"
          >
            Delete Receipt
          </Button>
        </div>
      </div>
    </div>
  );
}

function formatDate(dateString: string): string {
  try {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      weekday: "long",
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  } catch {
    return dateString;
  }
}
