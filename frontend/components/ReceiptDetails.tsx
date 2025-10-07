"use client";

import { Card, CardBody, CardHeader } from "@heroui/react";
import { Receipt } from "@/types/receipt";
import ReceiptTable from "./ReceiptTable";
import ReceiptSummary from "./ReceiptSummary";

interface ReceiptDetailsProps {
  receipt: Receipt;
}

export default function ReceiptDetails({ receipt }: ReceiptDetailsProps) {
  return (
    <Card className="bg-[#1a1d29] border-[#2a2d3a]">
      <CardHeader className="flex flex-col items-start gap-2 pb-4">
        <h2 className="text-xl font-semibold text-white">Receipt Details</h2>
        <p className="text-lg text-gray-300">
          Store: <span className="font-semibold text-blue-400">{receipt.store_name}</span>
        </p>
      </CardHeader>
      <CardBody className="space-y-6">
        <ReceiptTable items={receipt.items} />
        <ReceiptSummary
          subtotal={receipt.subtotal}
          discounts={receipt.discounts}
          totalAmount={receipt.total_amount}
        />
      </CardBody>
    </Card>
  );
}
