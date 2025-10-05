"use client";

import { Button, Card, CardBody, CardHeader } from "@heroui/react";

interface UploadFormProps {
  selectedFile: File | null;
  isProcessing: boolean;
  onFileChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onSubmit: (e: React.FormEvent) => void;
}

export default function UploadForm({
  selectedFile,
  isProcessing,
  onFileChange,
  onSubmit,
}: UploadFormProps) {
  return (
    <Card className="mb-8 bg-[#1a1d29] border-[#2a2d3a]">
      <CardHeader>
        <h2 className="text-xl font-semibold text-white">Upload Receipt</h2>
      </CardHeader>
      <CardBody>
        <form onSubmit={onSubmit} className="space-y-4">
          <div>
            <input
              type="file"
              accept=".jpg,.jpeg,.png"
              onChange={onFileChange}
              className="block w-full text-sm text-gray-300
                file:mr-4 file:py-2 file:px-4
                file:rounded-full file:border-0
                file:text-sm file:font-semibold
                file:bg-blue-600 file:text-white
                hover:file:bg-blue-700
                file:cursor-pointer
                cursor-pointer"
            />
            {selectedFile && (
              <p className="mt-2 text-sm text-gray-400">
                Selected: {selectedFile.name}
              </p>
            )}
          </div>

          <Button
            type="submit"
            color="primary"
            variant="shadow"
            size="lg"
            fullWidth
            isDisabled={!selectedFile || isProcessing}
            isLoading={isProcessing}
          >
            {isProcessing ? "Processing..." : "Process Receipt"}
          </Button>
        </form>
      </CardBody>
    </Card>
  );
}
