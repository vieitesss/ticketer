"use client";

import { Spinner } from "@heroui/react";

export default function LoadingSpinner() {
  return (
    <section className="flex justify-center items-center py-12" role="status" aria-live="polite">
      <Spinner size="lg" label="Processing receipt..." color="primary" />
    </section>
  );
}
