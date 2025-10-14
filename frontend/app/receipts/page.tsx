"use client";

import { useState, useEffect } from "react";
import { Button } from "@heroui/react";
import { api } from "@/lib/api";
import { Receipt, ReceiptListItem } from "@/types/receipt";
import ReceiptListSidebar from "@/components/ReceiptListSidebar";
import ReceiptDetailView from "@/components/ReceiptDetailView";
import UploadModal from "@/components/UploadModal";

export default function ReceiptsPage() {
  const [receipts, setReceipts] = useState<ReceiptListItem[]>([]);
  const [selectedReceiptId, setSelectedReceiptId] = useState<string | null>(
    null
  );
  const [selectedReceipt, setSelectedReceipt] = useState<Receipt | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingDetail, setIsLoadingDetail] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isUploadModalOpen, setIsUploadModalOpen] = useState(false);

  // Load receipts list
  useEffect(() => {
    loadReceipts();
  }, []);

  // Load selected receipt details
  useEffect(() => {
    if (selectedReceiptId) {
      loadReceiptDetail(selectedReceiptId);
    } else {
      setSelectedReceipt(null);
    }
  }, [selectedReceiptId]);

  const loadReceipts = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const data = await api.getReceipts();
      setReceipts(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load receipts");
    } finally {
      setIsLoading(false);
    }
  };

  const loadReceiptDetail = async (id: string) => {
    try {
      setIsLoadingDetail(true);
      setError(null);
      const data = await api.getReceipt(id);
      setSelectedReceipt(data);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to load receipt details"
      );
    } finally {
      setIsLoadingDetail(false);
    }
  };

  const handleReceiptClick = (id: string) => {
    setSelectedReceiptId(id);
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Are you sure you want to delete this receipt?")) {
      return;
    }

    try {
      await api.deleteReceipt(id);
      // Refresh receipts list
      await loadReceipts();
      // Clear selection if deleted receipt was selected
      if (selectedReceiptId === id) {
        setSelectedReceiptId(null);
      }
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to delete receipt");
    }
  };

  const handleUploadSuccess = () => {
    setIsUploadModalOpen(false);
    loadReceipts();
  };

  return (
    <div className="flex justify-center h-screen bg-[#0f1117]">
      <div className="flex w-full max-w-7xl h-screen">
        {/* Left Sidebar - Receipt List */}
        <div className="w-80 border-r border-[#2a2d3a] flex flex-col">
          <div className="p-4 border-b border-[#2a2d3a]">
            <h1 className="text-xl font-bold text-white mb-4">Receipts</h1>
            <Button
              className="w-full bg-[#3b82f6] text-white hover:bg-[#2563eb]"
              onPress={() => setIsUploadModalOpen(true)}
            >
              Upload Receipt
            </Button>
          </div>

          <ReceiptListSidebar
            receipts={receipts}
            selectedId={selectedReceiptId}
            isLoading={isLoading}
            error={error}
            onReceiptClick={handleReceiptClick}
          />
        </div>

        {/* Right Side - Receipt Detail */}
        <div className="flex-1 overflow-auto">
          <ReceiptDetailView
            receipt={selectedReceipt}
            isLoading={isLoadingDetail}
            onDelete={handleDelete}
            onUpdate={loadReceiptDetail}
          />
        </div>

        {/* Upload Modal */}
        <UploadModal
          isOpen={isUploadModalOpen}
          onClose={() => setIsUploadModalOpen(false)}
          onSuccess={handleUploadSuccess}
        />
      </div>
    </div>
  );
}
