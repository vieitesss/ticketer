"use client";

import {
  Table,
  TableHeader,
  TableColumn,
  TableBody,
  TableRow,
  TableCell,
} from "@heroui/react";
import { Item } from "@/types/receipt";

interface ReceiptTableProps {
  items: Item[];
}

export default function ReceiptTable({ items }: ReceiptTableProps) {
  return (
    <div>
      <h3 className="font-semibold mb-4 text-white text-lg">Items Purchased</h3>
      <Table
        aria-label="Receipt items table"
        className="dark"
        classNames={{
          wrapper: "bg-[#252936] border border-[#2a2d3a]",
          th: "bg-[#1a1d29] text-gray-300 font-semibold text-sm",
          td: "text-gray-200",
        }}
      >
        <TableHeader>
          <TableColumn>ITEM</TableColumn>
          <TableColumn align="center">QUANTITY</TableColumn>
          <TableColumn align="end">UNIT PRICE</TableColumn>
          <TableColumn align="end">SUBTOTAL</TableColumn>
        </TableHeader>
        <TableBody>
          {items.map((item) => (
            <TableRow key={item.id} className="hover:bg-[#2a2d3a] transition-colors">
              <TableCell className="font-medium">{item.product_name}</TableCell>
              <TableCell className="text-center">{item.quantity.toFixed(3)}</TableCell>
              <TableCell className="text-right">${item.price_paid.toFixed(2)}</TableCell>
              <TableCell className="text-right font-semibold">
                ${item.subtotal.toFixed(2)}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
