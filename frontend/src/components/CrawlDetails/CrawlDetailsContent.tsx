import React from "react";
import type { LinkData, HeadingData } from "../../types";
import { LinksTable } from "./LinksTable";
import { HeadingsTable } from "./HeadingsTable";

interface CrawlDetailsContentProps {
  detailLinks: LinkData[];
  detailHeadings: HeadingData[];
  detailLoading: boolean;
}

export const CrawlDetailsContent: React.FC<CrawlDetailsContentProps> = ({
  detailLinks,
  detailHeadings,
  detailLoading,
}) => {
  if (detailLoading) {
    return (
      <div className="h-full flex items-center justify-center text-gray-500">
        <div className="flex items-center">
          <svg
            className="animate-spin -ml-1 mr-3 h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            />
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
          Loading...
        </div>
      </div>
    );
  }

  return (
    <div className="h-full flex flex-col gap-6">
      <div className="h-1/2 min-h-0">
        <LinksTable links={detailLinks} />
      </div>
      <div className="h-1/2 min-h-0">
        <HeadingsTable headings={detailHeadings} />
      </div>
    </div>
  );
};
