import React, { useState } from "react";
import { Button } from "../ui/Button";
import { Input } from "../ui/Input";

interface UrlFormProps {
  onSubmit: (url: string) => Promise<void>;
  loading?: boolean;
}

export const UrlForm: React.FC<UrlFormProps> = ({
  onSubmit,
  loading = false,
}) => {
  const [url, setUrl] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!url.trim()) {
      setError("URL is required");
      return;
    }

    try {
      await onSubmit(url.trim());
      setUrl("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to submit URL");
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex gap-4 mb-4">
      <Input
        type="url"
        value={url}
        onChange={(e) => setUrl(e.target.value)}
        placeholder="Enter website URL (e.g., https://example.com)"
        required
        disabled={loading}
        error={error}
        className="flex-1"
      />
      <Button
        type="submit"
        loading={loading}
        disabled={loading}
        className="flex-shrink-0"
      >
        {loading ? "Adding..." : "Crawl"}
      </Button>
    </form>
  );
};
