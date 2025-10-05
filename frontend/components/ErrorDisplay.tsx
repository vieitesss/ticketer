"use client";

import { Card, CardBody } from "@heroui/react";

interface ErrorDisplayProps {
  message: string;
}

export default function ErrorDisplay({ message }: ErrorDisplayProps) {
  return (
    <Card className="border-red-600 bg-red-950/30" role="alert">
      <CardBody>
        <p className="text-red-400 font-semibold">Error: {message}</p>
      </CardBody>
    </Card>
  );
}
