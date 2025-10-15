import { Receipt, ReceiptListItem, UpdateItemRequest } from "@/types/receipt";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export class APIError extends Error {
  constructor(
    message: string,
    public status?: number
  ) {
    super(message);
    this.name = "APIError";
  }
}

export const api = {
  // Upload and process a receipt image
  async uploadReceipt(file: File): Promise<Receipt> {
    const formData = new FormData();
    formData.append("receipt", file);

    const response = await fetch(`${API_BASE_URL}/receipts/upload`, {
      method: "POST",
      body: formData,
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new APIError(
        errorText || "Failed to upload receipt",
        response.status
      );
    }

    return response.json();
  },

  // Get all receipts (for sidebar list)
  async getReceipts(limit = 50, offset = 0): Promise<ReceiptListItem[]> {
    const response = await fetch(
      `${API_BASE_URL}/receipts?limit=${limit}&offset=${offset}`
    );

    if (!response.ok) {
      throw new APIError("Failed to fetch receipts", response.status);
    }

    return response.json();
  },

  // Get a single receipt by ID (full details)
  async getReceipt(id: string): Promise<Receipt> {
    const response = await fetch(`${API_BASE_URL}/receipts/${id}`);

    if (!response.ok) {
      throw new APIError("Failed to fetch receipt", response.status);
    }

    return response.json();
  },

  // Delete a receipt
  async deleteReceipt(id: string): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/receipts/${id}`, {
      method: "DELETE",
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new APIError(
        errorText || "Failed to delete receipt",
        response.status
      );
    }
  },

  // Update an item's quantity and price
  async updateItem(
    itemId: string,
    data: UpdateItemRequest
  ): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/items/${itemId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new APIError(
        errorText || "Failed to update item",
        response.status
      );
    }
  },
};
