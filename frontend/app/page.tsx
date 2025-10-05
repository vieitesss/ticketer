"use client";

import { useReceiptUpload } from "@/hooks/useReceiptUpload";
import UploadForm from "@/components/UploadForm";
import LoadingSpinner from "@/components/LoadingSpinner";
import ErrorDisplay from "@/components/ErrorDisplay";
import ReceiptDetails from "@/components/ReceiptDetails";

export default function Home() {
  const {
    selectedFile,
    isProcessing,
    result,
    error,
    handleFileChange,
    handleSubmit,
  } = useReceiptUpload();

  return (
    <main className="min-h-screen bg-[#0f1117] p-8">
      <div className="max-w-4xl mx-auto">
        <header>
          <h1 className="text-4xl font-bold text-center mb-8 text-white">
            Receipt Processor
          </h1>
        </header>

        <UploadForm
          selectedFile={selectedFile}
          isProcessing={isProcessing}
          onFileChange={handleFileChange}
          onSubmit={handleSubmit}
        />

        {isProcessing && <LoadingSpinner />}

        {error && <ErrorDisplay message={error} />}

        {result && <ReceiptDetails receipt={result} />}
      </div>
    </main>
  );
}
