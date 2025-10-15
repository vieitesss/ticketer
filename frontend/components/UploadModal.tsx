import { useState } from "react";
import {
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Button,
  Card,
  CardBody,
} from "@heroui/react";
import { api, APIError } from "@/lib/api";

interface Props {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export default function UploadModal({ isOpen, onClose, onSuccess }: Props) {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      // Validate file type
      const validTypes = ["image/jpeg", "image/jpg", "image/png"];
      if (!validTypes.includes(file.type)) {
        setError("Please select a valid image file (JPG, JPEG, or PNG)");
        return;
      }

      // Validate file size (max 10MB)
      const maxSize = 10 * 1024 * 1024; // 10MB
      if (file.size > maxSize) {
        setError("File size must be less than 10MB");
        return;
      }

      setSelectedFile(file);
      setError(null);
    }
  };

  const handleUpload = async () => {
    if (!selectedFile) {
      setError("Please select a file");
      return;
    }

    try {
      setIsUploading(true);
      setError(null);
      await api.uploadReceipt(selectedFile);
      onSuccess();
      handleClose();
    } catch (err) {
      if (err instanceof APIError) {
        setError(err.message);
      } else {
        setError("Failed to upload receipt. Please try again.");
      }
    } finally {
      setIsUploading(false);
    }
  };

  const handleClose = () => {
    setSelectedFile(null);
    setError(null);
    onClose();
  };

  return (
    <Modal isOpen={isOpen} onClose={handleClose} size="lg" className="bg-[#16181f]">
      <ModalContent className="bg-[#16181f] border border-[#2a2d3a]">
        <ModalHeader className="border-b border-[#2a2d3a]">
          <h3 className="text-xl font-bold text-white">Upload Receipt</h3>
        </ModalHeader>
        <ModalBody className="py-6">
          <div className="space-y-4">
            {/* File Input */}
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                Select receipt image
              </label>
              <input
                type="file"
                accept="image/jpeg,image/jpg,image/png"
                onChange={handleFileChange}
                disabled={isUploading}
                className="block w-full text-sm text-gray-400
                  file:mr-4 file:py-2 file:px-4
                  file:rounded-lg file:border-0
                  file:text-sm file:font-semibold
                  file:bg-[#3b82f6] file:text-white
                  hover:file:bg-[#2563eb]
                  file:cursor-pointer
                  cursor-pointer"
              />
              <p className="mt-1 text-xs text-gray-500">
                Supported formats: JPG, JPEG, PNG (max 10MB)
              </p>
            </div>

            {/* Preview */}
            {selectedFile && (
              <Card className="bg-[#1e2330] border border-[#2a2d3a]">
                <CardBody className="p-3">
                  <div className="flex items-center gap-3">
                    <div className="flex-shrink-0">
                      <svg
                        className="w-8 h-8 text-[#3b82f6]"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                        />
                      </svg>
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-white truncate">
                        {selectedFile.name}
                      </p>
                      <p className="text-xs text-gray-400">
                        {(selectedFile.size / 1024).toFixed(2)} KB
                      </p>
                    </div>
                  </div>
                </CardBody>
              </Card>
            )}

            {/* Error */}
            {error && (
              <Card className="bg-red-900/20 border border-red-500">
                <CardBody className="p-3">
                  <p className="text-red-400 text-sm">{error}</p>
                </CardBody>
              </Card>
            )}
          </div>
        </ModalBody>
        <ModalFooter className="border-t border-[#2a2d3a]">
          <Button 
            variant="flat" 
            onPress={handleClose} 
            isDisabled={isUploading}
            className="bg-[#1e2330] text-gray-300 hover:bg-[#2a2d3a]"
          >
            Cancel
          </Button>
          <Button
            onPress={handleUpload}
            isLoading={isUploading}
            isDisabled={!selectedFile}
            className="bg-[#3b82f6] text-white hover:bg-[#2563eb]"
          >
            Upload
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}
