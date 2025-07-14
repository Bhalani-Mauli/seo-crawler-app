import React from "react";
import type { URLData } from "../../types";

interface CrawlDetailsHeaderProps {
  detailMeta: URLData | null;
  onClose: () => void;
}

export const CrawlDetailsHeader: React.FC<CrawlDetailsHeaderProps> = ({
  detailMeta,
  onClose,
}) => {
  return (
    <div className="p-6 border-b border-gray-200 flex-shrink-0">
      <div className="flex justify-between items-start">
        <div>
          <h3 className="text-2xl font-bold mb-2">Crawl Details</h3>
          {detailMeta && (
            <div>
              <div className="text-lg font-semibold text-gray-900">
                {detailMeta.url}
              </div>
              <div className="text-sm text-gray-500 mb-2">
                {detailMeta.crawl_data.title || "No title"}
              </div>
              <div className="flex flex-wrap gap-4 text-xs text-gray-700">
                <span>
                  Status: <b>{detailMeta.crawl_data.status}</b>
                </span>
                <span>
                  HTML Version: <b>{detailMeta.crawl_data.html_version}</b>
                </span>
                <span>
                  Internal Links: <b>{detailMeta.crawl_data.internal_links}</b>
                </span>
                <span>
                  External Links: <b>{detailMeta.crawl_data.external_links}</b>
                </span>
                <span>
                  Broken Links:{" "}
                  <b>{detailMeta.crawl_data.inaccessible_links}</b>
                </span>
                <span>
                  Login Form:{" "}
                  <b>{detailMeta.crawl_data.has_login_form ? "Yes" : "No"}</b>
                </span>
              </div>
            </div>
          )}
        </div>
        <button
          className="text-gray-500 hover:text-gray-700 text-2xl font-bold transition-colors"
          onClick={onClose}
          aria-label="Close modal"
        >
          &times;
        </button>
      </div>
    </div>
  );
};
